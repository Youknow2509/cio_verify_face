package impl

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	applicationErrors "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/application/errors"
	model "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/application/model"
	constants "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/constants"
	domainModel "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/domain/model"
	"github.com/youknow2509/cio_verify_face/server/service_analytic/internal/domain/repository"
	"github.com/youknow2509/cio_verify_face/server/service_analytic/internal/global"
	cacheutil "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/shared/utils/cache"
)

// AnalyticServiceImpl implements IAnalyticService
type AnalyticServiceImpl struct {
	repo repository.IAnalyticRepository
}

// NewAnalyticService creates a new analytics service instance
func NewAnalyticService(repo repository.IAnalyticRepository) *AnalyticServiceImpl {
	return &AnalyticServiceImpl{
		repo: repo,
	}
}

// GetDailyReport returns daily attendance report
func (s *AnalyticServiceImpl) GetDailyReport(ctx context.Context, input *model.DailyReportInput) (*model.DailyReportOutput, *applicationErrors.Error) {
	// Log entry
	if global.Logger != nil {
		global.Logger.Info("GetDailyReport start",
			"company_id", safeStrPtr(input.CompanyID),
			"date", input.Date.Format("2006-01-02"),
			"device_id", safeStrPtr(input.DeviceID),
		)
	}

	// Authorization check
	if err := s.checkAuthorization(input.Session, input.CompanyID); err != nil {
		return nil, err
	}

	// Parse company ID (required for ScyllaDB queries)
	if input.CompanyID == nil {
		return nil, applicationErrors.ErrInvalidInput.WithDetails("company_id is required")
	}

	companyID, err := uuid.Parse(*input.CompanyID)
	if err != nil {
		return nil, applicationErrors.ErrInvalidInput.WithDetails("invalid company_id")
	}

	// Optional parse deviceID to uuid for cache key
	var deviceUUIDPtr *uuid.UUID
	if input.DeviceID != nil && *input.DeviceID != "" {
		du, perr := uuid.Parse(*input.DeviceID)
		if perr != nil {
			return nil, applicationErrors.ErrInvalidInput.WithDetails("invalid device_id")
		}
		deviceUUIDPtr = &du
	}

	// Cache key for daily report (includes device filter when provided)
	cacheKey := cacheutil.BuildDailyByDateKey(companyID, input.Date, deviceUUIDPtr)

	// Try local cache first
	if v, ok := cacheutil.GetLocal(cacheKey); ok {
		if out, ok2 := v.(*model.DailyReportOutput); ok2 {
			if global.Logger != nil {
				global.Logger.Debug("GetDailyReport cache hit (local)", "key", cacheKey)
			}
			return out, nil
		}
	}

	var cachedOut model.DailyReportOutput
	if hit, _ := cacheutil.GetDistributed(ctx, cacheKey, &cachedOut); hit {
		// backfill local cache with slightly shorter TTL than distributed
		cacheutil.SetLocal(cacheKey, &cachedOut, localTTLFrom(constants.CacheTTLMidSeconds))
		if global.Logger != nil {
			global.Logger.Debug("GetDailyReport cache hit (redis)", "key", cacheKey)
		}
		return &cachedOut, nil
	}

	// Get daily summaries from ScyllaDB
	summaries, err := s.repo.GetDailySummariesByDate(ctx, companyID, input.Date)
	if err != nil {
		if global.Logger != nil {
			global.Logger.Error("GetDailySummariesByDate failed", "error", err.Error())
		}
		return nil, applicationErrors.ErrDatabaseError.WithDetails(err.Error())
	}

	// Optional: filter by device_id if provided
	if deviceUUIDPtr != nil {
		empIDs, rerr := s.repo.GetEmployeeIDsByDeviceAndDate(ctx, *deviceUUIDPtr, input.Date)
		if rerr != nil {
			if global.Logger != nil {
				global.Logger.Error("GetEmployeeIDsByDeviceAndDate failed", "error", rerr.Error())
			}
			return nil, applicationErrors.ErrDatabaseError.WithDetails(rerr.Error())
		}
		if len(empIDs) == 0 {
			// Short-circuit: no employees found for this device/date
			out := &model.DailyReportOutput{
				Date:                input.Date.Format("2006-01-02"),
				TotalEmployees:      0,
				PresentEmployees:    0,
				LateEmployees:       0,
				EarlyLeaveEmployees: 0,
				AbsentEmployees:     0,
				AttendanceRate:      0,
				Departments:         []model.DepartmentReport{},
				Shifts:              []model.ShiftReport{},
			}
			// set caches
			_ = cacheutil.SetDistributed(ctx, cacheKey, out, time.Duration(constants.CacheTTLMidSeconds)*time.Second)
			_ = cacheutil.SetLocal(cacheKey, out, localTTLFrom(constants.CacheTTLMidSeconds))
			if global.Logger != nil {
				global.Logger.Info("GetDailyReport computed (no employees for device)", "key", cacheKey)
			}
			return out, nil
		}
		// Build a set for quick lookup
		empSet := make(map[uuid.UUID]struct{}, len(empIDs))
		for _, id := range empIDs {
			empSet[id] = struct{}{}
		}
		// Filter summaries to those employees only
		filtered := make([]*domainModel.DailySummary, 0, len(summaries))
		for _, ssum := range summaries {
			if _, ok := empSet[ssum.EmployeeID]; ok {
				filtered = append(filtered, ssum)
			}
		}
		summaries = filtered
	}

	// Calculate statistics
	totalEmployees := len(summaries)
	presentEmployees := 0
	lateEmployees := 0
	earlyLeaveEmployees := 0
	absentEmployees := 0

	for _, summary := range summaries {
		switch summary.AttendanceStatus {
		case 0: // PRESENT
			presentEmployees++
		case 1: // LATE
			lateEmployees++
		case 2: // EARLY_LEAVE
			earlyLeaveEmployees++
		case 3: // ABSENT
			absentEmployees++
		}
	}

	attendanceRate := 0.0
	if totalEmployees > 0 {
		attendanceRate = float64(presentEmployees) / float64(totalEmployees) * 100
	}

	// Group by departments
	departments := s.groupByDepartment(ctx, summaries)

	// Group by shifts
	shifts := s.groupByShift(ctx, summaries)

	out := &model.DailyReportOutput{
		Date:                input.Date.Format("2006-01-02"),
		TotalEmployees:      totalEmployees,
		PresentEmployees:    presentEmployees,
		LateEmployees:       lateEmployees,
		EarlyLeaveEmployees: earlyLeaveEmployees,
		AbsentEmployees:     absentEmployees,
		AttendanceRate:      roundFloat(attendanceRate, 2),
		Departments:         departments,
		Shifts:              shifts,
	}
	// store in caches (local TTL < distributed TTL)
	_ = cacheutil.SetDistributed(ctx, cacheKey, out, time.Duration(constants.CacheTTLMidSeconds)*time.Second)
	_ = cacheutil.SetLocal(cacheKey, out, localTTLFrom(constants.CacheTTLMidSeconds))
	if global.Logger != nil {
		global.Logger.Info("GetDailyReport computed", "key", cacheKey, "total_employees", totalEmployees)
	}
	return out, nil
}

