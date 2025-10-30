package model

import (
	"time"

	"github.com/google/uuid"
)

// =================================================
// Shift application model
// =================================================

// For DeleteShift
type DeleteShiftInput struct {
	// User info
	UserId      uuid.UUID `json:"user_id"`
	UserSession uuid.UUID `json:"user_session"`
	Role        int       `json:"role"`
	//
	ShiftId uuid.UUID `json:"shift_id"`
}

// For EditShift
type EditShiftInput struct {
	// User info
	UserId      uuid.UUID `json:"user_id"`
	UserSession uuid.UUID `json:"user_session"`
	Role        int       `json:"role"`
	//
	ShiftId               uuid.UUID `json:"shift_id" validate:"required"`
	CompanyId             uuid.UUID `json:"company_id" validate:"required"`
	Name                  string    `json:"name" validate:"required"`
	Description           string    `json:"description"`
	StartTime             time.Time `json:"start_time" validate:"required"`
	EndTime               time.Time `json:"end_time" validate:"required"`
	BreakDurationMinutes  int       `json:"break_duration_minutes" validate:"required"`
	GracePeriodMinutes    int       `json:"grace_period_minutes" validate:"required"`
	EarlyDepartureMinutes int       `json:"early_departure_minutes" validate:"required"`
	WorkDays              []int     `json:"work_days" validate:"required"`
}

// For CreateShift
type CreateShiftInput struct {
	// User info
	UserId      uuid.UUID `json:"user_id"`
	UserSession uuid.UUID `json:"user_session"`
	Role        int       `json:"role"`
	//
	CompanyId             uuid.UUID `json:"company_id"`
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
