package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/youknow2509/cio_verify_face/server/service_analytic/internal/domain/model"
)

// IAnalyticRepository defines the interface for analytics data access (ScyllaDB + PostgreSQL)
type IAnalyticRepository interface {
	// ScyllaDB - Daily Summary queries
	GetDailySummariesByDate(ctx context.Context, companyID uuid.UUID, workDate time.Time) ([]*model.DailySummary, error)
	GetDailySummariesByMonth(ctx context.Context, companyID uuid.UUID, month string) ([]*model.DailySummary, error)
	GetDailySummariesByDateRange(ctx context.Context, companyID uuid.UUID, startDate, endDate time.Time) ([]*model.DailySummary, error)

	// ScyllaDB - Metrics queries
	GetDailyMetrics(ctx context.Context, companyID uuid.UUID, metricDate time.Time) (*model.AttendanceMetricsDaily, error)
	GetDailyMetricsByRange(ctx context.Context, companyID uuid.UUID, startDate, endDate time.Time) ([]*model.AttendanceMetricsDaily, error)
	GetHourlyMetrics(ctx context.Context, companyID uuid.UUID, metricDate time.Time) ([]*model.AttendanceMetricsHourly, error)

	// PostgreSQL - Employee queries (master data)
	GetEmployeeByID(ctx context.Context, employeeID uuid.UUID) (*model.Employee, error)
	GetEmployeesByCompany(ctx context.Context, companyID uuid.UUID) ([]*model.Employee, error)
	GetTotalEmployees(ctx context.Context, companyID *uuid.UUID) (int64, error)

	// PostgreSQL - User queries (master data)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*model.User, error)

	// PostgreSQL - Work shift queries (master data)
	GetWorkShiftByID(ctx context.Context, shiftID uuid.UUID) (*model.WorkShift, error)
	GetWorkShiftsByCompany(ctx context.Context, companyID uuid.UUID) ([]*model.WorkShift, error)

	// PostgreSQL - Company queries (master data)
	GetCompanyByID(ctx context.Context, companyID uuid.UUID) (*model.Company, error)

	// PostgreSQL - Attendance records queries (for device-level filtering)
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