// GetSummaryReport returns monthly summary report
func (s *AnalyticServiceImpl) GetSummaryReport(ctx context.Context, input *model.SummaryReportInput) (*model.SummaryReportOutput, *applicationErrors.Error) {
	// Authorization check
	if err := s.checkAuthorization(input.Session, input.CompanyID); err != nil {
		return nil, err
	}

	if global.Logger != nil {
		global.Logger.Info("GetSummaryReport start", "company_id", safeStrPtr(input.CompanyID), "month", input.Month)
	}
	startDate, endDate, err := parseMonth(input.Month)
	if err != nil {
		return nil, applicationErrors.ErrInvalidDateFormat.WithDetails("use YYYY-MM format")
	}
	if input.CompanyID == nil {
		return nil, applicationErrors.ErrInvalidInput.WithDetails("company_id is required")
	}
	companyID, err := uuid.Parse(*input.CompanyID)
	if err != nil {
		return nil, applicationErrors.ErrInvalidInput.WithDetails("invalid company_id")
	}
	monthKey := cacheutil.BuildDailyByMonthKey(companyID, input.Month)
	if v, ok := cacheutil.GetLocal(monthKey); ok {
		if out, ok2 := v.(*model.SummaryReportOutput); ok2 {
			if global.Logger != nil {
				global.Logger.Debug("GetSummaryReport cache hit (local)", "key", monthKey)
			}
			return out, nil
		}
	}
	var cachedMonth model.SummaryReportOutput
	if hit, _ := cacheutil.GetDistributed(ctx, monthKey, &cachedMonth); hit {
		cacheutil.SetLocal(monthKey, &cachedMonth, localTTLFrom(constants.CacheTTLLongSeconds))
		if global.Logger != nil {
			global.Logger.Debug("GetSummaryReport cache hit (redis)", "key", monthKey)
		}
		return &cachedMonth, nil
	}
	summaries, err := s.repo.GetDailySummariesByMonth(ctx, companyID, input.Month)
	if err != nil {
		if global.Logger != nil {
			global.Logger.Error("GetDailySummariesByMonth failed", "error", err.Error())
		}
		return nil, applicationErrors.ErrDatabaseError.WithDetails(err.Error())
	}
	totalWorkingDays := endDate.Day()
	totalEmployees64, err := s.repo.GetTotalEmployees(ctx, &companyID)
	if err != nil {
		if global.Logger != nil {
			global.Logger.Error("GetTotalEmployees failed", "error", err.Error())
		}
		return nil, applicationErrors.ErrDatabaseError.WithDetails(err.Error())
	}
	totalEmployees := int(totalEmployees64)
	totalPresentDays := 0
	totalWorkingMinutes := 0
	totalOvertimeMinutes := 0
	for _, summary := range summaries {
		if summary.AttendanceStatus == 0 {
			totalPresentDays++
		}
		totalWorkingMinutes += summary.TotalWorkMinutes
		totalOvertimeMinutes += summary.OvertimeMinutes
	}
	averageAttendanceRate := 0.0
	if totalEmployees > 0 && totalWorkingDays > 0 {
		averageAttendanceRate = float64(totalPresentDays) / float64(totalEmployees*totalWorkingDays) * 100
	}
	weeklySummary := s.calculateWeeklySummary(ctx, startDate, endDate, companyID, totalEmployees)
	topEmployees := s.getTopAttendanceEmployees(ctx, summaries, 10)
	lowEmployees := s.getLowAttendanceEmployees(ctx, summaries, 10)
	out := &model.SummaryReportOutput{
		Month:                  input.Month,
		TotalWorkingDays:       totalWorkingDays,
		TotalEmployees:         totalEmployees,
		AverageAttendanceRate:  roundFloat(averageAttendanceRate, 2),
		TotalWorkingHours:      totalWorkingMinutes / 60,
		TotalOvertimeHours:     totalOvertimeMinutes / 60,
		WeeklySummary:          weeklySummary,
		TopAttendanceEmployees: topEmployees,
		LowAttendanceEmployees: lowEmployees,
	}
	_ = cacheutil.SetDistributed(ctx, monthKey, out, time.Duration(constants.CacheTTLLongSeconds)*time.Second)
	_ = cacheutil.SetLocal(monthKey, out, localTTLFrom(constants.CacheTTLLongSeconds))
	if global.Logger != nil {
		global.Logger.Info("GetSummaryReport computed", "key", monthKey, "total_employees", totalEmployees)
	}
	return out, nil
}

