package health

import (
	"context"
	"errors"

	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/model"
)

/**
 * Interface health check
 */
type IHealthCheck interface {
	CheckSystemResource(ctx context.Context) *model.ComponentCheck
	CheckWebSocketServer(ctx context.Context) *model.ComponentCheck
	CheckDownstreamServices(ctx context.Context) *model.ComponentCheck
}

/**
 * Save instance
 */
var _vIHealthCheck IHealthCheck

/**
 * Getter and setter instance
 */
func GetHealthCheck() IHealthCheck {
	return _vIHealthCheck
}

func SetHealthCheck(data IHealthCheck) error {
	if data == nil {
		return errors.New("data init health check ")
	}
	if _vIHealthCheck != nil {
		return errors.New("data init health check ")
	}
	_vIHealthCheck = data
	return nil
}
