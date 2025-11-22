package model

import (
	"time"

	"github.com/google/uuid"
)

// For Delete daily attendance summary
type DeleteDailyAttendanceSummaryModel struct {
	Session        *SessionReq     `json:"session"`
	ServiceSession *ServiceSession `json:"service_session,omitempty"`
	//
	CompanyID    uuid.UUID `json:"company_id"`
	WorkDate     time.Time `json:"work_date"`
	SummaryMonth string    `json:"summary_month"` // Format: YYYY-MM
}

// For Delete attendance records no shift
type DeleteAttendanceRecordNoShiftModel struct {
	Session        *SessionReq     `json:"session"`
	ServiceSession *ServiceSession `json:"service_session,omitempty"`
	//
	CompanyID uuid.UUID `json:"company_id"`
	YearMonth string    `json:"year_month"` // Format: YYYY-MM
}

// For Delete attendace records
type DeleteAttendanceModel struct {
	Session        *SessionReq     `json:"session"`
	ServiceSession *ServiceSession `json:"service_session,omitempty"`
	//
	CompanyID  uuid.UUID `json:"company_id"`
	YearMonth  string    `json:"year_month"` // Format: YYYY-MM
	Time       time.Time `json:"time,omitempty"`
	EmployeeId uuid.UUID `json:"employee_id,omitempty"`
}

// ShiftTimeEmployee
type ShiftTimeEmployee struct {
	ShiftID               uuid.UUID  `json:"shift_id"`
	StartTime             time.Time  `json:"start_time"`
	EndTime               time.Time  `json:"end_time"`
	GracePeriodMinutes    int        `json:"grace_period_minutes"`
	EarlyDepartureMinutes int        `json:"early_departure_minutes"`
	WorkDays              []int32    `json:"work_days"`
	EffectiveFrom         time.Time  `json:"effective_from"`
	EffectiveTo           *time.Time `json:"effective_to"`
}

// For GetDailyAttendanceSummary
type GetDailyAttendanceSummaryModel struct {
	Session *SessionReq `json:"session"`
	//
	CompanyID    uuid.UUID `json:"company_id"`
	SummaryMonth string    `json:"summary_month"` // Format: YYYY-MM
	WorkDate     time.Time `json:"work_date"`
	PageSize     int       `json:"page_size,omitempty"`
	PageStage    []byte    `json:"page_stage,omitempty"`
}

type GetDailyAttendanceSummaryResultModel struct {
	Records       []DailySummariesCompanyInfo `json:"records"`
	PageStageNext string                      `json:"page_stage_next,omitempty"`
	PageSize      int                         `json:"page_size,omitempty"`
}

