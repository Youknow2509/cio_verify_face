package start

import (
	applicationService "github.com/youknow2509/cio_verify_face/server/service_device/internal/application/service"
	applicationServiceImpl "github.com/youknow2509/cio_verify_face/server/service_device/internal/application/service/impl"
)

// initialize application services
func initApplication() error {
	// Init ICoreAuthService
	coreAuthServiceImpl := applicationServiceImpl.NewCoreAuthService()
	if err := applicationService.SetCoreAuthService(coreAuthServiceImpl); err != nil {
		return err
	}
	return nil
}
