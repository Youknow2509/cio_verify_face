package model

import (
	"time"

	"github.com/google/uuid"
)

// =================================================
// Shift application model
// =================================================

// For ChangeStatusShift
type ChangeStatusShiftInput struct {
	// User info
	UserId      uuid.UUID `json:"user_id"`
	UserSession uuid.UUID `json:"user_session"`
	Role        int       `json:"role"`
	ClientIp    string    `json:"client_ip"`
	ClientAgent string    `json:"client_agent"`
	CompanyId   uuid.UUID `json:"company_id"`
	//
	CompanyIdReq uuid.UUID `json:"company_id_req"`
	ShiftId      uuid.UUID `json:"shift_id"`
	IsActive     bool      `json:"is_active"`
}

// For GetListShift
type GetListShiftInput struct {
	// User info
	UserId      uuid.UUID `json:"user_id"`
	UserSession uuid.UUID `json:"user_session"`
	Role        int       `json:"role"`
	ClientIp    string    `json:"client_ip"`
	ClientAgent string    `json:"client_agent"`
	CompanyId   uuid.UUID `json:"company_id"`
	//
	Page         int       `json:"page"`
}

// For DeleteShift
type DeleteShiftInput struct {
	// User info
	UserId      uuid.UUID `json:"user_id"`
	UserSession uuid.UUID `json:"user_session"`
	Role        int       `json:"role"`
	ClientIp    string    `json:"client_ip"`
	ClientAgent string    `json:"client_agent"`
	CompanyId   uuid.UUID `json:"company_id"`
	//
	ShiftId uuid.UUID `json:"shift_id"`
}

// For EditShift
type EditShiftInput struct {
	// User info
	UserId      uuid.UUID `json:"user_id"`
	UserSession uuid.UUID `json:"user_session"`
	Role        int       `json:"role"`
	ClientIp    string    `json:"client_ip"`
	ClientAgent string    `json:"client_agent"`
	CompanyId   uuid.UUID `json:"company_id"`
	//
	ShiftId               uuid.UUID `json:"shift_id"`
	CompanyIdReq          uuid.UUID `json:"company_id_req"`
	Name                  string    `json:"name"`
	Description           string    `json:"description"`
	StartTime             time.Time `json:"start_time"`
	EndTime               time.Time `json:"end_time"`
	BreakDurationMinutes  int       `json:"break_duration_minutes"`
	GracePeriodMinutes    int       `json:"grace_period_minutes"`
	EarlyDepartureMinutes int       `json:"early_departure_minutes"`
	WorkDays              []int     `json:"work_days"`
}

// For CreateShift
type CreateShiftInput struct {
	// User info
	UserId      uuid.UUID `json:"user_id"`
	UserSession uuid.UUID `json:"user_session"`
	Role        int       `json:"role"`
	ClientIp    string    `json:"client_ip"`
	ClientAgent string    `json:"client_agent"`
	CompanyId   uuid.UUID `json:"company_id"`
	//
	CompanyIdReq          uuid.UUID `json:"company_id_req"`
	Name                  string    `json:"name"`
	Description           string    `json:"description"`
	StartTime             time.Time `json:"start_time"`
	EndTime               time.Time `json:"end_time"`
	BreakDurationMinutes  int       `json:"break_duration_minutes"`
	GracePeriodMinutes    int       `json:"grace_period_minutes"`
	EarlyDepartureMinutes int       `json:"early_departure_minutes"`
	WorkDays              []int     `json:"work_days"`
}

type CreateShiftOutput struct {
	ShiftId string `json:"shift_id"`
}

// For GetDetailShift
type GetDetailShiftInput struct {
	// User info
	UserId      uuid.UUID `json:"user_id"`
	UserSession uuid.UUID `json:"user_session"`
	Role        int       `json:"role"`
	ClientIp    string    `json:"client_ip"`
	ClientAgent string    `json:"client_agent"`
	CompanyId   uuid.UUID `json:"company_id"`
	//
	ShiftId uuid.UUID `json:"shift_id"`
}

type GetDetailShiftOutput struct {
	ShiftId               string    `json:"shift_id"`
	CompanyId             string    `json:"company_id"`
	Name                  string    `json:"name"`
	Description           string    `json:"description"`
	StartTime             time.Time `json:"start_time"`
	EndTime               time.Time `json:"end_time"`
	BreakDurationMinutes  int       `json:"break_duration_minutes"`
	GracePeriodMinutes    int       `json:"grace_period_minutes"`
	EarlyDepartureMinutes int       `json:"early_departure_minutes"`
	WorkDays              []int     `json:"work_days"`
	IsActive              bool      `json:"is_active"`
}