type DailySummariesCompanyInfo struct {
	CompanyId         uuid.UUID `json:"company_id"`
	SummaryMonth      string    `json:"summary_month"`
	WorkDate          time.Time `json:"work_date"`
	EmployeeId        uuid.UUID `json:"employee_id"`
	ShiftId           uuid.UUID `json:"shift_id"`
	ActualCheckIn     time.Time `json:"actual_check_in"`
	ActualCheckOut    time.Time `json:"actual_check_out"`
	AttendanceStatus  string    `json:"attendance_status"`
	LateMinutes       int       `json:"late_minutes"`
	EarlyLeaveMinutes int       `json:"early_leave_minutes"`
	TotalWorkMinutes  int       `json:"total_work_minutes"`
	Notes             string    `json:"notes"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// For GetDailyAttendanceSummaryEmployee
type GetDailyAttendanceSummaryEmployeeResultModel struct {
	Records       []DailySummariesEmployeeInfo `json:"records"`
	PageStageNext string                       `json:"page_stage_next,omitempty"`
	PageSize      int                          `json:"page_size,omitempty"`
}

type DailySummariesEmployeeInfo struct {
	CompanyId         uuid.UUID `json:"company_id"`
	SummaryMonth      string    `json:"summary_month"`
	WorkDate          time.Time `json:"work_date"`
	EmployeeId        uuid.UUID `json:"employee_id"`
	ShiftId           uuid.UUID `json:"shift_id"`
	ActualCheckIn     time.Time `json:"actual_check_in"`
	ActualCheckOut    time.Time `json:"actual_check_out"`
	AttendanceStatus  string    `json:"attendance_status"`
	LateMinutes       int       `json:"late_minutes"`
	EarlyLeaveMinutes int       `json:"early_leave_minutes"`
	TotalWorkMinutes  int       `json:"total_work_minutes"`
	Notes             string    `json:"notes"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type GetDailyAttendanceSummaryEmployeeModel struct {
	Session *SessionReq `json:"session"`
	//
	CompanyID    uuid.UUID `json:"company_id"`
	EmployeeID   uuid.UUID `json:"employee_id"`
	SummaryMonth string    `json:"summary_month" ` // Format: YYYY-MM,
	PageSize     int       `json:"page_size"`
	PageStage    []byte    `json:"page_stage"`
}

type GetAttendanceRecordsCompanyResultModel struct {
	Records       []AttendanceRecordInfo `json:"records"`
	PageStageNext string                 `json:"page_stage_next,omitempty"`
	PageSize      int                    `json:"page_size,omitempty"`
}

type AttendanceRecordInfo struct {
	CompanyID           uuid.UUID         `json:"company_id"`
	YearMonth           string            `json:"year_month"`
	RecordTime          time.Time         `json:"record_time"`
	EmployeeID          uuid.UUID         `json:"employee_id"`
	DeviceID            uuid.UUID         `json:"device_id"`
	RecordType          int               `json:"record_type"`
	VerificationMethod  string            `json:"verification_method"`
	VerificationScore   float64           `json:"verification_score"`
	FaceImageURL        string            `json:"face_image_url"`
	LocationCoordinates string            `json:"location_coordinates"`
	Metadata            map[string]string `json:"metadata"`
	SyncStatus          string            `json:"sync_status"`
	CreatedAt           time.Time         `json:"created_at"`
}

// For GetAttendanceRecordsEmployee
type GetAttendanceRecordsEmployeeModel struct {
	//
	CompanyID  uuid.UUID `json:"company_id"`
	YearMonth  string    `json:"year_month" ` // Format: YYYY-MM
	EmployeeID uuid.UUID `json:"employee_id"`
	PageSize   int       `json:"page_size,omitempty"`
	PageStage  []byte    `json:"page_stage,omitempty"`
	//
	Session *SessionReq `json:"session"`
}

type GetAttendanceRecordsCompanyModel struct {
	//
	CompanyID uuid.UUID `json:"company_id"`
	YearMonth string    `json:"year_month" ` // Format: YYYY-MM
	PageSize  int       `json:"page_size"`
	PageStage []byte    `json:"page_stage"`
	//
	Session *SessionReq `json:"session"`
}

// For AddAttendanceModel
type AddAttendanceModel struct {
	// Info required for adding an attendance record
	CompanyID           uuid.UUID `json:"company_id"`
	EmployeeID          uuid.UUID `json:"employee_id"`
	DeviceID            uuid.UUID `json:"device_id"`
	RecordTime          time.Time `json:"record_time"`
	VerificationMethod  string    `json:"verification_method"`
	VerificationScore   float64   `json:"verification_score"`
	FaceImageURL        string    `json:"face_image_url"`
	LocationCoordinates string    `json:"location_coordinates"`
	// Session information
	Session        *SessionReq     `json:"session"`
	ServiceSession *ServiceSession `json:"service_session,omitempty"`
}

type ServiceSession struct {
	ServiceId   string `json:"service_id"`
	ServiceName string `json:"service_name"`
	ClientIp    string `json:"client_ip"`
	ClientAgent string `json:"client_agent"`
}

// Session input for adding attendance record
type SessionReq struct {
	SessionId   uuid.UUID `json:"session_id"`
	UserId      uuid.UUID `json:"user_id"`
	CompanyId   uuid.UUID `json:"company_id"`
	Role        int       `json:"role"`
	ClientIp    string    `json:"client_ip"`
	ClientAgent string    `json:"client_agent"`
}
