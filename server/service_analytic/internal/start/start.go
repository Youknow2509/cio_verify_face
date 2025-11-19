package start

import (
	"github.com/youknow2509/cio_verify_face/server/service_analytic/internal/global"
	infraCache "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/infrastructure/cache"
	infraConfig "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/infrastructure/config"
	infraConn "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/infrastructure/conn"
	infraLogger "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/infrastructure/logger"
	infraRepository "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/infrastructure/repository"
	domainRepository "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/domain/repository"
	applicationService "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/application/service"
	applicationServiceImpl "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/application/service/impl"
)

// StartService starts the analytics service
func StartService() error {
	// Load configuration
	config, err := infraConfig.LoadConfig()
	if err != nil {
		return err
	}
	global.SettingServer = *config

	// Initialize logger
	logger, zapLogger, err := infraLogger.NewZapLogger(&config.Logger)
	if err != nil {
		return err
	}
	global.Logger = logger
	global.ZapLogger = zapLogger
	logger.Info("Logger initialized")

	// Initialize local cache
	if err := infraCache.InitLocalCache(); err != nil {
		return err
	}
	cache, _ := infraCache.GetLocalCache()
	global.LocalCache = cache
	logger.Info("Local cache initialized")

	// Initialize PostgreSQL connection
	if err := infraConn.InitPostgresClient(&config.Postgres); err != nil {
		return err
	}
	pgPool, _ := infraConn.GetPostgresPool()
	global.PostgresPool = pgPool
	logger.Info("PostgreSQL connection initialized")

	// Initialize ScyllaDB connection
	if err := infraConn.InitScylladbClient(&config.Scylla); err != nil {
		return err
	}
	logger.Info("ScyllaDB connection initialized")

	// Initialize Redis connection
	if err := infraConn.InitRedisClient(&config.Redis); err != nil {
		return err
	}
	redisClient, _ := infraConn.GetRedisClient()
	global.RedisClient = redisClient
	logger.Info("Redis connection initialized")

	// Initialize auth gRPC client
	if err := initAuthGrpcClient(); err != nil {
		logger.Warn("Auth gRPC client initialization skipped", "error", err)
		// Don't fail if auth client can't be initialized
	}

	// Initialize repositories
	scyllaSession, _ := infraConn.GetScylladbClient()
	analyticRepo := infraRepository.NewAnalyticRepository(scyllaSession, pgPool)
	if err := domainRepository.SetAnalyticRepository(analyticRepo); err != nil {
		return err
	}
	logger.Info("Repositories initialized")

	// Initialize application services
	analyticService := applicationServiceImpl.NewAnalyticService(analyticRepo)
	if err := applicationService.SetAnalyticService(analyticService); err != nil {
		return err
	}
	logger.Info("Application services initialized")

	// Initialize gRPC server
	if err := initGrpcServer(); err != nil {
		logger.Warn("gRPC server initialization failed", "error", err)
		// Don't fail if gRPC can't be initialized
	}

	// Initialize HTTP server
	if err := initHTTPServer(&config.Server); err != nil {
		return err
	}
	logger.Info("HTTP server initialized")

	return nil
}
