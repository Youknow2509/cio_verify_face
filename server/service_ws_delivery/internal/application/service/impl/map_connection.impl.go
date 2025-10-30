package impl

import (
	"context"

	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/application/model"
	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/application/service"
	domainModel "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/model"
	domainRepository "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/repository"
	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/global"
	utilsCache "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/utils/cache"
)

// ================================================
// Service Map connection implementation
// ================================================
type MapConnectionService struct {
}

// RegisterConnection implements service.IMapConnectionService.
func (m *MapConnectionService) RegisterConnection(ctx context.Context, input *model.RegisterConnection) error {
	domainMapConnectionRepo := domainRepository.GetManagerConnectionRepository()
	if _, err := domainMapConnectionRepo.CreateConnection(
		ctx,
		&domainModel.CreateConnectionInput{
			DeviceConnectionsKey:  utilsCache.GetDeviceConnectionWsKey(input.DeviceId),
			ServiceConnectionsKey: utilsCache.GetServiceWsConnectionKey(global.ServerSetting.Id),
			ConnectionId:          input.ConnectionId,
			DeviceId:              input.DeviceId,
			ServiceId:             global.ServerSetting.Id,
			IpAddress:             input.IPAddress,
			ConnectedAt:           input.ConnectedAt,
			UserAgent:             input.UserAgent,
		},
	); err != nil {
		global.Logger.Error("MapConnectionService.RegisterConnection", "error", err)
		return err
	}
	return nil
}

// UnregisterConnection implements service.IMapConnectionService.
func (m *MapConnectionService) UnregisterConnection(ctx context.Context, input *model.UnregisterConnection) error {
	domainMapConnectionRepo := domainRepository.GetManagerConnectionRepository()
	if _, err := domainMapConnectionRepo.RemoveConnection(
		ctx,
		&domainModel.RemoveConnectionInput{
			DeviceId:              input.DeviceId,
			DeviceConnectionsKey:  utilsCache.GetDeviceConnectionWsKey(input.DeviceId),
			ServiceConnectionsKey: utilsCache.GetServiceWsConnectionKey(global.ServerSetting.Id),
		}); err != nil {
		global.Logger.Error("MapConnectionService.UnregisterConnection", "error", err)
		return err
	}
	return nil
}

// New service and implement
func NewMapConnectionService() service.IMapConnectionService {
	return &MapConnectionService{}
}
