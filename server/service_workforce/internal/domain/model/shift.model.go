package model

// =========================================================
// Shift model for repository
// ==================================================

import (
	"time"

	"github.com/google/uuid"
)

// Shift entity (domain-facing)
type Shift struct {
	ShiftID               uuid.UUID
	CompanyID             uuid.UUID
	Name                  string
	Description           string
	StartTime             time.Time
	EndTime               time.Time
	BreakDurationMinutes  int32
	GracePeriodMinutes    int32
	EarlyDepartureMinutes int32
	WorkDays              []int32
	IsFlexible            bool
	OvertimeAfterMinutes  int32
	IsActive              bool
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

// CreateShiftInput carries data to create a shift
type CreateShiftInput struct {
	CompanyID             uuid.UUID
	Name                  string
	Description           string
	StartTime             time.Time
	EndTime               time.Time
	BreakDurationMinutes  int32
	GracePeriodMinutes    int32
	EarlyDepartureMinutes int32
	WorkDays              []int32
}

// ListShiftsInput filters and paginates shifts
type ListShiftsInput struct {
	CompanyID uuid.UUID
	IsActive  bool
	Limit     int32
	Offset    int32
}

// UpdateTimeShiftInput updates time-related fields of a shift
type UpdateTimeShiftInput struct {
	ShiftID               uuid.UUID
	StartTime             time.Time
	EndTime               time.Time
	BreakDurationMinutes  int32
	GracePeriodMinutes    int32
	EarlyDepartureMinutes int32
	WorkDays              []int32
}
