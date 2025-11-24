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
// Table: daily_summaries
// PRIMARY KEY ((company_id, summary_month), work_date, employee_id)
type DailySummary struct {
	CompanyID         uuid.UUID  `db:"company_id"`
	SummaryMonth      string     `db:"summary_month"` // YYYY-MM format
	WorkDate          time.Time  `db:"work_date"`
	EmployeeID        uuid.UUID  `db:"employee_id"`
	ShiftID           uuid.UUID  `db:"shift_id"`
	ActualCheckIn     *time.Time `db:"actual_check_in"`
	ActualCheckOut    *time.Time `db:"actual_check_out"`
	AttendanceStatus  int        `db:"attendance_status"`
	LateMinutes       int        `db:"late_minutes"`
	EarlyLeaveMinutes int        `db:"early_leave_minutes"`
	TotalWorkMinutes  int        `db:"total_work_minutes"`
	Notes             string     `db:"notes"`
	UpdatedAt         time.Time  `db:"updated_at"`

	// Calculated fields (not in ScyllaDB, computed on demand)
	OvertimeMinutes      int     `db:"-"` // Calculated field
	AttendancePercentage float64 `db:"-"` // Calculated field
}

// AttendanceMetricsDaily represents daily rollup metrics from ScyllaDB
// Table: attendance_metrics_daily
type AttendanceMetricsDaily struct {
	CompanyID              uuid.UUID       `db:"company_id"`
	MetricDate             time.Time       `db:"metric_date"`
	TotalAttendanceRecords int             `db:"total_attendance_records"`
	UniqueEmployeesCount   int             `db:"unique_employees_count"`
	PresentCount           int             `db:"present_count"`
	LateCount              int             `db:"late_count"`
	AbsentCount            int             `db:"absent_count"`
	AvgWorkHours           NullableFloat32 `db:"avg_work_hours"`
	TotalOvertimeMinutes   int             `db:"total_overtime_minutes"`
	AttendanceRate         NullableFloat32 `db:"attendance_rate"`
	PunctualityRate        NullableFloat32 `db:"punctuality_rate"`
	CreatedAt              time.Time       `db:"created_at"`
}

// AttendanceMetricsHourly represents hourly metrics from ScyllaDB
// Table: attendance_metrics_hourly
type AttendanceMetricsHourly struct {
	CompanyID            uuid.UUID       `db:"company_id"`
	MetricDate           time.Time       `db:"metric_date"`
	MetricHour           int             `db:"metric_hour"`
	TotalCheckins        int             `db:"total_checkins"`
	TotalCheckouts       int             `db:"total_checkouts"`
	UniqueEmployees      int             `db:"unique_employees"`
	ActiveDevices        int             `db:"active_devices"`
	AvgVerificationScore NullableFloat32 `db:"avg_verification_score"`
	PeakConcurrentUsers  int             `db:"peak_concurrent_users"`
	CreatedAt            time.Time       `db:"created_at"`
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
	EmployeeID   uuid.UUID       `db:"employee_id"`
	CompanyID    uuid.UUID       `db:"company_id"`
	EmployeeCode string          `db:"employee_code"`
	Department   *string         `db:"department"`
	Position     *string         `db:"position"`
	HireDate     *time.Time      `db:"hire_date"`
	Salary       NullableFloat64 `db:"salary"`
	Status       int             `db:"status"` // 0: active, 1: inactive, 2: on leave
	CreatedAt    time.Time       `db:"created_at"`
	UpdatedAt    time.Time       `db:"updated_at"`
}

