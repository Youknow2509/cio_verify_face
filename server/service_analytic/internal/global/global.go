package global

import (
	"sync"

	"github.com/dgraph-io/ristretto"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/youknow2509/cio_verify_face/server/service_analytic/internal/domain/config"
	"github.com/youknow2509/cio_verify_face/server/service_analytic/internal/domain/logger"
	"go.uber.org/zap"
)

var (
	// Server settings
	SettingServer config.Config

	// Logger
	Logger logger.ILogger

	// Database connections
	PostgresPool *pgxpool.Pool

	// Cache
	RedisClient  *redis.Client
	LocalCache   *ristretto.Cache

	// Wait group for graceful shutdown
	WaitGroup *sync.WaitGroup

	// Zap logger instance
	ZapLogger *zap.Logger
)
