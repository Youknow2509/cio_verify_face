package model

import (
	"time"

	"github.com/google/uuid"
)

// AttendanceStatus represents the status of attendance
type AttendanceStatus int

const (
	AttendanceStatusPresent    AttendanceStatus = 0
	AttendanceStatusLate       AttendanceStatus = 1
	AttendanceStatusEarlyLeave AttendanceStatus = 2
	AttendanceStatusAbsent     AttendanceStatus = 3
)

// DailySummary represents the daily attendance summary model from ScyllaDB
// Table: daily_summaries_by_company_month
type DailySummary struct {
	CompanyID            uuid.UUID        `db:"company_id"`
	SummaryMonth         string           `db:"summary_month"` // YYYY-MM format
	WorkDate             time.Time        `db:"work_date"`
	EmployeeID           uuid.UUID        `db:"employee_id"`
	ShiftID              uuid.UUID        `db:"shift_id"`
	ScheduledIn          time.Time        `db:"scheduled_in"`
	ScheduledOut         time.Time        `db:"scheduled_out"`
	ActualCheckIn        *time.Time       `db:"actual_check_in"`
	ActualCheckOut       *time.Time       `db:"actual_check_out"`
	TotalWorkMinutes     int              `db:"total_work_minutes"`
	BreakMinutes         int              `db:"break_minutes"`
	OvertimeMinutes      int              `db:"overtime_minutes"`
	LateMinutes          int              `db:"late_minutes"`
	EarlyLeaveMinutes    int              `db:"early_leave_minutes"`
	AttendanceStatus     int              `db:"attendance_status"` // 0: PRESENT, 1: LATE, 2: EARLY_LEAVE, 3: ABSENT
	AttendancePercentage float64          `db:"attendance_percentage"`
	Notes                string           `db:"notes"`
	ApprovedBy           *uuid.UUID       `db:"approved_by"`
	ApprovedAt           *time.Time       `db:"approved_at"`
	CreatedAt            time.Time        `db:"created_at"`
	UpdatedAt            time.Time        `db:"updated_at"`
}

// AttendanceMetricsDaily represents daily rollup metrics from ScyllaDB
// Table: attendance_metrics_daily
type AttendanceMetricsDaily struct {
	CompanyID               uuid.UUID `db:"company_id"`
	MetricDate              time.Time `db:"metric_date"`
	TotalAttendanceRecords  int       `db:"total_attendance_records"`
	UniqueEmployeesCount    int       `db:"unique_employees_count"`
	PresentCount            int       `db:"present_count"`
	LateCount               int       `db:"late_count"`
	AbsentCount             int       `db:"absent_count"`
	AvgWorkHours            float64   `db:"avg_work_hours"`
	TotalOvertimeMinutes    int       `db:"total_overtime_minutes"`
	AttendanceRate          float64   `db:"attendance_rate"`
	PunctualityRate         float64   `db:"punctuality_rate"`
	CreatedAt               time.Time `db:"created_at"`
}

// AttendanceMetricsHourly represents hourly metrics from ScyllaDB
// Table: attendance_metrics_hourly
type AttendanceMetricsHourly struct {
	CompanyID              uuid.UUID `db:"company_id"`
	MetricDate             time.Time `db:"metric_date"`
	MetricHour             int       `db:"metric_hour"`
	TotalCheckins          int       `db:"total_checkins"`
	TotalCheckouts         int       `db:"total_checkouts"`
	UniqueEmployees        int       `db:"unique_employees"`
	ActiveDevices          int       `db:"active_devices"`
	AvgVerificationScore   float64   `db:"avg_verification_score"`
	PeakConcurrentUsers    int       `db:"peak_concurrent_users"`
	CreatedAt              time.Time `db:"created_at"`
}