// ExportReport exports attendance report to file
func (s *AnalyticServiceImpl) ExportReport(ctx context.Context, input *model.ExportReportInput) (*model.ExportReportOutput, *applicationErrors.Error) {
	// Authorization check
	if err := s.checkAuthorization(input.Session, input.CompanyID); err != nil {
		return nil, err
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		return nil, applicationErrors.ErrInvalidDateFormat.WithDetails("start_date must be YYYY-MM-DD")
	}

	endDate, err := time.Parse("2006-01-02", input.EndDate)
	if err != nil {
		return nil, applicationErrors.ErrInvalidDateFormat.WithDetails("end_date must be YYYY-MM-DD")
	}

	if startDate.After(endDate) {
		return nil, applicationErrors.ErrInvalidDateRange.WithDetails("start_date must be before end_date")
	}
	// Company required
	if input.CompanyID == nil || *input.CompanyID == "" {
		return nil, applicationErrors.ErrInvalidInput.WithDetails("company_id is required")
	}
	companyID, err := uuid.Parse(*input.CompanyID)
	if err != nil {
		return nil, applicationErrors.ErrInvalidInput.WithDetails("invalid company_id")
	}

	// Normalize format (excel -> csv)
	exportFormat := input.Format
	if exportFormat == "excel" {
		exportFormat = "csv"
	}

	// Build cache key for export
	exportKey := cacheutil.BuildExportKey(companyID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), exportFormat)

	// Attempt cache hit: if object storage used previously, re-presign fresh URL
	type exportCacheEntry struct {
		Storage   string `json:"storage"`              // "object" or "local"
		ObjectKey string `json:"object_key,omitempty"` // for object storage
		LocalPath string `json:"local_path,omitempty"` // for local fallback
	}
	// Try local cache first
	if v, ok := cacheutil.GetLocal(exportKey); ok {
		if e, ok2 := v.(*exportCacheEntry); ok2 {
			// Serve from cache
			if e.Storage == "object" {
				objCfg := global.SettingServer.ObjectStorage
				if objCfg.Endpoint != "" && objCfg.Bucket != "" {
					cli, cerr := minio.New(objCfg.Endpoint, &minio.Options{
						Creds:  credentials.NewStaticV4(objCfg.AccessKey, objCfg.SecretKey, ""),
						Secure: objCfg.UseSSL,
						Region: objCfg.Region,
					})
					if cerr == nil {
						expireMinutes := objCfg.PresignExpireMinutes
						if expireMinutes <= 0 {
							expireMinutes = 60
						}
						presigned, perr := cli.PresignedGetObject(ctx, objCfg.Bucket, e.ObjectKey, time.Duration(expireMinutes)*time.Minute, nil)
						if perr == nil {
							download := presigned.String()
							if input.Email != nil && *input.Email != "" {
								if nerr := s.publishExportEmail(ctx, *input.Email, download, exportFormat, startDate, endDate, *input.CompanyID); nerr != nil {
									return nil, applicationErrors.ErrExportFailed.WithDetails("notify email failed: " + nerr.Error())
								}
							}
							return &model.ExportReportOutput{JobID: "cached", Status: "completed", Message: "Export served from cache", DownloadURL: &download}, nil
						}
					}
				}
			} else if e.Storage == "local" {
				if e.LocalPath != "" {
					if _, statErr := os.Stat(e.LocalPath); statErr == nil {
						download := e.LocalPath
						if input.Email != nil && *input.Email != "" {
							if nerr := s.publishExportEmail(ctx, *input.Email, download, exportFormat, startDate, endDate, *input.CompanyID); nerr != nil {
								return nil, applicationErrors.ErrExportFailed.WithDetails("notify email failed: " + nerr.Error())
							}
						}
						return &model.ExportReportOutput{JobID: "cached", Status: "completed", Message: "Export served from cache", DownloadURL: &download}, nil
					}
				}
			}
		}
	}
	// Try distributed cache (Redis)
	var ec exportCacheEntry
	if hit, _ := cacheutil.GetDistributed(ctx, exportKey, &ec); hit {
		// backfill local cache with shorter TTL
		cacheutil.SetLocal(exportKey, &ec, localTTLFrom(constants.CacheTTLLongSeconds))
		if ec.Storage == "object" {
			objCfg := global.SettingServer.ObjectStorage
			if objCfg.Endpoint != "" && objCfg.Bucket != "" {
				cli, cerr := minio.New(objCfg.Endpoint, &minio.Options{
					Creds:  credentials.NewStaticV4(objCfg.AccessKey, objCfg.SecretKey, ""),
					Secure: objCfg.UseSSL,
					Region: objCfg.Region,
				})
				if cerr == nil {
					expireMinutes := objCfg.PresignExpireMinutes
					if expireMinutes <= 0 {
						expireMinutes = 60
					}
					presigned, perr := cli.PresignedGetObject(ctx, objCfg.Bucket, ec.ObjectKey, time.Duration(expireMinutes)*time.Minute, nil)
					if perr == nil {
						download := presigned.String()
						if input.Email != nil && *input.Email != "" {
							if nerr := s.publishExportEmail(ctx, *input.Email, download, exportFormat, startDate, endDate, *input.CompanyID); nerr != nil {
								return nil, applicationErrors.ErrExportFailed.WithDetails("notify email failed: " + nerr.Error())
							}
						}
						return &model.ExportReportOutput{JobID: "cached", Status: "completed", Message: "Export served from cache", DownloadURL: &download}, nil
					}
				}
			}
		} else if ec.Storage == "local" && ec.LocalPath != "" {
			if _, statErr := os.Stat(ec.LocalPath); statErr == nil {
				download := ec.LocalPath
				if input.Email != nil && *input.Email != "" {
					if nerr := s.publishExportEmail(ctx, *input.Email, download, exportFormat, startDate, endDate, *input.CompanyID); nerr != nil {
						return nil, applicationErrors.ErrExportFailed.WithDetails("notify email failed: " + nerr.Error())
					}
				}
				return &model.ExportReportOutput{JobID: "cached", Status: "completed", Message: "Export served from cache", DownloadURL: &download}, nil
			}
		}
	}

	// Fetch data from ScyllaDB (cache miss)
	summaries, derr := s.repo.GetDailySummariesByDateRange(ctx, companyID, startDate, endDate)
	if derr != nil {
		return nil, applicationErrors.ErrDatabaseError.WithDetails(derr.Error())
	}

	// Ensure export directory exists
	exportDir := "exports"
	if err := os.MkdirAll(exportDir, 0o755); err != nil {
		return nil, applicationErrors.ErrExportFailed.WithDetails("failed to create export directory")
	}
	jobID := fmt.Sprintf("export_%d_%s", time.Now().Unix(), uuid.New().String()[:8])
	baseName := fmt.Sprintf("%s_%s_to_%s", jobID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

	var filePath string
	// Note: we only remove the local file when uploaded to object storage
	switch input.Format {
	case "csv":
		filePath = filepath.Join(exportDir, baseName+".csv")
		if werr := s.writeCSV(filePath, summaries); werr != nil {
			return nil, applicationErrors.ErrExportFailed.WithDetails(werr.Error())
		}
	case "excel":
		// Default to CSV when Excel lib not available; still generate .csv
		filePath = filepath.Join(exportDir, baseName+".csv")
		if werr := s.writeCSV(filePath, summaries); werr != nil {
			return nil, applicationErrors.ErrExportFailed.WithDetails(werr.Error())
		}
	case "pdf":
		return nil, applicationErrors.ErrInvalidInput.WithDetails("pdf export not supported yet")
	default:
		return nil, applicationErrors.ErrInvalidInput.WithDetails("unsupported format; use excel or csv")
	}

	// Upload to object storage (MinIO/S3-compatible)
	objCfg := global.SettingServer.ObjectStorage
	if objCfg.Endpoint != "" && objCfg.Bucket != "" {
		cli, cerr := minio.New(objCfg.Endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(objCfg.AccessKey, objCfg.SecretKey, ""),
			Secure: objCfg.UseSSL,
			Region: objCfg.Region,
		})
		if cerr != nil {
			return nil, applicationErrors.ErrExportFailed.WithDetails("init object storage client failed: " + cerr.Error())
		}

		// Ensure bucket exists
		exists, eerr := cli.BucketExists(ctx, objCfg.Bucket)
		if eerr != nil {
			return nil, applicationErrors.ErrExportFailed.WithDetails("check bucket failed: " + eerr.Error())
		}
		if !exists {
			if cberr := cli.MakeBucket(ctx, objCfg.Bucket, minio.MakeBucketOptions{Region: objCfg.Region}); cberr != nil {
				return nil, applicationErrors.ErrExportFailed.WithDetails("create bucket failed: " + cberr.Error())
			}
		}

		objectKey := fmt.Sprintf("reports/%s/%s_%s_%s.%s", time.Now().Format("2006/01/02"), companyID.String(), startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), exportFormat)
		contentType := "text/csv"
		if _, uperr := cli.FPutObject(ctx, objCfg.Bucket, objectKey, filePath, minio.PutObjectOptions{ContentType: contentType}); uperr != nil {
			return nil, applicationErrors.ErrExportFailed.WithDetails("upload to object storage failed: " + uperr.Error())
		}

		// Generate presigned URL with expiry
		expireMinutes := objCfg.PresignExpireMinutes
		if expireMinutes <= 0 {
			expireMinutes = 60
		}
		presigned, perr := cli.PresignedGetObject(ctx, objCfg.Bucket, objectKey, time.Duration(expireMinutes)*time.Minute, nil)
		if perr != nil {
			return nil, applicationErrors.ErrExportFailed.WithDetails("presign url failed: " + perr.Error())
		}

		download := presigned.String()

		// Cache the export location (object storage)
		entry := &exportCacheEntry{Storage: "object", ObjectKey: objectKey}
		_ = cacheutil.SetDistributed(ctx, exportKey, entry, time.Duration(constants.CacheTTLLongSeconds)*time.Second)
		_ = cacheutil.SetLocal(exportKey, entry, localTTLFrom(constants.CacheTTLLongSeconds))

		// If email provided, publish Kafka notification message
		if input.Email != nil && *input.Email != "" {
			if nerr := s.publishExportEmail(ctx, *input.Email, download, input.Format, startDate, endDate, *input.CompanyID); nerr != nil {
				return nil, applicationErrors.ErrExportFailed.WithDetails("notify email failed: " + nerr.Error())
			}
		}
		// Remove local temp after successful upload
		_ = os.Remove(filePath)
		return &model.ExportReportOutput{
			JobID:       jobID,
			Status:      "completed",
			Message:     fmt.Sprintf("Exported %d rows to object storage", len(summaries)),
			DownloadURL: &download,
		}, nil
	}

	// Fallback: local file path
	// Fallback: local file path; cache the local path for reuse
	download := filePath
	entry := &exportCacheEntry{Storage: "local", LocalPath: filePath}
	_ = cacheutil.SetDistributed(ctx, exportKey, entry, time.Duration(constants.CacheTTLMidSeconds)*time.Second)
	_ = cacheutil.SetLocal(exportKey, entry, localTTLFrom(constants.CacheTTLMidSeconds))
	if input.Email != nil && *input.Email != "" {
		if nerr := s.publishExportEmail(ctx, *input.Email, download, input.Format, startDate, endDate, *input.CompanyID); nerr != nil {
			return nil, applicationErrors.ErrExportFailed.WithDetails("notify email failed: " + nerr.Error())
		}
	}
	return &model.ExportReportOutput{
		JobID:       jobID,
		Status:      "completed",
		Message:     fmt.Sprintf("Exported %d rows to %s", len(summaries), filePath),
		DownloadURL: &download,
	}, nil
}

