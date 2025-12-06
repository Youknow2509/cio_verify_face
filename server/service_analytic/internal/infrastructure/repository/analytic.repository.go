package repository

import (
	"context"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/youknow2509/cio_verify_face/server/service_analytic/internal/domain/model"
	domainRepo "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/domain/repository"
	database "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/infrastructure/gen"
	"time"
)

// AnalyticRepositoryImpl implements IAnalyticRepository
type AnalyticRepositoryImpl struct {
	scyllaSession *gocql.Session
	pgPool        *pgxpool.Pool
	queries       *database.Queries
}

// GetDailySummariesByDatePage implements repository.IAnalyticRepository.
func (r *AnalyticRepositoryImpl) GetDailySummariesByDatePage(ctx context.Context, companyID uuid.UUID, workDate time.Time, pageState []byte, limit int) ([]*model.DailySummary, []byte, error) {
	month := workDate.Format("2006-01")
	query := `SELECT company_id, summary_month, work_date, employee_id, shift_id,
		actual_check_in, actual_check_out, attendance_status, late_minutes,
		early_leave_minutes, total_work_minutes, notes, updated_at
		FROM daily_summaries
		WHERE company_id = ? AND summary_month = ? AND work_date = ?`
	if pageState != nil {
		// Get first page
		iter := r.scyllaSession.Query(query, uuidToGocql(companyID), month, workDate).PageSize(limit).Iter()
		nextPage := iter.PageState()
		summaries, err := scanDailySummaries(iter)
		return summaries, nextPage, err
	} else {
		// Get subsequent pages
		iter := r.scyllaSession.Query(query, uuidToGocql(companyID), month, workDate).PageSize(limit).PageState(pageState).Iter()
		nextPage := iter.PageState()
		summaries, err := scanDailySummaries(iter)
		return summaries, nextPage, err
	}
}

// NewAnalyticRepository creates a new analytics repository instance
func NewAnalyticRepository(scyllaSession *gocql.Session, pgPool *pgxpool.Pool) domainRepo.IAnalyticRepository {
	return &AnalyticRepositoryImpl{
		scyllaSession: scyllaSession,
		pgPool:        pgPool,
		queries:       database.New(pgPool),
	}
}

// ============================================
// ScyllaDB - Attendance Records queries
// ============================================

// GetAttendanceRecords retrieves attendance records by company and month
func (r *AnalyticRepositoryImpl) GetAttendanceRecords(ctx context.Context, companyID uuid.UUID, yearMonth string, limit int) ([]*model.AttendanceRecord, error) {
	query := `SELECT company_id, year_month, record_time, employee_id, device_id, record_type,
		verification_method, verification_score, face_image_url, location_coordinates,
		metadata, sync_status, created_at
		FROM attendance_records
		WHERE company_id = ? AND year_month = ?
		LIMIT ?`

	iter := r.scyllaSession.Query(query, uuidToGocql(companyID), yearMonth, limit).Iter()
	return scanAttendanceRecords(iter)
}

// GetAttendanceRecordsByTimeRange retrieves attendance records within a time range
func (r *AnalyticRepositoryImpl) GetAttendanceRecordsByTimeRange(ctx context.Context, companyID uuid.UUID, yearMonth string, startTime, endTime time.Time) ([]*model.AttendanceRecord, error) {
	query := `SELECT company_id, year_month, record_time, employee_id, device_id, record_type,
		verification_method, verification_score, face_image_url, location_coordinates,
		metadata, sync_status, created_at
		FROM attendance_records
		WHERE company_id = ? AND year_month = ? AND record_time >= ? AND record_time <= ?`

	iter := r.scyllaSession.Query(query, uuidToGocql(companyID), yearMonth, startTime, endTime).Iter()
	return scanAttendanceRecords(iter)
}

// GetAttendanceRecordsByEmployee retrieves attendance records for a specific employee
func (r *AnalyticRepositoryImpl) GetAttendanceRecordsByEmployee(ctx context.Context, companyID uuid.UUID, yearMonth string, employeeID uuid.UUID) ([]*model.AttendanceRecord, error) {
	query := `SELECT company_id, year_month, record_time, employee_id, device_id, record_type,
		verification_method, verification_score, face_image_url, location_coordinates,
		metadata, sync_status, created_at
		FROM attendance_records_by_user
		WHERE company_id = ? AND year_month = ? AND employee_id = ?`

	iter := r.scyllaSession.Query(query, uuidToGocql(companyID), yearMonth, uuidToGocql(employeeID)).Iter()
	return scanAttendanceRecords(iter)
}

// CreateAttendanceRecord creates a new attendance record

// GetAttendanceRecordsByUser retrieves attendance records indexed by user
func (r *AnalyticRepositoryImpl) GetAttendanceRecordsByUser(ctx context.Context, companyID, employeeID uuid.UUID, yearMonth string, limit int) ([]*model.AttendanceRecordByUser, error) {
	query := `SELECT company_id, employee_id, year_month, record_time, device_id, record_type,
		verification_method, verification_score, face_image_url, location_coordinates,
		metadata, sync_status, created_at
		FROM attendance_records_by_user
		WHERE company_id = ? AND employee_id = ? AND year_month = ?
		LIMIT ?`

	iter := r.scyllaSession.Query(query, uuidToGocql(companyID), uuidToGocql(employeeID), yearMonth, limit).Iter()
	return scanAttendanceRecordsByUser(iter)
}

