package dto

// =======================================
// Shift DTO
// =======================================

// Edit shift
type EditShiftReq struct {
	ShiftId               string `json:"shift_id" validate:"required"`
	CompanyId             string `json:"company_id" validate:"required"`
	Name                  string `json:"name" validate:"required"`
	Description           string `json:"description"`
	StartTime             int64  `json:"start_time" validate:"required"`
	EndTime               int64  `json:"end_time" validate:"required"`
	BreakDurationMinutes  int    `json:"break_duration_minutes" validate:"required"`
	GracePeriodMinutes    int    `json:"grace_period_minutes" validate:"required"`
	EarlyDepartureMinutes int    `json:"early_departure_minutes" validate:"required"`
	WorkDays              []int  `json:"work_days" validate:"required"`
}

// Get detail shift
type GetDetailShift struct {
}

// Create shift
type CreateShiftReq struct {
	CompanyId             string `json:"company_id" validate:"required"`
	Name                  string `json:"name" validate:"required"`
	Description           string `json:"description"`
	StartTime             int64  `json:"start_time" validate:"required"`
	EndTime               int64  `json:"end_time" validate:"required"`
	BreakDurationMinutes  int    `json:"break_duration_minutes" validate:"required"`
	GracePeriodMinutes    int    `json:"grace_period_minutes" validate:"required"`
	EarlyDepartureMinutes int    `json:"early_departure_minutes" validate:"required"`
	WorkDays              []int  `json:"work_days" validate:"required"`
}

// Get shift info base (for user)
type GetShiftInfoBase struct {
}
