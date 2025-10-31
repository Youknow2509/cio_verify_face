package impl

import (
	"context"
	"errors"

	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/application/model"
	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/application/service"
	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/global"
	wsCore "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/interfaces/ws/core"
)

/**
 * Client service struct
 */
type ClientService struct{}

// SendMessageToClient implements service.IClientService.
func (c *ClientService) SendMessageToClient(ctx context.Context, input *model.SendMessageToClientInput) error {
	// Get instance use
	var (
		hub    *wsCore.Hub
		client *wsCore.Client
	)
	hub = wsCore.GetHub()
	client, ok := hub.GetClient(input.ConnectionId)
	if !ok {
		return errors.New("client not found")
	}
	// Send message to client
	err := client.Send(input.Message)
	if err != nil {
		global.Logger.Warn("Failed to send message to client", err)
		return err
	}
	return nil
}

/**
 * New send event service and implementation
 */
func NewClientService() service.IClientService {
	return &ClientService{}
}