// GetAttendanceRecordsByUserTimeRange retrieves attendance records by user within a time range
func (r *AnalyticRepositoryImpl) GetAttendanceRecordsByUserTimeRange(ctx context.Context, companyID, employeeID uuid.UUID, yearMonth string, startTime, endTime time.Time) ([]*model.AttendanceRecordByUser, error) {
	query := `SELECT company_id, employee_id, year_month, record_time, device_id, record_type,
		verification_method, verification_score, face_image_url, location_coordinates,
		metadata, sync_status, created_at
		FROM attendance_records_by_user
		WHERE company_id = ? AND employee_id = ? AND year_month = ? AND record_time >= ? AND record_time <= ?`

	iter := r.scyllaSession.Query(query, uuidToGocql(companyID), uuidToGocql(employeeID), yearMonth, startTime, endTime).Iter()
	return scanAttendanceRecordsByUser(iter)
}

// ============================================
// ScyllaDB - Daily Summary queries
// ============================================

// GetDailySummariesByDate retrieves daily summaries for a specific date
func (r *AnalyticRepositoryImpl) GetDailySummariesByDate(ctx context.Context, companyID uuid.UUID, workDate time.Time) ([]*model.DailySummary, error) {
	month := workDate.Format("2006-01")

	query := `SELECT company_id, summary_month, work_date, employee_id, shift_id,
		actual_check_in, actual_check_out, attendance_status, late_minutes,
		early_leave_minutes, total_work_minutes, notes, updated_at
		FROM daily_summaries
		WHERE company_id = ? AND summary_month = ? AND work_date = ?`

	iter := r.scyllaSession.Query(query, uuidToGocql(companyID), month, workDate).Iter()
	return scanDailySummaries(iter)
}

// GetDailySummariesByMonth retrieves all daily summaries for a month
func (r *AnalyticRepositoryImpl) GetDailySummariesByMonth(ctx context.Context, companyID uuid.UUID, month string) ([]*model.DailySummary, error) {
	query := `SELECT company_id, summary_month, work_date, employee_id, shift_id,
		actual_check_in, actual_check_out, attendance_status, late_minutes,
		early_leave_minutes, total_work_minutes, notes, updated_at
		FROM daily_summaries
		WHERE company_id = ? AND summary_month = ?`

	iter := r.scyllaSession.Query(query, uuidToGocql(companyID), month).Iter()
	return scanDailySummaries(iter)
}

// GetDailySummariesByDateRange retrieves daily summaries for a date range
func (r *AnalyticRepositoryImpl) GetDailySummariesByDateRange(ctx context.Context, companyID uuid.UUID, startDate, endDate time.Time) ([]*model.DailySummary, error) {
	var allSummaries []*model.DailySummary

	// Generate list of months in the range
	currentMonth := time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, time.UTC)
	endMonth := time.Date(endDate.Year(), endDate.Month(), 1, 0, 0, 0, 0, time.UTC)

	for currentMonth.Before(endMonth) || currentMonth.Equal(endMonth) {
		month := currentMonth.Format("2006-01")

		summaries, err := r.GetDailySummariesByMonth(ctx, companyID, month)
		if err != nil {
			return nil, err
		}

		// Filter by date range
		for _, summary := range summaries {
			if (summary.WorkDate.After(startDate) || summary.WorkDate.Equal(startDate)) &&
				(summary.WorkDate.Before(endDate) || summary.WorkDate.Equal(endDate)) {
				allSummaries = append(allSummaries, summary)
			}
		}

		currentMonth = currentMonth.AddDate(0, 1, 0)
	}

	return allSummaries, nil
}

// GetDailySummariesByEmployeeMonth retrieves all daily summaries for a specific employee in a given month
func (r *AnalyticRepositoryImpl) GetDailySummariesByEmployeeMonth(ctx context.Context, companyID, employeeID uuid.UUID, month string) ([]*model.DailySummary, error) {
	query := `SELECT company_id, summary_month, work_date, employee_id, shift_id,
		actual_check_in, actual_check_out, attendance_status, late_minutes,
		early_leave_minutes, total_work_minutes, notes, updated_at
		FROM daily_summaries_by_user
		WHERE company_id = ? AND summary_month = ? AND employee_id = ?`

	iter := r.scyllaSession.Query(query, uuidToGocql(companyID), month, uuidToGocql(employeeID)).Iter()
	return scanDailySummaries(iter)
}

