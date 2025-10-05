package service

import (
	"context"
	"errors"

	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/application/model"
)

// ===================================
// Health Check Service
// ===================================
type IHealthCheckService interface {
	SystemDetails(ctx context.Context, input *model.SystemDetailsInput) *model.SystemDetails
}

// Save instance interface
var (
	_vIHealthCheckService IHealthCheckService
)

// ================================================
//
//	Getter and setter for instance
//
// ================================================
func GetHealthCheckService() IHealthCheckService {
	return _vIHealthCheckService
}

func SetHealthCheckService(service IHealthCheckService) error {
	if service == nil {
		return errors.New("service is nil")
	}
	if _vIHealthCheckService != nil {
		return errors.New("service already set")
	}
	_vIHealthCheckService = service
	return nil
}
