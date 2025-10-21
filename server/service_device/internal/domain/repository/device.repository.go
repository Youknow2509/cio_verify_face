package repository

import (
	"context"
	"errors"

	"github.com/youknow2509/cio_verify_face/server/service_device/internal/domain/model"
)

/**
 * Interface for device repository
 */
type IDeviceRepository interface {
	DeviceExist(ctx context.Context, input *model.DeviceExistInput) (bool, error)
	CreateNewDevice(ctx context.Context, input *model.NewDevice) error
	DeviceInfoBase(ctx context.Context, input *model.DeviceInfoBaseInput) (*model.DeviceInfoBaseOutput, error)
	DeviceInfo(ctx context.Context, input *model.DeviceInfoInput) (*model.DeviceInfoOutput, error)
	ListDeviceInCompany(ctx context.Context, input *model.ListDeviceInCompanyInput) (*model.ListDeviceInCompanyOutput, error)
	DeleteDevice(ctx context.Context, input *model.DeleteDeviceInput) error
	DisableDevice(ctx context.Context, input *model.DisableDeviceInput) error
	EnableDevice(ctx context.Context, input *model.EnableDeviceInput) error
	UpdateDeviceName(ctx context.Context, input *model.UpdateDeviceNameInput) error
	UpdateDeviceLocation(ctx context.Context, input *model.UpdateDeviceLocationInput) error
	UpdateDeviceInfo(ctx context.Context, input *model.UpdateDeviceInfoInput) error
}

/**
 * Variable for device repository instance
 */
var _vDeviceRepository IDeviceRepository

/**
 * Set the Device repository instance
 */
func SetDeviceRepository(v IDeviceRepository) error {
	if _vDeviceRepository != nil {
		return errors.New("device repository initialization failed, not nil")
	}
	_vDeviceRepository = v
	return nil
}

/**
 * Get the device repository instance
 */
func GetDeviceRepository() (IDeviceRepository, error) {
	if _vDeviceRepository == nil {
		return nil, errors.New("device repository not initialized")
	}
	return _vDeviceRepository, nil
}
