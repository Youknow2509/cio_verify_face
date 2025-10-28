package service

import (
	// "context"
	"context"
	applicationError "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/application/error"
	model "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/application/model"
	service "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/application/service"
)

// =================================================
// ShiftEmployee service implementation interface
// =================================================
type ShiftEmployeeService struct {
}

// AddShiftEmployee implements service.IShiftEmployeeService.
func (s *ShiftEmployeeService) AddShiftEmployee(ctx context.Context, input *model.AddShiftEmployeeInput) (**model.AddShiftEmployeeOutput, *applicationError.Error) {
	panic("unimplemented")
}

// DeleteShiftUser implements service.IShiftEmployeeService.
func (s *ShiftEmployeeService) DeleteShiftUser(ctx context.Context, input *model.DeleteShiftUserInput) *applicationError.Error {
	panic("unimplemented")
}

// DisableShiftUser implements service.IShiftEmployeeService.
func (s *ShiftEmployeeService) DisableShiftUser(ctx context.Context, input *model.DisableShiftUserInput) *applicationError.Error {
	panic("unimplemented")
}

// EditShiftForUserWithEffectiveDate implements service.IShiftEmployeeService.
func (s *ShiftEmployeeService) EditShiftForUserWithEffectiveDate(ctx context.Context, input *model.EditShiftForUserWithEffectiveDateInput) *applicationError.Error {
	panic("unimplemented")
}

// EnableShiftUser implements service.IShiftEmployeeService.
func (s *ShiftEmployeeService) EnableShiftUser(ctx context.Context, input *model.EnableShiftUserInput) *applicationError.Error {
	panic("unimplemented")
}

// GetShiftForUserWithEffectiveDate implements service.IShiftEmployeeService.
func (s *ShiftEmployeeService) GetShiftForUserWithEffectiveDate(ctx context.Context, input *model.GetShiftForUserWithEffectiveDateInput) (*model.GetShiftForUserWithEffectiveDateOutput, *applicationError.Error) {
	panic("unimplemented")
}

// New instance
func NewShiftEmployeeService() service.IShiftEmployeeService {
	return &ShiftEmployeeService{}
}
