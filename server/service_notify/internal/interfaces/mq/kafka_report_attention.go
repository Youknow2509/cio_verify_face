package mq

import (
	"context"
	"encoding/json"

	applicationModel "github.com/youknow2509/cio_verify_face/server/service_notify/internal/application/model"
	applicationService "github.com/youknow2509/cio_verify_face/server/service_notify/internal/application/service"
	domainMq "github.com/youknow2509/cio_verify_face/server/service_notify/internal/domain/mq"
	"github.com/youknow2509/cio_verify_face/server/service_notify/internal/global"
	"github.com/youknow2509/cio_verify_face/server/service_notify/internal/interfaces/dto"
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

				msg, err := mq.ReadMessageAutoCommit(ctxUsed, k.Topic)
				if err != nil {
					if ctxUsed.Err() != nil {
						global.Logger.Info("Report attention Kafka read aborted by context", "topic", k.Topic, "thread", thread)
						return
					}
					global.Logger.Warn("Report attention Kafka read message error", "error", err)
					continue
				}

				// Parse message to ReportAttentionNotificationEvent
				var event dto.ReportAttentionNotificationEvent
				if err := json.Unmarshal(msg, &event); err != nil {
					global.Logger.Warn("Report attention Kafka unmarshal message error", "error", err, "message", string(msg))
					continue
				}

				// Validate message
				if err := global.Validator.Struct(event); err != nil {
					global.Logger.Warn("Report attention Kafka validate message error", "error", err)
					continue
				}

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
				if err := applicationService.GetMailService().SendReportAttentionNotification(ctxUsed, input); err != nil {
					global.Logger.Warn("Report attention send mail error", "error", err)
					continue
				}

				global.Logger.Info("Report attention notification sent successfully", "to", input.Email, "company_id", input.CompanyID)
			}
		}(i)
	}
	return nil
}
