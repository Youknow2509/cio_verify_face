package service

import (
	"context"
	"errors"
	applicationErrors "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/application/errors"
	model "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/application/model"
)

// ============================================
// Attendance Service Interfaces
// ============================================
type IAttendanceService interface {
	GetDailyAttendanceSummaryForCompany(ctx context.Context, req *model.GetDailyAttendanceSummaryModel) (*model.GetDailyAttendanceSummaryResultModel, *applicationErrors.Error)
	GetDailyAttendanceSummaryEmployeeForCompany(ctx context.Context, req *model.GetDailyAttendanceSummaryEmployeeModel) (*model.GetDailyAttendanceSummaryEmployeeResultModel, *applicationErrors.Error)
	GetAttendanceRecordsEmployeeForConpany(ctx context.Context, req *model.GetAttendanceRecordsEmployeeModel) (*model.GetAttendanceRecordsCompanyResultModel, *applicationErrors.Error)
	GetAttendanceRecordsCompany(ctx context.Context, req *model.GetAttendanceRecordsCompanyModel) (*model.GetAttendanceRecordsCompanyResultModel, *applicationErrors.Error)
	AddAttendance(ctx context.Context, req *model.AddAttendanceModel) *applicationErrors.Error
	DeleteAttendanceRecord(ctx context.Context, req *model.DeleteAttendanceModel) *applicationErrors.Error
	DeleteAttendanceEmployeeBeforeTime(ctx context.Context, req *model.DeleteAttendanceModel) *applicationErrors.Error
	DeleteAttendanceNoShift(ctx context.Context, req *model.DeleteAttendanceRecordNoShiftModel) *applicationErrors.Error
	DeleteDailyAttendanceSummary(ctx context.Context, req *model.DeleteDailyAttendanceSummaryModel) *applicationErrors.Error
}

// Manager instance of attendance service
var _vIAttendanceService IAttendanceService

// Getter for attendance service instance
func GetAttendanceService() IAttendanceService {
	return _vIAttendanceService
}

// Setter for attendance service instance
func SetAttendanceService(service IAttendanceService) error {
	if service == nil {
		return errors.New("attendance service set is nil")
	}
	if _vIAttendanceService != nil {
		return errors.New("attendance service is already set")
	}
	_vIAttendanceService = service
	return nil
}
