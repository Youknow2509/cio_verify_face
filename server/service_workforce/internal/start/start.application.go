package start

import (
	applicationService "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/application/service"
	applicationServiceImpl "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/application/service/impl"
)

// initialize application services
func initApplication() error {
	// Init IDeviceService
	deviceServiceImpl := applicationServiceImpl.NewDeviceService()
	if err := applicationService.SetDeviceService(deviceServiceImpl); err != nil {
		return err
	}
	return nil
}
