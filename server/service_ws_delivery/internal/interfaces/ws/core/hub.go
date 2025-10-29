package core

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/constants"
	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/global"
	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/interfaces/ws/model"
)

// =============================================
//
//	Hub WS Structure
//
// =============================================
type Hub struct {
	// Client management
	clients      map[uuid.UUID]*Client
	clientsMutex sync.RWMutex
	register     chan *Client
	unregister   chan *Client
	// Worker Manager
	numWorkers     int
	RegisterChan   chan model.ClientInfo
	UnregisterChan chan model.ClientInfo
	HandlerReceive chan model.ClientWriterData
	SendData       chan model.DataSend
}

// Manager instance hub
var _vHub *Hub

// Getter hub
func GetHub() *Hub {
	return _vHub
}

// Setter hub
func SetHub(hub *Hub) error {
	if hub == nil {
		return errors.New("hub is nil")
	}
	if _vHub != nil {
		return errors.New("hub already set")
	}
	_vHub = hub
	return nil
}

// NewHub được cập nhật
func NewHub(
	numWorkers int,
	sizeRegister int,
	sizeUnregister int,
	sizeReceive int,
	sizeSend int,
) *Hub {
	return &Hub{
		// Client Management
		clients:    make(map[uuid.UUID]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		// Worker Manager
		numWorkers:     numWorkers,
		RegisterChan:   make(chan model.ClientInfo, sizeRegister),
		UnregisterChan: make(chan model.ClientInfo, sizeUnregister),
		HandlerReceive: make(chan model.ClientWriterData, sizeReceive),
		SendData:       make(chan model.DataSend, sizeSend),
	}
}

// Run hub service
func (h *Hub) Run(ctx context.Context) {
	defer func() {
		global.WaitGroup.Wait()
	}()
	// Run worker
	for i := 1; i <= h.numWorkers; i++ {
		go func(id int) {
			defer func() {
				global.WaitGroup.Done()
			}()
			global.WaitGroup.Add(1)
			NewWorker(ctx, id).Run()
		}(i)
	}
	//
	for {
		select {
		case client := <-h.register:
			h.clientsMutex.Lock()
			if len(h.clients) > global.ServerWsSetting.MaxConnectionSystem {
				h.clientsMutex.Unlock()
				global.Logger.Warn(fmt.Sprintf("Max connection limit reached: %d", global.ServerWsSetting.MaxConnectionSystem))
				continue
			}
			h.clients[client.ConnId] = client
			h.clientsMutex.Unlock()
			h.RegisterChan <- model.ClientInfo{
				DeviceId:     client.DeviceId,
				ConnectionId: client.ConnId,
				UserAgent:    client.ClientUserAgent,
				IpAddress:    client.ClientIpAddress,
			}
			global.WaitGroup.Add(1)

		case client := <-h.unregister:
			h.clientsMutex.Lock()
			if _, ok := h.clients[client.ConnId]; ok {
				delete(h.clients, client.ConnId)
				h.UnregisterChan <- model.ClientInfo{
					DeviceId:     client.DeviceId,
					ConnectionId: client.ConnId,
					UserAgent:    client.ClientUserAgent,
					IpAddress:    client.ClientIpAddress,
				}
			}
			h.clientsMutex.Unlock()

		case <-ctx.Done():
			close(h.register)
			close(h.unregister)
			return
		}
	}
}

// Register client
func (h *Hub) RegisterClient(client *Client) {
	select {
	case h.register <- client:
		global.Logger.Warn(fmt.Sprintf("Client %s registered successfully", client.ConnId))
	default:
		global.Logger.Warn(fmt.Sprintf("Failed to register client %s: hub is shutting down", client.ConnId))
	}
}

// Unregister client
func (h *Hub) UnregisterClient(client *Client) {
	select {
	case h.unregister <- client:
		global.Logger.Warn(fmt.Sprintf("Client %s unregistered successfully", client.ConnId))
	default:
		global.Logger.Warn(fmt.Sprintf("Failed to unregister client %s: hub is shutting down", client.ConnId))
	}
}

// Send data to client
func (h *Hub) SendDataToClient(ctx context.Context, input model.DataSend) error {
	timeOut := time.Duration(constants.TIME_OUT_SEND_MSG_TO_CHAN * float64(time.Second))
	select {
	case h.SendData <- input:
		return nil
	case <-time.After(timeOut):
		return errors.New("send data to client timeout: channel is full")
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Get client
func (h *Hub) GetClient(connId uuid.UUID) (*Client, bool) {
	h.clientsMutex.RLock()
	defer h.clientsMutex.RUnlock()
	client, ok := h.clients[connId]
	return client, ok
}

// Num client
func (h *Hub) NumClients() int {
	h.clientsMutex.RLock()
	defer h.clientsMutex.RUnlock()
	return len(h.clients)
}
