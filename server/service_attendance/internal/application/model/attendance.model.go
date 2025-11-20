package model

import (
	"time"

	"github.com/google/uuid"
)

// For GetDailyAttendanceSummary
type GetDailyAttendanceSummaryModel struct {
	Session *SessionReq `json:"session"`
	//
	CompanyID    uuid.UUID `json:"company_id"`
	SummaryMonth string    `json:"summary_month"` // Format: YYYY-MM
	WorkDate     time.Time `json:"work_date"`
	PageSize     int       `json:"page_size" omitempty`
	PageStage    []byte    `json:"page_stage" omitempty`
}

type GetDailyAttendanceSummaryResultModel struct {
	Records       []DailySummariesCompanyInfo `json:"records"`
	PageStageNext string                      `json:"page_stage_next" omitempty`
	PageSize      int                         `json:"page_size" omitempty`
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
	PageStageNext string                       `json:"page_stage_next" omitempty`
	PageSize      int                          `json:"page_size" omitempty`
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
	CompanyID  uuid.UUID `json:"company_id"`
	EmployeeID uuid.UUID `json:"employee_id"`
	YearMonth  string    `json:"year_month" ` // Format: YYYY-MM,
	PageSize   int       `json:"page_size"`
	PageStage  []byte    `json:"page_stage"`
}

type GetAttendanceRecordsCompanyResultModel struct {
	Records       []AttendanceRecordInfo `json:"records"`
	PageStageNext string                 `json:"page_stage_next" omitempty`
	PageSize      int                    `json:"page_size" omitempty`
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
	EmployeeID uuid.UUID `json:"employee_id""`
	PageSize   int       `json:"page_size"`
	PageStage  []byte    `json:"page_stage"`
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
	Session *SessionReq `json:"session"`
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
