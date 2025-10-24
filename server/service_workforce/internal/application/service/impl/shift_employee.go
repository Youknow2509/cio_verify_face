package service

import (
	// "context"
	// applicationError "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/application/error"
	service "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/application/service"
)

// =================================================
// ShiftEmployee service implementation interface
// =================================================
type ShiftEmployeeService struct {
}

// New instance
func NewShiftEmployeeService() service.IShiftEmployeeService {
	return &ShiftEmployeeService{}
}
