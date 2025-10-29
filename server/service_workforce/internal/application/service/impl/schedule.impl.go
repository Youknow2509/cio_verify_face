package service

import (
	// "context"
	// applicationError "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/application/error"
	service "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/application/service"
)

// =================================================
// Schedule service implementation interface
// =================================================
type ScheduleService struct {
}

// New instance
func NewScheduleService() service.IScheduleService {
	return &ScheduleService{}
}
