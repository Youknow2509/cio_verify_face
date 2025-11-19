package model

import (
	"time"

	"github.com/google/uuid"
)

// ============================================
// Attendance Record Model
// ============================================
type DailySummariesEmployeeOutput struct {
	Records       []DailySummariesEmployeeInfo `json:"records"`
	PageStageNext []byte                       `json:"page_stage_next" omitempty`
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

type DailySummariesCompanyOutput struct {
	Records       []DailySummariesCompanyInfo `json:"records"`
	PageStageNext []byte                      `json:"page_stage_next" omitempty`
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

type GetDailySummariesCompanyForEmployeeInput struct {
	// company_id = uuid_company AND employee_id = uuid_employee AND summary_month = '2023-10';
	CompanyID    uuid.UUID `json:"company_id"`
	EmployeeID   uuid.UUID `json:"employee_id"`
	SummaryMonth string    `json:"summary_month"`
	PageSize     int       `json:"page_size" omitempty`
	PageStage    []byte    `json:"page_stage" omitempty`
}

type GetDailySummariesCompanyInput struct {
	// company_id = uuid_company AND summary_month = '2023-10' AND work_date = '2023-10-25';
	CompanyID    uuid.UUID `json:"company_id"`
	SummaryMonth string    `json:"summary_month"`
	WorkDate     time.Time `json:"work_date"`
	PageSize     int       `json:"page_size" omitempty`
	PageStage    []byte    `json:"page_stage" omitempty`
}

type UpdateDailySummariesEmployeeInput struct {
	CompanyID        uuid.UUID `json:"company_id"`
	SummaryMonth     string    `json:"summary_month"`
	WorkDate         time.Time `json:"work_date"`
	EmployeeID       uuid.UUID `json:"employee_id"`
	Notes            string    `json:"notes"`
	AttendanceStatus int       `json:"attendance_status"`
}

type DeleteDailySummariesInput struct {
	CompanyID    uuid.UUID `json:"company_id"`
	SummaryMonth string    `json:"summary_month"`
	WorkDate     time.Time `json:"work_date"`
	EmployeeID   uuid.UUID `json:"employee_id"`
}

type DeleteDailySummariesEmployeeInput struct {
	CompanyID    uuid.UUID `json:"company_id"`
	SummaryMonth string    `json:"summary_month"`
	EmployeeID   uuid.UUID `json:"employee_id"`
	WorkDate     time.Time `json:"work_date"`
}

type DeleteAttendanceRecordInput struct {
	CompanyID  uuid.UUID `json:"company_id"`
	YearMonth  string    `json:"year_month"`
	RecordTime time.Time `json:"record_time"`
	EmployeeID uuid.UUID `json:"employee_id"`
}

type GetAttendanceRecordCompanyForEmployeeInput struct {
	CompanyID  uuid.UUID `json:"company_id"`
	YearMonth  string    `json:"year_month"`
	EmployeeID uuid.UUID `json:"employee_id"`
	PageSize   int       `json:"page_size" omitempty`
	PageStage  []byte    `json:"page_stage" omitempty`
}

type AttendanceRecordOutput struct {
	Records       []AttendanceRecordInfo `json:"records"`
	PageStageNext []byte                 `json:"page_stage_next" omitempty`
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

type GetAttendanceRecordCompanyInput struct {
	CompanyID uuid.UUID `json:"company_id"`
	YearMonth string    `json:"year_month"`
	PageSize  int       `json:"page_size" omitempty`
	PageStage []byte    `json:"page_stage" omitempty`
}

type AddDailySummariesInput struct {
	CompanyID         uuid.UUID `json:"company_id"`
	SummaryMonth      string    `json:"summary_month"`
	WorkDate          time.Time `json:"work_date"`
	EmployeeID        uuid.UUID `json:"employee_id"`
	ShiftID           uuid.UUID `json:"shift_id"`
	ActualCheckIn     time.Time `json:"actual_check_in"`
	ActualCheckOut    time.Time `json:"actual_check_out"`
	AttendanceStatus  int       `json:"attendance_status"`
	LateMinutes       int       `json:"late_minutes"`
	EarlyLeaveMinutes int       `json:"early_leave_minutes"`
	TotalWorkMinutes  int       `json:"total_work_minutes"`
	Notes             string    `json:"notes"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type AddAttendanceRecordInput struct {
	CompanyID           uuid.UUID         `json:"company_id"`
	EmployeeID          uuid.UUID         `json:"employee_id"`
	YearMonth           string            `json:"year_month"`
	RecordTime          time.Time         `json:"record_time"`
	DeviceID            uuid.UUID         `json:"device_id"`
	RecordType          int               `json:"record_type"`
	VerificationMethod  string            `json:"verification_method"`
	VerificationScore   float64           `json:"verification_score"`
	FaceImageURL        string            `json:"face_image_url"`
	LocationCoordinates string            `json:"location_coordinates"`
	Metadata            map[string]string `json:"metadata"`
}
