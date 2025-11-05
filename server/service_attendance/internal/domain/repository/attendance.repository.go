package repository

import (
	"context"
	"errors"
	model "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/model"
)

// ============================================
// Interface for Attendance repository
// ============================================
type IAttendanceRepository interface {
	AddCheckInRecord(ctx context.Context, input *model.AddCheckInRecordInput) error
	AddCheckOutRecord(ctx context.Context, input *model.AddCheckOutRecordInput) error
	GetAttendanceRecordRangeTime(ctx context.Context, input *model.GetAttendanceRecordRangeTimeInput) ([]*model.AttendanceRecord, error)
}

// ============================================
// Variable for Attendance repository instance
// ============================================
var _vAttendanceRepository IAttendanceRepository

// ============================================
// Set the Attendance repository instance
// ============================================
func SetAttendanceRepository(v IAttendanceRepository) error {
	if v == nil {
		return errors.New("Attendance repository initialization failed, nil value")
	}
	if _vAttendanceRepository != nil {
		return errors.New("Attendance repository initialization failed, not nil")
	}
	_vAttendanceRepository = v
	return nil
}

// ============================================
// Get the Attendance repository instance
// ============================================
func GetAttendanceRepository() IAttendanceRepository {
	return _vAttendanceRepository
}