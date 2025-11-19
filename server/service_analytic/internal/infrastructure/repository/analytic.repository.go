package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/youknow2509/cio_verify_face/server/service_analytic/internal/domain/model"
	database "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/infrastructure/gen"
)

// AnalyticRepositoryImpl implements IAnalyticRepository
type AnalyticRepositoryImpl struct {
	scyllaSession *gocql.Session
	pgPool        *pgxpool.Pool
	queries       *database.Queries
}

// NewAnalyticRepository creates a new analytics repository instance
func NewAnalyticRepository(scyllaSession *gocql.Session, pgPool *pgxpool.Pool) *AnalyticRepositoryImpl {
	return &AnalyticRepositoryImpl{
		scyllaSession: scyllaSession,
		pgPool:        pgPool,
		queries:       database.New(pgPool),
	}
}

// ============================================
// ScyllaDB - Daily Summary queries
// ============================================

// GetDailySummariesByDate retrieves daily summaries for a specific date
func (r *AnalyticRepositoryImpl) GetDailySummariesByDate(ctx context.Context, companyID uuid.UUID, workDate time.Time) ([]*model.DailySummary, error) {
	month := workDate.Format("2006-01")

	query := `SELECT company_id, summary_month, work_date, employee_id, shift_id, 
		scheduled_in, scheduled_out, actual_check_in, actual_check_out,
		total_work_minutes, break_minutes, overtime_minutes, late_minutes, early_leave_minutes,
		attendance_status, attendance_percentage, notes, approved_by, approved_at, created_at, updated_at
		FROM daily_summaries_by_company_month 
		WHERE company_id = ? AND summary_month = ? AND work_date = ?`

	iter := r.scyllaSession.Query(query, uuidToGocql(companyID), month, workDate).Iter()

	var summaries []*model.DailySummary
	var summary model.DailySummary

	for iter.Scan(
		&summary.CompanyID, &summary.SummaryMonth, &summary.WorkDate, &summary.EmployeeID, &summary.ShiftID,
		&summary.ScheduledIn, &summary.ScheduledOut, &summary.ActualCheckIn, &summary.ActualCheckOut,
		&summary.TotalWorkMinutes, &summary.BreakMinutes, &summary.OvertimeMinutes, &summary.LateMinutes, &summary.EarlyLeaveMinutes,
		&summary.AttendanceStatus, &summary.AttendancePercentage, &summary.Notes, &summary.ApprovedBy, &summary.ApprovedAt,
		&summary.CreatedAt, &summary.UpdatedAt,
	) {
		summaryCopy := summary
		summaries = append(summaries, &summaryCopy)
	}

	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("failed to get daily summaries: %w", err)
	}

	return summaries, nil
}

// GetDailySummariesByMonth retrieves all daily summaries for a month
func (r *AnalyticRepositoryImpl) GetDailySummariesByMonth(ctx context.Context, companyID uuid.UUID, month string) ([]*model.DailySummary, error) {
	query := `SELECT company_id, summary_month, work_date, employee_id, shift_id, 
		scheduled_in, scheduled_out, actual_check_in, actual_check_out,
		total_work_minutes, break_minutes, overtime_minutes, late_minutes, early_leave_minutes,
		attendance_status, attendance_percentage, notes, approved_by, approved_at, created_at, updated_at
		FROM daily_summaries_by_company_month 
		WHERE company_id = ? AND summary_month = ?`

	iter := r.scyllaSession.Query(query, uuidToGocql(companyID), month).Iter()

	var summaries []*model.DailySummary
	var summary model.DailySummary

	for iter.Scan(
		&summary.CompanyID, &summary.SummaryMonth, &summary.WorkDate, &summary.EmployeeID, &summary.ShiftID,
		&summary.ScheduledIn, &summary.ScheduledOut, &summary.ActualCheckIn, &summary.ActualCheckOut,
		&summary.TotalWorkMinutes, &summary.BreakMinutes, &summary.OvertimeMinutes, &summary.LateMinutes, &summary.EarlyLeaveMinutes,
		&summary.AttendanceStatus, &summary.AttendancePercentage, &summary.Notes, &summary.ApprovedBy, &summary.ApprovedAt,
		&summary.CreatedAt, &summary.UpdatedAt,
	) {
		summaryCopy := summary
		summaries = append(summaries, &summaryCopy)
	}

	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("failed to get daily summaries by month: %w", err)
	}

	return summaries, nil
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

// pgtypeToTimePtr converts pgtype.Timestamptz to *time.Time
func pgtypeToTimePtr(ts pgtype.Timestamptz) *time.Time {
	if !ts.Valid {
		return nil
	}
	return &ts.Time
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

// pgtypeToFloat64Ptr converts pgtype.Numeric to *float64
func pgtypeToFloat64Ptr(num pgtype.Numeric) *float64 {
	if !num.Valid {
		return nil
	}
	f64, _ := num.Float64Value()
	return &f64.Float64
}

// pgtypeInt4ToIntPtr converts pgtype.Int4 to *int
func pgtypeInt4ToIntPtr(i pgtype.Int4) *int {
	if !i.Valid {
		return nil
	}
	val := int(i.Int32)
	return &val
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
		Salary:       pgtypeToFloat64Ptr(emp.Salary),
		Status:       int(emp.Status),
		CreatedAt:    pgtypeToTime(emp.CreatedAt),
		UpdatedAt:    pgtypeToTime(emp.UpdatedAt),
	}
}