// GetHealthCheck returns service health status
func (s *AnalyticServiceImpl) GetHealthCheck(ctx context.Context) (*model.HealthCheckOutput, *applicationErrors.Error) {
	// Check database connectivity
	dbStatus := "healthy"
	_, err := s.repo.GetTotalEmployees(ctx, nil)
	if err != nil {
		dbStatus = "unhealthy"
	}

	return &model.HealthCheckOutput{
		Status:  "healthy",
		Version: "1.0.0",
		Services: map[string]interface{}{
			"database": dbStatus,
			"cache":    "healthy",
		},
	}, nil
}

// Helper functions

func (s *AnalyticServiceImpl) groupByDepartment(ctx context.Context, summaries []*domainModel.DailySummary) []model.DepartmentReport {
	departmentMap := make(map[string]*model.DepartmentReport)

	for _, summary := range summaries {
		// Get employee details
		employee, err := s.repo.GetEmployeeByID(ctx, summary.EmployeeID)
		if err != nil {
			continue
		}

		deptName := "Unknown"
		if employee.Department != nil {
			deptName = *employee.Department
		}

		if _, exists := departmentMap[deptName]; !exists {
			departmentMap[deptName] = &model.DepartmentReport{
				DepartmentName:   deptName,
				TotalEmployees:   0,
				PresentEmployees: 0,
			}
		}

		dept := departmentMap[deptName]
		dept.TotalEmployees++
		if summary.AttendanceStatus == 0 { // PRESENT
			dept.PresentEmployees++
		}
	}

	// Calculate attendance rate for each department
	result := make([]model.DepartmentReport, 0, len(departmentMap))
	for _, dept := range departmentMap {
		if dept.TotalEmployees > 0 {
			dept.AttendanceRate = roundFloat(float64(dept.PresentEmployees)/float64(dept.TotalEmployees)*100, 2)
		}
		result = append(result, *dept)
	}

	return result
}

