package service

import (
	"context"
	"errors"

	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/application/model"
)

/**
 * Send event service to Kafka
 */
type ISendEventService interface {
	SendDataVerifyFace(ctx context.Context, input *model.SendDataVerifyFace) error
}

/**
 * Manager instance
 */
var _vISendEventService ISendEventService

/**
 * Getter and setter instance
 */
func GetSendEventService() ISendEventService {
	return _vISendEventService
}

func SetSendEventService(service ISendEventService) error {
	if service == nil {
		return errors.New("service cannot be nil")
	}
	if _vISendEventService != nil {
		return errors.New("service already set")
	}
	_vISendEventService = service
	return nil
}
