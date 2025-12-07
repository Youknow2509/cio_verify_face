package mq

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	applicationModel "github.com/youknow2509/cio_verify_face/server/service_notify/internal/application/model"
	applicationService "github.com/youknow2509/cio_verify_face/server/service_notify/internal/application/service"
	"github.com/youknow2509/cio_verify_face/server/service_notify/internal/constants"
	domainMq "github.com/youknow2509/cio_verify_face/server/service_notify/internal/domain/mq"
	"github.com/youknow2509/cio_verify_face/server/service_notify/internal/global"
	"github.com/youknow2509/cio_verify_face/server/service_notify/internal/interfaces/dto"
	"go.opentelemetry.io/otel/attribute"
)

/**
 * Kafka listener data
 */
type KafkaListenerData struct {
	Topic     string
	NumThread int
}

/**
 * New Kafka listener data
 */
func NewKafkaListenerData(topic string, numThread int) *KafkaListenerData {
	return &KafkaListenerData{
		Topic:     topic,
		NumThread: numThread,
	}
}

/**
 * Listener interface for Kafka
 */
func (k *KafkaListenerData) Listener(ctx context.Context) error {
	mq, err := domainMq.GetKafkaReadService()
	if err != nil {
		return err
	}

	// nếu global.ContextSystem được thiết lập thì ưu tiên dùng để quản lý dừng/bắt đầu toàn hệ thống
	ctxUsed := ctx
	if global.ContextSystem != nil {
		ctxUsed = global.ContextSystem
	}

	for i := 0; i < k.NumThread; i++ {
		global.WaitGroup.Add(1)
		// Start a goroutine for each thread
		go func(thread int) {
			defer global.WaitGroup.Done()

			for {
				// nếu context bị huỷ thì thoát goroutine
				if ctxUsed.Err() != nil {
					global.Logger.Warn("Kafka listener stopping", "topic", k.Topic, "thread", thread)
					return
				}

				start := time.Now()
				msg, err := mq.ReadMessageAutoCommit(ctxUsed, k.Topic)
				if err != nil {
					if ctxUsed.Err() != nil {
						global.Logger.Info("Kafka read aborted by context", "topic", k.Topic, "thread", thread)
						return
					}
					recordKafkaError(k.Topic, "read", time.Since(start).Seconds())
					global.Logger.Warn("Kafka read message error", "error", err)
					continue
				}
				// Parse message
				var event dto.KafkaEvent
				if err := json.Unmarshal(msg, &event); err != nil {
					recordKafkaError(k.Topic, "unmarshal", time.Since(start).Seconds())
					global.Logger.Warn("Kafka unmarshal message error", "error", err)
					continue
				}
				// Validate message
				if err := global.Validator.Struct(event); err != nil {
					recordKafkaError(k.Topic, "validate", time.Since(start).Seconds())
					// log validation error và tiếp tục
					global.Logger.Warn("Kafka validate message error", "error", err)
					continue
				}
				ctxSpan, span := startKafkaSpan(ctxUsed, k.Topic, thread)
				span.SetAttributes(attribute.String("kafka.event_type", strconv.Itoa(event.EventType)))

				// Send message to application handler
				switch event.EventType {
				case constants.KAFKA_EVENT_TYPE_SEND_TOKEN_RESET_PASSWORD:
					// chuyển payload về struct mục tiêu an toàn (payload có thể là map[string]interface{})
					payloadBytes, perr := json.Marshal(event.Payload)
					if perr != nil {
						span.RecordError(perr)
						recordKafkaError(k.Topic, "marshal_payload", time.Since(start).Seconds())
						global.Logger.Warn("Kafka marshal payload error", "error", perr)
						continue
					}
					var input applicationModel.MailForgotPassword
					if perr := json.Unmarshal(payloadBytes, &input); perr != nil {
						span.RecordError(perr)
						recordKafkaError(k.Topic, "unmarshal_payload", time.Since(start).Seconds())
						global.Logger.Warn("Kafka unmarshal payload to MailForgotPassword error", "error", perr)
						continue
					}
					if err := global.Validator.Struct(input); err != nil {
						span.RecordError(err)
						recordKafkaError(k.Topic, "validate_payload", time.Since(start).Seconds())
						global.Logger.Warn("Kafka validate data error", "error", err)
						continue
					}
					if err := applicationService.GetMailService().SendMessageForgotPassword(ctxSpan, input); err != nil {
						span.RecordError(err)
						recordKafkaError(k.Topic, "handle", time.Since(start).Seconds())
						global.Logger.Warn("Kafka send mail forgot password error", "error", err)
						continue
					}
				default:
					global.Logger.Warn("Kafka unknown event type", "event_type", event.EventType)
				}

				recordKafkaSuccess(k.Topic, time.Since(start).Seconds())
				span.End()
			}
		}(i)
	}
	return nil
}
