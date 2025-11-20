package impl

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/youknow2509/cio_verify_face/server/service_attendance/internal/application/errors"
	"github.com/youknow2509/cio_verify_face/server/service_attendance/internal/application/model"
	"github.com/youknow2509/cio_verify_face/server/service_attendance/internal/application/service"
	"github.com/youknow2509/cio_verify_face/server/service_attendance/internal/constants"
	domainCache "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/cache"
	domainLogger "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/logger"
	domainModel "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/model"
	domainRepo "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/repository"
	utilsCache "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/shared/utils/cache"
	utilsCrypto "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/shared/utils/crypto"
)

// =================================================
// Attendance Service
// =================================================
type AttendanceService struct {
	attendanceRepo   domainRepo.IAttendanceRepository
	userRepo         domainRepo.IUserRepository
	logger           domainLogger.ILogger
	localCache       domainCache.ILocalCache
	distributedCache domainCache.IDistributedCache
}

// GetAttendanceRecordsEmployeeForConpany implements service.IAttendanceService.
func (a *AttendanceService) GetAttendanceRecordsEmployeeForConpany(ctx context.Context, req *model.GetAttendanceRecordsEmployeeModel) (*model.GetAttendanceRecordsCompanyResultModel, *errors.Error) {
	// 1. Check permission

	// 2. Check cache -> get from DB if not exist

	// 3. Return result
	panic("unimplemented")
}

// GetDailyAttendanceSummaryEmployeeForCompany implements service.IAttendanceService.
func (a *AttendanceService) GetDailyAttendanceSummaryEmployeeForCompany(ctx context.Context, req *model.GetDailyAttendanceSummaryEmployeeModel) (*model.GetDailyAttendanceSummaryEmployeeResultModel, *errors.Error) {
	// 1. Check permission

	// 2. Check cache -> get from DB if not exist

	// 3. Return result
	panic("unimplemented")
}

// GetDailyAttendanceSummaryForCompany implements service.IAttendanceService.
func (a *AttendanceService) GetDailyAttendanceSummaryForCompany(ctx context.Context, req *model.GetDailyAttendanceSummaryModel) (*model.GetDailyAttendanceSummaryResultModel, *errors.Error) {
	// 1. Check permission

	// 2. Check cache -> get from DB if not exist

	// 3. Return result
	panic("unimplemented")
}

// GetAttendanceRecordsCompany implements service.IAttendanceService.
func (a *AttendanceService) GetAttendanceRecordsCompany(ctx context.Context, req *model.GetAttendanceRecordsCompanyModel) (*model.GetAttendanceRecordsCompanyResultModel, *errors.Error) {
	// 1. Check permission

	// 2. Check cache -> get from DB if not exist

	// 3. Return result
	panic("unimplemented")
}

