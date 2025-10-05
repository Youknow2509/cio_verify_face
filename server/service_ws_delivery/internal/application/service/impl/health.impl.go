package impl

import (
	"context"

	model "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/application/model"
	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/application/service"
	domainHealth "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/health"
	domainModel "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/model"
)

/**
 * HealthCheck service struct
 */
type HealthCheckService struct{}

// SystemDetails implements service.IHealthCheckService.
func (h *HealthCheckService) SystemDetails(ctx context.Context, input *model.SystemDetailsInput) *model.SystemDetails {
	// Get instance use
	service := domainHealth.GetHealthCheck()
	if service == nil {
		return nil
	}
	// Get data
	checkResource := service.CheckSystemResource(ctx)
	checkWs := service.CheckWebSocketServer(ctx)
	checkDownStream := service.CheckDownstreamServices(ctx)
	return &model.SystemDetails{
		SystemResources:    convertComponentCheck(checkResource),
		WebSocketServer:    convertComponentCheck(checkWs),
		DownstreamServices: convertComponentCheck(checkDownStream),
	}
}

/**
 * New send event service and implementation
 */
func NewHealthCheckService() service.IHealthCheckService {
	return &HealthCheckService{}
}

// convertComponentCheck converts a domain model ComponentCheck to an application model ComponentCheck.
func convertComponentCheck(domainCheck *domainModel.ComponentCheck) *model.ComponentCheck {
	if domainCheck == nil {
		return nil
	}
	return &model.ComponentCheck{
		Status:  string(domainCheck.Status),
		Details: domainCheck.Details,
	}
}
