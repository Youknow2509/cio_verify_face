package model

import (
	"time"

	"github.com/google/uuid"
)

// =================================================
// Shift Employee application model
// =================================================

// For remove shift employee list
type RemoveShiftEmployeeListInput struct {
	// User info
	UserId      uuid.UUID `json:"user_id"`
	SessionId   uuid.UUID `json:"session_id"`
	Role        int       `json:"role"`
	ClientIp    string    `json:"client_ip"`
	ClientAgent string    `json:"client_agent"`
	CompanyId   uuid.UUID `json:"company_id"`
	//
	ShiftId     uuid.UUID   `json:"shift_id"`
	EmployeeIDs []uuid.UUID `json:"employee_ids"`
}

// For GetListEmployeeInShift
type GetListEmployeeShiftInput struct {
	// User info
	UserId      uuid.UUID `json:"user_id"`
	SessionId   uuid.UUID `json:"session_id"`
	Role        int       `json:"role"`
	ClientIp    string    `json:"client_ip"`
	ClientAgent string    `json:"client_agent"`
	CompanyId   uuid.UUID `json:"company_id"`
	//
	ShiftId uuid.UUID `json:"shift_id"`
	Page    int       `json:"page"`
}

type GetListEmployeeShiftOutput struct {
	Total     int                        `json:"total"`
	Size      int                        `json:"size"`
	Page      int                        `json:"page"`
	Employees []*EmployeeInfoInShiftBase `json:"employees"`
}

// For GetListShiftForEmployee (self-serve listing)
type GetListShiftForEmployeeInput struct {
	EmployeeID  uuid.UUID `json:"employee_id"`
	CompanyID   uuid.UUID `json:"company_id"`
	Page        int       `json:"page"`
	Size        int       `json:"size"`
	ClientIp    string    `json:"client_ip"`
	ClientAgent string    `json:"client_agent"`
}

type ShiftInfoForEmployee struct {
	ShiftId       uuid.UUID `json:"shift_id"`
	ShiftName     string    `json:"shift_name"`
	ShiftStart    string    `json:"shift_start"`
	ShiftEnd      string    `json:"shift_end"`
	EffectiveFrom time.Time `json:"effective_from"`
	EffectiveTo   time.Time `json:"effective_to"`
	IsActive      bool      `json:"is_active"`
}

type GetListShiftForEmployeeOutput struct {
	Total  int                     `json:"total"`
	Page   int                     `json:"page"`
	Size   int                     `json:"size"`
	Shifts []*ShiftInfoForEmployee `json:"shifts"`
}

type EmployeeInfoInShiftBase struct {
	EmployeeId          uuid.UUID `json:"employee_id"`
	EmployeeName        string    `json:"employee_name"`
	EmployeeCode        string    `json:"employee_code"`
	EmployeeShiftName   string    `json:"employee_shift_name"`
	EmployeeShiftActive bool      `json:"employee_shift_active"`
	ShiftEffectiveFrom  time.Time `json:"shift_effective_from,omitempty"`
	ShiftEffectiveTo    time.Time `json:"shift_effective_to,omitempty"`
}

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
	EmployeeIds []uuid.UUID `json:"employee_ids"`
	ShiftId     uuid.UUID   `json:"shift_id"`
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

type ShiftOfUserItem struct {
	ShiftId       uuid.UUID `json:"shift_id"`
	EffectiveFrom time.Time `json:"effective_from"`
	EffectiveTo   time.Time `json:"effective_to"`
	IsActive      bool      `json:"is_active"`
}

type GetShiftForUserWithEffectiveDateOutput struct {
	Total  int                `json:"total"`
	Page   int                `json:"page"`
	Size   int                `json:"size"`
	Shifts []*ShiftOfUserItem `json:"shifts"`
}
