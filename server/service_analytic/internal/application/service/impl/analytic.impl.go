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
	"github.com/youknow2509/cio_verify_face/server/service_analytic/internal/application/service"
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

// ExportDailyReportDetail implements service.IAnalyticService.
func (s *AnalyticServiceImpl) ExportDailyReportDetail(ctx context.Context, input *model.ExportDailyReportDetailInput) (*model.ExportDailyReportDetailOutput, *applicationErrors.Error) {
	// Authorization check: Require CompanyAdmin role minimum, no employee self-access for company-wide reports
	requestedCompanyID := input.CompanyID.String()
	if err := s.checkAuthorization(input.Session, &requestedCompanyID, domainModel.RoleCompanyAdmin, false); err != nil {
		return nil, err
	}

	// Validate date
	if input.Date.IsZero() {
		return nil, applicationErrors.ErrInvalidInput.WithDetails("date is required")
	}

	// Validate MinIO configuration (required for export)
	objCfg := global.SettingServer.ObjectStorage
	if objCfg.Endpoint == "" || objCfg.Bucket == "" {
		if global.Logger != nil {
			global.Logger.Error("ExportDailyReportDetail: Object storage not configured")
		}
		return nil, applicationErrors.ErrExportFailed.WithDetails("Object storage (MinIO/S3) must be configured for export functionality")
	}

	// Normalize format (excel -> csv, pdf stays as pdf)
	exportFormat := input.Format
	if exportFormat == "excel" {
		exportFormat = "csv"
	}

	// Build cache key for export
	exportKey := cacheutil.BuildExportKey(input.CompanyID, input.Date.Format("2006-01-02"), input.Date.Format("2006-01-02"), exportFormat+"_detail")

	type exportCacheEntry struct {
		JobID     string `json:"job_id,omitempty"`
		Status    string `json:"status"`  // processing | completed
		Storage   string `json:"storage"` // object | local
		ObjectKey string `json:"object_key,omitempty"`
		LocalPath string `json:"local_path,omitempty"`
		Format    string `json:"format"`
		Rows      int    `json:"rows"`
	}

	// Check Redis cache first
	var ec exportCacheEntry
	if hit, _ := cacheutil.GetDistributedOnly(ctx, exportKey, &ec); hit {
		if ec.Status == "processing" {
			jobID := ec.JobID
			if jobID == "" {
				jobID = "processing"
			}
			return &model.ExportDailyReportDetailOutput{JobID: jobID, Status: "processing", Message: "Export is processing"}, nil
		}
		if ec.Status == "completed" {
			if ec.Storage == "object" && ec.ObjectKey != "" {
				objCfg := global.SettingServer.ObjectStorage
				if objCfg.Endpoint != "" && objCfg.Bucket != "" {
					cli, cerr := minio.New(objCfg.Endpoint, &minio.Options{Creds: credentials.NewStaticV4(objCfg.AccessKey, objCfg.SecretKey, ""), Secure: objCfg.UseSSL, Region: objCfg.Region})
					if cerr == nil {
						// Generate 7-day presigned URL
						expireDuration := 7 * 24 * time.Hour // 7 days
						presigned, perr := cli.PresignedGetObject(ctx, objCfg.Bucket, ec.ObjectKey, expireDuration, nil)
						if perr == nil {
							urlStr := presigned.String()
							return &model.ExportDailyReportDetailOutput{
								JobID:       ec.JobID,
								Status:      "completed",
								Message:     fmt.Sprintf("Export completed with %d rows (cached, expires in 7 days)", ec.Rows),
								DownloadURL: &urlStr,
							}, nil
						}
					}
				}
			}
		}
	}

	// Prepare object key
	objectKey := fmt.Sprintf("reports/%s/%s_daily_detail_%s.%s", time.Now().Format("2006/01/02"), input.CompanyID.String(), input.Date.Format("2006-01-02"), exportFormat)

	// Check if object exists in storage
	cli, cerr := minio.New(objCfg.Endpoint, &minio.Options{Creds: credentials.NewStaticV4(objCfg.AccessKey, objCfg.SecretKey, ""), Secure: objCfg.UseSSL, Region: objCfg.Region})
	if cerr != nil {
		if global.Logger != nil {
			global.Logger.Error("ExportDailyReportDetail: Failed to connect to object storage", "error", cerr.Error())
		}
		return nil, applicationErrors.ErrExportFailed.WithDetails(fmt.Sprintf("Failed to connect to object storage: %v", cerr))
	}

	// Ensure bucket exists
	exists, eerr := cli.BucketExists(ctx, objCfg.Bucket)
	if eerr != nil {
		if global.Logger != nil {
			global.Logger.Error("ExportDailyReportDetail: Failed to check bucket", "error", eerr.Error())
		}
		return nil, applicationErrors.ErrExportFailed.WithDetails(fmt.Sprintf("Failed to check bucket: %v", eerr))
	}
	if !exists {
		if err := cli.MakeBucket(ctx, objCfg.Bucket, minio.MakeBucketOptions{Region: objCfg.Region}); err != nil {
			if global.Logger != nil {
				global.Logger.Error("ExportDailyReportDetail: Failed to create bucket", "error", err.Error())
			}
			return nil, applicationErrors.ErrExportFailed.WithDetails(fmt.Sprintf("Failed to create bucket: %v", err))
		}
	}

	// Try stat existing object to avoid regeneration
	if _, statErr := cli.StatObject(ctx, objCfg.Bucket, objectKey, minio.StatObjectOptions{}); statErr == nil {
		// Object exists, generate 7-day presigned URL
		expireDuration := 7 * 24 * time.Hour // 7 days
		presigned, perr := cli.PresignedGetObject(ctx, objCfg.Bucket, objectKey, expireDuration, nil)
		if perr == nil {
			jobID := fmt.Sprintf("export_%d_%s", time.Now().Unix(), uuid.New().String()[:8])
			urlStr := presigned.String()
			entry := &exportCacheEntry{JobID: jobID, Status: "completed", Storage: "object", ObjectKey: objectKey, Format: exportFormat}
			_ = cacheutil.SetDistributedOnly(ctx, exportKey, entry, time.Duration(constants.CacheTTLMidSeconds)*time.Second)
			if global.Logger != nil {
				global.Logger.Info("ExportDailyReportDetail: Returning existing object", "job_id", jobID, "object_key", objectKey)
			}
			return &model.ExportDailyReportDetailOutput{
				JobID:       jobID,
				Status:      "completed",
				Message:     "Export completed (existing file, expires in 7 days)",
				DownloadURL: &urlStr,
			}, nil
		}
	}

	// Decide async vs sync: async if email provided
	isAsync := input.Email != nil && *input.Email != ""
	if isAsync {
		lockAcquired, lockErr := cacheutil.AcquireLock(ctx, exportKey, 15*time.Minute)
		if lockErr != nil {
			if global.Logger != nil {
				global.Logger.Error("ExportDailyReportDetail: Failed to acquire lock", "error", lockErr.Error())
			}
			return nil, applicationErrors.ErrExportFailed.WithDetails("failed to acquire export lock")
		}

		if !lockAcquired {
			// Another instance is already processing this export
			var ecAsync exportCacheEntry
			if hitAsync, _ := cacheutil.GetDistributedOnly(ctx, exportKey, &ecAsync); hitAsync && ecAsync.Status == "processing" {
				jobID := ecAsync.JobID
				if jobID == "" {
					jobID = "processing"
				}
				return &model.ExportDailyReportDetailOutput{JobID: jobID, Status: "processing", Message: "Export is processing"}, nil
			}
			// Fallback: lock exists but no status found
			return &model.ExportDailyReportDetailOutput{JobID: "processing", Status: "processing", Message: "Export is processing"}, nil
		}

		// Lock acquired, proceed with async export
		jobID := fmt.Sprintf("job_%s", uuid.New().String())
		procEntry := &exportCacheEntry{JobID: jobID, Status: "processing", Storage: "", Format: exportFormat}
		_ = cacheutil.SetDistributedOnly(ctx, exportKey, procEntry, 15*time.Minute)

		go func(job string, compID uuid.UUID, date time.Time, exportFmt, objectK, exportK string, emailPtr *string) {
			bgCtx := context.Background()
			// Ensure lock is released when goroutine exits
			defer func() {
				cacheutil.ReleaseLock(bgCtx, exportK)
			}()

			// Fetch data from ScyllaDB
			summaries, derr := s.repo.GetDailySummariesByDate(bgCtx, compID, date)
			if derr != nil {
				failEntry := &exportCacheEntry{JobID: job, Status: "failed", Storage: "local", LocalPath: "", Format: exportFmt, Rows: 0}
				_ = cacheutil.SetDistributedOnly(bgCtx, exportK, failEntry, 2*time.Minute)
				return
			}

			// Ensure export directory exists
			exportDir := "exports"
			_ = os.MkdirAll(exportDir, 0o755)
			baseName := fmt.Sprintf("%s_daily_detail_%s", job, date.Format("2006-01-02"))
			filePath := filepath.Join(exportDir, baseName+".csv")

			if werr := s.writeCSV(filePath, summaries); werr != nil {
				failEntry := &exportCacheEntry{JobID: job, Status: "failed", Storage: "local", LocalPath: "", Format: exportFmt, Rows: 0}
				_ = cacheutil.SetDistributedOnly(bgCtx, exportK, failEntry, 2*time.Minute)
				return
			}

			// Upload to object storage (required)
			objCfgA := global.SettingServer.ObjectStorage
			cliA, cerrA := minio.New(objCfgA.Endpoint, &minio.Options{Creds: credentials.NewStaticV4(objCfgA.AccessKey, objCfgA.SecretKey, ""), Secure: objCfgA.UseSSL, Region: objCfgA.Region})
			if cerrA != nil {
				if global.Logger != nil {
					global.Logger.Error("ExportDailyReportDetail async: Failed to connect to object storage", "error", cerrA.Error())
				}
				failEntry := &exportCacheEntry{JobID: job, Status: "failed", Storage: "", LocalPath: "", Format: exportFmt, Rows: 0}
				_ = cacheutil.SetDistributedOnly(bgCtx, exportK, failEntry, 2*time.Minute)
				return
			}

			// Upload file to MinIO
			_, uerr := cliA.FPutObject(bgCtx, objCfgA.Bucket, objectK, filePath, minio.PutObjectOptions{ContentType: "text/csv"})
			if uerr != nil {
				if global.Logger != nil {
					global.Logger.Error("ExportDailyReportDetail async: Failed to upload to object storage", "error", uerr.Error())
				}
				failEntry := &exportCacheEntry{JobID: job, Status: "failed", Storage: "", LocalPath: "", Format: exportFmt, Rows: 0}
				_ = cacheutil.SetDistributedOnly(bgCtx, exportK, failEntry, 2*time.Minute)
				return
			}

			// Generate 7-day presigned URL
			expireDuration := 7 * 24 * time.Hour // 7 days
			presigned, perr := cliA.PresignedGetObject(bgCtx, objCfgA.Bucket, objectK, expireDuration, nil)
			if perr != nil {
				if global.Logger != nil {
					global.Logger.Error("ExportDailyReportDetail async: Failed to generate presigned URL", "error", perr.Error())
				}
				failEntry := &exportCacheEntry{JobID: job, Status: "failed", Storage: "", LocalPath: "", Format: exportFmt, Rows: 0}
				_ = cacheutil.SetDistributedOnly(bgCtx, exportK, failEntry, 2*time.Minute)
				return
			}

			download := presigned.String()
			// Remove local file after successful upload
			_ = os.Remove(filePath)
			entry := &exportCacheEntry{JobID: job, Status: "completed", Storage: "object", ObjectKey: objectK, Format: exportFmt, Rows: len(summaries)}
			_ = cacheutil.SetDistributedOnly(bgCtx, exportK, entry, time.Duration(constants.CacheTTLMidSeconds)*time.Second)

			// Send email notification via Kafka
			if emailPtr != nil && *emailPtr != "" {
				_ = s.publishExportEmail(bgCtx, *emailPtr, download, exportFmt, date, date, compID.String())
			}
			if global.Logger != nil {
				global.Logger.Info("ExportDailyReportDetail async: Completed", "job_id", job, "rows", len(summaries))
			}
		}(jobID, input.CompanyID, input.Date, exportFormat, objectKey, exportKey, input.Email)

		return &model.ExportDailyReportDetailOutput{JobID: jobID, Status: "processing", Message: "Export scheduled. Email will be sent when completed."}, nil
	}

	// Sync path below
	if global.Logger != nil {
		global.Logger.Info("ExportDailyReportDetail: Starting sync export",
			"company_id", input.CompanyID.String(),
			"date", input.Date.Format("2006-01-02"),
			"format", exportFormat)
	}

	// Generate job ID for sync export
	syncJobID := fmt.Sprintf("export_%d_%s", time.Now().Unix(), uuid.New().String()[:8])

	lockAcquired, lockErr := cacheutil.AcquireLock(ctx, exportKey, 5*time.Minute)
	if lockErr != nil {
		if global.Logger != nil {
			global.Logger.Error("ExportDailyReportDetail: Failed to acquire lock", "error", lockErr.Error())
		}
		return nil, applicationErrors.ErrExportFailed.WithDetails("failed to acquire export lock")
	}
	if !lockAcquired {
		// Another instance is processing, check status
		var ecSync exportCacheEntry
		if hitSync, _ := cacheutil.GetDistributedOnly(ctx, exportKey, &ecSync); hitSync && ecSync.Status == "processing" {
			jobID := ecSync.JobID
			if jobID == "" {
				jobID = "processing"
			}
			return &model.ExportDailyReportDetailOutput{JobID: jobID, Status: "processing", Message: "Export is processing"}, nil
		}
	}
	defer cacheutil.ReleaseLock(ctx, exportKey)

	procEntry := &exportCacheEntry{JobID: syncJobID, Status: "processing", Storage: "", Format: exportFormat}
	_ = cacheutil.SetDistributedOnly(ctx, exportKey, procEntry, 5*time.Minute)

	// Fetch data from ScyllaDB synchronously
	summaries, derr := s.repo.GetDailySummariesByDate(ctx, input.CompanyID, input.Date)
	if derr != nil {
		if global.Logger != nil {
			global.Logger.Error("GetDailySummariesByDate failed", "error", derr.Error())
		}
		return nil, applicationErrors.ErrDatabaseError.WithDetails(derr.Error())
	}

	if global.Logger != nil {
		global.Logger.Info("ExportDailyReportDetail: Fetched data", "rows", len(summaries))
	}

	// Ensure export directory exists
	exportDir := "exports"
	if err := os.MkdirAll(exportDir, 0o755); err != nil {
		if global.Logger != nil {
			global.Logger.Error("Failed to create export directory", "error", err.Error())
		}
		return nil, applicationErrors.ErrExportFailed.WithDetails("failed to create export directory")
	}

	jobID := fmt.Sprintf("export_%d_%s", time.Now().Unix(), uuid.New().String()[:8])
	baseName := fmt.Sprintf("%s_daily_detail_%s", jobID, input.Date.Format("2006-01-02"))

	var filePath string
	switch input.Format {
	case "csv", "excel":
		filePath = filepath.Join(exportDir, baseName+".csv")
		if err := s.writeCSV(filePath, summaries); err != nil {
			if global.Logger != nil {
				global.Logger.Error("Failed to write CSV", "error", err.Error())
			}
			return nil, applicationErrors.ErrExportFailed.WithDetails("failed to write CSV file")
		}
	case "pdf":
		// For PDF, you would implement a writePDF method similar to writeCSV
		return nil, applicationErrors.ErrInvalidInput.WithDetails("PDF export not yet implemented")
	default:
		return nil, applicationErrors.ErrInvalidInput.WithDetails("unsupported format")
	}

	if global.Logger != nil {
		global.Logger.Info("ExportDailyReportDetail: File created",
			"path", filePath,
			"format", exportFormat,
			"rows", len(summaries))
	}

	// Upload to object storage (MinIO/S3-compatible) - REQUIRED
	_, uerr := cli.FPutObject(ctx, objCfg.Bucket, objectKey, filePath, minio.PutObjectOptions{ContentType: "text/csv"})
	if uerr != nil {
		if global.Logger != nil {
			global.Logger.Error("ExportDailyReportDetail: Failed to upload to object storage", "error", uerr.Error())
		}
		// Clean up local file
		_ = os.Remove(filePath)
		return nil, applicationErrors.ErrExportFailed.WithDetails(fmt.Sprintf("Failed to upload to object storage: %v", uerr))
	}

	if global.Logger != nil {
		global.Logger.Info("ExportDailyReportDetail: Uploaded to object storage", "bucket", objCfg.Bucket, "key", objectKey)
	}

	// Generate 7-day presigned URL
	expireDuration := 7 * 24 * time.Hour // 7 days
	presigned, perr := cli.PresignedGetObject(ctx, objCfg.Bucket, objectKey, expireDuration, nil)
	if perr != nil {
		if global.Logger != nil {
			global.Logger.Error("ExportDailyReportDetail: Failed to generate presigned URL", "error", perr.Error())
		}
		// Clean up local file
		_ = os.Remove(filePath)
		return nil, applicationErrors.ErrExportFailed.WithDetails(fmt.Sprintf("Failed to generate presigned URL: %v", perr))
	}

	urlStr := presigned.String()
	// Remove local file after successful upload
	_ = os.Remove(filePath)

	// Cache the result
	entry := &exportCacheEntry{JobID: jobID, Status: "completed", Storage: "object", ObjectKey: objectKey, Format: exportFormat, Rows: len(summaries)}
	_ = cacheutil.SetDistributedOnly(ctx, exportKey, entry, time.Duration(constants.CacheTTLMidSeconds)*time.Second)

	// Send email notification via Kafka if email provided
	if input.Email != nil && *input.Email != "" {
		_ = s.publishExportEmail(ctx, *input.Email, urlStr, exportFormat, input.Date, input.Date, input.CompanyID.String())
	}

	if global.Logger != nil {
		global.Logger.Info("ExportDailyReportDetail: Completed", "job_id", jobID, "rows", len(summaries), "expires", "7 days")
	}

	return &model.ExportDailyReportDetailOutput{
		JobID:       jobID,
		Status:      "completed",
		Message:     fmt.Sprintf("Exported %d rows (download link expires in 7 days)", len(summaries)),
		DownloadURL: &urlStr,
	}, nil
}

