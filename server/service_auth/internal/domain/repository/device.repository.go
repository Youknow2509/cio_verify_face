package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	model "github.com/youknow2509/cio_verify_face/server/service_auth/internal/domain/model"
)

/**
 * Device repository interface
 */
type IDeviceRepository interface {
	CheckTokenDevice(ctx context.Context, input *model.CheckTokenDeviceInput) (bool, string, error)
	DeviceExist(ctx context.Context, deviceID uuid.UUID) (bool, error)
	CreateDeviceToken(ctx context.Context, input *model.CreateDeviceTokenInput) error
	BlockDeviceToken(ctx context.Context, deviceToken uuid.UUID) error
}

/**
 * Manager for device repository
 */
var _vIDeviceRepository IDeviceRepository

// SetDeviceRepository sets the device repository implementation
func SetDeviceRepository(repo IDeviceRepository) error {
	if repo == nil {
		return errors.New("invalid device repository")
	}
	if _vIDeviceRepository != nil {
		return errors.New("device repository already set")
	}
	_vIDeviceRepository = repo
	return nil
}

// GetDeviceRepository gets the device repository implementation
func GetDeviceRepository() IDeviceRepository {
	return _vIDeviceRepository
}
