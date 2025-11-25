package cache

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/constants"
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/domain/cache"
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/domain/config"
)

// RedisDistributedCache implements IDistributedCache using Redis
type RedisDistributedCache struct {
	client redis.UniversalClient
}

// NewRedisDistributedCache creates a new RedisDistributedCache instance
func NewRedisDistributedCache(cfg *config.RedisSetting) (cache.IDistributedCache, error) {
	var client redis.UniversalClient

	// Configure TLS if enabled
	var tlsConfig *tls.Config
	if cfg.UseTLS {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: true, // For development; in production, load proper certs
		}
	}

	switch cfg.Type {
	case constants.RedisTypeStandalone:
		client = redis.NewClient(&redis.Options{
			Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			Password:     cfg.Password,
			DB:           cfg.DB,
			PoolSize:     cfg.PoolSize,
			MinIdleConns: cfg.MinIdleConns,
			MaxRetries:   cfg.MaxRetries,
			TLSConfig:    tlsConfig,
		})

	case constants.RedisTypeSentinel:
		client = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:       cfg.MasterName,
			SentinelAddrs:    cfg.SentinelAddrs,
			Password:         cfg.Password,
			DB:               cfg.DB,
			PoolSize:         cfg.PoolSize,
			MinIdleConns:     cfg.MinIdleConns,
			MaxRetries:       cfg.MaxRetries,
			TLSConfig:        tlsConfig,
		})

	case constants.RedisTypeCluster:
		client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:          cfg.Address,
			Password:       cfg.Password,
			PoolSize:       cfg.PoolSize,
			MinIdleConns:   cfg.MinIdleConns,
			MaxRetries:     cfg.MaxRetries,
			RouteByLatency: cfg.RouteByLatency,
			RouteRandomly:  cfg.RouteRandomly,
			TLSConfig:      tlsConfig,
		})

	default:
		// Default to standalone
		client = redis.NewClient(&redis.Options{
			Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			Password:     cfg.Password,
			DB:           cfg.DB,
			PoolSize:     cfg.PoolSize,
			MinIdleConns: cfg.MinIdleConns,
			MaxRetries:   cfg.MaxRetries,
			TLSConfig:    tlsConfig,
		})
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisDistributedCache{client: client}, nil
}

func (r *RedisDistributedCache) Get(ctx context.Context, key string) (string, error) {
	result, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("failed to get key %s: %w", key, err)
	}
	return result, nil
}

func (r *RedisDistributedCache) Set(ctx context.Context, key string, value interface{}) error {
	valueStr, err := toString(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, valueStr, 0).Err()
}

func (r *RedisDistributedCache) SetTTL(ctx context.Context, key string, value interface{}, ttl int64) error {
	valueStr, err := toString(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, valueStr, time.Duration(ttl)*time.Second).Err()
}

func (r *RedisDistributedCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	return r.client.TTL(ctx, key).Result()
}

func (r *RedisDistributedCache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisDistributedCache) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

func (r *RedisDistributedCache) Increment(ctx context.Context, key string, delta int64) (int64, error) {
	return r.client.IncrBy(ctx, key, delta).Result()
}

func (r *RedisDistributedCache) Decrement(ctx context.Context, key string, delta int64) (int64, error) {
	return r.client.DecrBy(ctx, key, delta).Result()
}

func (r *RedisDistributedCache) HSet(ctx context.Context, key string, field string, value interface{}) error {
	valueStr, err := toString(value)
	if err != nil {
		return err
	}
	return r.client.HSet(ctx, key, field, valueStr).Err()
}

func (r *RedisDistributedCache) HGet(ctx context.Context, key string, field string) (string, error) {
	result, err := r.client.HGet(ctx, key, field).Result()
	if err == redis.Nil {
		return "", nil
	}
	return result, err
}

func (r *RedisDistributedCache) HDel(ctx context.Context, key string, field string) error {
	return r.client.HDel(ctx, key, field).Err()
}

func (r *RedisDistributedCache) HGetAll(ctx context.Context, key string) (map[string]interface{}, error) {
	result, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	converted := make(map[string]interface{})
	for field, value := range result {
		var jsonResult interface{}
		if err := json.Unmarshal([]byte(value), &jsonResult); err == nil {
			converted[field] = jsonResult
		} else {
			converted[field] = value
		}
	}
	return converted, nil
}

func (r *RedisDistributedCache) HSETExpire(ctx context.Context, key string, field string, value interface{}, ttl int64) error {
	if err := r.HSet(ctx, key, field, value); err != nil {
		return err
	}
	return r.client.Expire(ctx, key, time.Duration(ttl)*time.Second).Err()
}

func (r *RedisDistributedCache) SAdd(ctx context.Context, key string, value interface{}) error {
	valueStr, err := toString(value)
	if err != nil {
		return err
	}
	return r.client.SAdd(ctx, key, valueStr).Err()
}

func (r *RedisDistributedCache) SRem(ctx context.Context, key string, value interface{}) error {
	valueStr, err := toString(value)
	if err != nil {
		return err
	}
	return r.client.SRem(ctx, key, valueStr).Err()
}

func (r *RedisDistributedCache) SIsMember(ctx context.Context, key string, value interface{}) (bool, error) {
	valueStr, err := toString(value)
	if err != nil {
		return false, err
	}
	return r.client.SIsMember(ctx, key, valueStr).Result()
}

func (r *RedisDistributedCache) LuaScript(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error) {
	return r.client.Eval(ctx, script, keys, args...).Result()
}

// Helper function to convert interface to string
func toString(value interface{}) (string, error) {
	if str, ok := value.(string); ok {
		return str, nil
	}
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return "", fmt.Errorf("failed to marshal value: %w", err)
	}
	return string(jsonBytes), nil
}
