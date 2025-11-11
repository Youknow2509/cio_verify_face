package model

// =========================================================
// Shift User model for repository
// ==================================================

import (
	"time"

	"github.com/google/uuid"
)

// For get list employees in a shift
type GetListEmployyeShiftInput struct {
	ShiftID   uuid.UUID
	CompanyID uuid.UUID
	Limit     int32
	Offset    int32
}

type GetListEmployyeShiftOutput struct {
	EmployeeIDs []*EmployeeShiftInfoBase
	Total       int32
	PageSize    int32
}

type EmployeeShiftInfoBase struct {
	EmployeeId          uuid.UUID
	EmployeeName        string
	EmployeeCode        string
	EmployeeShiftName   string
	EmployeeShiftActive bool
}

// For rm list shift assignments from multiple employees
type RemoveListShiftForEmployeesInput struct {
	ShiftID     uuid.UUID
	EmployeeIDs []uuid.UUID
}

// For IsUserManagetShift
type IsUserManagetShiftInput struct {
	CompanyUserID uuid.UUID
	ShiftID       uuid.UUID
}

// For enable employee shift assignment
type EnableEmployeeShiftIInput struct {
	EmployeeID uuid.UUID
	ShiftID    uuid.UUID
}

// For disable employee shift assignment
type DisableEmployeeShiftInput struct {
	EmployeeID uuid.UUID
	ShiftID    uuid.UUID
}

// For delete employee shift assignment
type DeleteEmployeeShiftInput struct {
	EmployeeID uuid.UUID
	ShiftId    uuid.UUID
}

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
	EmployeeID    uuid.UUID
	ShiftID       uuid.UUID
	EffectiveFrom time.Time
	EffectiveTo   time.Time
	IsActive      bool
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
	EmployeeID    uuid.UUID
	ShiftID       uuid.UUID
	EffectiveFrom time.Time
	EffectiveTo   time.Time
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