// GetDailySummariesByEmployeeDateRange retrieves daily summaries for a specific employee within a date range
func (r *AnalyticRepositoryImpl) GetDailySummariesByEmployeeDateRange(ctx context.Context, companyID, employeeID uuid.UUID, startDate, endDate time.Time) ([]*model.DailySummary, error) {
	var allSummaries []*model.DailySummary
	// This is not the most efficient way to query a date range in ScyllaDB without month partitioning,
	// but it's a simple approach for this case.
	currentMonth := time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, time.UTC)
	endMonth := time.Date(endDate.Year(), endDate.Month(), 1, 0, 0, 0, 0, time.UTC)

	for currentMonth.Before(endMonth) || currentMonth.Equal(endMonth) {
		month := currentMonth.Format("2006-01")
		summaries, err := r.GetDailySummariesByEmployeeMonth(ctx, companyID, employeeID, month)
		if err != nil {
			return nil, err
		}

		for _, summary := range summaries {
			if (summary.WorkDate.After(startDate) || summary.WorkDate.Equal(startDate)) &&
				(summary.WorkDate.Before(endDate) || summary.WorkDate.Equal(endDate)) {
				allSummaries = append(allSummaries, summary)
			}
		}
		currentMonth = currentMonth.AddDate(0, 1, 0)
	}
	return allSummaries, nil
}

// GetDailySummaryByEmployeeDate retrieves a specific daily summary
func (r *AnalyticRepositoryImpl) GetDailySummaryByEmployeeDate(ctx context.Context, companyID uuid.UUID, month string, workDate time.Time, employeeID uuid.UUID) (*model.DailySummary, error) {
	query := `SELECT company_id, summary_month, work_date, employee_id, shift_id,
		actual_check_in, actual_check_out, attendance_status, late_minutes,
		early_leave_minutes, total_work_minutes, notes, updated_at
		FROM daily_summaries_by_user
		WHERE company_id = ? AND summary_month = ? AND work_date = ? AND employee_id = ?`

	var summary model.DailySummary
	err := r.scyllaSession.Query(query, uuidToGocql(companyID), month, workDate, uuidToGocql(employeeID)).
		Scan(&summary.CompanyID, &summary.SummaryMonth, &summary.WorkDate, &summary.EmployeeID, &summary.ShiftID,
			&summary.ActualCheckIn, &summary.ActualCheckOut, &summary.AttendanceStatus, &summary.LateMinutes,
			&summary.EarlyLeaveMinutes, &summary.TotalWorkMinutes, &summary.Notes, &summary.UpdatedAt)

	if err == gocql.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get daily summary: %w", err)
	}

	return &summary, nil
}

// GetDailySummariesByUser retrieves daily summaries indexed by user
func (r *AnalyticRepositoryImpl) GetDailySummariesByUser(ctx context.Context, companyID, employeeID uuid.UUID, month string) ([]*model.DailySummaryByUser, error) {
	query := `SELECT company_id, employee_id, summary_month, work_date, shift_id,
		actual_check_in, actual_check_out, attendance_status, late_minutes,
		early_leave_minutes, total_work_minutes, notes, updated_at
		FROM daily_summaries_by_user
		WHERE company_id = ? AND employee_id = ? AND summary_month = ?`

	iter := r.scyllaSession.Query(query, uuidToGocql(companyID), uuidToGocql(employeeID), month).Iter()
	return scanDailySummariesByUser(iter)
}

// GetDailySummaryByUserDate retrieves a specific daily summary by user
func (r *AnalyticRepositoryImpl) GetDailySummaryByUserDate(ctx context.Context, companyID, employeeID uuid.UUID, month string, workDate time.Time) (*model.DailySummaryByUser, error) {
	query := `SELECT company_id, employee_id, summary_month, work_date, shift_id,
		actual_check_in, actual_check_out, attendance_status, late_minutes,
		early_leave_minutes, total_work_minutes, notes, updated_at
		FROM daily_summaries_by_user
		WHERE company_id = ? AND employee_id = ? AND summary_month = ? AND work_date = ?`

	var summary model.DailySummaryByUser
	err := r.scyllaSession.Query(query, uuidToGocql(companyID), uuidToGocql(employeeID), month, workDate).
		Scan(&summary.CompanyID, &summary.EmployeeID, &summary.SummaryMonth, &summary.WorkDate, &summary.ShiftID,
			&summary.ActualCheckIn, &summary.ActualCheckOut, &summary.AttendanceStatus, &summary.LateMinutes,
			&summary.EarlyLeaveMinutes, &summary.TotalWorkMinutes, &summary.Notes, &summary.UpdatedAt)

	if err == gocql.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get daily summary by user: %w", err)
	}

	return &summary, nil
}

// ============================================
// ScyllaDB - Audit Logs queries
// ============================================

// GetAuditLogs retrieves audit logs by company and month
func (r *AnalyticRepositoryImpl) GetAuditLogs(ctx context.Context, companyID uuid.UUID, yearMonth string, limit int) ([]*model.AuditLog, error) {
	query := `SELECT company_id, year_month, created_at, actor_id, action_category,
		action_name, resource_type, resource_id, details, ip_address, user_agent, status
		FROM audit_logs
		WHERE company_id = ? AND year_month = ?
		LIMIT ?`

	iter := r.scyllaSession.Query(query, uuidToGocql(companyID), yearMonth, limit).Iter()
	return scanAuditLogs(iter)
}

