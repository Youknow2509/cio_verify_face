package service

import (
	// "context"
	"errors"
	// applicationError "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/application/error"
	// "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/application/model"
)

// =================================================
// Schedule application interface service
// =================================================
type IScheduleService interface {
}

/**
 * Managet instance
 */
var _vIScheduleService IScheduleService

/**
 * Getter and setter instance
 */
func GetScheduleService() IScheduleService {
	return _vIScheduleService
}
func SetScheduleService(s IScheduleService) error {
	if s == nil {
		return errors.New("invalid schedule service")
	}
	if _vIScheduleService != nil {
		return errors.New("schedule service already set")
	}
	_vIScheduleService = s
	return nil
}
