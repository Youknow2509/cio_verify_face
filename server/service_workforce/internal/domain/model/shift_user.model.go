package model

// =========================================================
// Shift User model for repository
// ==================================================

import (
	"time"

	"github.com/google/uuid"
)

// For add list shift assignments to multiple employees
type AddListShiftForEmployeesInput struct {
	ShiftID       uuid.UUID
	CompanyID     uuid.UUID
	EffectiveFrom time.Time
	EffectiveTo   time.Time
	EmployeeIDs   []uuid.UUID
}

// Row returned from GetShiftEmployeeWithEffectiveDate
type EmployeeShiftRow struct {
	EmployeeShiftID uuid.UUID
	ShiftID         uuid.UUID
	EffectiveFrom   time.Time
	EffectiveTo     time.Time
	IsActive        bool
}

// Input for querying shifts of an employee effective on a date
type GetShiftEmployeeWithEffectiveDateInput struct {
	EmployeeID    uuid.UUID
	EffectiveFrom time.Time
	Limit         int32
	Offset        int32
}

// Update effective window for an employee shift
type EditEffectiveShiftForEmployeeInput struct {
	EmployeeShiftID uuid.UUID
	EffectiveFrom   time.Time
	EffectiveTo     time.Time
}

// Add a shift assignment to an employee
type AddShiftForEmployeeInput struct {
	EmployeeID    uuid.UUID
	ShiftID       uuid.UUID
	EffectiveFrom time.Time
	EffectiveTo   time.Time
}

// Check whether a user has an active shift within a range
type CheckUserExistShiftInput struct {
	EmployeeID    uuid.UUID
	EffectiveFrom time.Time
	EffectiveTo   time.Time
	Limit         int32
	Offset        int32
}