// GetAuditLogsByTimeRange retrieves audit logs within a time range
func (r *AnalyticRepositoryImpl) GetAuditLogsByTimeRange(ctx context.Context, companyID uuid.UUID, yearMonth string, startTime, endTime time.Time) ([]*model.AuditLog, error) {
	query := `SELECT company_id, year_month, created_at, actor_id, action_category,
		action_name, resource_type, resource_id, details, ip_address, user_agent, status
		FROM audit_logs
		WHERE company_id = ? AND year_month = ? AND created_at >= ? AND created_at <= ?`

	iter := r.scyllaSession.Query(query, uuidToGocql(companyID), yearMonth, startTime, endTime).Iter()
	return scanAuditLogs(iter)
}

// GetAuditLogsByActor retrieves audit logs for a specific actor
func (r *AnalyticRepositoryImpl) GetAuditLogsByActor(ctx context.Context, companyID uuid.UUID, yearMonth string, actorID uuid.UUID) ([]*model.AuditLog, error) {
	query := `SELECT company_id, year_month, created_at, actor_id, action_category,
		action_name, resource_type, resource_id, details, ip_address, user_agent, status
		FROM audit_logs
		WHERE company_id = ? AND year_month = ? AND actor_id = ?`

	iter := r.scyllaSession.Query(query, uuidToGocql(companyID), yearMonth, uuidToGocql(actorID)).Iter()
	return scanAuditLogs(iter)
}

