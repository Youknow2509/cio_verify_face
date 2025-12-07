package mq

import (
	"context"
	"encoding/json"
	"time"

	applicationModel "github.com/youknow2509/cio_verify_face/server/service_notify/internal/application/model"
	applicationService "github.com/youknow2509/cio_verify_face/server/service_notify/internal/application/service"
	domainMq "github.com/youknow2509/cio_verify_face/server/service_notify/internal/domain/mq"
	"github.com/youknow2509/cio_verify_face/server/service_notify/internal/global"
	"github.com/youknow2509/cio_verify_face/server/service_notify/internal/interfaces/dto"
)

/**
 * Password Reset Notification Kafka listener
 */
type PasswordResetKafkaListener struct {
	Topic     string
	NumThread int
}

/**
 * New Password Reset Kafka listener
 */
func NewPasswordResetKafkaListener(topic string, numThread int) *PasswordResetKafkaListener {
	return &PasswordResetKafkaListener{
		Topic:     topic,
		NumThread: numThread,
	}
}

/**
 * Listener for password reset notifications
 */
func (k *PasswordResetKafkaListener) Listener(ctx context.Context) error {
	mq, err := domainMq.GetKafkaReadService()
	if err != nil {
		return err
	}

	ctxUsed := ctx
	if global.ContextSystem != nil {
		ctxUsed = global.ContextSystem
	}

	for i := 0; i < k.NumThread; i++ {
		global.WaitGroup.Add(1)
		go func(thread int) {
			defer global.WaitGroup.Done()

			for {
				if ctxUsed.Err() != nil {
					global.Logger.Warn("Password reset Kafka listener stopping", "topic", k.Topic, "thread", thread)
					return
				}

				start := time.Now()
				msg, err := mq.ReadMessageAutoCommit(ctxUsed, k.Topic)
				if err != nil {
					if ctxUsed.Err() != nil {
						global.Logger.Info("Password reset Kafka read aborted by context", "topic", k.Topic, "thread", thread)
						return
					}
					recordKafkaError(k.Topic, "read", time.Since(start).Seconds())
					global.Logger.Warn("Password reset Kafka read message error", "error", err)
					continue
				}

				// Parse message to PasswordResetNotificationEvent
				var event dto.PasswordResetNotificationEvent
				if err := json.Unmarshal(msg, &event); err != nil {
					recordKafkaError(k.Topic, "unmarshal", time.Since(start).Seconds())
					global.Logger.Warn("Password reset Kafka unmarshal message error", "error", err, "message", string(msg))
					continue
				}

				// Validate message
				if err := global.Validator.Struct(event); err != nil {
					recordKafkaError(k.Topic, "validate", time.Since(start).Seconds())
					global.Logger.Warn("Password reset Kafka validate message error", "error", err)
					continue
				}

				ctxSpan, span := startKafkaSpan(ctxUsed, k.Topic, thread)

				// Convert to application model
				input := applicationModel.PasswordResetNotification{
					To:        event.Payload.To,
					FullName:  event.Payload.FullName,
					ResetURL:  event.Payload.ResetURL,
					ExpiresIn: event.Payload.ExpiresIn,
				}

				// Send notification
				if err := applicationService.GetMailService().SendPasswordResetNotification(ctxSpan, input); err != nil {
					span.RecordError(err)
					recordKafkaError(k.Topic, "handle", time.Since(start).Seconds())
					global.Logger.Warn("Password reset send mail error", "error", err)
					continue
				}

				global.Logger.Info("Password reset notification sent successfully", "to", input.To, "user_id", event.Metadata.UserID)
				recordKafkaSuccess(k.Topic, time.Since(start).Seconds())
				span.End()
			}
		}(i)
	}
	return nil
}
