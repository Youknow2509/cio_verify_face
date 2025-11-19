package conn

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	domainConfig "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/domain/config"
)

// redis client variables
var (
	vRedisClient *redis.Client
)

// InitRedisClient initializes the Redis client
func InitRedisClient(redisConfig *domainConfig.RedisConfig) error {
	if vRedisClient != nil {
		return errors.New("Redis client is already initialized")
	}

	// Create Redis client based on type
	var client *redis.Client

	switch redisConfig.Type {
	case 1: // Standalone
		client = redis.NewClient(&redis.Options{
			Addr:         fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
			Password:     redisConfig.Password,
			DB:           redisConfig.DB,
			PoolSize:     redisConfig.PoolSize,
			MinIdleConns: redisConfig.MinIdleConns,
			MaxRetries:   redisConfig.MaxRetries,
		})
	default:
		return fmt.Errorf("unsupported redis type: %d", redisConfig.Type)
	}

	// Ping to verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to ping redis: %w", err)
	}

	vRedisClient = client
	return nil
}

// GetRedisClient returns the Redis client
func GetRedisClient() (*redis.Client, error) {
	if vRedisClient == nil {
		return nil, errors.New("Redis client is not initialized, please call InitRedisClient first")
	}
	return vRedisClient, nil
}

// CloseRedisClient closes the Redis connection
func CloseRedisClient() error {
	if vRedisClient != nil {
		err := vRedisClient.Close()
		vRedisClient = nil
		return err
	}
	return nil
}