// AddAttendance implements service.IAttendanceService.
func (a *AttendanceService) AddAttendance(ctx context.Context, req *model.AddAttendanceModel) *errors.Error {
	// 1. Check permission
	if err := checkPermisionManager(
		ctx,
		*req.Session,
		req.CompanyID,
	); err != nil {
		a.logger.Warn("Permission Denied", "session", req.Session)
		return err
	}
	// 2. Get shift info -> check in(0) or check out(1)
	// - Cache local cache, distributed cache(Redis)
	// - if not exist, get from DB and set to cache
	keyListShiftEmployee := utilsCache.GetKeyListShiftTimeEmployee(
		utilsCrypto.GetHash(req.EmployeeID.String()),
	)
	cacheData := ""
	if data, err := a.localCache.Get(ctx, keyListShiftEmployee); err == nil {
		cacheData = data
		a.logger.Info("Get shift time from local cache", "key", keyListShiftEmployee)
	} else if data, err := a.distributedCache.Get(ctx, keyListShiftEmployee); err == nil {
		cacheData = data
		a.logger.Info("Get shift time from distributed cache", "key", keyListShiftEmployee)
		// Set to local cache
		_ = a.localCache.SetTTL(ctx, keyListShiftEmployee, cacheData, getTTLTimeCacheLocal(constants.TTL_CACHE_DEFAULT))
	}
	var shiftTimes []model.ShiftTimeEmployee
	if err := json.Unmarshal([]byte(cacheData), &shiftTimes); err != nil  && cacheData != "" {
		a.logger.Warn("Failed to unmarshal shift times", "error", err, "cache_data", cacheData)
	} else {
		respListShiftEmployee, err := a.userRepo.GetListTimeShiftEmployee(
			ctx,
			&domainModel.GetListTimeShiftEmployeeInput{
				EmployeeID: req.EmployeeID,
				CompanyID:  req.CompanyID,
			},
		)
		if err != nil {
			a.logger.Error("Failed to get list shift time employee from DB", "error", err)
			return &errors.Error{
				ErrorSystem: err,
				ErrorClient: "InternalError",
			}
		}
		// mapping data
		for _, item := range respListShiftEmployee {
			shiftTimes = append(shiftTimes, model.ShiftTimeEmployee{
				StartTime:             item.StartTime,
				EndTime:               item.EndTime,
				GracePeriodMinutes:    item.GracePeriodMinutes,
				EarlyDepartureMinutes: item.EarlyDepartureMinutes,
				WorkDays:              item.WorkDays,
				EffectiveFrom:         item.EffectiveFrom,
				EffectiveTo:           item.EffectiveTo,
			})
		}
		// Set to caches
		marshaledData, _ := json.Marshal(shiftTimes)
		cacheValue := string(marshaledData)
		_ = a.distributedCache.SetTTL(ctx, keyListShiftEmployee, cacheValue, constants.TTL_CACHE_DEFAULT)
		_ = a.localCache.SetTTL(ctx, keyListShiftEmployee, cacheValue, getTTLTimeCacheLocal(constants.TTL_CACHE_DEFAULT))
	}
	// 3. Check req is check in or check out
	var (
		isCheckIn       = true
		foundValidShift = false
		// Khởi tạo chênh lệch thời gian tối đa (12 giờ) để tìm ca làm việc gần nhất
		minDiffMinutes = float64(12 * 60)
	)

	// Chuyển đổi ngày trong tuần của Go (0=Chủ nhật) sang quy ước ISO 8601 (1=Thứ hai...7=Chủ nhật)
	currentWeekday := int32(req.RecordTime.Weekday())
	if currentWeekday == 0 {
		currentWeekday = 7
	}

	for _, shift := range shiftTimes {
		// A. Kiểm tra ngày hiệu lực của ca làm việc
		if req.RecordTime.Before(shift.EffectiveFrom) {
			continue
		}
		if shift.EffectiveTo != nil && req.RecordTime.After(*shift.EffectiveTo) {
			continue
		}

		// B. Kiểm tra xem hôm nay có phải là ngày làm việc theo lịch không
		isWorkDay := false
		for _, day := range shift.WorkDays {
			if day == currentWeekday {
				isWorkDay = true
				break
			}
		}
		if !isWorkDay {
			continue
		}

		// C. Chuẩn hóa thời gian bắt đầu/kết thúc ca về cùng ngày với ngày chấm công
		year, month, day := req.RecordTime.Date()
		shiftStart := time.Date(year, month, day, shift.StartTime.Hour(), shift.StartTime.Minute(), 0, 0, req.RecordTime.Location())
		shiftEnd := time.Date(year, month, day, shift.EndTime.Hour(), shift.EndTime.Minute(), 0, 0, req.RecordTime.Location())

		// Xử lý ca đêm (ví dụ: 22:00 - 06:00)
		// Nếu giờ kết thúc nhỏ hơn giờ bắt đầu, nghĩa là ca làm việc kết thúc vào ngày hôm sau
		if shiftEnd.Before(shiftStart) {
			// Nếu thời gian chấm công là sáng sớm (ví dụ: 05:00), so sánh với ca bắt đầu từ hôm qua
			if req.RecordTime.Hour() < 12 {
				shiftStart = shiftStart.Add(-24 * time.Hour)
			} else {
				// Nếu thời gian chấm công là tối muộn (ví dụ: 23:00), ca làm việc sẽ kết thúc vào ngày mai
				shiftEnd = shiftEnd.Add(24 * time.Hour)
			}
		}

		// D. Tính toán độ chênh lệch (tính theo phút) để xác định là Check-in hay Check-out
		diffStart := req.RecordTime.Sub(shiftStart).Minutes()
		if diffStart < 0 {
			diffStart = -diffStart
		}

		diffEnd := req.RecordTime.Sub(shiftEnd).Minutes()
		if diffEnd < 0 {
			diffEnd = -diffEnd
		}

		// Tìm ca làm việc có thời gian gần nhất với thời gian chấm công
		localMin := diffStart
		if diffEnd < localMin {
			localMin = diffEnd
		}

		if localMin < minDiffMinutes {
			minDiffMinutes = localMin
			foundValidShift = true
			// Nếu gần thời gian bắt đầu hơn -> Check In, ngược lại -> Check Out
			isCheckIn = diffStart <= diffEnd
		}
	}

	if !foundValidShift {
		a.logger.Warn("No scheduled shift found for the record time", "employee_id", req.EmployeeID, "record_time", req.RecordTime)
		return &errors.Error{
			ErrorClient: "NoScheduledShift",
			ErrorSystem: nil,
		}
	}
	// 4. Add attendance record
	inputAddAttendanceRecord := &domainModel.AddAttendanceRecordInput{
		CompanyID:  req.CompanyID,
		EmployeeID: req.EmployeeID,
		YearMonth:  req.RecordTime.Format("2006-01"),
		RecordTime: req.RecordTime,
		DeviceID:   req.DeviceID,
		RecordType: func() int {
			if isCheckIn {
				return 0
			} else {
				return 1
			}
		}(),
		VerificationMethod:  req.VerificationMethod,
		VerificationScore:   req.VerificationScore,
		FaceImageURL:        req.FaceImageURL,
		LocationCoordinates: req.LocationCoordinates,
		Metadata: map[string]string{
			"client_ip":    req.Session.ClientIp,
			"client_agent": req.Session.ClientAgent,
			"session_id":   req.Session.SessionId.String(),
		},
	}
	if err := a.attendanceRepo.AddAttendanceRecord(
		ctx,
		inputAddAttendanceRecord,
	); err != nil {
		a.logger.Error("Failed to add attendance record", "error", err)
		return &errors.Error{
			ErrorSystem: err,
			ErrorClient: "InternalError",
		}
	}
	return nil
}

