package core

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/global"
	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/interfaces/dto"
	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/interfaces/ws/handler"
)

/**
 * Worker structure for WebSocket
 */
type Worker struct {
	ctx context.Context
	id  int
}

// Worker run
func (h *Worker) Run() {
	for {
		select {
		case data := <-GetHub().RegisterChan:
			if err := handler.GetWorkerHandler().RegisterClient(
				h.ctx,
				data,
			); err != nil {
				global.Logger.Warn(fmt.Sprintf("Failed to register connection: %v", err))
			}
		case data := <-GetHub().UnregisterChan:
			if err := handler.GetWorkerHandler().UnregisterClient(
				h.ctx,
				data,
			); err != nil {
				global.Logger.Warn(fmt.Sprintf("Failed to unregister connection: %v", err))
			}
		case data := <-GetHub().HandlerReceive:
			var dataObj dto.DataServerReceive
			if err := json.Unmarshal(data.Data, &dataObj); err != nil {
				global.Logger.Warn(fmt.Sprintf("Failed to unmarshal data: %v", err))
				continue
			}
			if err := handler.GetWorkerHandler().HandleDataReceive(
				h.ctx,
				data.ClientInfo,
				dataObj.Type,
				dataObj.Payload,
			); err != nil {
				global.Logger.Warn(fmt.Sprintf("Failed to handle data receive: %v", err))
				continue
			}
		case data := <-GetHub().SendData:
			if err := handler.GetWorkerHandler().SendDataToClient(
				h.ctx,
				data.ConnectionId,
				data.Data,
			); err != nil {
				global.Logger.Warn(fmt.Sprintf("Failed to send data to client %s: %v", data.ConnectionId, err))
			}
		case <-h.ctx.Done():
			global.Logger.Warn(fmt.Sprintf("Handler %d stopped", h.id))
			return
		}
	}
}

func NewWorker(ctx context.Context, id int) *Worker {
	return &Worker{
		ctx: ctx,
		id:  id,
	}
}
