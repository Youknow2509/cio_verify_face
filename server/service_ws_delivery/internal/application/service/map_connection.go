package service

import (
	"context"
	"errors"

	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/application/model"
)

// ================================================
//
//	Map connection service
//
// ================================================
type IMapConnectionService interface {
	// Register connection
	RegisterConnection(ctx context.Context, input *model.RegisterConnection) error
	// Unregister connection
	UnregisterConnection(ctx context.Context, input *model.UnregisterConnection) error
}

// Save instance interface
var (
	_vIMapConnectionService IMapConnectionService
)

// ================================================
//
//	Getter and setter for instance
//
// ================================================
func GetMapConnectionService() IMapConnectionService {
	return _vIMapConnectionService
}

func SetMapConnectionService(service IMapConnectionService) error {
	if service == nil {
		return errors.New("service is nil")
	}
	if _vIMapConnectionService != nil {
		return errors.New("service already set")
	}
	_vIMapConnectionService = service
	return nil
}