// User represents the user model
type User struct {
	UserID    uuid.UUID `db:"user_id"`
	CompanyID uuid.UUID `db:"company_id"`
	FullName  string    `db:"full_name"`
	Email     string    `db:"email"`
	Phone     *string   `db:"phone"`
	Avatar    *string   `db:"avatar"`
	Role      int       `db:"role"`   // 0: employee, 1: company_admin, 2: system_admin
	Status    int       `db:"status"` // 0: active, 1: inactive, 2: locked
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// WorkShift represents the work shift model
type WorkShift struct {
	ShiftID       uuid.UUID `db:"shift_id"`
	CompanyID     uuid.UUID `db:"company_id"`
	Name          string    `db:"name"`
	StartTime     string    `db:"start_time"`
	EndTime       string    `db:"end_time"`
	BreakDuration int       `db:"break_duration"` // in minutes
	LateThreshold int       `db:"late_threshold"` // in minutes
	IsActive      bool      `db:"is_active"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

// Company represents the company model
type Company struct {
	CompanyID uuid.UUID `db:"company_id"`
	Name      string    `db:"name"`
	Code      string    `db:"code"`
	Email     *string   `db:"email"`
	Phone     *string   `db:"phone"`
	Address   *string   `db:"address"`
	Status    int       `db:"status"` // 0: active, 1: inactive, 2: trial
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// AttendanceRecord represents individual check-in/out records from ScyllaDB
// Table: attendance_records
// PRIMARY KEY ((company_id, year_month), record_time, employee_id)
type AttendanceRecord struct {
	CompanyID           uuid.UUID         `db:"company_id"`
	YearMonth           string            `db:"year_month"` // YYYY-MM format
	RecordTime          time.Time         `db:"record_time"`
	EmployeeID          uuid.UUID         `db:"employee_id"`
	DeviceID            uuid.UUID         `db:"device_id"`
	RecordType          int               `db:"record_type"` // 0: CHECK_IN, 1: CHECK_OUT
	VerificationMethod  string            `db:"verification_method"`
	VerificationScore   NullableFloat32   `db:"verification_score"`
	FaceImageURL        string            `db:"face_image_url"`
	LocationCoordinates string            `db:"location_coordinates"`
	Metadata            map[string]string `db:"metadata"`
	SyncStatus          string            `db:"sync_status"`
	CreatedAt           time.Time         `db:"created_at"`
}

// AttendanceRecordByUser represents attendance records indexed by user from ScyllaDB
// Table: attendance_records_by_user
// PRIMARY KEY ((company_id, employee_id, year_month), record_time)
type AttendanceRecordByUser struct {
	CompanyID           uuid.UUID         `db:"company_id"`
	EmployeeID          uuid.UUID         `db:"employee_id"`
	YearMonth           string            `db:"year_month"` // YYYY-MM format
	RecordTime          time.Time         `db:"record_time"`
	DeviceID            uuid.UUID         `db:"device_id"`
	RecordType          int               `db:"record_type"` // 0: CHECK_IN, 1: CHECK_OUT
	VerificationMethod  string            `db:"verification_method"`
	VerificationScore   NullableFloat32   `db:"verification_score"`
	FaceImageURL        string            `db:"face_image_url"`
	LocationCoordinates string            `db:"location_coordinates"`
	Metadata            map[string]string `db:"metadata"`
	SyncStatus          string            `db:"sync_status"`
	CreatedAt           time.Time         `db:"created_at"`
}

// DailySummaryByUser represents daily summaries indexed by user from ScyllaDB
// Table: daily_summaries_by_user
// PRIMARY KEY ((company_id, employee_id, summary_month), work_date)
type DailySummaryByUser struct {
	CompanyID         uuid.UUID  `db:"company_id"`
	EmployeeID        uuid.UUID  `db:"employee_id"`
	SummaryMonth      string     `db:"summary_month"` // YYYY-MM format
	WorkDate          time.Time  `db:"work_date"`
	ShiftID           uuid.UUID  `db:"shift_id"`
	ActualCheckIn     *time.Time `db:"actual_check_in"`
	ActualCheckOut    *time.Time `db:"actual_check_out"`
	AttendanceStatus  int        `db:"attendance_status"`
	LateMinutes       int        `db:"late_minutes"`
	EarlyLeaveMinutes int        `db:"early_leave_minutes"`
	TotalWorkMinutes  int        `db:"total_work_minutes"`
	Notes             string     `db:"notes"`
	UpdatedAt         time.Time  `db:"updated_at"`
}

// AuditLog represents audit log entries from ScyllaDB
// Table: audit_logs
// PRIMARY KEY ((company_id, year_month), created_at, actor_id)
type AuditLog struct {
	CompanyID      uuid.UUID         `db:"company_id"`
	YearMonth      string            `db:"year_month"` // YYYY-MM format
	CreatedAt      time.Time         `db:"created_at"`
	ActorID        uuid.UUID         `db:"actor_id"`
	ActionCategory string            `db:"action_category"`
	ActionName     string            `db:"action_name"`
	ResourceType   string            `db:"resource_type"`
	ResourceID     string            `db:"resource_id"`
	Details        map[string]string `db:"details"`
	IPAddress      string            `db:"ip_address"`
	UserAgent      string            `db:"user_agent"`
	Status         string            `db:"status"`
}

// FaceEnrollmentLog represents face enrollment log entries from ScyllaDB
// Table: face_enrollment_logs
// PRIMARY KEY ((company_id, year_month), created_at, employee_id)
type FaceEnrollmentLog struct {
	CompanyID     uuid.UUID         `db:"company_id"`
	YearMonth     string            `db:"year_month"` // YYYY-MM format
	CreatedAt     time.Time         `db:"created_at"`
	EmployeeID    uuid.UUID         `db:"employee_id"`
	ActionType    string            `db:"action_type"` // e.g., "enrollment", "update", "deletion"
	Status        string            `db:"status"`      // e.g., "success", "failure"
	ImageURL      string            `db:"image_url"`
	FailureReason string            `db:"failure_reason"`
	Metadata      map[string]string `db:"metadata"`
}

// AttendanceRecordNoShift represents attendance records without shift information from ScyllaDB
// Table: attendance_records_no_shift
// PRIMARY KEY ((company_id, year_month), record_time)
type AttendanceRecordNoShift struct {
	CompanyID           uuid.UUID       `db:"company_id"`
	YearMonth           string          `db:"year_month"` // YYYY-MM format
	RecordTime          time.Time       `db:"record_time"`
	EmployeeID          uuid.UUID       `db:"employee_id"`
	DeviceID            uuid.UUID       `db:"device_id"`
	VerificationMethod  string          `db:"verification_method"`
	VerificationScore   NullableFloat32 `db:"verification_score"`
	FaceImageURL        string          `db:"face_image_url"`
	LocationCoordinates string          `db:"location_coordinates"`
	CreatedAt           time.Time       `db:"created_at"`
}
