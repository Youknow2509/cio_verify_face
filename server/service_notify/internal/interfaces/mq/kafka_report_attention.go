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
	"go.opentelemetry.io/otel/attribute"
)

/**
 * Report Attention Notification Kafka listener
 */
type ReportAttentionKafkaListener struct {
	Topic     string
	NumThread int
}

/**
 * New Report Attention Kafka listener
 */
func NewReportAttentionKafkaListener(topic string, numThread int) *ReportAttentionKafkaListener {
	return &ReportAttentionKafkaListener{
		Topic:     topic,
		NumThread: numThread,
	}
}

/**
 * Listener for report attention notifications
 */
func (k *ReportAttentionKafkaListener) Listener(ctx context.Context) error {
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
					global.Logger.Warn("Report attention Kafka listener stopping", "topic", k.Topic, "thread", thread)
					return
				}

				start := time.Now()
				msg, err := mq.ReadMessageAutoCommit(ctxUsed, k.Topic)
				if err != nil {
					if ctxUsed.Err() != nil {
						global.Logger.Info("Report attention Kafka read aborted by context", "topic", k.Topic, "thread", thread)
						return
					}
					recordKafkaError(k.Topic, "read", time.Since(start).Seconds())
					global.Logger.Warn("Report attention Kafka read message error", "error", err)
					continue
				}

				// Parse message to ReportAttentionNotificationEvent
				var event dto.ReportAttentionNotificationEvent
				if err := json.Unmarshal(msg, &event); err != nil {
					recordKafkaError(k.Topic, "unmarshal", time.Since(start).Seconds())
					global.Logger.Warn("Report attention Kafka unmarshal message error", "error", err, "message", string(msg))
					continue
				}

				// Validate message
				if err := global.Validator.Struct(event); err != nil {
					recordKafkaError(k.Topic, "validate", time.Since(start).Seconds())
					global.Logger.Warn("Report attention Kafka validate message error", "error", err)
					continue
				}

				ctxSpan, span := startKafkaSpan(ctxUsed, k.Topic, thread)
				span.SetAttributes(attribute.String("report_attention.type", event.Type), attribute.String("report_attention.format", event.Format))

				// Convert to application model
				input := applicationModel.ReportAttentionNotification{
					Email:       event.Email,
					CompanyID:   event.CompanyID,
					DownloadURL: event.DownloadURL,
					Type:        event.Type,
					Format:      event.Format,
					StartDate:   event.StartDate,
					EndDate:     event.EndDate,
					CreatedAt:   event.CreatedAt,
				}

				// Send notification
				if err := applicationService.GetMailService().SendReportAttentionNotification(ctxSpan, input); err != nil {
					span.RecordError(err)
					recordKafkaError(k.Topic, "handle", time.Since(start).Seconds())
					global.Logger.Warn("Report attention send mail error", "error", err)
					continue
				}

				global.Logger.Info("Report attention notification sent successfully", "to", input.Email, "company_id", input.CompanyID)
				recordKafkaSuccess(k.Topic, time.Since(start).Seconds())
				span.End()
			}
		}(i)
	}
	return nil
}