func (s *AnalyticServiceImpl) groupByShift(ctx context.Context, summaries []*domainModel.DailySummary) []model.ShiftReport {
	shiftMap := make(map[uuid.UUID]*model.ShiftReport)

	for _, summary := range summaries {
		if summary.ShiftID == uuid.Nil {
			continue
		}

		if _, exists := shiftMap[summary.ShiftID]; !exists {
			// Get shift details
			shift, err := s.repo.GetWorkShiftByID(ctx, summary.ShiftID)
			if err != nil {
				continue
			}

			shiftMap[summary.ShiftID] = &model.ShiftReport{
				ShiftName:        shift.Name,
				StartTime:        shift.StartTime,
				EndTime:          shift.EndTime,
				TotalEmployees:   0,
				PresentEmployees: 0,
			}
		}

		shiftRpt := shiftMap[summary.ShiftID]
		shiftRpt.TotalEmployees++
		if summary.AttendanceStatus == 0 { // PRESENT
			shiftRpt.PresentEmployees++
		}
	}

	// Calculate attendance rate for each shift
	result := make([]model.ShiftReport, 0, len(shiftMap))
	for _, shift := range shiftMap {
		if shift.TotalEmployees > 0 {
			shift.AttendanceRate = roundFloat(float64(shift.PresentEmployees)/float64(shift.TotalEmployees)*100, 2)
		}
		result = append(result, *shift)
	}

	return result
}

