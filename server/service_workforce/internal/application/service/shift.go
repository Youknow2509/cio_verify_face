package service

import (
	// "context"
	"context"
	"errors"

	applicationError "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/application/error"
	model "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/application/model"
	// "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/application/model"
)

// =================================================
// Shift application interface service
// =================================================
type IShiftService interface {
	DeleteShift(ctx context.Context, input *model.DeleteShiftInput) *applicationError.Error
	EditShift(ctx context.Context, input *model.EditShiftInput) *applicationError.Error
	CreateShift(ctx context.Context, input *model.CreateShiftInput) (*model.CreateShiftOutput, *applicationError.Error)
	GetDetailShift(ctx context.Context, input *model.GetDetailShiftInput) (*model.GetDetailShiftOutput, *applicationError.Error)
}

/**
 * Managet instance
 */
var _vIShiftService IShiftService

/**
 * Getter and setter instance
 */
func GetShiftService() IShiftService {
	return _vIShiftService
}
func SetShiftService(s IShiftService) error {
	if s == nil {
		return errors.New("invalid shift service")
	}
	if _vIShiftService != nil {
		return errors.New("shift service already set")
	}
	_vIShiftService = s
	return nil
}
