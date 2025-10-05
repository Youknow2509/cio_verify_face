package mq

import (
	"context"
	"errors"

	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/model"
)

// =========================================
// Server send event handle to kafka
// =========================================
type ISendEventToKafka interface {
	WriteNewMessage(ctx context.Context, input model.WriteNewMessage) error
	UpgradeStatusTypingUser(ctx context.Context, input model.UpgradeStatusTypingUser) error
	UserReadMessageStatus(ctx context.Context, input model.UserReadMessageStatus) error
	UserEditMessage(ctx context.Context, input model.UserEditMessage) error
	UserDeleteMessage(ctx context.Context, input model.UserDeleteMessage) error
	UserReactMessage(ctx context.Context, input model.UserReactMessage) error
	UserCallOfferInitilize(ctx context.Context, input model.UserCallOfferInitilize) error
	UserCallAnswer(ctx context.Context, input model.UserCallAnswer) error
	UserCallIceCandidate(ctx context.Context, input model.UserCallIceCandidate) error
	UserCallEnd(ctx context.Context, input model.UserCallEnd) error
}

/**
 * Save instance
 */
var _vISendEventToKafka ISendEventToKafka

/**
 * Getter and setter instance
 */
func GetSendEventToKafka() ISendEventToKafka {
	return _vISendEventToKafka
}

func SetSendEventToKafka(data ISendEventToKafka) error {
	if data == nil {
		return errors.New("data init send event ")
	}
	if _vISendEventToKafka != nil {
		return errors.New("data init send event ")
	}
	_vISendEventToKafka = data
	return nil
}
