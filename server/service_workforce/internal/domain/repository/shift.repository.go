package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	model "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/domain/model"
)

/**
 * Interface for Shift repository
 */
type IShiftRepository interface {
	CreateShift(ctx context.Context, input *model.CreateShiftInput) (uuid.UUID, error)
	ListShifts(ctx context.Context, input *model.ListShiftsInput) ([]*model.Shift, error)
	GetShiftByID(ctx context.Context, shiftID uuid.UUID) (*model.Shift, error)
	UpdateTimeShift(ctx context.Context, input *model.UpdateTimeShiftInput) error
	DeleteShift(ctx context.Context, shiftID uuid.UUID) error
	GetShiftsIdForCompany(ctx context.Context, companyID uuid.UUID, limit, offset int32) ([]uuid.UUID, error)
	DisableShiftWithId(ctx context.Context, input *model.DisableShiftInput) error
	EnableShiftWithId(ctx context.Context, input *model.EnableShiftInput) error
}

/**
 * Variable for Shift repository instance
 */
var _vShiftRepository IShiftRepository

/**
 * Set the Shift repository instance
 */
func SetShiftRepository(v IShiftRepository) error {
	if _vShiftRepository != nil {
		return errors.New("shift repository initialization failed, not nil")
	}
	_vShiftRepository = v
	return nil
}

/**
 * Get the Shift repository instance
 */
func GetShiftRepository() (IShiftRepository, error) {
	if _vShiftRepository == nil {
		return nil, errors.New("shift repository not initialized")
	}
	return _vShiftRepository, nil
}
