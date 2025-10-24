package service

import (
	// "context"
	"errors"
	// applicationError "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/application/error"
	// "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/application/model"
)

// =================================================
// ShiftEmployee application interface service
// =================================================
type IShiftEmployeeService interface {
}

/**
 * Managet instance
 */
var _vIShiftEmployeeService IShiftEmployeeService

/**
 * Getter and setter instance
 */
func GetShiftEmployeeService() IShiftEmployeeService {
	return _vIShiftEmployeeService
}
func SetShiftEmployeeService(s IShiftEmployeeService) error {
	if s == nil {
		return errors.New("invalid shiftEmployee service")
	}
	if _vIShiftEmployeeService != nil {
		return errors.New("shiftEmployee service already set")
	}
	_vIShiftEmployeeService = s
	return nil
}
