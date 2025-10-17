package service

import (
	"context"
	applicationError "github.com/youknow2509/cio_verify_face/server/service_device/internal/application/error"
	"github.com/youknow2509/cio_verify_face/server/service_device/internal/application/model"
	"github.com/youknow2509/cio_verify_face/server/service_device/internal/application/service"
)

// =================================================
// Device application service
// =================================================
type DeviceService struct{}

// UpdateInfoDevice implements service.IDeviceService.
func (d *DeviceService) UpdateInfoDevice(ctx context.Context, input *model.UpdateInfoDeviceInput) *applicationError.Error {
	panic("unimplemented")
}

// UpdateLocationDevice implements service.IDeviceService.
func (d *DeviceService) UpdateLocationDevice(ctx context.Context, input *model.UpdateLocationDeviceInput) *applicationError.Error {
	panic("unimplemented")
}

// UpdateNameDevice implements service.IDeviceService.
func (d *DeviceService) UpdateNameDevice(ctx context.Context, input *model.UpdateNameDeviceInput) *applicationError.Error {
	panic("unimplemented")
}

// CreateNewDevice implements service.IDeviceService.
func (d *DeviceService) CreateNewDevice(ctx context.Context, input *model.CreateNewDeviceInput) (*model.CreateNewDeviceOutput, *applicationError.Error) {
	panic("unimplemented")
}

// DeleteDeviceById implements service.IDeviceService.
func (d *DeviceService) DeleteDeviceById(ctx context.Context, input *model.DeleteDeviceInput) *applicationError.Error {
	panic("unimplemented")
}

// GetDeviceById implements service.IDeviceService.
func (d *DeviceService) GetDeviceById(ctx context.Context, input *model.GetDeviceByIdInput) (*model.GetDeviceByIdOutput, *applicationError.Error) {
	panic("unimplemented")
}

// GetListDevices implements service.IDeviceService.
func (d *DeviceService) GetListDevices(ctx context.Context, input *model.ListDevicesInput) (*model.ListDevicesOutput, *applicationError.Error) {
	panic("unimplemented")
}

// UpdateDeviceById implements service.IDeviceService.
func (d *DeviceService) UpdateDeviceById(ctx context.Context, input *model.UpdateDeviceInput) (*model.UpdateDeviceOutput, *applicationError.Error) {
	panic("unimplemented")
}

// NewDeviceService create new instance and implement IDeviceService
func NewDeviceService() service.IDeviceService {
	return &DeviceService{}
}
