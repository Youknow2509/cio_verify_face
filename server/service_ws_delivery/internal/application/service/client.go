package service

import (
	"context"
	"errors"

	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/application/model"
)

// =======================================================
// Client service interface
// =======================================================
type IClientService interface {
	SendMessageToClient(ctx context.Context, input *model.SendMessageToClientInput) error
}

// Save instance interface
var (
	_vIClientService IClientService
)

// ================================================
//
//	Getter and setter for instance
//
// ================================================
func GetClientService() IClientService {
	return _vIClientService
}

func SetClientService(service IClientService) error {
	if service == nil {
		return errors.New("service is nil")
	}
	if _vIClientService != nil {
		return errors.New("service already set")
	}
	_vIClientService = service
	return nil
}
