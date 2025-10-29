package start

import (
	applicationService "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/application/service"
	applicationServiceImpl "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/application/service/impl"
)

// initialize application services
func initApplication() error {
	// Init IShiftService
	shiftServiceImpl := applicationServiceImpl.NewShiftService()
	if err := applicationService.SetShiftService(shiftServiceImpl); err != nil {
		return err
	}
	// Init IScheduleService
	scheduleServiceImpl := applicationServiceImpl.NewScheduleService()
	if err := applicationService.SetScheduleService(scheduleServiceImpl); err != nil {
		return err
	}
	return nil
}
