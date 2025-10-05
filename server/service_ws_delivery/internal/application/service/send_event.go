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
	NewMessageText(ctx context.Context, input model.NewMessageText) error
	UserTypingStatus(ctx context.Context, input model.UserTypingStatus) error
	MessageReadStatus(ctx context.Context, input model.MessageReadStatus) error
	EditMessage(ctx context.Context, input model.EditMessage) error
	DeleteMessage(ctx context.Context, input model.DeleteMessage) error
	ReactMessage(ctx context.Context, input model.ReactMessage) error
	CallOffer(ctx context.Context, input model.CallOffer) error
	CallAnswer(ctx context.Context, input model.CallAnswer) error
	CallIceCandidate(ctx context.Context, input model.CallIceCandidate) error
	CallEnd(ctx context.Context, input model.CallEnd) error
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
