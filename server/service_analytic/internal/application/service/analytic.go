package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	applicationErrors "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/application/errors"
	model "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/application/model"
	domainModel "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/domain/model"
)

// IAnalyticService interface defines the analytics service operations
type IAnalyticService interface {
	// ============================================
	// Existing report methods
	// ============================================
	
	// GetDailyReport returns daily attendance report
	GetDailyReport(ctx context.Context, input *model.DailyReportInput) (*model.DailyReportOutput, *applicationErrors.Error)
	GetDailyReportDetail(ctx context.Context, input *model.DailyDetailReportInput) (*model.DailyReportDetailOutput, *applicationErrors.Error)
	ExportDailyReportDetail(ctx context.Context, input *model.ExportDailyReportDetailInput) (*model.ExportDailyReportDetailOutput, *applicationErrors.Error)
	// GetSummaryReport returns monthly summary report
	GetSummaryReport(ctx context.Context, input *model.SummaryReportInput) (*model.SummaryReportOutput, *applicationErrors.Error)
	
	// ExportReport exports attendance report to file
	ExportReport(ctx context.Context, input *model.ExportReportInput) (*model.ExportReportOutput, *applicationErrors.Error)
	
	// GetHealthCheck returns service health status
	GetHealthCheck(ctx context.Context) (*model.HealthCheckOutput, *applicationErrors.Error)

	// ============================================
	// Attendance Records methods
	// ============================================
	
	// GetAttendanceRecords retrieves attendance records for a company and month
	GetAttendanceRecords(ctx context.Context, companyID uuid.UUID, yearMonth string, limit int) ([]*domainModel.AttendanceRecord, error)
	
	// GetAttendanceRecordsByTimeRange retrieves attendance records within a time range
	GetAttendanceRecordsByTimeRange(ctx context.Context, companyID uuid.UUID, yearMonth string, startTime, endTime time.Time) ([]*domainModel.AttendanceRecord, error)
	
	// GetAttendanceRecordsByEmployee retrieves attendance records for a specific employee
	GetAttendanceRecordsByEmployee(ctx context.Context, companyID uuid.UUID, yearMonth string, employeeID uuid.UUID) ([]*domainModel.AttendanceRecord, error)
	
	// GetAttendanceRecordsByUser retrieves attendance records indexed by user
	GetAttendanceRecordsByUser(ctx context.Context, companyID, employeeID uuid.UUID, yearMonth string, limit int) ([]*domainModel.AttendanceRecordByUser, error)
	
	// GetAttendanceRecordsByUserTimeRange retrieves attendance records for a user within a time range
	GetAttendanceRecordsByUserTimeRange(ctx context.Context, companyID, employeeID uuid.UUID, yearMonth string, startTime, endTime time.Time) ([]*domainModel.AttendanceRecordByUser, error)
	
	// ============================================
	// Daily Summary methods
	// ============================================
	
	// GetDailySummaries retrieves daily summaries for a month
	GetDailySummaries(ctx context.Context, companyID uuid.UUID, month string) ([]*domainModel.DailySummary, error)
	
	// GetDailySummaryByEmployeeDate retrieves a specific daily summary
	GetDailySummaryByEmployeeDate(ctx context.Context, companyID uuid.UUID, month string, workDate time.Time, employeeID uuid.UUID) (*domainModel.DailySummary, error)
	
	// GetDailySummariesByUser retrieves daily summaries for a user
	GetDailySummariesByUser(ctx context.Context, companyID, employeeID uuid.UUID, month string) ([]*domainModel.DailySummaryByUser, error)
	
	// GetDailySummaryByUserDate retrieves a specific daily summary for a user and date
	GetDailySummaryByUserDate(ctx context.Context, companyID, employeeID uuid.UUID, month string, workDate time.Time) (*domainModel.DailySummaryByUser, error)
	
	// ============================================
	// Audit Logs methods
	// ============================================
	
	// GetAuditLogs retrieves audit logs for a company and month
	GetAuditLogs(ctx context.Context, companyID uuid.UUID, yearMonth string, limit int) ([]*domainModel.AuditLog, error)
	
	// GetAuditLogsByTimeRange retrieves audit logs within a time range
	GetAuditLogsByTimeRange(ctx context.Context, companyID uuid.UUID, yearMonth string, startTime, endTime time.Time) ([]*domainModel.AuditLog, error)
	
	// CreateAuditLog creates a new audit log
	CreateAuditLog(ctx context.Context, log *domainModel.AuditLog) error
	
	// ============================================
	// Face Enrollment Logs methods
	// ============================================
	
	// GetFaceEnrollmentLogs retrieves face enrollment logs for a company and month
	GetFaceEnrollmentLogs(ctx context.Context, companyID uuid.UUID, yearMonth string, limit int) ([]*domainModel.FaceEnrollmentLog, error)
	
	// GetFaceEnrollmentLogsByEmployee retrieves face enrollment logs for a specific employee
	GetFaceEnrollmentLogsByEmployee(ctx context.Context, companyID uuid.UUID, yearMonth string, employeeID uuid.UUID) ([]*domainModel.FaceEnrollmentLog, error)
	
	// ============================================
	// Attendance Records No Shift methods
	// ============================================
	
	// GetAttendanceRecordsNoShift retrieves attendance records without shift
	GetAttendanceRecordsNoShift(ctx context.Context, companyID uuid.UUID, yearMonth string, limit int) ([]*domainModel.AttendanceRecordNoShift, error)
}

// Manager instance of analytic service
var _vIAnalyticService IAnalyticService

// GetAnalyticService returns the singleton instance
func GetAnalyticService() IAnalyticService {
	return _vIAnalyticService
}

// SetAnalyticService sets the singleton instance
func SetAnalyticService(service IAnalyticService) error {
	if service == nil {
		return errors.New("analytic service set is nil")
	}
	if _vIAnalyticService != nil {
		return errors.New("analytic service is already set")
	}
	_vIAnalyticService = service
	return nil
}
