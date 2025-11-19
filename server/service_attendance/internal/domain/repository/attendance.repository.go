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
	// Add Attendance Record
	AddAttendanceRecord(ctx context.Context, input *model.AddAttendanceRecordInput) error
	AddDailySummaries(ctx context.Context, input *model.AddDailySummariesInput) error
	// Get
	GetAttendanceRecordCompany(ctx context.Context, input *model.GetAttendanceRecordCompanyInput) (*model.AttendanceRecordOutput, error)
	GetAttendanceRecordCompanyForEmployee(ctx context.Context, input *model.GetAttendanceRecordCompanyForEmployeeInput) (*model.AttendanceRecordOutput, error)
	GetDailySummarieCompany(ctx context.Context, input *model.GetDailySummariesCompanyInput) (*model.DailySummariesCompanyOutput, error)
	GetDailySummarieCompanyForEmployee(ctx context.Context, input *model.GetDailySummariesCompanyForEmployeeInput) (*model.DailySummariesEmployeeOutput, error)
	// Delete
	DeleteAttendanceRecordBeforeTimestamp(ctx context.Context, input *model.DeleteAttendanceRecordInput) error
	DeleteDailySummariesCompanyBeforeDate(ctx context.Context, input *model.DeleteDailySummariesInput) error
	DeleteDailySummariesEmployeeBeforeDate(ctx context.Context, input *model.DeleteDailySummariesEmployeeInput) error
	DeleteAttendanceRecord(ctx context.Context, input *model.DeleteAttendanceRecordInput) error
	DeleteDailySummariesCompany(ctx context.Context, input *model.DeleteDailySummariesInput) error
	DeleteDailySummariesEmployee(ctx context.Context, input *model.DeleteDailySummariesEmployeeInput) error
	// Update
	UpdateDailySummariesEmployee(ctx context.Context, input *model.UpdateDailySummariesEmployeeInput) error
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
