package start

import "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/global"

/**
 * Start service
 */
func StartService() error {
	// Load configuration
	setting, err := loadConfig()
	if err != nil {
		return err
	}
	global.ServerSetting = setting.Server
	global.ServerWsSetting = setting.WsServer
	global.ServerGrpcSetting = setting.GrpcServer

	// Initialize logger
	if err := initLogger(&setting.Logger); err != nil {
		return err
	}
	// Initialize observability (Prometheus metrics and Jaeger tracing)
	if err := initObservability(&setting.Observability, setting.Server.Name); err != nil {
		return err
	}
	// Initialize local cache
	if err := initLocalCache(); err != nil {
		return err
	}
	// Initialize rate limit
	if err := initRateLimit(setting.RateLimitPolicies); err != nil {
		return err
	}
	// Initialize connection to infrastructure
	if err := initConnectionToInfrastructure(setting); err != nil {
		return err
	}
	// Initialize domain
	if err := initDomain(); err != nil {
		return err
	}
	// Initialize application
	if err := initApplication(); err != nil {
		return err
	}
	// Initialize validator
	if err := initValidator(); err != nil {
		return err
	}
	// Initialize WS
	if err := initWebSocketServer(); err != nil {
		return err
	}
	// Initialize gRPC server
	if err := initGrpcServer(); err != nil {
		return err
	}
	// Initialize Gin Engine
	if err := initGinRouter(); err != nil {
		return err
	}
	return nil
}