// DailyAttendanceSummary represents the daily attendance summary model (for PostgreSQL compatibility)
type DailyAttendanceSummary struct {
	SummaryID            uuid.UUID        `db:"summary_id"`
	EmployeeID           uuid.UUID        `db:"employee_id"`
	ShiftID              *uuid.UUID       `db:"shift_id"`
	WorkDate             time.Time        `db:"work_date"`
	ScheduledIn          *time.Time       `db:"scheduled_in"`
	ScheduledOut         *time.Time       `db:"scheduled_out"`
	ActualCheckIn        *time.Time       `db:"actual_check_in"`
	ActualCheckOut       *time.Time       `db:"actual_check_out"`
	TotalWorkMinutes     int              `db:"total_work_minutes"`
	BreakMinutes         int              `db:"break_minutes"`
	OvertimeMinutes      int              `db:"overtime_minutes"`
	LateMinutes          int              `db:"late_minutes"`
	EarlyLeaveMinutes    int              `db:"early_leave_minutes"`
	Status               AttendanceStatus `db:"status"`
	AttendancePercentage float64          `db:"attendance_percentage"`
	Notes                *string          `db:"notes"`
	ApprovedBy           *uuid.UUID       `db:"approved_by"`
	ApprovedAt           *time.Time       `db:"approved_at"`
	CreatedAt            time.Time        `db:"created_at"`
	UpdatedAt            time.Time        `db:"updated_at"`
}

// Employee represents the employee model
type Employee struct {
	EmployeeID   uuid.UUID  `db:"employee_id"`
	CompanyID    uuid.UUID  `db:"company_id"`
	EmployeeCode string     `db:"employee_code"`
	Department   *string    `db:"department"`
	Position     *string    `db:"position"`
	HireDate     *time.Time `db:"hire_date"`
	Salary       *float64   `db:"salary"`
	Status       int        `db:"status"` // 0: active, 1: inactive, 2: on leave
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
}

// User represents the user model
type User struct {
	UserID    uuid.UUID  `db:"user_id"`
	CompanyID uuid.UUID  `db:"company_id"`
	FullName  string     `db:"full_name"`
	Email     string     `db:"email"`
	Phone     *string    `db:"phone"`
	Avatar    *string    `db:"avatar"`
	Role      int        `db:"role"` // 0: employee, 1: company_admin, 2: system_admin
	Status    int        `db:"status"` // 0: active, 1: inactive, 2: locked
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
}

// WorkShift represents the work shift model
type WorkShift struct {
	ShiftID       uuid.UUID  `db:"shift_id"`
	CompanyID     uuid.UUID  `db:"company_id"`
	Name          string     `db:"name"`
	StartTime     string     `db:"start_time"`
	EndTime       string     `db:"end_time"`
	BreakDuration int        `db:"break_duration"` // in minutes
	LateThreshold int        `db:"late_threshold"` // in minutes
	IsActive      bool       `db:"is_active"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at"`
}

// Company represents the company model
type Company struct {
	CompanyID   uuid.UUID  `db:"company_id"`
	Name        string     `db:"name"`
	Code        string     `db:"code"`
	Email       *string    `db:"email"`
	Phone       *string    `db:"phone"`
	Address     *string    `db:"address"`
	Status      int        `db:"status"` // 0: active, 1: inactive, 2: trial
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}

// AttendanceRecord represents individual check-in/out records
type AttendanceRecord struct {
	RecordID           uuid.UUID  `db:"record_id"`
	EmployeeID         uuid.UUID  `db:"employee_id"`
	DeviceID           uuid.UUID  `db:"device_id"`
	Timestamp          time.Time  `db:"timestamp"`
	RecordType         int        `db:"record_type"` // 0: CHECK_IN, 1: CHECK_OUT
	VerificationMethod string     `db:"verification_method"`
	VerificationScore  *float64   `db:"verification_score"`
	FaceImageURL       *string    `db:"face_image_url"`
	SyncStatus         int        `db:"sync_status"` // 0: SYNCED, 1: PENDING, 2: FAILED
	CreatedAt          time.Time  `db:"created_at"`
}
