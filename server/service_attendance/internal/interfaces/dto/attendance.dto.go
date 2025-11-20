package dto

// ============================================
// Attendance DTOs
// ============================================
type GetDailyAttendanceSummaryRequest struct {
	CompanyID    string `json:"company_id" validate:"required"`
	SummaryMonth string `json:"summary_month" validate:"required ,len=7"` // Format: YYYY-MM
	WorkDate     int64  `json:"work_date" validate:"omitempty"`           // Unix timestamp in milliseconds
	PageSize     int    `json:"page_size" omitempty`
	PageStage    string `json:"page_stage" omitempty`
}

type GetDailyAttendanceSummaryEmployeeRequest struct {
	CompanyID    string `json:"company_id" validate:"required"`
	EmployeeID   string `json:"employee_id" validate:"omitempty"`
	SummaryMonth string `json:"summary_month" validate:"required ,len=7"` // Format: YYYY-MM
	PageSize     int    `json:"page_size" validate:"omitempty"`
	PageStage    string `json:"page_stage" validate:"omitempty"`
}

type GetAttendanceRecordsEmployeeRequest struct {
	CompanyID  string `json:"company_id" validate:"required"`
	YearMonth  string `json:"year_month" validate:"required,len=7"` // Format: YYYY-MM
	EmployeeID string `json:"employee_id" validate:"omitempty"`
	PageSize   int    `json:"page_size" validate:"omitempty"`
	PageStage  string `json:"page_stage" validate:"omitempty"`
}

type GetAttendanceRecordsRequest struct {
	CompanyID string `json:"company_id" validate:"required"`
	YearMonth string `json:"year_month" validate:"required,len=7"` // Format: YYYY-MM
	PageSize  int    `json:"page_size" validate:"omitempty"`
	PageStage string `json:"page_stage" validate:"omitempty"`
}

type AddAttendanceRequest struct {
	CompanyID           string  `json:"company_id" validate:"required"`
	EmployeeID          string  `json:"employee_id" validate:"required"`
	DeviceID            string  `json:"device_id" validate:"required"`
	RecordTime          int64   `json:"record_time" validate:"required,gt=0"`
	VerificationMethod  string  `json:"verification_method" validate:"required"`
	VerificationScore   float64 `json:"verification_score" validate:"required,gte=0,lte=1"`
	FaceImageURL        string  `json:"face_image_url" validate:"required"`
	LocationCoordinates string  `json:"location_coordinates" validate:"required"`
}