// GetDailyReportDetail implements service.IAnalyticService.
func (s *AnalyticServiceImpl) GetDailyReportDetail(ctx context.Context, input *model.DailyDetailReportInput) (*model.DailyReportDetailOutput, *applicationErrors.Error) {
	// Log entry
	if global.Logger != nil {
		global.Logger.Info("GetDailyReport start",
			"company_id", input.CompanyID.String(),
			"date", input.Date.Format("2006-01-02"),
		)
	}

	// Authorization check: Require CompanyAdmin role minimum, no employee self-access for company-wide reports
	requestedCompanyID := input.CompanyID.String()
	if err := s.checkAuthorization(input.Session, &requestedCompanyID, domainModel.RoleCompanyAdmin, false); err != nil {
		return nil, err
	}
	companyID, err := uuid.Parse(requestedCompanyID)
	if err != nil {
		return nil, applicationErrors.ErrInvalidInput.WithDetails("invalid company_id")
	}

	// Cache key for daily report (includes device filter when provided)
	cacheKey := cacheutil.BuildDailyDetailByDateKey(companyID, input.Date)

	// Try local cache first
	if v, ok := cacheutil.GetLocal(cacheKey); ok {
		if out, ok2 := v.(*model.DailyReportDetailOutput); ok2 {
			if global.Logger != nil {
				global.Logger.Debug("GetDailyReport cache hit (local)", "key", cacheKey)
			}
			return out, nil
		}
	}

	var cachedOut model.DailyReportDetailOutput
	if hit, _ := cacheutil.GetDistributed(ctx, cacheKey, &cachedOut); hit {
		// backfill local cache with slightly shorter TTL than distributed
		cacheutil.SetLocal(cacheKey, &cachedOut, localTTLFrom(constants.CacheTTLMidSeconds))
		if global.Logger != nil {
			global.Logger.Debug("GetDailyReport cache hit (redis)", "key", cacheKey)
		}
		return &cachedOut, nil
	}

	// Validate data querying parameters
	if input.Limit == nil || *input.Limit <= 0 {
		return nil, applicationErrors.ErrInvalidInput.WithDetails("limit must be greater than 0")
	}
	resp, nextPage, err := s.repo.GetDailySummariesByDatePage(ctx, companyID, input.Date, input.PageToken, *input.Limit)
	if err != nil {
		if global.Logger != nil {
			global.Logger.Error("GetDailySummariesByDatePage failed", "error", err.Error())
		}
		return nil, applicationErrors.ErrDatabaseError.WithDetails(err.Error())
	}
	if resp == nil {
		return nil, applicationErrors.ErrNotFound.WithDetails("no daily summaries found")
	}
	out := &model.DailyReportDetailOutput{
		Total:    len(resp),
		Items:    make([]model.DailyReportDetailEmployeeRow, 0, len(resp)),
		NextPage: nextPage,
	}
	// Map domain models to output models
	for _, item := range resp {
		row := model.DailyReportDetailEmployeeRow{
			CompanyID:            item.CompanyID,
			SummaryMonth:         item.SummaryMonth,
			WorkDate:             item.WorkDate,
			EmployeeID:           item.EmployeeID,
			ShiftID:              item.ShiftID,
			ActualCheckIn:        item.ActualCheckIn,
			ActualCheckOut:       item.ActualCheckOut,
			AttendanceStatus:     item.AttendanceStatus,
			LateMinutes:          item.LateMinutes,
			EarlyLeaveMinutes:    item.EarlyLeaveMinutes,
			TotalWorkMinutes:     item.TotalWorkMinutes,
			Notes:                item.Notes,
			UpdatedAt:            item.UpdatedAt,
			OvertimeMinutes:      item.OvertimeMinutes,
			AttendancePercentage: item.AttendancePercentage,
		}
		out.Items = append(out.Items, row)
	}

	// Cache the result
	_ = cacheutil.SetDistributed(ctx, cacheKey, out, time.Duration(constants.CacheTTLMidSeconds)*time.Second)
	_ = cacheutil.SetLocal(cacheKey, out, localTTLFrom(constants.CacheTTLMidSeconds))
	if global.Logger != nil {
		global.Logger.Info("GetDailyReportDetail computed", "key", cacheKey, "total_rows", len(resp))
	}
	return out, nil
}

