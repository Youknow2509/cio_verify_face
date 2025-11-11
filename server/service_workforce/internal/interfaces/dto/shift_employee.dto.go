package dto

// =======================================
// Shift employee DTO
// =======================================

// GetInfoEmployeeInShiftReq request
type GetInfoEmployeeInShiftReq struct {
	ShiftId string `json:"shift_id" validate:"required"`
	Page    int    `json:"page" validate:"gte=1"`
}

// GetInfoEmployeeDonotInShiftReq request
type GetInfoEmployeeDonotInShiftReq struct {
	ShiftId string `json:"shift_id" validate:"required"`
	Page    int    `json:"page" validate:"gte=1"`
}

// DisableShiftUserReq request
type DisableShiftUserReq struct {
	ShiftId    string `json:"shift_id" validate:"required"`
	EmployeeId string `json:"employee_id" validate:"required"`
}

// DeleteShiftUserReq request
type DeleteShiftUserReq struct {
	ShiftId    string `json:"shift_id" validate:"required"`
	EmployeeId string `json:"employee_id" validate:"required"`
}

// Add shift employee list request
type AddShiftEmployeeListReq struct {
	EmployeeIDs   []string `json:"employee_ids" validate:"required,dive,required"`
	EffectiveFrom int64    `json:"effective_from" validate:"required"`
	EffectiveTo   int64    `json:"effective_to" validate:"required"`
	CompanyId     string   `json:"company_id" validate:"required"`
	ShiftId       string   `json:"shift_id" validate:"required"`
}

// List Employee id
type ListEmployeeIDReq struct {
	EmployeeIDs []string `json:"employee_ids" validate:"required,dive,required"`
}

// Enable shift for user request
type EnableShiftUserReq struct {
	ShiftId    string `json:"shift_id" validate:"required"`
	EmployeeId string `json:"employee_id" validate:"required"`
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
	ShiftId          string `json:"shift_id" validate:"required"`
	EmployeeID       string `json:"employee_id" validate:"required"`
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
