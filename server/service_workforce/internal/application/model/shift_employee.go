package model

import (
	"time"

	"github.com/google/uuid"
)

// =================================================
// Shift Employee application model
// =================================================

// For GetInfoEmployeeInShift
type GetInfoEmployeeInShiftInput struct {
	// User info
	UserId      uuid.UUID `json:"user_id"`
	SessionId   uuid.UUID `json:"session_id"`
	Role        int       `json:"role"`
	ClientIp    string    `json:"client_ip"`
	ClientAgent string    `json:"client_agent"`
	CompanyId   uuid.UUID `json:"company_id"`
	//
	ShiftId uuid.UUID `json:"shift_id"`
}

type EmployeeInfoInShiftBase struct {
	UserId         string `json:"user_id"`
	NumberEmployee string `json:"number_employee"`
	Name           string `json:"user_name"`
	CurrentShift   bool   `json:"current_shift"`
	ShiftActive    string `json:"shift_active"`
}

// For add shift employee list
type AddShiftEmployeeListInput struct {
	// User info
	UserId      uuid.UUID `json:"user_id"`
	SessionId   uuid.UUID `json:"session_id"`
	Role        int       `json:"role"`
	ClientIp    string    `json:"client_ip"`
	ClientAgent string    `json:"client_agent"`
	CompanyId   uuid.UUID `json:"company_id"`
	//
	CompanyRequestId uuid.UUID   `json:"company_request_id"`
	ShiftId          uuid.UUID   `json:"shift_id"`
	EffectiveFrom    time.Time   `json:"effective_from"`
	EffectiveTo      time.Time   `json:"effective_to"`
	EmployeeIDs      []uuid.UUID `json:"employee_ids"`
}

// For AddShiftEmployee
type AddShiftEmployeeInput struct {
	// User info
	UserId      uuid.UUID `json:"user_id"`
	SessionId   uuid.UUID `json:"session_id"`
	Role        int       `json:"role"`
	ClientIp    string    `json:"client_ip"`
	ClientAgent string    `json:"client_agent"`
	CompanyId   uuid.UUID `json:"company_id"`
	//
	EmployeeId    uuid.UUID `json:"employee_id"`
	ShiftId       uuid.UUID `json:"shift_id"`
	EffectiveFrom time.Time `json:"effective_from"`
	EffectiveTo   time.Time `json:"effective_to"`
}
type AddShiftEmployeeOutput struct {
	ShiftUserId string `json:"shift_user_id"`
}

// For DeleteShiftUser
type DeleteShiftUserInput struct {
	// User info
	UserId      uuid.UUID `json:"user_id"`
	SessionId   uuid.UUID `json:"session_id"`
	Role        int       `json:"role"`
	ClientIp    string    `json:"client_ip"`
	ClientAgent string    `json:"client_agent"`
	CompanyId   uuid.UUID `json:"company_id"`
	//
	UserIdReq uuid.UUID `json:"user_id_req"`
	ShiftId   uuid.UUID `json:"shift_id"`
}

// For DisableShiftUser
type DisableShiftUserInput struct {
	// User info
	UserId      uuid.UUID `json:"user_id"`
	SessionId   uuid.UUID `json:"session_id"`
	Role        int       `json:"role"`
	ClientIp    string    `json:"client_ip"`
	ClientAgent string    `json:"client_agent"`
	CompanyId   uuid.UUID `json:"company_id"`
	//
	UserIdReq uuid.UUID `json:"user_id_req"`
	ShiftId   uuid.UUID `json:"shift_id"`
}

// For EditShiftForUserWithEffectiveDate
type EditShiftForUserWithEffectiveDateInput struct {
	// User info
	UserId      uuid.UUID `json:"user_id"`
	SessionId   uuid.UUID `json:"session_id"`
	Role        int       `json:"role"`
	ClientIp    string    `json:"client_ip"`
	ClientAgent string    `json:"client_agent"`
	CompanyId   uuid.UUID `json:"company_id"`
	//
	UserIdReq        uuid.UUID `json:"user_id_req"`
	ShiftId          uuid.UUID `json:"shift_id"`
	NewEffectiveFrom time.Time `json:"new_effective_from"`
	NewEffectiveTo   time.Time `json:"new_effective_to"`
}

// For EnableShiftUser
type EnableShiftUserInput struct {
	// User info
	UserId      uuid.UUID `json:"user_id"`
	SessionId   uuid.UUID `json:"session_id"`
	Role        int       `json:"role"`
	ClientIp    string    `json:"client_ip"`
	ClientAgent string    `json:"client_agent"`
	CompanyId   uuid.UUID `json:"company_id"`
	//
	UserIdReq uuid.UUID `json:"user_id_req"`
	ShiftId   uuid.UUID `json:"shift_id"`
}

// For GetShiftForUserWithEffectiveDate
type GetShiftForUserWithEffectiveDateInput struct {
	// User info
	UserId      uuid.UUID `json:"user_id"`
	SessionId   uuid.UUID `json:"session_id"`
	Role        int       `json:"role"`
	ClientIp    string    `json:"client_ip"`
	ClientAgent string    `json:"client_agent"`
	CompanyId   uuid.UUID `json:"company_id"`
	//
	EffectiveFrom time.Time `json:"effective_from"`
	EffectiveTo   time.Time `json:"effective_to"`
	Page          int       `json:"page"`
	Size          int       `json:"size"`
}
type GetShiftForUserWithEffectiveDateOutput struct{}