// NewAttendanceService creates a new instance of AttendanceService
func NewAttendanceService() service.IAttendanceService {
	// Get dependencies
	attendanceRepo := domainRepo.GetAttendanceRepository()
	userRepo := domainRepo.GetUserRepository()
	logger := domainLogger.GetLogger()
	localCache, _ := domainCache.GetLocalCache()
	distributedCache, _ := domainCache.GetDistributedCache()

	// Create service instance
	return &AttendanceService{
		attendanceRepo:   attendanceRepo,
		userRepo:         userRepo,
		logger:           logger,
		localCache:       localCache,
		distributedCache: distributedCache,
	}
}

// =================================================
// Helper Functions
// =================================================
// handler ttl time cache local
func getTTLTimeCacheLocal(ttlDistributed int64) int64 {
	ttlLocal := ttlDistributed / 3
	if ttlLocal <= 0 {
		ttlLocal = 10 // default 10 seconds
	}
	return ttlLocal
}

// Check perrmission manager helper function
func checkPermisionManager(ctx context.Context, session model.SessionReq, companyReq uuid.UUID) *errors.Error {
	if session.Role == domainModel.RoleAdmin {
		return nil
	}
	if session.Role == domainModel.RoleManager && session.CompanyId == companyReq {
		return nil
	}
	return &errors.Error{
		ErrorClient: "PermissionDenied",
	}
}

// check permission employee helper function
func checkPermissionEmployee(ctx context.Context, session model.SessionReq, employeeReq uuid.UUID) *errors.Error {
	if err := checkPermisionManager(ctx, session, employeeReq); err == nil {
		return nil
	}
	if session.UserId == employeeReq {
		return nil
	}
	return &errors.Error{
		ErrorClient: "PermissionDenied",
	}
}
