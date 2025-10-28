package dto

// =======================================
// Shift employee DTO
// =======================================

// Enable shift for user request
type EnableShiftForUserReq struct {
	ShiftUserId string `json:"shift_user_id" validate:"required"`
}

// Disable shift for user request
type DisableShiftForUserReq struct {
	ShiftUserId string `json:"shift_user_id" validate:"required"`
}

// Add shift employee request
type AddShiftEmployeeReq struct {
	EmployeeId    string `json:"employee_id" validate:"required"`
	ShiftId       string `json:"shift_id" validate:"required"`
	EffectiveFrom int64  `json:"effective_from" validate:"required"`
	EffectiveTo   int64  `json:"effective_to" validate:"required"`
}

// Edit shift employee effective date request
type ShiftEmployeeEditEffectiveDateReq struct {
	ShiftUserId      string `json:"shift_user_id" validate:"required"`
	NewEffectiveFrom int64  `json:"new_effective_from" validate:"required"`
	NewEffectiveTo   int64  `json:"new_effective_to" validate:"required"`
}

// Get shift for user with effective date
type ShiftEmployeeEffectiveDateReq struct {
	UserId        string `json:"user_id"`
	EffectiveFrom int64  `json:"effective_from" validate:"required"`
	EffectiveTo   int64  `json:"effective_to" validate:"required"`
	Page          int    `json:"page" validate:"required"`
	Size          int    `json:"size" validate:"required"`
}
