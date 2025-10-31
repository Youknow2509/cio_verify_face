package start

import (
	"context"
	"sync"

	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/global"
	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/interfaces/ws/core"
)

/**
 * WebSocket server startup
 */
func initWebSocketServer() error {
	// Create ws context
	global.WsContext = context.Background()
	// Create wait group
	global.WaitGroup = &sync.WaitGroup{}
	// Create hub
	hub := core.NewHub(
		global.ServerWsSetting.NumWorkers,
		global.ServerWsSetting.SizeBufferChan,
		global.ServerWsSetting.SizeBufferChan,
		global.ServerWsSetting.SizeBufferChan,
		global.ServerWsSetting.SizeBufferChan,
	)
	core.SetHub(hub)
	// run hub
	go hub.Run(global.WsContext)
	return nil
}
