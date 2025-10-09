package start

import (
	appService "github.com/youknow2509/cio_verify_face/server/service_notify/internal/application/service"
	implService "github.com/youknow2509/cio_verify_face/server/service_notify/internal/application/service/impl"
)

// initialize application services
func initApplication() error {
	// Initialize mail service
	mailServiceImpl := implService.NewMailService()
	if err := appService.SetMailService(mailServiceImpl); err != nil {
		return err
	}
	return nil
}
