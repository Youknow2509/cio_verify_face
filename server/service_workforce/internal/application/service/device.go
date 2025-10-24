package service

import (
	"context"
	"errors"

	applicationError "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/application/error"
	"github.com/youknow2509/cio_verify_face/server/service_workforce/internal/application/model"
)

// =================================================
// Device application interface service
// =================================================
type IDeviceService interface {
	CreateNewDevice(ctx context.Context, input *model.CreateNewDeviceInput) (*model.CreateNewDeviceOutput, *applicationError.Error)
	GetListDevices(ctx context.Context, input *model.ListDevicesInput) (*model.ListDevicesOutput, *applicationError.Error)
	GetDeviceById(ctx context.Context, input *model.GetDeviceByIdInput) (*model.GetDeviceByIdOutput, *applicationError.Error)
	UpdateDeviceById(ctx context.Context, input *model.UpdateDeviceInput) (*model.UpdateDeviceOutput, *applicationError.Error)
	DeleteDeviceById(ctx context.Context, input *model.DeleteDeviceInput) *applicationError.Error
	UpdateLocationDevice(ctx context.Context, input *model.UpdateLocationDeviceInput) *applicationError.Error
	UpdateNameDevice(ctx context.Context, input *model.UpdateNameDeviceInput) *applicationError.Error
	UpdateInfoDevice(ctx context.Context, input *model.UpdateInfoDeviceInput) *applicationError.Error
}

/**
 * Managet instance
 */
var _vIDeviceService IDeviceService

/**
 * Getter and setter instance
 */
func GetDeviceService() IDeviceService {
	return _vIDeviceService
}
func SetDeviceService(s IDeviceService) error {
	if s == nil {
		return errors.New("invalid device service")
	}
	if _vIDeviceService != nil {
		return errors.New("device service already set")
	}
	_vIDeviceService = s
	return nil
}
