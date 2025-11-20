package impl

import (
	"context"
	"github.com/youknow2509/cio_verify_face/server/service_attendance/internal/application/errors"
	"github.com/youknow2509/cio_verify_face/server/service_attendance/internal/application/model"
	"github.com/youknow2509/cio_verify_face/server/service_attendance/internal/application/service"
	domainCache "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/cache"
	domainLogger "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/logger"
	domainRepo "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/repository"
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

	// 2. Get shift info -> check in(0) or check out(1)
	// - Cache local cache, distributed cache(Redis)
	// - if not exist, get from DB and set to cache

	// 3. Add attendance record
	panic("unimplemented")
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
