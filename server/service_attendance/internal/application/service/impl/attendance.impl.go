package impl

import (
	"context"
	applicationErrors "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/application/errors"
	model "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/application/model"
	service "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/application/service"
)

// ============================================
// Attendance Service
// ============================================
type AttendanceService struct{}

// CheckInUser implements service.IAttendanceService.
func (a *AttendanceService) CheckInUser(ctx context.Context, input *model.CheckInInput) *applicationErrors.Error {
	panic("unimplemented")
}

// CheckOutUser implements service.IAttendanceService.
func (a *AttendanceService) CheckOutUser(ctx context.Context, input *model.CheckOutInput) *applicationErrors.Error {
	panic("unimplemented")
}

// GetMyRecords implements service.IAttendanceService.
func (a *AttendanceService) GetMyRecords(ctx context.Context, input *model.GetMyRecordsInput) ([]*model.GetMyRecordsOutput, *applicationErrors.Error) {
	panic("unimplemented")
}

// GetRecordByID implements service.IAttendanceService.
func (a *AttendanceService) GetRecordByID(ctx context.Context, input *model.GetAttendanceRecordByIDInput) (*model.AttendanceRecordInfo, *applicationErrors.Error) {
	panic("unimplemented")
}

// GetRecords implements service.IAttendanceService.
func (a *AttendanceService) GetRecords(ctx context.Context, input *model.GetAttendanceRecordsInput) ([]*model.AttendanceRecordOutput, *applicationErrors.Error) {
	panic("unimplemented")
}

// New Attendance Service instance and impl interface
func NewAttendanceService() service.IAttendanceService {
	return &AttendanceService{}
}
