package clients

import (
	"context"
	"crypto/tls"
	"errors"
	"strconv"

	"github.com/redis/go-redis/v9"
	domainConfig "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/config"
	"github.com/youknow2509/cio_verify_face/server/service_attendance/internal/constants"
)

// redisClient variable holds the Redis client instance.
var (
	vRedisClient *redis.Client
)

/**
 * Initializes the Redis client with the provided options.
 * @param *config.RedisSetting redisSetting - The Redis configuration settings.
 * @return error - Returns an error if the initialization fails.
 */
func InitRedisClient(redisSetting *domainConfig.RedisSetting) error {
	if redisSetting == nil {
		return errors.New("redis settings are nil")
	}

	var client redis.UniversalClient

	switch redisSetting.Type {
	case constants.RedisTypeStandalone: // Standalone
		opt := &redis.Options{
			Addr:         redisSetting.Host + ":" + strconv.Itoa(redisSetting.Port),
			Password:     redisSetting.Password,
			DB:           redisSetting.DB,
			PoolSize:     redisSetting.PoolSize,
			MinIdleConns: redisSetting.MinIdleConns,
			MaxRetries:   redisSetting.MaxRetries,
		}
		if redisSetting.UseTLS {
			opt.TLSConfig = &tls.Config{
				InsecureSkipVerify: false,
			}
		}
		client = redis.NewClient(opt)

	case constants.RedisTypeSentinel: // Sentinel
		opt := &redis.FailoverOptions{
			MasterName:    redisSetting.MasterName,
			SentinelAddrs: redisSetting.SentinelAddrs,
			Password:      redisSetting.Password,
			DB:            redisSetting.DB,
			PoolSize:      redisSetting.PoolSize,
			MinIdleConns:  redisSetting.MinIdleConns,
			MaxRetries:    redisSetting.MaxRetries,
		}
		if redisSetting.UseTLS {
			opt.TLSConfig = &tls.Config{
				InsecureSkipVerify: false,
			}
		}
		client = redis.NewFailoverClient(opt)

	case constants.RedisTypeCluster: // Cluster
		opt := &redis.ClusterOptions{
			Addrs:          redisSetting.Address,
			Password:       redisSetting.Password,
			PoolSize:       redisSetting.PoolSize,
			MinIdleConns:   redisSetting.MinIdleConns,
			MaxRetries:     redisSetting.MaxRetries,
			RouteByLatency: redisSetting.RouteByLatency,
			RouteRandomly:  redisSetting.RouteRandomly,
		}
		if redisSetting.UseTLS {
			opt.TLSConfig = &tls.Config{
				InsecureSkipVerify: false,
			}
		}
		client = redis.NewClusterClient(opt)

	default:
		return errors.New("unsupported Redis type")
	}

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return errors.New("failed to connect to Redis: " + err.Error())
	}

	// Cast to *redis.Client if standalone/sentinel, else keep as UniversalClient
	switch c := client.(type) {
	case *redis.Client:
		vRedisClient = c
	case *redis.ClusterClient:
		// You may want to handle cluster client separately if needed
		// For now, do nothing
	}
	// set pool configurations
	vRedisClient.Options().PoolSize = redisSetting.PoolSize
	vRedisClient.Options().MinIdleConns = redisSetting.MinIdleConns
	vRedisClient.Options().MaxRetries = redisSetting.MaxRetries

	return nil
}

/**
 * Get the Redis client instance.
 * @return (*redis.Client, error) - Returns the Redis client instance and an error if any.
 */
func GetRedisClient() (*redis.Client, error) {
	if vRedisClient == nil {
		return nil, errors.New("redis client is not initialized, please call InitRedisClient first")
	}
	return vRedisClient, nil
}

// ===========================================
// 		helper functions for Redis client
// ===========================================
