package start

import "github.com/youknow2509/cio_verify_face/server/service_auth/internal/global"

/**
 * Start service
 */
func StartService() error {
	// Load configuration
	setting, err := loadConfig()
	if err != nil {
		return err
	}
	global.SettingServer = *setting
	// Initialize logger
	if err := initLogger(&setting.Logger); err != nil {
		return err
	}
	// Initialize observability (Prometheus metrics and Jaeger tracing)
	if err := initObservability(&setting.Observability, setting.Server.Name); err != nil {
		return err
	}
	// Inittialize Ristretto - in-memory cache
	if err := initLocalCache(); err != nil {
		return err
	}
	// Initialize connection to infrastructure
	if err := initConnectionToInfrastructure(setting); err != nil {
		return err
	}
	// Initialize Domain
	if err := initDomain(); err != nil {
		return err
	}
	// Initialize application
	if err := initApplication(); err != nil {
		return err
	}
	// Initialize Grpc Server
	if err := initServerGrpc(); err != nil {
		return err
	}
	// Initialize Gin Engine
	if err := initGinRouter(&setting.Server); err != nil {
		return err
	}
	return nil
}