// NewAnalyticService creates a new analytics service instance
func NewAnalyticService(repo repository.IAnalyticRepository) service.IAnalyticService {
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

	// Authorization check: Require CompanyAdmin role minimum, no employee self-access for company-wide reports
	if err := s.checkAuthorization(input.Session, input.CompanyID, domainModel.RoleCompanyAdmin, false); err != nil {
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
	// Authorization check: Require CompanyAdmin role, but allow employees to view their own summary
	if err := s.checkAuthorization(input.Session, input.CompanyID, domainModel.RoleCompanyAdmin, true); err != nil {
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

	// Check if user is an employee - if so, return only their own summary
	userRole := domainModel.Role(input.Session.Role)
	if userRole == domainModel.RoleEmployee {
		// Employee can only view their own data
		employeeID, parseErr := uuid.Parse(input.Session.UserID)
		if parseErr != nil {
			return nil, applicationErrors.ErrInvalidInput.WithDetails("invalid user_id in session")
		}
		return s.getEmployeeSummaryReport(ctx, companyID, employeeID, input.Month, startDate, endDate)
	}

	// For CompanyAdmin and SystemAdmin: return full company summary
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

// getEmployeeSummaryReport returns monthly summary report for a specific employee
func (s *AnalyticServiceImpl) getEmployeeSummaryReport(ctx context.Context, companyID, employeeID uuid.UUID, month string, startDate, endDate time.Time) (*model.SummaryReportOutput, *applicationErrors.Error) {
	// Fetch employee's daily summaries for the month
	summaries, err := s.repo.GetDailySummariesByEmployeeMonth(ctx, companyID, employeeID, month)
	if err != nil {
		if global.Logger != nil {
			global.Logger.Error("GetDailySummariesByEmployeeMonth failed", "error", err.Error(), "employee_id", employeeID.String())
		}
		return nil, applicationErrors.ErrDatabaseError.WithDetails(err.Error())
	}

	totalWorkingDays := endDate.Day()
	totalPresentDays := 0
	totalWorkingMinutes := 0
	totalOvertimeMinutes := 0

	for _, summary := range summaries {
		if summary.AttendanceStatus == 0 { // PRESENT
			totalPresentDays++
		}
		totalWorkingMinutes += summary.TotalWorkMinutes
		totalOvertimeMinutes += summary.OvertimeMinutes
	}

	// Calculate attendance rate for this employee
	averageAttendanceRate := 0.0
	if totalWorkingDays > 0 {
		averageAttendanceRate = float64(totalPresentDays) / float64(totalWorkingDays) * 100
	}

	// For employee view, weekly summary is calculated only for this employee
	weeklySummary := s.calculateWeeklySummaryForEmployee(ctx, startDate, endDate, companyID, employeeID)

	// Top/Low attendance lists don't make sense for single employee, return empty
	out := &model.SummaryReportOutput{
		Month:                  month,
		TotalWorkingDays:       totalWorkingDays,
		TotalEmployees:         1, // Only this employee
		AverageAttendanceRate:  roundFloat(averageAttendanceRate, 2),
		TotalWorkingHours:      totalWorkingMinutes / 60,
		TotalOvertimeHours:     totalOvertimeMinutes / 60,
		WeeklySummary:          weeklySummary,
		TopAttendanceEmployees: []model.EmployeeAttendanceStat{}, // Not applicable for single employee
		LowAttendanceEmployees: []model.EmployeeAttendanceStat{}, // Not applicable for single employee
	}

	if global.Logger != nil {
		global.Logger.Info("GetEmployeeSummaryReport computed", "employee_id", employeeID.String(), "month", month, "present_days", totalPresentDays)
	}

	return out, nil
}

// calculateWeeklySummaryForEmployee calculates weekly summary for a specific employee
func (s *AnalyticServiceImpl) calculateWeeklySummaryForEmployee(ctx context.Context, startDate, endDate time.Time, companyID, employeeID uuid.UUID) []model.WeeklySummary {
	weeklySummaries := []model.WeeklySummary{}
	currentDate := startDate
	weekNumber := 1

	for currentDate.Before(endDate) || currentDate.Equal(endDate) {
		weekStart := currentDate
		weekEnd := currentDate.AddDate(0, 0, 6)
		if weekEnd.After(endDate) {
			weekEnd = endDate
		}

		// Get attendance data for this week for this specific employee
		weekSummaries, err := s.repo.GetDailySummariesByEmployeeDateRange(ctx, companyID, employeeID, weekStart, weekEnd)
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

		// For employee, attendance rate is based on actual working days in the week
		daysInWeek := int(weekEnd.Sub(weekStart).Hours()/24) + 1
		attendanceRate := 0.0
		if daysInWeek > 0 {
			attendanceRate = float64(totalPresentDays) / float64(daysInWeek) * 100
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

// ExportReport exports attendance report to file
func (s *AnalyticServiceImpl) ExportReport(ctx context.Context, input *model.ExportReportInput) (*model.ExportReportOutput, *applicationErrors.Error) {
	// Authorization check: Require CompanyAdmin role, but allow employees to export their own data
	if err := s.checkAuthorization(input.Session, input.CompanyID, domainModel.RoleCompanyAdmin, true); err != nil {
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

	// Check if user is an employee - if so, filter to their own data only
	userRole := domainModel.Role(input.Session.Role)
	var employeeFilterID *uuid.UUID
	if userRole == domainModel.RoleEmployee {
		employeeID, parseErr := uuid.Parse(input.Session.UserID)
		if parseErr != nil {
			return nil, applicationErrors.ErrInvalidInput.WithDetails("invalid user_id in session")
		}
		employeeFilterID = &employeeID
	}

	// Normalize format (excel -> csv)
	exportFormat := input.Format
	if exportFormat == "excel" {
		exportFormat = "csv"
	}

	// Build cache key for export
	exportKey := cacheutil.BuildExportKey(companyID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), exportFormat)

	type exportCacheEntry struct {
		JobID     string `json:"job_id,omitempty"`
		Status    string `json:"status"`  // processing | completed
		Storage   string `json:"storage"` // object | local
		ObjectKey string `json:"object_key,omitempty"`
		LocalPath string `json:"local_path,omitempty"`
		Format    string `json:"format"`
		Rows      int    `json:"rows"`
	}

	// Only check Redis to avoid stale data across multiple instances
	var ec exportCacheEntry
	if hit, _ := cacheutil.GetDistributedOnly(ctx, exportKey, &ec); hit {
		// Do NOT backfill to local cache for mutable export status
		if ec.Status == "processing" {
			jobID := ec.JobID
			if jobID == "" {
				jobID = "processing"
			}
			return &model.ExportReportOutput{JobID: jobID, Status: "processing", Message: "Export is processing"}, nil
		}
		if ec.Status == "completed" {
			if ec.Storage == "object" && ec.ObjectKey != "" {
				objCfg := global.SettingServer.ObjectStorage
				if objCfg.Endpoint != "" && objCfg.Bucket != "" {
					cli, cerr := minio.New(objCfg.Endpoint, &minio.Options{Creds: credentials.NewStaticV4(objCfg.AccessKey, objCfg.SecretKey, ""), Secure: objCfg.UseSSL, Region: objCfg.Region})
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
	}

	// Prepare object key early (includes employee scope when present)
	var scopePart string
	if employeeFilterID != nil {
		scopePart = employeeFilterID.String()
	} else {
		scopePart = "all"
	}
	objectKey := fmt.Sprintf("reports/%s/%s_%s_%s_%s.%s", time.Now().Format("2006/01/02"), companyID.String(), scopePart, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), exportFormat)

	objCfg := global.SettingServer.ObjectStorage
	if objCfg.Endpoint != "" && objCfg.Bucket != "" {
		cli, cerr := minio.New(objCfg.Endpoint, &minio.Options{Creds: credentials.NewStaticV4(objCfg.AccessKey, objCfg.SecretKey, ""), Secure: objCfg.UseSSL, Region: objCfg.Region})
		if cerr == nil {
			exists, eerr := cli.BucketExists(ctx, objCfg.Bucket)
			if eerr == nil && !exists {
				_ = cli.MakeBucket(ctx, objCfg.Bucket, minio.MakeBucketOptions{Region: objCfg.Region})
			}
			// Try stat existing object to avoid regeneration
			if _, statErr := cli.StatObject(ctx, objCfg.Bucket, objectKey, minio.StatObjectOptions{}); statErr == nil {
				expireMinutes := objCfg.PresignExpireMinutes
				if expireMinutes <= 0 {
					expireMinutes = 60
				}
				presigned, perr := cli.PresignedGetObject(ctx, objCfg.Bucket, objectKey, time.Duration(expireMinutes)*time.Minute, nil)
				if perr == nil {
					download := presigned.String()
					jobIDCached := fmt.Sprintf("cached_%s", uuid.New().String()[:8])
					entry := &exportCacheEntry{JobID: jobIDCached, Status: "completed", Storage: "object", ObjectKey: objectKey, Format: exportFormat, Rows: 0}
					_ = cacheutil.SetDistributedOnly(ctx, exportKey, entry, time.Duration(expireMinutes)*time.Minute)
					if input.Email != nil && *input.Email != "" {
						if nerr := s.publishExportEmail(ctx, *input.Email, download, exportFormat, startDate, endDate, *input.CompanyID); nerr != nil {
							return nil, applicationErrors.ErrExportFailed.WithDetails("notify email failed: " + nerr.Error())
						}
					}
					return &model.ExportReportOutput{JobID: "cached", Status: "completed", Message: "Export served from existing object", DownloadURL: &download}, nil
				}
			}
		}
	}

	// Decide async vs sync: async if email provided
	isAsync := input.Email != nil && *input.Email != ""
	if isAsync {
		lockAcquired, lockErr := cacheutil.AcquireLock(ctx, exportKey, 15*time.Minute)
		if lockErr != nil {
			if global.Logger != nil {
				global.Logger.Error("ExportReport: Failed to acquire lock", "error", lockErr.Error())
			}
			return nil, applicationErrors.ErrExportFailed.WithDetails("failed to acquire export lock")
		}

		if !lockAcquired {
			// Another instance is already processing this export
			var ecAsync exportCacheEntry
			if hitAsync, _ := cacheutil.GetDistributedOnly(ctx, exportKey, &ecAsync); hitAsync && ecAsync.Status == "processing" {
				jobID := ecAsync.JobID
				if jobID == "" {
					jobID = "processing"
				}
				if global.Logger != nil {
					global.Logger.Info("ExportReport: Export already processing on another instance", "job_id", jobID)
				}
				return &model.ExportReportOutput{JobID: jobID, Status: "processing", Message: "Export is being processed by another instance"}, nil
			}
			// Fallback: lock exists but no status found
			return &model.ExportReportOutput{JobID: "processing", Status: "processing", Message: "Export is processing"}, nil
		}

		// Lock acquired, proceed with export
		jobID := fmt.Sprintf("job_%s", uuid.New().String())
		procEntry := &exportCacheEntry{JobID: jobID, Status: "processing", Storage: "", Format: exportFormat}
		_ = cacheutil.SetDistributedOnly(ctx, exportKey, procEntry, 15*time.Minute)
		go func(job string, compID uuid.UUID, start, end time.Time, employeeID *uuid.UUID, exportFmt, objectK, exportK string, emailPtr *string) {
			bgCtx := context.Background()
			// Ensure lock is released when goroutine exits
			defer func() {
				if rerr := cacheutil.ReleaseLock(bgCtx, exportK); rerr != nil && global.Logger != nil {
					global.Logger.Warn("ExportReport: Failed to release lock", "error", rerr.Error())
				}
			}()
			var summaries []*domainModel.DailySummary
			var derr2 error
			if employeeID != nil {
				summaries, derr2 = s.repo.GetDailySummariesByEmployeeDateRange(bgCtx, compID, *employeeID, start, end)
			} else {
				summaries, derr2 = s.repo.GetDailySummariesByDateRange(bgCtx, compID, start, end)
			}
			if derr2 != nil {
				failEntry := &exportCacheEntry{JobID: job, Status: "failed", Storage: "local", LocalPath: "", Format: exportFmt, Rows: 0}
				_ = cacheutil.SetDistributedOnly(bgCtx, exportK, failEntry, 2*time.Minute)
				return
			}
			exportDir := "exports"
			_ = os.MkdirAll(exportDir, 0o755)
			baseName := fmt.Sprintf("%s_%s_to_%s", job, start.Format("2006-01-02"), end.Format("2006-01-02"))
			filePath := filepath.Join(exportDir, baseName+".csv")
			if werr := s.writeCSV(filePath, summaries); werr != nil {
				_ = cacheutil.SetDistributedOnly(bgCtx, exportK, &exportCacheEntry{JobID: job, Status: "failed", Storage: "local", LocalPath: "", Format: exportFmt, Rows: 0}, 2*time.Minute)
				return
			}
			objCfgA := global.SettingServer.ObjectStorage
			var download string
			if objCfgA.Endpoint != "" && objCfgA.Bucket != "" {
				cliA, cerrA := minio.New(objCfgA.Endpoint, &minio.Options{Creds: credentials.NewStaticV4(objCfgA.AccessKey, objCfgA.SecretKey, ""), Secure: objCfgA.UseSSL, Region: objCfgA.Region})
				if cerrA == nil {
					existsA, eerrA := cliA.BucketExists(bgCtx, objCfgA.Bucket)
					if eerrA == nil && !existsA {
						_ = cliA.MakeBucket(bgCtx, objCfgA.Bucket, minio.MakeBucketOptions{Region: objCfgA.Region})
					}
					if _, uperrA := cliA.FPutObject(bgCtx, objCfgA.Bucket, objectK, filePath, minio.PutObjectOptions{ContentType: "text/csv"}); uperrA == nil {
						expireMinutesA := objCfgA.PresignExpireMinutes
						if expireMinutesA <= 0 {
							expireMinutesA = 60
						}
						presignedA, perrA := cliA.PresignedGetObject(bgCtx, objCfgA.Bucket, objectK, time.Duration(expireMinutesA)*time.Minute, nil)
						if perrA == nil {
							download = presignedA.String()
							entryA := &exportCacheEntry{JobID: job, Status: "completed", Storage: "object", ObjectKey: objectK, Format: exportFmt, Rows: len(summaries)}
							_ = cacheutil.SetDistributedOnly(bgCtx, exportK, entryA, time.Duration(expireMinutesA)*time.Minute)
							_ = os.Remove(filePath)
						}
					}
				}
			}
			if download == "" {
				entryLocal := &exportCacheEntry{JobID: job, Status: "completed", Storage: "local", LocalPath: filePath, Format: exportFmt, Rows: len(summaries)}
				_ = cacheutil.SetDistributedOnly(bgCtx, exportK, entryLocal, time.Duration(constants.CacheTTLMidSeconds)*time.Second)
				download = filePath
			}
			if emailPtr != nil && *emailPtr != "" {
				_ = s.publishExportEmail(bgCtx, *emailPtr, download, exportFmt, start, end, compID.String())
			}
		}(jobID, companyID, startDate, endDate, employeeFilterID, exportFormat, objectKey, exportKey, input.Email)
		return &model.ExportReportOutput{JobID: jobID, Status: "processing", Message: "Export scheduled. Email will be sent when completed."}, nil
	}

	// Sync path below
	if global.Logger != nil {
		global.Logger.Info("ExportReport: Starting sync export",
			"company_id", companyID.String(),
			"start_date", startDate.Format("2006-01-02"),
			"end_date", endDate.Format("2006-01-02"),
			"format", exportFormat,
			"has_employee_filter", employeeFilterID != nil)
	}

	// Generate job ID for sync export and set processing status
	syncJobID := fmt.Sprintf("export_%d_%s", time.Now().Unix(), uuid.New().String()[:8])

	lockAcquired, lockErr := cacheutil.AcquireLock(ctx, exportKey, 5*time.Minute)
	if lockErr != nil {
		if global.Logger != nil {
			global.Logger.Error("ExportReport: Failed to acquire sync export lock", "error", lockErr.Error())
		}
		return nil, applicationErrors.ErrExportFailed.WithDetails("failed to acquire export lock")
	}
	if !lockAcquired {
		// Another instance is processing, check status
		var ecSync exportCacheEntry
		if hitSync, _ := cacheutil.GetDistributedOnly(ctx, exportKey, &ecSync); hitSync {
			if ecSync.Status == "processing" {
				return &model.ExportReportOutput{JobID: ecSync.JobID, Status: "processing", Message: "Export is being processed by another instance"}, nil
			}
		}
	}
	defer cacheutil.ReleaseLock(ctx, exportKey) // Release lock after sync export completes

	procEntry := &exportCacheEntry{JobID: syncJobID, Status: "processing", Storage: "", Format: exportFormat}
	_ = cacheutil.SetDistributedOnly(ctx, exportKey, procEntry, 5*time.Minute) // Fetch data from ScyllaDB (cache miss) synchronously
	var summaries []*domainModel.DailySummary
	var derr error
	if employeeFilterID != nil {
		summaries, derr = s.repo.GetDailySummariesByEmployeeDateRange(ctx, companyID, *employeeFilterID, startDate, endDate)
	} else {
		summaries, derr = s.repo.GetDailySummariesByDateRange(ctx, companyID, startDate, endDate)
	}
	if derr != nil {
		if global.Logger != nil {
			global.Logger.Error("ExportReport: Failed to fetch data", "error", derr.Error())
		}
		return nil, applicationErrors.ErrDatabaseError.WithDetails(derr.Error())
	}

	if global.Logger != nil {
		global.Logger.Info("ExportReport: Fetched data from database", "rows", len(summaries))
	}

	// Ensure export directory exists
	exportDir := "exports"
	if err := os.MkdirAll(exportDir, 0o755); err != nil {
		if global.Logger != nil {
			global.Logger.Error("ExportReport: Failed to create export directory", "error", err.Error())
		}
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
			if global.Logger != nil {
				global.Logger.Error("ExportReport: Failed to write CSV", "error", werr.Error())
			}
			return nil, applicationErrors.ErrExportFailed.WithDetails(werr.Error())
		}
	case "excel":
		// Default to CSV when Excel lib not available; still generate .csv
		filePath = filepath.Join(exportDir, baseName+".csv")
		if werr := s.writeCSV(filePath, summaries); werr != nil {
			if global.Logger != nil {
				global.Logger.Error("ExportReport: Failed to write CSV", "error", werr.Error())
			}
			return nil, applicationErrors.ErrExportFailed.WithDetails(werr.Error())
		}
	case "pdf":
		return nil, applicationErrors.ErrInvalidInput.WithDetails("pdf export not supported yet")
	default:
		return nil, applicationErrors.ErrInvalidInput.WithDetails("unsupported format; use excel or csv")
	}

	if global.Logger != nil {
		global.Logger.Info("ExportReport: File created successfully", "path", filePath, "size_bytes", func() int64 {
			if fi, err := os.Stat(filePath); err == nil {
				return fi.Size()
			}
			return 0
		}())
	}

	// Upload to object storage (MinIO/S3-compatible) using precomputed objectKey
	objCfg = global.SettingServer.ObjectStorage
	if objCfg.Endpoint != "" && objCfg.Bucket != "" {
		if global.Logger != nil {
			global.Logger.Info("ExportReport: Attempting to upload to object storage",
				"endpoint", objCfg.Endpoint,
				"bucket", objCfg.Bucket,
				"object_key", objectKey)
		}
		cli, cerr := minio.New(objCfg.Endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(objCfg.AccessKey, objCfg.SecretKey, ""),
			Secure: objCfg.UseSSL,
			Region: objCfg.Region,
		})
		if cerr != nil {
			if global.Logger != nil {
				global.Logger.Error("ExportReport: Failed to init MinIO client", "error", cerr.Error())
			}
			// Don't return error, fallback to local file
		} else {
			// Ensure bucket exists
			exists, eerr := cli.BucketExists(ctx, objCfg.Bucket)
			if eerr != nil {
				if global.Logger != nil {
					global.Logger.Error("ExportReport: Failed to check bucket", "error", eerr.Error())
				}
				// Don't return error, fallback to local file
			} else {
				if !exists {
					if global.Logger != nil {
						global.Logger.Info("ExportReport: Creating bucket", "bucket", objCfg.Bucket)
					}
					if cberr := cli.MakeBucket(ctx, objCfg.Bucket, minio.MakeBucketOptions{Region: objCfg.Region}); cberr != nil {
						if global.Logger != nil {
							global.Logger.Error("ExportReport: Failed to create bucket", "error", cberr.Error())
						}
						// Don't return error, fallback to local file
					}
				}

				contentType := "text/csv"
				if global.Logger != nil {
					global.Logger.Info("ExportReport: Uploading file to MinIO", "file", filePath, "object_key", objectKey)
				}
				if _, uperr := cli.FPutObject(ctx, objCfg.Bucket, objectKey, filePath, minio.PutObjectOptions{ContentType: contentType}); uperr != nil {
					if global.Logger != nil {
						global.Logger.Error("ExportReport: Failed to upload to MinIO", "error", uperr.Error())
					}
					// Don't return error, fallback to local file
				} else {
					if global.Logger != nil {
						global.Logger.Info("ExportReport: File uploaded successfully, generating presigned URL")
					}

					// Generate presigned URL with expiry
					expireMinutes := objCfg.PresignExpireMinutes
					if expireMinutes <= 0 {
						expireMinutes = 60
					}
					presigned, perr := cli.PresignedGetObject(ctx, objCfg.Bucket, objectKey, time.Duration(expireMinutes)*time.Minute, nil)
					if perr != nil {
						if global.Logger != nil {
							global.Logger.Error("ExportReport: Failed to generate presigned URL", "error", perr.Error())
						}
						// Don't return error, fallback to local file
					} else {
						download := presigned.String()

						if global.Logger != nil {
							global.Logger.Info("ExportReport: Presigned URL generated",
								"url_length", len(download),
								"expire_minutes", expireMinutes)
						}

						// Cache completed export with TTL equal to presign expiry
						entry := &exportCacheEntry{JobID: jobID, Status: "completed", Storage: "object", ObjectKey: objectKey, Format: exportFormat, Rows: len(summaries)}
						_ = cacheutil.SetDistributedOnly(ctx, exportKey, entry, time.Duration(expireMinutes)*time.Minute) // If email provided, publish Kafka notification message
						if input.Email != nil && *input.Email != "" {
							if nerr := s.publishExportEmail(ctx, *input.Email, download, input.Format, startDate, endDate, *input.CompanyID); nerr != nil {
								if global.Logger != nil {
									global.Logger.Error("ExportReport: Failed to send email notification", "error", nerr.Error())
								}
								return nil, applicationErrors.ErrExportFailed.WithDetails("notify email failed: " + nerr.Error())
							}
						}
						// Remove local temp after successful upload
						_ = os.Remove(filePath)
						if global.Logger != nil {
							global.Logger.Info("ExportReport: Sync export completed successfully",
								"job_id", jobID,
								"rows", len(summaries),
								"storage", "object")
						}
						return &model.ExportReportOutput{
							JobID:       jobID,
							Status:      "completed",
							Message:     fmt.Sprintf("Exported %d rows to object storage", len(summaries)),
							DownloadURL: &download,
						}, nil
					}
				}
			}
		}
	}

	// Fallback: local file - construct download URL using API endpoint
	if global.Logger != nil {
		global.Logger.Warn("ExportReport: Object storage not configured or upload failed, using local storage",
			"file_path", filePath)
	}

	// Extract just the filename from the path
	filename := filepath.Base(filePath)

	// Build download URL through API endpoint
	// Use localhost and server port from config
	downloadURL := fmt.Sprintf("http://localhost:%d/api/v1/reports/download/%s", global.SettingServer.Server.Port, filename)

	entry := &exportCacheEntry{JobID: jobID, Status: "completed", Storage: "local", LocalPath: filePath, Format: exportFormat, Rows: len(summaries)}
	_ = cacheutil.SetDistributedOnly(ctx, exportKey, entry, time.Duration(constants.CacheTTLMidSeconds)*time.Second)
	if input.Email != nil && *input.Email != "" {
		if nerr := s.publishExportEmail(ctx, *input.Email, downloadURL, input.Format, startDate, endDate, *input.CompanyID); nerr != nil {
			if global.Logger != nil {
				global.Logger.Error("ExportReport: Failed to send email notification", "error", nerr.Error())
			}
			return nil, applicationErrors.ErrExportFailed.WithDetails("notify email failed: " + nerr.Error())
		}
	}
	if global.Logger != nil {
		global.Logger.Info("ExportReport: Sync export completed with local storage",
			"job_id", jobID,
			"rows", len(summaries),
			"file_path", filePath,
			"download_url", downloadURL)
	}
	return &model.ExportReportOutput{
		JobID:       jobID,
		Status:      "completed",
		Message:     fmt.Sprintf("Exported %d rows to local storage", len(summaries)),
		DownloadURL: &downloadURL,
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
// Returns the effective role that determines data access scope
// minRole: minimum role required to access this endpoint (RoleSystemAdmin=0, RoleCompanyAdmin=1, RoleEmployee=2)
// allowEmployeeSelfAccess: if true, employees can access their own data even if minRole is higher
func (s *AnalyticServiceImpl) checkAuthorization(session *model.SessionInfo, requestedCompanyID *string, minRole domainModel.Role, allowEmployeeSelfAccess bool) *applicationErrors.Error {
	if session == nil {
		return applicationErrors.ErrUnauthorized.WithDetails("session info required")
	}

	if requestedCompanyID == nil {
		return applicationErrors.ErrInvalidInput.WithDetails("company_id is required")
	}

	userRole := domainModel.Role(session.Role)

	// 1. System admin has full access to all companies
	if userRole == domainModel.RoleSystemAdmin {
		return nil
	}

	// 2. Check company_id match for non-system-admin users
	// Both CompanyAdmin and Employee must belong to the requested company
	if session.CompanyID != *requestedCompanyID {
		if global.Logger != nil {
			global.Logger.Warn("Authorization failed: company mismatch",
				"user_id", session.UserID,
				"session_company_id", session.CompanyID,
				"requested_company_id", *requestedCompanyID)
		}
		return applicationErrors.ErrForbidden.WithDetails("access denied: you can only access your own company data")
	}

	// 3. Check if user's role meets the minimum required role
	// Lower role value = higher privilege (SystemAdmin=0, CompanyAdmin=1, Employee=2)
	if userRole > minRole {
		// User doesn't have sufficient privileges
		// Check if employee self-access is allowed for this endpoint
		if !allowEmployeeSelfAccess || userRole != domainModel.RoleEmployee {
			if global.Logger != nil {
				global.Logger.Warn("Authorization failed: insufficient role",
					"user_id", session.UserID,
					"user_role", userRole,
					"min_role", minRole,
					"allow_employee_self_access", allowEmployeeSelfAccess)
			}
			return applicationErrors.ErrForbidden.WithDetails("access denied: insufficient permissions")
		}
		// If allowEmployeeSelfAccess is true and user is employee,
		// authorization passes but business logic must filter data to employee's own records
	}

	// 4. CompanyAdmin has passed: they can see all company data
	// 5. Employee with allowEmployeeSelfAccess has passed: business logic must filter to their own data
	return nil
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

// ============================================
// Attendance Records methods implementation
// ============================================

// GetAttendanceRecords retrieves attendance records for a company and month
func (s *AnalyticServiceImpl) GetAttendanceRecords(ctx context.Context, companyID uuid.UUID, yearMonth string, limit int) ([]*domainModel.AttendanceRecord, error) {
	return s.repo.GetAttendanceRecords(ctx, companyID, yearMonth, limit)
}

// GetAttendanceRecordsByTimeRange retrieves attendance records within a time range
func (s *AnalyticServiceImpl) GetAttendanceRecordsByTimeRange(ctx context.Context, companyID uuid.UUID, yearMonth string, startTime, endTime time.Time) ([]*domainModel.AttendanceRecord, error) {
	return s.repo.GetAttendanceRecordsByTimeRange(ctx, companyID, yearMonth, startTime, endTime)
}

// GetAttendanceRecordsByEmployee retrieves attendance records for a specific employee
func (s *AnalyticServiceImpl) GetAttendanceRecordsByEmployee(ctx context.Context, companyID uuid.UUID, yearMonth string, employeeID uuid.UUID) ([]*domainModel.AttendanceRecord, error) {
	return s.repo.GetAttendanceRecordsByEmployee(ctx, companyID, yearMonth, employeeID)
}

// GetAttendanceRecordsByUser retrieves attendance records indexed by user
func (s *AnalyticServiceImpl) GetAttendanceRecordsByUser(ctx context.Context, companyID, employeeID uuid.UUID, yearMonth string, limit int) ([]*domainModel.AttendanceRecordByUser, error) {
	return s.repo.GetAttendanceRecordsByUser(ctx, companyID, employeeID, yearMonth, limit)
}

// ============================================
// Daily Summary methods implementation
// ============================================

// GetDailySummaries retrieves daily summaries for a month
func (s *AnalyticServiceImpl) GetDailySummaries(ctx context.Context, companyID uuid.UUID, month string) ([]*domainModel.DailySummary, error) {
	return s.repo.GetDailySummariesByMonth(ctx, companyID, month)
}

// GetDailySummaryByEmployeeDate retrieves a specific daily summary
func (s *AnalyticServiceImpl) GetDailySummaryByEmployeeDate(ctx context.Context, companyID uuid.UUID, month string, workDate time.Time, employeeID uuid.UUID) (*domainModel.DailySummary, error) {
	return s.repo.GetDailySummaryByEmployeeDate(ctx, companyID, month, workDate, employeeID)
}

// GetDailySummariesByUser retrieves daily summaries for a user
func (s *AnalyticServiceImpl) GetDailySummariesByUser(ctx context.Context, companyID, employeeID uuid.UUID, month string) ([]*domainModel.DailySummaryByUser, error) {
	return s.repo.GetDailySummariesByUser(ctx, companyID, employeeID, month)
}

// ============================================
// Audit Logs methods implementation
// ============================================

// GetAuditLogs retrieves audit logs for a company and month
func (s *AnalyticServiceImpl) GetAuditLogs(ctx context.Context, companyID uuid.UUID, yearMonth string, limit int) ([]*domainModel.AuditLog, error) {
	return s.repo.GetAuditLogs(ctx, companyID, yearMonth, limit)
}

// GetAuditLogsByTimeRange retrieves audit logs within a time range
func (s *AnalyticServiceImpl) GetAuditLogsByTimeRange(ctx context.Context, companyID uuid.UUID, yearMonth string, startTime, endTime time.Time) ([]*domainModel.AuditLog, error) {
	return s.repo.GetAuditLogsByTimeRange(ctx, companyID, yearMonth, startTime, endTime)
}

// CreateAuditLog creates a new audit log
func (s *AnalyticServiceImpl) CreateAuditLog(ctx context.Context, log *domainModel.AuditLog) error {
	// Ensure YearMonth is set
	if log.YearMonth == "" {
		log.YearMonth = log.CreatedAt.Format("2006-01")
	}
	return s.repo.CreateAuditLog(ctx, log)
}

// ============================================
// Face Enrollment Logs methods implementation
// ============================================

// GetFaceEnrollmentLogs retrieves face enrollment logs for a company and month
func (s *AnalyticServiceImpl) GetFaceEnrollmentLogs(ctx context.Context, companyID uuid.UUID, yearMonth string, limit int) ([]*domainModel.FaceEnrollmentLog, error) {
	return s.repo.GetFaceEnrollmentLogs(ctx, companyID, yearMonth, limit)
}

// GetFaceEnrollmentLogsByEmployee retrieves face enrollment logs for a specific employee
func (s *AnalyticServiceImpl) GetFaceEnrollmentLogsByEmployee(ctx context.Context, companyID uuid.UUID, yearMonth string, employeeID uuid.UUID) ([]*domainModel.FaceEnrollmentLog, error) {
	return s.repo.GetFaceEnrollmentLogsByEmployee(ctx, companyID, yearMonth, employeeID)
}

// ============================================
// Attendance Records No Shift methods implementation
// ============================================

// GetAttendanceRecordsNoShift retrieves attendance records without shift
func (s *AnalyticServiceImpl) GetAttendanceRecordsNoShift(ctx context.Context, companyID uuid.UUID, yearMonth string, limit int) ([]*domainModel.AttendanceRecordNoShift, error) {
	return s.repo.GetAttendanceRecordsNoShift(ctx, companyID, yearMonth, limit)
}

// ============================================
// Additional helper methods
// ============================================

// GetAttendanceRecordsByUserTimeRange retrieves attendance records for a user within a time range
func (s *AnalyticServiceImpl) GetAttendanceRecordsByUserTimeRange(ctx context.Context, companyID, employeeID uuid.UUID, yearMonth string, startTime, endTime time.Time) ([]*domainModel.AttendanceRecordByUser, error) {
	return s.repo.GetAttendanceRecordsByUserTimeRange(ctx, companyID, employeeID, yearMonth, startTime, endTime)
}

// GetDailySummaryByUserDate retrieves a specific daily summary for a user and date
func (s *AnalyticServiceImpl) GetDailySummaryByUserDate(ctx context.Context, companyID, employeeID uuid.UUID, month string, workDate time.Time) (*domainModel.DailySummaryByUser, error) {
	return s.repo.GetDailySummaryByUserDate(ctx, companyID, employeeID, month, workDate)
}
