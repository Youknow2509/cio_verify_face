package mq

import (
	"context"
	"encoding/json"

	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/constants"
	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/model"
	domainMq "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/mq"
)

/**
 *
 */
type SendEventToKafka struct{}

// SendDataVerify implements mq.ISendEventToKafka.
func (s *SendEventToKafka) SendDataVerify(ctx context.Context, input model.KafkaAttendanceVerifyReceived) error {
	kafkaWrite, err := domainMq.GetKafkaWriteService()
	if err != nil {
		return err
	}
	dataBytes, err := json.Marshal(&input)
	if err != nil {
		return err
	}
	if err := kafkaWrite.WriteMessageRequireAllAck(
		ctx,
		constants.KAFKA_TOPIC_ATTENDANCE_VERIFY,
		"",
		dataBytes,
	); err != nil {
		return err
	}
	
	return nil
}

/**
 * New struct and implementation for sending events to Kafka
 */
func NewSendEventToKafka() domainMq.ISendEventToKafka {
	return &SendEventToKafka{}
}
