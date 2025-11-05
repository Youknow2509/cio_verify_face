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
	CheckInUser(ctx context.Context, input *model.CheckInInput) *applicationErrors.Error
	CheckOutUser(ctx context.Context, input *model.CheckOutInput) *applicationErrors.Error
	GetRecords(ctx context.Context, input *model.GetAttendanceRecordsInput) ([]*model.AttendanceRecordOutput, *applicationErrors.Error)	
	GetMyRecords(ctx context.Context, input *model.GetMyRecordsInput) ([]*model.GetMyRecordsOutput, *applicationErrors.Error)
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