func (s *AnalyticServiceImpl) calculateWeeklySummary(ctx context.Context, startDate, endDate time.Time, companyID uuid.UUID, totalEmployees int) []model.WeeklySummary {
	weeklySummaries := []model.WeeklySummary{}
	currentDate := startDate
	weekNumber := 1

	for currentDate.Before(endDate) || currentDate.Equal(endDate) {
		weekStart := currentDate
		weekEnd := currentDate.AddDate(0, 0, 6)
		if weekEnd.After(endDate) {
			weekEnd = endDate
		}

		// Get attendance data for this week from ScyllaDB
		weekSummaries, err := s.repo.GetDailySummariesByDateRange(ctx, companyID, weekStart, weekEnd)
		if err != nil {
			currentDate = currentDate.AddDate(0, 0, 7)
			weekNumber++
			continue
		}

		totalPresentDays := 0
		totalMinutes := 0
		for _, summary := range weekSummaries {
			if summary.AttendanceStatus == 0 { // PRESENT
				totalPresentDays++
			}
			totalMinutes += summary.TotalWorkMinutes
		}

		attendanceRate := 0.0
		if totalEmployees > 0 {
			attendanceRate = float64(totalPresentDays) / float64(totalEmployees*7) * 100
		}

		weeklySummaries = append(weeklySummaries, model.WeeklySummary{
			Week:           weekNumber,
			StartDate:      weekStart.Format("2006-01-02"),
			EndDate:        weekEnd.Format("2006-01-02"),
			AttendanceRate: roundFloat(attendanceRate, 2),
			TotalHours:     totalMinutes / 60,
		})

		currentDate = currentDate.AddDate(0, 0, 7)
		weekNumber++
	}

	return weeklySummaries
}