// CreateAuditLog creates a new audit log
func (r *AnalyticRepositoryImpl) CreateAuditLog(ctx context.Context, log *model.AuditLog) error {
	query := `INSERT INTO audit_logs (company_id, year_month, created_at, actor_id,
		action_category, action_name, resource_type, resource_id, details, ip_address, user_agent, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	return r.scyllaSession.Query(query,
		uuidToGocql(log.CompanyID), log.YearMonth, log.CreatedAt, uuidToGocql(log.ActorID),
		log.ActionCategory, log.ActionName, log.ResourceType, log.ResourceID,
		log.Details, log.IPAddress, log.UserAgent, log.Status,
	).Exec()
}

// ============================================
// ScyllaDB - Face Enrollment Logs queries
// ============================================

// GetFaceEnrollmentLogs retrieves face enrollment logs by company and month
func (r *AnalyticRepositoryImpl) GetFaceEnrollmentLogs(ctx context.Context, companyID uuid.UUID, yearMonth string, limit int) ([]*model.FaceEnrollmentLog, error) {
	query := `SELECT company_id, year_month, created_at, employee_id, action_type,
		status, image_url, failure_reason, metadata
		FROM face_enrollment_logs
		WHERE company_id = ? AND year_month = ?
		LIMIT ?`

	iter := r.scyllaSession.Query(query, uuidToGocql(companyID), yearMonth, limit).Iter()
	return scanFaceEnrollmentLogs(iter)
}

// GetFaceEnrollmentLogsByTimeRange retrieves face enrollment logs within a time range
func (r *AnalyticRepositoryImpl) GetFaceEnrollmentLogsByTimeRange(ctx context.Context, companyID uuid.UUID, yearMonth string, startTime, endTime time.Time) ([]*model.FaceEnrollmentLog, error) {
	query := `SELECT company_id, year_month, created_at, employee_id, action_type,
		status, image_url, failure_reason, metadata
		FROM face_enrollment_logs
		WHERE company_id = ? AND year_month = ? AND created_at >= ? AND created_at <= ?`

	iter := r.scyllaSession.Query(query, uuidToGocql(companyID), yearMonth, startTime, endTime).Iter()
	return scanFaceEnrollmentLogs(iter)
}

// GetFaceEnrollmentLogsByEmployee retrieves face enrollment logs for a specific employee
func (r *AnalyticRepositoryImpl) GetFaceEnrollmentLogsByEmployee(ctx context.Context, companyID uuid.UUID, yearMonth string, employeeID uuid.UUID) ([]*model.FaceEnrollmentLog, error) {
	query := `SELECT company_id, year_month, created_at, employee_id, action_type,
		status, image_url, failure_reason, metadata
		FROM face_enrollment_logs
		WHERE company_id = ? AND year_month = ? AND employee_id = ?`

	iter := r.scyllaSession.Query(query, uuidToGocql(companyID), yearMonth, uuidToGocql(employeeID)).Iter()
	return scanFaceEnrollmentLogs(iter)
}

// ============================================
// ScyllaDB - Attendance Records No Shift queries
// ============================================

// GetAttendanceRecordsNoShift retrieves attendance records without shift by company and month
func (r *AnalyticRepositoryImpl) GetAttendanceRecordsNoShift(ctx context.Context, companyID uuid.UUID, yearMonth string, limit int) ([]*model.AttendanceRecordNoShift, error) {
	query := `SELECT company_id, year_month, record_time, employee_id, device_id,
		verification_method, verification_score, face_image_url, location_coordinates, created_at
		FROM attendance_records_no_shift
		WHERE company_id = ? AND year_month = ?
		LIMIT ?`

	iter := r.scyllaSession.Query(query, uuidToGocql(companyID), yearMonth, limit).Iter()
	return scanAttendanceRecordsNoShift(iter)
}

// GetAttendanceRecordsNoShiftByTimeRange retrieves attendance records without shift within a time range
func (r *AnalyticRepositoryImpl) GetAttendanceRecordsNoShiftByTimeRange(ctx context.Context, companyID uuid.UUID, yearMonth string, startTime, endTime time.Time) ([]*model.AttendanceRecordNoShift, error) {
	query := `SELECT company_id, year_month, record_time, employee_id, device_id,
		verification_method, verification_score, face_image_url, location_coordinates, created_at
		FROM attendance_records_no_shift
		WHERE company_id = ? AND year_month = ? AND record_time >= ? AND record_time <= ?`

	iter := r.scyllaSession.Query(query, uuidToGocql(companyID), yearMonth, startTime, endTime).Iter()
	return scanAttendanceRecordsNoShift(iter)
}

// ============================================
// ScyllaDB - Metrics queries
// ============================================

// GetDailyMetrics retrieves daily metrics for a specific date
func (r *AnalyticRepositoryImpl) GetDailyMetrics(ctx context.Context, companyID uuid.UUID, metricDate time.Time) (*model.AttendanceMetricsDaily, error) {
	query := `SELECT company_id, metric_date, total_attendance_records, unique_employees_count,
		present_count, late_count, absent_count, avg_work_hours, total_overtime_minutes,
		attendance_rate, punctuality_rate, created_at
		FROM attendance_metrics_daily 
		WHERE company_id = ? AND metric_date = ?`

	var metrics model.AttendanceMetricsDaily

	err := r.scyllaSession.Query(query, uuidToGocql(companyID), metricDate).Scan(
		&metrics.CompanyID, &metrics.MetricDate, &metrics.TotalAttendanceRecords, &metrics.UniqueEmployeesCount,
		&metrics.PresentCount, &metrics.LateCount, &metrics.AbsentCount, &metrics.AvgWorkHours, &metrics.TotalOvertimeMinutes,
		&metrics.AttendanceRate, &metrics.PunctualityRate, &metrics.CreatedAt,
	)

	if err != nil {
		if err == gocql.ErrNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get daily metrics: %w", err)
	}

	return &metrics, nil
}

// GetDailyMetricsByRange retrieves daily metrics for a date range
func (r *AnalyticRepositoryImpl) GetDailyMetricsByRange(ctx context.Context, companyID uuid.UUID, startDate, endDate time.Time) ([]*model.AttendanceMetricsDaily, error) {
	query := `SELECT company_id, metric_date, total_attendance_records, unique_employees_count,
		present_count, late_count, absent_count, avg_work_hours, total_overtime_minutes,
		attendance_rate, punctuality_rate, created_at
		FROM attendance_metrics_daily 
		WHERE company_id = ? AND metric_date >= ? AND metric_date <= ?`

	iter := r.scyllaSession.Query(query, uuidToGocql(companyID), startDate, endDate).Iter()

	var metricsList []*model.AttendanceMetricsDaily
	var metrics model.AttendanceMetricsDaily

	for iter.Scan(
		&metrics.CompanyID, &metrics.MetricDate, &metrics.TotalAttendanceRecords, &metrics.UniqueEmployeesCount,
		&metrics.PresentCount, &metrics.LateCount, &metrics.AbsentCount, &metrics.AvgWorkHours, &metrics.TotalOvertimeMinutes,
		&metrics.AttendanceRate, &metrics.PunctualityRate, &metrics.CreatedAt,
	) {
		metricsCopy := metrics
		metricsList = append(metricsList, &metricsCopy)
	}

	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("failed to get daily metrics by range: %w", err)
	}

	return metricsList, nil
}

// GetHourlyMetrics retrieves hourly metrics for a specific date
func (r *AnalyticRepositoryImpl) GetHourlyMetrics(ctx context.Context, companyID uuid.UUID, metricDate time.Time) ([]*model.AttendanceMetricsHourly, error) {
	query := `SELECT company_id, metric_date, metric_hour, total_checkins, total_checkouts,
		unique_employees, active_devices, avg_verification_score, peak_concurrent_users, created_at
		FROM attendance_metrics_hourly 
		WHERE company_id = ? AND metric_date = ?`

	iter := r.scyllaSession.Query(query, uuidToGocql(companyID), metricDate).Iter()

	var metricsList []*model.AttendanceMetricsHourly
	var metrics model.AttendanceMetricsHourly

	for iter.Scan(
		&metrics.CompanyID, &metrics.MetricDate, &metrics.MetricHour, &metrics.TotalCheckins, &metrics.TotalCheckouts,
		&metrics.UniqueEmployees, &metrics.ActiveDevices, &metrics.AvgVerificationScore, &metrics.PeakConcurrentUsers, &metrics.CreatedAt,
	) {
		metricsCopy := metrics
		metricsList = append(metricsList, &metricsCopy)
	}

	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("failed to get hourly metrics: %w", err)
	}

	return metricsList, nil
}

// ============================================
// PostgreSQL - Master Data queries
// ============================================

// GetEmployeeByID retrieves employee by ID from PostgreSQL
func (r *AnalyticRepositoryImpl) GetEmployeeByID(ctx context.Context, employeeID uuid.UUID) (*model.Employee, error) {
	employee, err := r.queries.GetEmployeeByID(ctx, uuidToPgtype(employeeID))
	if err != nil {
		return nil, fmt.Errorf("failed to get employee: %w", err)
	}

	return convertEmployeeToModel(&employee), nil
}

// GetEmployeesByCompany retrieves all employees for a company
func (r *AnalyticRepositoryImpl) GetEmployeesByCompany(ctx context.Context, companyID uuid.UUID) ([]*model.Employee, error) {
	employees, err := r.queries.GetEmployeesByCompany(ctx, uuidToPgtype(companyID))
	if err != nil {
		return nil, fmt.Errorf("failed to get employees: %w", err)
	}

	result := make([]*model.Employee, 0, len(employees))
	for i := range employees {
		result = append(result, convertEmployeeToModel(&employees[i]))
	}

	return result, nil
}

// GetTotalEmployees returns total employee count
func (r *AnalyticRepositoryImpl) GetTotalEmployees(ctx context.Context, companyID *uuid.UUID) (int64, error) {
	var pgCompanyID pgtype.UUID
	if companyID != nil {
		pgCompanyID = uuidToPgtype(*companyID)
	}

	count, err := r.queries.GetTotalEmployeesCount(ctx, pgCompanyID)
	if err != nil {
		return 0, fmt.Errorf("failed to get total employees: %w", err)
	}

	return count, nil
}

// GetUserByID retrieves user by ID from PostgreSQL
func (r *AnalyticRepositoryImpl) GetUserByID(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	user, err := r.queries.GetUserByID(ctx, uuidToPgtype(userID))
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &model.User{
		UserID:   pgtypeToUUID(user.UserID),
		FullName: user.FullName,
		Email:    user.Email,
		Role:     int(user.Role),
		// Note: sqlc query returns limited fields, full model may need schema update
	}, nil
}

// GetWorkShiftByID retrieves work shift by ID from PostgreSQL
func (r *AnalyticRepositoryImpl) GetWorkShiftByID(ctx context.Context, shiftID uuid.UUID) (*model.WorkShift, error) {
	shift, err := r.queries.GetWorkShiftByID(ctx, uuidToPgtype(shiftID))
	if err != nil {
		return nil, fmt.Errorf("failed to get work shift: %w", err)
	}

	return &model.WorkShift{
		ShiftID:   pgtypeToUUID(shift.ShiftID),
		CompanyID: pgtypeToUUID(shift.CompanyID),
		Name:      shift.Name,
		StartTime: shift.StartTime,
		EndTime:   shift.EndTime,
		// Note: sqlc query returns limited fields, full model may need schema update
	}, nil
}

// GetWorkShiftsByCompany retrieves all work shifts for a company
func (r *AnalyticRepositoryImpl) GetWorkShiftsByCompany(ctx context.Context, companyID uuid.UUID) ([]*model.WorkShift, error) {
	shifts, err := r.queries.GetWorkShiftsByCompany(ctx, uuidToPgtype(companyID))
	if err != nil {
		return nil, fmt.Errorf("failed to get work shifts: %w", err)
	}

	result := make([]*model.WorkShift, 0, len(shifts))
	for i := range shifts {
		result = append(result, &model.WorkShift{
			ShiftID:   pgtypeToUUID(shifts[i].ShiftID),
			CompanyID: pgtypeToUUID(shifts[i].CompanyID),
			Name:      shifts[i].Name,
			StartTime: shifts[i].StartTime,
			EndTime:   shifts[i].EndTime,
		})
	}

	return result, nil
}

// GetCompanyByID retrieves company by ID from PostgreSQL
func (r *AnalyticRepositoryImpl) GetCompanyByID(ctx context.Context, companyID uuid.UUID) (*model.Company, error) {
	company, err := r.queries.GetCompanyByID(ctx, uuidToPgtype(companyID))
	if err != nil {
		return nil, fmt.Errorf("failed to get company: %w", err)
	}

	address := pgtypeTextToStringPtr(company.Address)
	return &model.Company{
		CompanyID: pgtypeToUUID(company.CompanyID),
		Name:      company.Name,
		Address:   address,
		// Note: sqlc query returns limited fields, full model may need schema update
	}, nil
}

// GetEmployeeIDsByDeviceAndDate returns distinct employee IDs with records from a device on a date
func (r *AnalyticRepositoryImpl) GetEmployeeIDsByDeviceAndDate(ctx context.Context, deviceID uuid.UUID, date time.Time) ([]uuid.UUID, error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.AddDate(0, 0, 1)

	const q = `
		SELECT DISTINCT employee_id
		FROM attendance_records
		WHERE device_id = $1
		  AND timestamp >= $2 AND timestamp < $3
	`

	rows, err := r.pgPool.Query(ctx, q, uuidToPgtype(deviceID), start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to query employee IDs by device/date: %w", err)
	}
	defer rows.Close()

	ids := make([]uuid.UUID, 0)
	for rows.Next() {
		var id pgtype.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan employee_id: %w", err)
		}
		if id.Valid {
			ids = append(ids, id.Bytes)
		}
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("row iteration error: %w", rows.Err())
	}
	return ids, nil
}

// ============================================
// Helper functions for type conversion
// ============================================

// uuidToPgtype converts uuid.UUID to pgtype.UUID
func uuidToPgtype(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: id,
		Valid: true,
	}
}

// pgtypeToUUID converts pgtype.UUID to uuid.UUID
func pgtypeToUUID(pgID pgtype.UUID) uuid.UUID {
	if !pgID.Valid {
		return uuid.Nil
	}
	return pgID.Bytes
}

// pgtypeTextToStringPtr converts pgtype.Text to *string
func pgtypeTextToStringPtr(text pgtype.Text) *string {
	if !text.Valid {
		return nil
	}
	str := text.String
	return &str
}

// pgtypeToTime converts pgtype.Timestamptz to time.Time
func pgtypeToTime(ts pgtype.Timestamptz) time.Time {
	if !ts.Valid {
		return time.Time{}
	}
	return ts.Time
}

// pgtypeDateToTimePtr converts pgtype.Date to *time.Time
func pgtypeDateToTimePtr(date pgtype.Date) *time.Time {
	if !date.Valid {
		return nil
	}
	return &date.Time
}

// pgtypeToNullableFloat64Ptr converts pgtype.Numeric to model.NullableFloat64
func pgtypeToNullableFloat64Ptr(num pgtype.Numeric) model.NullableFloat64 {
	if !num.Valid {
		return model.NullableFloat64{Float: nil}
	}
	f64, err := num.Float64Value()
	if err != nil {
		return model.NullableFloat64{Float: nil}
	}
	val := f64.Float64
	return model.NullableFloat64{Float: &val}
}

// uuidToGocql converts uuid.UUID (google/uuid) to gocql.UUID for ScyllaDB queries
func uuidToGocql(id uuid.UUID) gocql.UUID {
	return gocql.UUID(id)
}

// convertEmployeeToModel converts database.Employee to model.Employee
func convertEmployeeToModel(emp *database.Employee) *model.Employee {
	return &model.Employee{
		EmployeeID:   pgtypeToUUID(emp.EmployeeID),
		CompanyID:    pgtypeToUUID(emp.CompanyID),
		EmployeeCode: emp.EmployeeCode,
		Department:   pgtypeTextToStringPtr(emp.Department),
		Position:     pgtypeTextToStringPtr(emp.Position),
		HireDate:     pgtypeDateToTimePtr(emp.HireDate),
		Salary:       pgtypeToNullableFloat64Ptr(emp.Salary),
		Status:       int(emp.Status),
		CreatedAt:    pgtypeToTime(emp.CreatedAt),
		UpdatedAt:    pgtypeToTime(emp.UpdatedAt),
	}
}

// ============================================
// ScyllaDB scan helper functions
// ============================================

// scanAttendanceRecords scans attendance records from iterator
func scanAttendanceRecords(iter *gocql.Iter) ([]*model.AttendanceRecord, error) {
	var records []*model.AttendanceRecord
	var record model.AttendanceRecord
	// Scan UUIDs into gocql.UUID first to avoid unmarshal errors
	var companyUUID, employeeUUID, deviceUUID gocql.UUID

	for iter.Scan(
		&companyUUID, &record.YearMonth, &record.RecordTime, &employeeUUID,
		&deviceUUID, &record.RecordType, &record.VerificationMethod, &record.VerificationScore,
		&record.FaceImageURL, &record.LocationCoordinates, &record.Metadata, &record.SyncStatus, &record.CreatedAt,
	) {
		// Convert gocql.UUID to uuid.UUID
		record.CompanyID = uuid.UUID(companyUUID)
		record.EmployeeID = uuid.UUID(employeeUUID)
		record.DeviceID = uuid.UUID(deviceUUID)
		recordCopy := record
		records = append(records, &recordCopy)
	}

	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("failed to scan attendance records: %w", err)
	}

	return records, nil
}

// scanAttendanceRecordsByUser scans attendance records by user from iterator
func scanAttendanceRecordsByUser(iter *gocql.Iter) ([]*model.AttendanceRecordByUser, error) {
	var records []*model.AttendanceRecordByUser
	var record model.AttendanceRecordByUser
	// Scan UUIDs into gocql.UUID first to avoid unmarshal errors
	var companyUUID, employeeUUID, deviceUUID gocql.UUID

	for iter.Scan(
		&companyUUID, &employeeUUID, &record.YearMonth, &record.RecordTime,
		&deviceUUID, &record.RecordType, &record.VerificationMethod, &record.VerificationScore,
		&record.FaceImageURL, &record.LocationCoordinates, &record.Metadata, &record.SyncStatus, &record.CreatedAt,
	) {
		// Convert gocql.UUID to uuid.UUID
		record.CompanyID = uuid.UUID(companyUUID)
		record.EmployeeID = uuid.UUID(employeeUUID)
		record.DeviceID = uuid.UUID(deviceUUID)
		recordCopy := record
		records = append(records, &recordCopy)
	}

	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("failed to scan attendance records by user: %w", err)
	}

	return records, nil
}

// scanDailySummaries scans daily summaries from iterator
func scanDailySummaries(iter *gocql.Iter) ([]*model.DailySummary, error) {
	var summaries []*model.DailySummary
	var summary model.DailySummary
	// mashal data special
	companyUuid := gocql.UUID{}
	employeeUuid := gocql.UUID{}
	shiftUuid := gocql.UUID{}
	for iter.Scan(
		&companyUuid, &summary.SummaryMonth, &summary.WorkDate, &employeeUuid, &shiftUuid,
		&summary.ActualCheckIn, &summary.ActualCheckOut, &summary.AttendanceStatus, &summary.LateMinutes,
		&summary.EarlyLeaveMinutes, &summary.TotalWorkMinutes, &summary.Notes, &summary.UpdatedAt,
	) {
		summary.CompanyID = uuid.UUID(companyUuid)
		summary.EmployeeID = uuid.UUID(employeeUuid)
		summary.ShiftID = uuid.UUID(shiftUuid)
		summaryCopy := summary
		summaries = append(summaries, &summaryCopy)
	}

	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("failed to scan daily summaries: %w", err)
	}

	return summaries, nil
}

// scanDailySummariesByUser scans daily summaries by user from iterator
func scanDailySummariesByUser(iter *gocql.Iter) ([]*model.DailySummaryByUser, error) {
	var summaries []*model.DailySummaryByUser
	var summary model.DailySummaryByUser
	// Scan UUIDs into gocql.UUID first to avoid unmarshal errors
	var companyUUID, employeeUUID, shiftUUID gocql.UUID

	for iter.Scan(
		&companyUUID, &employeeUUID, &summary.SummaryMonth, &summary.WorkDate, &shiftUUID,
		&summary.ActualCheckIn, &summary.ActualCheckOut, &summary.AttendanceStatus, &summary.LateMinutes,
		&summary.EarlyLeaveMinutes, &summary.TotalWorkMinutes, &summary.Notes, &summary.UpdatedAt,
	) {
		// Convert gocql.UUID to uuid.UUID
		summary.CompanyID = uuid.UUID(companyUUID)
		summary.EmployeeID = uuid.UUID(employeeUUID)
		summary.ShiftID = uuid.UUID(shiftUUID)
		summaryCopy := summary
		summaries = append(summaries, &summaryCopy)
	}

	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("failed to scan daily summaries by user: %w", err)
	}

	return summaries, nil
}

// scanAuditLogs scans audit logs from iterator
func scanAuditLogs(iter *gocql.Iter) ([]*model.AuditLog, error) {
	var logs []*model.AuditLog
	var log model.AuditLog
	// Scan UUIDs into gocql.UUID first to avoid unmarshal errors
	var companyUUID, actorUUID gocql.UUID

	for iter.Scan(
		&companyUUID, &log.YearMonth, &log.CreatedAt, &actorUUID, &log.ActionCategory,
		&log.ActionName, &log.ResourceType, &log.ResourceID, &log.Details, &log.IPAddress,
		&log.UserAgent, &log.Status,
	) {
		// Convert gocql.UUID to uuid.UUID (ResourceID is string, not UUID)
		log.CompanyID = uuid.UUID(companyUUID)
		log.ActorID = uuid.UUID(actorUUID)
		logCopy := log
		logs = append(logs, &logCopy)
	}

	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("failed to scan audit logs: %w", err)
	}

	return logs, nil
}

// scanFaceEnrollmentLogs scans face enrollment logs from iterator
func scanFaceEnrollmentLogs(iter *gocql.Iter) ([]*model.FaceEnrollmentLog, error) {
	var logs []*model.FaceEnrollmentLog
	var log model.FaceEnrollmentLog
	// Scan UUIDs into gocql.UUID first to avoid unmarshal errors
	var companyUUID, employeeUUID gocql.UUID

	for iter.Scan(
		&companyUUID, &log.YearMonth, &log.CreatedAt, &employeeUUID, &log.ActionType,
		&log.Status, &log.ImageURL, &log.FailureReason, &log.Metadata,
	) {
		// Convert gocql.UUID to uuid.UUID
		log.CompanyID = uuid.UUID(companyUUID)
		log.EmployeeID = uuid.UUID(employeeUUID)
		logCopy := log
		logs = append(logs, &logCopy)
	}

	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("failed to scan face enrollment logs: %w", err)
	}

	return logs, nil
}

// scanAttendanceRecordsNoShift scans attendance records no shift from iterator
func scanAttendanceRecordsNoShift(iter *gocql.Iter) ([]*model.AttendanceRecordNoShift, error) {
	var records []*model.AttendanceRecordNoShift
	var record model.AttendanceRecordNoShift
	// Scan UUIDs into gocql.UUID first to avoid unmarshal errors
	var companyUUID, employeeUUID, deviceUUID gocql.UUID

	for iter.Scan(
		&companyUUID, &record.YearMonth, &record.RecordTime, &employeeUUID, &deviceUUID,
		&record.VerificationMethod, &record.VerificationScore, &record.FaceImageURL,
		&record.LocationCoordinates, &record.CreatedAt,
	) {
		// Convert gocql.UUID to uuid.UUID
		record.CompanyID = uuid.UUID(companyUUID)
		record.EmployeeID = uuid.UUID(employeeUUID)
		record.DeviceID = uuid.UUID(deviceUUID)
		recordCopy := record
		records = append(records, &recordCopy)
	}

	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("failed to scan attendance records no shift: %w", err)
	}

	return records, nil
}
