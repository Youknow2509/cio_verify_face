package service

import (
	// "context"
	"context"
	applicationError "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/application/error"
	model "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/application/model"
	service "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/application/service"
)

// =================================================
// Shift service implementation interface
// =================================================
type ShiftService struct {
}

// CreateShift implements service.IShiftService.
func (s *ShiftService) CreateShift(ctx context.Context, input *model.CreateShiftInput) (*model.CreateShiftOutput, *applicationError.Error) {
	panic("unimplemented")
}

// DeleteShift implements service.IShiftService.
func (s *ShiftService) DeleteShift(ctx context.Context, input *model.DeleteShiftInput) *applicationError.Error {
	panic("unimplemented")
}

// EditShift implements service.IShiftService.
func (s *ShiftService) EditShift(ctx context.Context, input *model.EditShiftInput) *applicationError.Error {
	panic("unimplemented")
}

// GetDetailShift implements service.IShiftService.
func (s *ShiftService) GetDetailShift(ctx context.Context, input *model.GetDetailShiftInput) (*model.GetDetailShiftOutput, *applicationError.Error) {
	panic("unimplemented")
}

// New instance
func NewShiftService() service.IShiftService {
	return &ShiftService{}
}