func (s *AnalyticServiceImpl) getTopAttendanceEmployees(ctx context.Context, summaries []*domainModel.DailySummary, limit int) []model.EmployeeAttendanceStat {
	employeeMap := make(map[uuid.UUID]*model.EmployeeAttendanceStat)

	for _, summary := range summaries {
		if _, exists := employeeMap[summary.EmployeeID]; !exists {
			employee, err := s.repo.GetEmployeeByID(ctx, summary.EmployeeID)
			if err != nil {
				continue
			}

			user, err := s.repo.GetUserByID(ctx, summary.EmployeeID)
			if err != nil {
				continue
			}

			employeeMap[summary.EmployeeID] = &model.EmployeeAttendanceStat{
				EmployeeCode:   employee.EmployeeCode,
				FullName:       user.FullName,
				PresentDays:    0,
				AttendanceRate: 0,
				TotalHours:     0,
			}
		}

		emp := employeeMap[summary.EmployeeID]
		if summary.AttendanceStatus == 0 { // PRESENT
			emp.PresentDays++
		}
		emp.TotalHours += summary.TotalWorkMinutes / 60
	}

	// Convert to slice and sort by present days (descending)
	result := make([]model.EmployeeAttendanceStat, 0, len(employeeMap))
	for _, emp := range employeeMap {
		emp.AttendanceRate = roundFloat(float64(emp.PresentDays)/30*100, 2) // Assuming 30 days
		result = append(result, *emp)
	}

	// Sort by present days descending
	for i := 0; i < len(result)-1; i++ {
		for j := i + 1; j < len(result); j++ {
			if result[j].PresentDays > result[i].PresentDays {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	if len(result) > limit {
		result = result[:limit]
	}

	return result
}

func (s *AnalyticServiceImpl) getLowAttendanceEmployees(ctx context.Context, summaries []*domainModel.DailySummary, limit int) []model.EmployeeAttendanceStat {
	topEmployees := s.getTopAttendanceEmployees(ctx, summaries, len(summaries))

	// Sort by present days ascending
	for i := 0; i < len(topEmployees)-1; i++ {
		for j := i + 1; j < len(topEmployees); j++ {
			if topEmployees[j].PresentDays < topEmployees[i].PresentDays {
				topEmployees[i], topEmployees[j] = topEmployees[j], topEmployees[i]
			}
		}
	}

	if len(topEmployees) > limit {
		topEmployees = topEmployees[:limit]
	}

	return topEmployees
}

// Utility functions

func parseMonth(month string) (time.Time, time.Time, error) {
	startDate, err := time.Parse("2006-01", month)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	// Last day of month
	endDate := startDate.AddDate(0, 1, -1)
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 0, endDate.Location())

	return startDate, endDate, nil
}

func roundFloat(val float64, precision int) float64 {
	ratio := 1.0
	for i := 0; i < precision; i++ {
		ratio *= 10
	}
	return float64(int(val*ratio+0.5)) / ratio
}

// writeCSV writes daily summaries to a CSV file
func (s *AnalyticServiceImpl) writeCSV(path string, summaries []*domainModel.DailySummary) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	// Header
	if err := w.Write([]string{
		"work_date",
		"employee_id",
		"shift_id",
		"total_work_minutes",
		"overtime_minutes",
		"late_minutes",
		"early_leave_minutes",
		"attendance_status",
		"attendance_percentage",
	}); err != nil {
		return err
	}

	// Rows
	for _, s := range summaries {
		rec := []string{
			s.WorkDate.Format("2006-01-02"),
			s.EmployeeID.String(),
			s.ShiftID.String(),
			fmt.Sprintf("%d", s.TotalWorkMinutes),
			fmt.Sprintf("%d", s.OvertimeMinutes),
			fmt.Sprintf("%d", s.LateMinutes),
			fmt.Sprintf("%d", s.EarlyLeaveMinutes),
			fmt.Sprintf("%d", s.AttendanceStatus),
			fmt.Sprintf("%.2f", s.AttendancePercentage),
		}
		if err := w.Write(rec); err != nil {
			return err
		}
	}
	return nil
}

// publishExportEmail publishes a message to Kafka for service_notify to send email
func (s *AnalyticServiceImpl) publishExportEmail(ctx context.Context, email, url, format string, start, end time.Time, companyID string) error {
	kcfg := global.SettingServer.Kafka
	if len(kcfg.Brokers) == 0 || kcfg.NotifyTopic == "" {
		return fmt.Errorf("kafka not configured")
	}

	cfg := sarama.NewConfig()
	cfg.ClientID = kcfg.ClientID
	cfg.Producer.Return.Successes = true
	cfg.Version = sarama.V2_8_0_0
	if kcfg.SASLEnabled {
		cfg.Net.SASL.Enable = true
		cfg.Net.SASL.User = kcfg.SASLUser
		cfg.Net.SASL.Password = kcfg.SASLPassword
		cfg.Net.SASL.Mechanism = sarama.SASLTypePlaintext
	}
	if kcfg.TLSEnabled {
		cfg.Net.TLS.Enable = true
	}

	prod, err := sarama.NewSyncProducer(kcfg.Brokers, cfg)
	if err != nil {
		return err
	}
	defer prod.Close()

	// Build notify payload
	payload := map[string]interface{}{
		"type":         "export_report",
		"email":        email,
		"download_url": url,
		"format":       format,
		"company_id":   companyID,
		"start_date":   start.Format("2006-01-02"),
		"end_date":     end.Format("2006-01-02"),
		"created_at":   time.Now().UTC().Format(time.RFC3339),
	}
	b, _ := jsonMarshal(payload)

	msg := &sarama.ProducerMessage{
		Topic: kcfg.NotifyTopic,
		Value: sarama.ByteEncoder(b),
	}
	_, _, err = prod.SendMessage(msg)
	return err
}

// checkAuthorization checks if the user has permission to access the requested company data
func (s *AnalyticServiceImpl) checkAuthorization(session *model.SessionInfo, requestedCompanyID *string) *applicationErrors.Error {
	if session == nil {
		return applicationErrors.ErrUnauthorized.WithDetails("session info required")
	}

	if requestedCompanyID == nil {
		return applicationErrors.ErrInvalidInput.WithDetails("company_id is required")
	}

	// System admin (root) has full access
	if session.Role == int32(domainModel.RoleSystemAdmin) {
		return nil
	}

	// Company admin/manager can only access their own company data
	if session.Role == int32(domainModel.RoleCompanyAdmin) {
		// Check if company_id in session matches requested company_id
		if session.CompanyID != *requestedCompanyID {
			if global.Logger != nil {
				global.Logger.Warn("Authorization failed: company mismatch",
					"user_id", session.UserID,
					"session_company_id", session.CompanyID,
					"requested_company_id", *requestedCompanyID)
			}
			return applicationErrors.ErrForbidden.WithDetails("access denied: you can only access your own company data")
		}
		return nil
	}

	// Employees (role 0) should not access analytics
	return applicationErrors.ErrForbidden.WithDetails("access denied: insufficient permissions")
}

// jsonMarshal keeps minimal deps by using stdlib encoding/json via indirection
func jsonMarshal(v interface{}) ([]byte, error) {
	type jm = interface{}
	return json.Marshal(jm(v))
}

// safeStrPtr returns empty string when pointer is nil
func safeStrPtr(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

// localTTLFrom returns a TTL strictly less than the distributed TTL
func localTTLFrom(distSeconds int) time.Duration {
	dist := time.Duration(distSeconds) * time.Second
	skew := 10 * time.Second
	if dist > skew {
		return dist - skew
	}
	if dist > time.Second {
		return dist / 2
	}
	return time.Second
}
