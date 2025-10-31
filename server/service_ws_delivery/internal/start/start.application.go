package start

import (
	applicationService "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/application/service"
	applicationServiceImpl "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/application/service/impl"
)

// initialize application services
func initApplication() error {
	// Initialize map connection service
	_mapConnectionService := applicationServiceImpl.NewMapConnectionService()
	if err := applicationService.SetMapConnectionService(_mapConnectionService); err != nil {
		return err
	}
	_sendEventService := applicationServiceImpl.NewSendEventService()
	if err := applicationService.SetSendEventService(_sendEventService); err != nil {
		return err
	}
	_clientService := applicationServiceImpl.NewClientService()
	if err := applicationService.SetClientService(_clientService); err != nil {
		return err
	}
	_healthCheckService := applicationServiceImpl.NewHealthCheckService()
	if err := applicationService.SetHealthCheckService(_healthCheckService); err != nil {
		return err
	}
	// v.v
	return nil
}
