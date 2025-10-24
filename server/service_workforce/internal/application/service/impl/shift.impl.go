package service

import (
	// "context"
	// applicationError "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/application/error"
	service "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/application/service"
)

// =================================================
// Shift service implementation interface
// =================================================
type ShiftService struct {
}

// New instance
func NewShiftService() service.IShiftService {
	return &ShiftService{}
}
