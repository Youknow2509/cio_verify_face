package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	applicationModel "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/application/model"
	applicationService "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/application/service"
	domainModel "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/model"
	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/global"
	dto "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/interfaces/dto"
	model "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/interfaces/ws/model"
)

/**
 * Interface worker handle
 */
type WorkerHandler struct {
}

/**
 * Get WorkerHandler instance
 */

func GetWorkerHandler() *WorkerHandler {
	return &WorkerHandler{}
}

/**
 * Register new client
 */
func (wh *WorkerHandler) RegisterClient(ctx context.Context, input model.ClientInfo) error {
	service := applicationService.GetMapConnectionService()
	if err := service.RegisterConnection(
		ctx,
		&applicationModel.RegisterConnection{
			ConnectionId: input.ConnectionId.String(),
			DeviceId:     input.DeviceId.String(),
			IPAddress:    input.IpAddress,
			ConnectedAt:  time.Now().Format("2006-01-02 15:04:05"),
			UserAgent:    input.UserAgent,
		},
	); err != nil {
		return fmt.Errorf("failed to register connection: %v", err)
	}
	return nil
}

/**
 * Unregister client
 */
func (wh *WorkerHandler) UnregisterClient(ctx context.Context, input model.ClientInfo) error {
	service := applicationService.GetMapConnectionService()
	if err := service.UnregisterConnection(
		ctx,
		&applicationModel.UnregisterConnection{
			DeviceId: input.DeviceId.String(),
		},
	); err != nil {
		return fmt.Errorf("failed to unregister connection: %v", err)
	}
	return nil
}

/**
 * Handle data receive
 */
func (wh *WorkerHandler) HandleDataReceive(ctx context.Context, clientInfo model.ClientInfo, eventType int, data []byte) error {
	switch eventType {
	case int(domainModel.WSEventReceivedAttendance):
		// Unmarshal data
		var attendanceData dto.SendDataVerifyFace
		if err := json.Unmarshal(data, &attendanceData); err != nil {
			wh.SendDataToClient(
				ctx,
				clientInfo.ConnectionId,
				[]byte("invalid attendance data format"),
			)
			return fmt.Errorf("failed to unmarshal attendance data: %v", err)
		}
		// Validate data
		if err := global.Validate.Struct(attendanceData); err != nil {
            errs := err.(validator.ValidationErrors)
            var messages []string
            for _, e := range errs {
                messages = append(messages, fmt.Sprintf("%s: %s", e.Field(), e.Tag()))
            }
			err_str := strings.Join(messages, ", ")
			wh.SendDataToClient(
				ctx,
				clientInfo.ConnectionId,
				[]byte(err_str),
			)
            return errors.New(err_str)
        }
		// Send data to verify face service
		if err := applicationService.GetSendEventService().SendDataVerifyFace(
			ctx,
			&applicationModel.SendDataVerifyFace{
				DeviceId:  clientInfo.DeviceId,
				DataUrl:   attendanceData.DataUrl,
				Metadata:  attendanceData.Metadata,
				Timestamp: time.Now().Unix(),
			},
		); err != nil {
			global.Logger.Warn("failed to send data verify face: %v", err)
			return errors.New("failed to send data verify face")
		}
		return nil
	default:
		errorBytes, _ := json.Marshal(fmt.Sprintf("Unknown event type: %d", eventType))
		err := wh.SendDataToClient(ctx, clientInfo.ConnectionId, errorBytes)
		if err != nil {
			return err
		}
	}
	return nil
}

/**
 * Send data to client
 */
func (wh *WorkerHandler) SendDataToClient(ctx context.Context, clientId uuid.UUID, data []byte) error {
	err := applicationService.GetClientService().SendMessageToClient(
		ctx,
		&applicationModel.SendMessageToClientInput{
			ConnectionId: clientId,
			Message:      data,
		},
	)
	return err
}
