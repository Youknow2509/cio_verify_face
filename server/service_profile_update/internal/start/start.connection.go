package start

import (
	"fmt"

	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/domain/cache"
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/domain/config"
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/domain/mq"
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/global"
	infraCache "github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/infrastructure/cache"
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/infrastructure/conn"
	infraMq "github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/infrastructure/mq"
)

func initConnectionToInfrastructure(setting *config.Setting) error {
	// Initialize Redis connection
	redisCache, err := infraCache.NewRedisDistributedCache(&setting.Redis)
	if err != nil {
		return fmt.Errorf("failed to initialize Redis: %w", err)
	}

	if err := cache.SetDistributedCache(redisCache); err != nil {
		return fmt.Errorf("failed to set distributed cache: %w", err)
	}

	global.Logger.Info("Redis connection initialized successfully")

	// Initialize PostgreSQL connection
	if err := conn.InitPostgresqlClient(&setting.Postgres); err != nil {
		return fmt.Errorf("failed to initialize PostgreSQL: %w", err)
	}

	global.Logger.Info("PostgreSQL connection initialized successfully")

	// Initialize Kafka writer
	kafkaWriter := infraMq.NewKafkaWriterService(&setting.Kafka)
	if err := mq.SetKafkaWriter(kafkaWriter); err != nil {
		return fmt.Errorf("failed to set kafka writer: %w", err)
	}

	global.Logger.Info("Kafka writer initialized successfully")

	return nil
}
