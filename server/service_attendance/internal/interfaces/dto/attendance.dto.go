package dto

// ============================================
// Attendance DTOs
// ============================================

// GetAttendanceRecordsRequest represents the request body for retrieving attendance records.
type GetAttendanceRecordsRequest struct {
	Page      int    `json:"page" validate:"omitempty,min=1"`
	PageSize  int    `json:"page_size" validate:"omitempty,min=1,max=100"`
	StartDate string `json:"start_date" validate:"omitempty,datetime=2006-01-02"`
	EndDate   string `json:"end_date" validate:"omitempty,datetime=2006-01-02"`
	DeviceId  string `json:"device_id" validate:"omitempty,uuid4"`
	UserID    string `json:"user_id" validate:"omitempty,uuid4"`
}

// GetMyAttendanceRecordsRequest represents the request body for retrieving a user's attendance records.
type GetMyAttendanceRecordsRequest struct {
	Page      int    `json:"page" validate:"omitempty,min=1"`
	PageSize  int    `json:"page_size" validate:"omitempty,min=1,max=100"`
	StartDate string `json:"start_date" validate:"omitempty,datetime=2006-01-02"`
	EndDate   string `json:"end_date" validate:"omitempty,datetime=2006-01-02"`
}

// CheckInRequest represents the request body for checking in.
type CheckInRequest struct {
	UserID    string `json:"user_id" validate:"required"`
	Timestamp string `json:"timestamp" validate:"required"`
	Location  string `json:"location" validate:"required"`
	DeviceId  string `json:"device_id" validate:"required"`
}

// CheckOutRequest represents the request body for checking out.
type CheckOutRequest struct {
	UserID    string `json:"user_id" validate:"required"`
	Timestamp string `json:"timestamp" validate:"required"`
	Location  string `json:"location" validate:"required"`
	DeviceId  string `json:"device_id" validate:"required"`
}
