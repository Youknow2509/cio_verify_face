package repository

import (
	"context"
	"errors"

	model "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/domain/model"
)

/**
 * Interface for ShiftUser repository
 */
type IShiftUserRepository interface {
	GetShiftEmployeeWithEffectiveDate(ctx context.Context, input *model.GetShiftEmployeeWithEffectiveDateInput) ([]*model.EmployeeShiftRow, error)
	EditEffectiveShiftForEmployee(ctx context.Context, input *model.EditEffectiveShiftForEmployeeInput) error
	DeleteEmployeeShift(ctx context.Context, input *model.DeleteEmployeeShiftInput) error
	DisableEmployeeShift(ctx context.Context, input *model.DisableEmployeeShiftInput) error
	EnableEmployeeShift(ctx context.Context, input *model.EnableEmployeeShiftIInput) error
	AddShiftForEmployee(ctx context.Context, input *model.AddShiftForEmployeeInput) error
	CheckUserExistShift(ctx context.Context, input *model.CheckUserExistShiftInput) (bool, error)
	AddListShiftForEmployees(ctx context.Context, input *model.AddListShiftForEmployeesInput) error
}

/**
 * Variable for ShiftUser repository instance
 */
var _vShiftUserRepository IShiftUserRepository

/**
 * Set the ShiftUser repository instance
 */
func SetShiftUserRepository(v IShiftUserRepository) error {
	if _vShiftUserRepository != nil {
		return errors.New("shift user repository initialization failed, not nil")
	}
	_vShiftUserRepository = v
	return nil
}

/**
 * Get the ShiftUser repository instance
 */
func GetShiftUserRepository() (IShiftUserRepository, error) {
	if _vShiftUserRepository == nil {
		return nil, errors.New("shift user repository not initialized")
	}
	return _vShiftUserRepository, nil
}
