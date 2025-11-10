package service

import (
	// "context"
	"context"
	"errors"
	applicationError "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/application/error"
	model "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/application/model"
)

// =================================================
// ShiftEmployee application interface service
// =================================================
type IShiftEmployeeService interface {
	AddShiftEmployee(ctx context.Context, input *model.AddShiftEmployeeInput) *applicationError.Error
	DeleteShiftUser(ctx context.Context, input *model.DeleteShiftUserInput) *applicationError.Error
	DisableShiftUser(ctx context.Context, input *model.DisableShiftUserInput) *applicationError.Error
	EditShiftForUserWithEffectiveDate(ctx context.Context, input *model.EditShiftForUserWithEffectiveDateInput) *applicationError.Error
	EnableShiftUser(ctx context.Context, input *model.EnableShiftUserInput) *applicationError.Error
	GetShiftForUserWithEffectiveDate(ctx context.Context, input *model.GetShiftForUserWithEffectiveDateInput) (*model.GetShiftForUserWithEffectiveDateOutput, *applicationError.Error)
	AddListShiftEmployee(ctx context.Context, input *model.AddShiftEmployeeListInput) *applicationError.Error
}

/**
 * Managet instance
 */
var _vIShiftEmployeeService IShiftEmployeeService

/**
 * Getter and setter instance
 */
func GetShiftEmployeeService() IShiftEmployeeService {
	return _vIShiftEmployeeService
}
func SetShiftEmployeeService(s IShiftEmployeeService) error {
	if s == nil {
		return errors.New("invalid shiftEmployee service")
	}
	if _vIShiftEmployeeService != nil {
		return errors.New("shiftEmployee service already set")
	}
	_vIShiftEmployeeService = s
	return nil
}
