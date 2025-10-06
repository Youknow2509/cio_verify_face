package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	applicationModel "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/application/model"
	applicationService "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/application/service"
	domainModel "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/model"
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
			DeviceId:  input.DeviceId.String(),
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
	case int(domainModel.WSEventMessageText):
		panic("Not implemented yet") // TODO: Implement this
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
