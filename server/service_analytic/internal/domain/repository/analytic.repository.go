package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/youknow2509/cio_verify_face/server/service_analytic/internal/domain/model"
)

// IAnalyticRepository defines the interface for analytics data access (ScyllaDB + PostgreSQL)
type IAnalyticRepository interface {
	// ============================================
	// ScyllaDB - Attendance Records queries
	// ============================================

	// GetAttendanceRecords retrieves attendance records by company and month
	GetAttendanceRecords(ctx context.Context, companyID uuid.UUID, yearMonth string, limit int) ([]*model.AttendanceRecord, error)
	GetAttendanceRecordsByTimeRange(ctx context.Context, companyID uuid.UUID, yearMonth string, startTime, endTime time.Time) ([]*model.AttendanceRecord, error)
	GetAttendanceRecordsByEmployee(ctx context.Context, companyID uuid.UUID, yearMonth string, employeeID uuid.UUID) ([]*model.AttendanceRecord, error)

	// GetAttendanceRecordsByUser retrieves attendance records indexed by user
	GetAttendanceRecordsByUser(ctx context.Context, companyID, employeeID uuid.UUID, yearMonth string, limit int) ([]*model.AttendanceRecordByUser, error)
	GetAttendanceRecordsByUserTimeRange(ctx context.Context, companyID, employeeID uuid.UUID, yearMonth string, startTime, endTime time.Time) ([]*model.AttendanceRecordByUser, error)

	// ============================================
	// ScyllaDB - Daily Summary queries
	// ============================================

	GetDailySummariesByDate(ctx context.Context, companyID uuid.UUID, workDate time.Time) ([]*model.DailySummary, error)
	GetDailySummariesByDatePage(ctx context.Context, companyID uuid.UUID, workDate time.Time, pageState []byte, limit int) ([]*model.DailySummary, []byte, error)
	GetDailySummariesByMonth(ctx context.Context, companyID uuid.UUID, month string) ([]*model.DailySummary, error)
	GetDailySummariesByDateRange(ctx context.Context, companyID uuid.UUID, startDate, endDate time.Time) ([]*model.DailySummary, error)
	GetDailySummariesByEmployeeDateRange(ctx context.Context, companyID, employeeID uuid.UUID, startDate, endDate time.Time) ([]*model.DailySummary, error)
	GetDailySummaryByEmployeeDate(ctx context.Context, companyID uuid.UUID, month string, workDate time.Time, employeeID uuid.UUID) (*model.DailySummary, error)
	GetDailySummariesByEmployeeMonth(ctx context.Context, companyID, employeeID uuid.UUID, month string) ([]*model.DailySummary, error)

	// GetDailySummariesByUser retrieves daily summaries indexed by user
	GetDailySummariesByUser(ctx context.Context, companyID, employeeID uuid.UUID, month string) ([]*model.DailySummaryByUser, error)
	GetDailySummaryByUserDate(ctx context.Context, companyID, employeeID uuid.UUID, month string, workDate time.Time) (*model.DailySummaryByUser, error)

	// ============================================
	// ScyllaDB - Audit Logs queries
	// ============================================

	GetAuditLogs(ctx context.Context, companyID uuid.UUID, yearMonth string, limit int) ([]*model.AuditLog, error)
	GetAuditLogsByTimeRange(ctx context.Context, companyID uuid.UUID, yearMonth string, startTime, endTime time.Time) ([]*model.AuditLog, error)
	GetAuditLogsByActor(ctx context.Context, companyID uuid.UUID, yearMonth string, actorID uuid.UUID) ([]*model.AuditLog, error)
	CreateAuditLog(ctx context.Context, log *model.AuditLog) error

	// ============================================
	// ScyllaDB - Face Enrollment Logs queries
	// ============================================

	GetFaceEnrollmentLogs(ctx context.Context, companyID uuid.UUID, yearMonth string, limit int) ([]*model.FaceEnrollmentLog, error)
	GetFaceEnrollmentLogsByTimeRange(ctx context.Context, companyID uuid.UUID, yearMonth string, startTime, endTime time.Time) ([]*model.FaceEnrollmentLog, error)
	GetFaceEnrollmentLogsByEmployee(ctx context.Context, companyID uuid.UUID, yearMonth string, employeeID uuid.UUID) ([]*model.FaceEnrollmentLog, error)

	// ============================================
	// ScyllaDB - Attendance Records No Shift queries
	// ============================================

	GetAttendanceRecordsNoShift(ctx context.Context, companyID uuid.UUID, yearMonth string, limit int) ([]*model.AttendanceRecordNoShift, error)
	GetAttendanceRecordsNoShiftByTimeRange(ctx context.Context, companyID uuid.UUID, yearMonth string, startTime, endTime time.Time) ([]*model.AttendanceRecordNoShift, error)

	// ============================================
	// ScyllaDB - Metrics queries (existing)
	// ============================================

	GetDailyMetrics(ctx context.Context, companyID uuid.UUID, metricDate time.Time) (*model.AttendanceMetricsDaily, error)
	GetDailyMetricsByRange(ctx context.Context, companyID uuid.UUID, startDate, endDate time.Time) ([]*model.AttendanceMetricsDaily, error)
	GetHourlyMetrics(ctx context.Context, companyID uuid.UUID, metricDate time.Time) ([]*model.AttendanceMetricsHourly, error)

	// ============================================
	// PostgreSQL - Master Data queries (existing)
	// ============================================

	GetEmployeeByID(ctx context.Context, employeeID uuid.UUID) (*model.Employee, error)
	GetEmployeesByCompany(ctx context.Context, companyID uuid.UUID) ([]*model.Employee, error)
	GetTotalEmployees(ctx context.Context, companyID *uuid.UUID) (int64, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*model.User, error)
	GetWorkShiftByID(ctx context.Context, shiftID uuid.UUID) (*model.WorkShift, error)
	GetWorkShiftsByCompany(ctx context.Context, companyID uuid.UUID) ([]*model.WorkShift, error)
	GetCompanyByID(ctx context.Context, companyID uuid.UUID) (*model.Company, error)
	GetEmployeeIDsByDeviceAndDate(ctx context.Context, deviceID uuid.UUID, date time.Time) ([]uuid.UUID, error)
}

// Manager instance of analytic repository
var _vIAnalyticRepository IAnalyticRepository

// GetAnalyticRepository returns the singleton instance
func GetAnalyticRepository() IAnalyticRepository {
	return _vIAnalyticRepository
}

// SetAnalyticRepository sets the singleton instance
func SetAnalyticRepository(repo IAnalyticRepository) error {
	if repo == nil {
		return ErrRepositoryNil
	}
	if _vIAnalyticRepository != nil {
		return ErrRepositoryAlreadySet
	}
	_vIAnalyticRepository = repo
	return nil
}

var (
	ErrRepositoryNil        = &RepositoryError{Code: "REPO_NIL", Message: "repository instance is nil"}
	ErrRepositoryAlreadySet = &RepositoryError{Code: "REPO_ALREADY_SET", Message: "repository instance already set"}
)

// RepositoryError represents repository-specific errors
type RepositoryError struct {
	Code    string
	Message string
}

func (e *RepositoryError) Error() string {
	return e.Message
}
