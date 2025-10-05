package core

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	libsDomainModel "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/model"
	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/global"
	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/interfaces/ws/model"
)

// TODO: Handle protobuf message
// =============================================
//
//	Client WS Structure
//
// =============================================
type Client struct {
	// Info Client
	SessionId       uuid.UUID
	UserId          uuid.UUID
	ConnId          uuid.UUID
	ClientIpAddress string
	ClientUserAgent string
	// Connection
	conn *websocket.Conn
	// Message queue for outgoing messages, protected by mutex.
	sendQueue      [][]byte
	sendQueueMutex sync.Mutex
	// Config Client
	maxMessageSize   int64
	maxSendQueueSize int
	readWait         time.Duration
	writeWait        time.Duration
	pingPeriod       time.Duration
	// Context management lifecycle
	ctx    context.Context
	cancel context.CancelFunc
}

// Send handle add message to the queue safely
func (c *Client) Send(message []byte) error {
	c.sendQueueMutex.Lock()
	defer c.sendQueueMutex.Unlock()
	// Disconnection connection client if send queue is full - Full memory
	if len(c.sendQueue) >= c.maxSendQueueSize {
		c.cancel()
		return fmt.Errorf("Client %s send queue is full, disconnecting client", c.ConnId)
	}
	c.sendQueue = append(c.sendQueue, message)
	return nil
}

// readPump - listen for messages from the client.
func (c *Client) ReadPump() {
	defer func() {
		GetHub().UnregisterClient(c)
		c.conn.Close()
		global.WaitGroup.Done()
	}()
	// Set read limit and deadlines
	c.conn.SetReadLimit(c.maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(c.readWait))
	c.conn.SetPongHandler(func(string) error {
		_ = c.conn.SetReadDeadline(time.Now().Add(c.readWait))
		return nil
	})
	//
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			global.Logger.Warn(fmt.Sprintf("Error reading message from client %s: %v", c.ConnId, err))
			c.cancel()
			return
		}
		// Rate limit check
		verdict, err := global.RateLimitWsRead.Allow(c.ctx, c.ConnId.String())
		if err != nil {
			global.Logger.Error(fmt.Sprintf("Error rate limit ws read for client %s: %v", c.ConnId, err))
			c.cancel()
			return
		}
		if verdict == libsDomainModel.Denied {
			c.Send([]byte("Error: Too many requests, please slow down"))
			time.Sleep(100 * time.Millisecond) // Give some time for the message to be sent
			global.Logger.Warn(fmt.Sprintf("Client %s disconnected due to rate limiting", c.ConnId))
			c.cancel()
			return
		}
		// Message handling
		select {
		case GetHub().HandlerReceive <- model.ClientWriterData{
			ClientInfo: model.ClientInfo{
				ConnectionId: c.ConnId,
				UserId:       c.UserId,
				SessionId:    c.SessionId,
				IpAddress:    c.ClientIpAddress,
				UserAgent:    c.ClientUserAgent,
			},
			Data: msg,
		}:
		case <-c.ctx.Done():
			return
		}
	}
}

// writePump send batch messages to the websocket.
func (c *Client) WritePump() {
	ticker := time.NewTicker(c.pingPeriod)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(c.writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.cancel()
				return
			}
		case <-c.ctx.Done():
			return
		default:
			// Get batch messages from sendQueue
			c.sendQueueMutex.Lock()
			if len(c.sendQueue) == 0 {
				c.sendQueueMutex.Unlock()
				// Wait a bit before checking again to avoid busy-loop
				time.Sleep(50 * time.Millisecond)
				continue
			}
			// Copy the queue to a local variable and clear the original queue
			messages := make([][]byte, len(c.sendQueue))
			copy(messages, c.sendQueue)
			c.sendQueue = c.sendQueue[:0]
			c.sendQueueMutex.Unlock()
			// Send all messages that have been retrieved
			for _, message := range messages {
				_ = c.conn.SetWriteDeadline(time.Now().Add(c.writeWait))
				if err := c.conn.WriteMessage(websocket.BinaryMessage, message); err != nil {
					global.Logger.Warn(fmt.Sprintf("Error sending message to client %s: %v", c.ConnId, err))
					return
				}
			}
		}
	}
}

// NewClientWS
func NewClientWS(
	ctx context.Context,
	sessionId uuid.UUID,
	userId uuid.UUID,
	clientIpAddress string,
	clientUserAgent string,
	conn *websocket.Conn,
	maxMessageSize int64,
	readWait time.Duration,
	writeWait time.Duration,
	pingPeriod time.Duration,
	maxSendQueueSize int,
) *Client {
	clientCtx, cancel := context.WithCancel(ctx)
	return &Client{
		SessionId:        sessionId,
		UserId:           userId,
		ClientIpAddress:  clientIpAddress,
		ClientUserAgent:  clientUserAgent,
		ConnId:           uuid.New(),
		conn:             conn,
		sendQueue:        make([][]byte, 0, maxSendQueueSize),
		maxSendQueueSize: maxSendQueueSize,
		maxMessageSize:   maxMessageSize,
		readWait:         readWait,
		writeWait:        writeWait,
		pingPeriod:       pingPeriod,
		ctx:              clientCtx,
		cancel:           cancel,
	}
}
