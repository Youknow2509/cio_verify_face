package mq

import (
	"context"
	"errors"

	model "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/model"
)

// =========================================
// Server send event handle to kafka
// =========================================
type ISendEventToKafka interface {
	SendDataVerify(ctx context.Context, input model.KafkaAttendanceVerifyReceived) error
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
