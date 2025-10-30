package router

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/global"
	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/interfaces/ws/core"
)

/**
 * Base router struct
 */
type BaseRouter struct{}

func GetBaseRouter() *BaseRouter {
	return &BaseRouter{}
}

/**
 * Init base router struct
 */
func (b *BaseRouter) Initialize(g *gin.Engine) {
	group := g.Group("/ws")
	group.GET("", func(c *gin.Context) {
		// create upgrade connection
		upgradeObj := websocket.Upgrader{
			HandshakeTimeout: time.Second * time.Duration(global.ServerWsSetting.HandshakeTimeout),
			ReadBufferSize:   global.ServerWsSetting.ReadBufferSize,
			WriteBufferSize:  global.ServerWsSetting.WriteBufferSize,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			EnableCompression: global.ServerWsSetting.EnableCompression,
		}
		// create connection handshake
		conn, err := upgradeObj.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			global.Logger.Warn("Failed to upgrade connection", "error", err)
			c.Error(err)
			return
		}
		// Get data user after auth
		deviceId, _ := uuid.NewUUID() // TODO: fix in prod
		// Create client info
		clientInfo := core.NewClientWS(
			c,
			deviceId,
			c.ClientIP(),
			c.Request.UserAgent(),
			conn,
			global.ServerWsSetting.MaxMessageSize,
			time.Second*time.Duration(global.ServerWsSetting.ReadWait),
			time.Second*time.Duration(global.ServerWsSetting.WriteWait),
			time.Second*time.Duration(global.ServerWsSetting.PingPeriod),
			global.ServerWsSetting.MaxSendQueueSize,
		)
		// register with hub
		core.GetHub().RegisterClient(clientInfo)
		// read, write
		go clientInfo.ReadPump()
		go clientInfo.WritePump()
	})
}
