package impl

import (
	"context"

	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/application/model"
	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/application/service"
	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/global"

	domainModel "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/model"
	domainMq "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/mq"
)

/**
 * Send event service struct
 */
type SendEventService struct{}

// SendDataVerifyFace implements service.ISendEventService.
func (s *SendEventService) SendDataVerifyFace(ctx context.Context, input *model.SendDataVerifyFace) error {
	serviceMq := domainMq.GetSendEventToKafka()
	if serviceMq == nil {
		global.Logger.Panic("SendEventService.SendDataVerifyFace: failed to get SendEventToKafka service")
	}
	err := serviceMq.SendDataVerify(
		ctx,
		domainModel.KafkaAttendanceVerifyReceived{
			ServiceId: global.ServerSetting.Id,
			DeviceId: input.DeviceId.String(),
			DataUrl: input.DataUrl,
			Metadata: input.Metadata,
			Timestamp: input.Timestamp,
		},
	)
	return err
}

/**
 * New send event service and implementation
 */
func NewSendEventService() service.ISendEventService {
	return &SendEventService{}
}
