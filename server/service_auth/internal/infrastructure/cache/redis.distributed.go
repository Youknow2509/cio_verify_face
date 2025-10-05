package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	infraConn "github.com/youknow2509/cio_verify_face/server/service_auth/internal/infrastructure/conn"
	"github.com/youknow2509/cio_verify_face/server/service_auth/internal/domain/config"
	domainCache "github.com/youknow2509/cio_verify_face/server/service_auth/internal/domain/cache"
)

/**
 * Redis service interface implementation distributed cache
 */
type RedisDistributedCache struct {
	client *redis.Client
}

// TTL implements cache.IDistributedCache.
func (r *RedisDistributedCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	result, err := r.client.TTL(ctx, key).Result()
	if handleErrorRedis(err) != nil {
		return 0, fmt.Errorf("failed to get TTL for key %s: %w", key, err)
	}
	return result, nil
}

// Decrement implements cache.IDistributedCache.
func (r *RedisDistributedCache) Decrement(ctx context.Context, key string, delta int64) (int64, error) {
	result, err := r.client.DecrBy(ctx, key, delta).Result()
	if handleErrorRedis(err) != nil {
		return 0, fmt.Errorf("failed to decrement key %s: %w", key, err)
	}
	return result, nil
}

// Delete implements cache.IDistributedCache.
func (r *RedisDistributedCache) Delete(ctx context.Context, key string) error {
	err := r.client.Del(ctx, key).Err()
	if handleErrorRedis(err) != nil {
		return fmt.Errorf("failed to delete key %s: %w", key, err)
	}
	return nil
}

// Exists implements cache.IDistributedCache.
func (r *RedisDistributedCache) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, key).Result()
	if handleErrorRedis(err) != nil {
		return false, fmt.Errorf("failed to check existence of key %s: %w", key, err)
	}
	return result > 0, nil
}

// Get implements cache.IDistributedCache.
func (r *RedisDistributedCache) Get(ctx context.Context, key string) (string, error) {
	result, err := r.client.Get(ctx, key).Result()
	if handleErrorRedis(err) != nil {
		return "", fmt.Errorf("failed to get key %s: %w", key, err)
	}

	return result, nil
}

// HDel implements cache.IDistributedCache.
func (r *RedisDistributedCache) HDel(ctx context.Context, key string, field string) error {
	err := r.client.HDel(ctx, key, field).Err()
	if handleErrorRedis(err) != nil {
		return fmt.Errorf("failed to delete hash field %s from key %s: %w", field, key, err)
	}
	return nil
}

// HDelAll implements cache.IDistributedCache.
func (r *RedisDistributedCache) HDelAll(ctx context.Context, key string) error {
	err := r.client.Del(ctx, key).Err()
	if handleErrorRedis(err) != nil {
		return fmt.Errorf("failed to delete all hash fields from key %s: %w", key, err)
	}
	return nil
}

// HGet implements cache.IDistributedCache.
func (r *RedisDistributedCache) HGet(ctx context.Context, key string, field string) (string, error) {
	result, err := r.client.HGet(ctx, key, field).Result()
	if handleErrorRedis(err) != nil {
		return "", fmt.Errorf("failed to get hash field %s from key %s: %w", field, key, err)
	}

	return result, nil
}

// HGetAll implements cache.IDistributedCache.
func (r *RedisDistributedCache) HGetAll(ctx context.Context, key string) (map[string]interface{}, error) {
	result, err := r.client.HGetAll(ctx, key).Result()
	if handleErrorRedis(err) != nil {
		return nil, fmt.Errorf("failed to get all hash fields from key %s: %w", key, err)
	}

	// Convert map[string]string to map[string]interface{}
	converted := make(map[string]interface{})
	for field, value := range result {
		// Try to parse as JSON first, if it fails keep as string
		var jsonResult interface{}
		if err := json.Unmarshal([]byte(value), &jsonResult); err == nil {
			converted[field] = jsonResult
		} else {
			converted[field] = value
		}
	}
	return converted, nil
}

// HSETExpire implements cache.IDistributedCache.
func (r *RedisDistributedCache) HSETExpire(ctx context.Context, key string, field string, value interface{}, ttl int64) error {
	// First set the hash field
	if err := r.HSet(ctx, key, field, value); err != nil {
		return err
	}

	// Then set the TTL for the entire hash key
	err := r.client.Expire(ctx, key, time.Duration(ttl)*time.Second).Err()
	if handleErrorRedis(err) != nil {
		return fmt.Errorf("failed to set TTL for key %s: %w", key, err)
	}
	return nil
}

// HSet implements cache.IDistributedCache.
func (r *RedisDistributedCache) HSet(ctx context.Context, key string, field string, value interface{}) error {
	// Convert value to JSON string if it's not already a string
	var valueStr string
	if str, ok := value.(string); ok {
		valueStr = str
	} else {
		jsonBytes, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value for hash field %s in key %s: %w", field, key, err)
		}
		valueStr = string(jsonBytes)
	}

	err := r.client.HSet(ctx, key, field, valueStr).Err()
	if handleErrorRedis(err) != nil {
		return fmt.Errorf("failed to set hash field %s in key %s: %w", field, key, err)
	}
	return nil
}

// Increment implements cache.IDistributedCache.
func (r *RedisDistributedCache) Increment(ctx context.Context, key string, delta int64) (int64, error) {
	result, err := r.client.IncrBy(ctx, key, delta).Result()
	if handleErrorRedis(err) != nil {
		return 0, fmt.Errorf("failed to increment key %s: %w", key, err)
	}
	return result, nil
}

// LLength implements cache.IDistributedCache.
func (r *RedisDistributedCache) LLength(ctx context.Context, key string) (int64, error) {
	result, err := r.client.LLen(ctx, key).Result()
	if handleErrorRedis(err) != nil {
		return 0, fmt.Errorf("failed to get list length for key %s: %w", key, err)
	}
	return result, nil
}

// LMove implements cache.IDistributedCache.
func (r *RedisDistributedCache) LMove(ctx context.Context, sourceKey string, destinationKey string, value interface{}) error {
	// Convert value to JSON string if it's not already a string
	var valueStr string
	if str, ok := value.(string); ok {
		valueStr = str
	} else {
		jsonBytes, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value for list move: %w", err)
		}
		valueStr = string(jsonBytes)
	}

	// Remove from source list
	err := r.client.LRem(ctx, sourceKey, 1, valueStr).Err()
	if handleErrorRedis(err) != nil {
		return fmt.Errorf("failed to remove value from source list %s: %w", sourceKey, err)
	}

	// Add to destination list
	err = r.client.LPush(ctx, destinationKey, valueStr).Err()
	if handleErrorRedis(err) != nil {
		return fmt.Errorf("failed to add value to destination list %s: %w", destinationKey, err)
	}

	return nil
}

// LPop implements cache.IDistributedCache.
func (r *RedisDistributedCache) LPop(ctx context.Context, key string) (string, error) {
	result, err := r.client.LPop(ctx, key).Result()
	if handleErrorRedis(err) != nil {
		return "", fmt.Errorf("failed to pop from list %s: %w", key, err)
	}

	return result, nil
}

// LPush implements cache.IDistributedCache.
func (r *RedisDistributedCache) LPush(ctx context.Context, key string, value interface{}) error {
	// Convert value to JSON string if it's not already a string
	var valueStr string
	if str, ok := value.(string); ok {
		valueStr = str
	} else {
		jsonBytes, err := json.Marshal(value)
		if (err) != nil {
			return fmt.Errorf("failed to marshal value for list push: %w", err)
		}
		valueStr = string(jsonBytes)
	}

	err := r.client.LPush(ctx, key, valueStr).Err()
	if handleErrorRedis(err) != nil {
		return fmt.Errorf("failed to push to list %s: %w", key, err)
	}
	return nil
}

// LRange implements cache.IDistributedCache.
func (r *RedisDistributedCache) LRange(ctx context.Context, key string, start int64, stop int64) ([]interface{}, error) {
	result, err := r.client.LRange(ctx, key, start, stop).Result()
	if handleErrorRedis(err) != nil {
		return nil, fmt.Errorf("failed to get range from list %s: %w", key, err)
	}

	// Convert []string to []interface{}
	converted := make([]interface{}, len(result))
	for i, value := range result {
		// Try to parse as JSON first, if it fails keep as string
		var jsonResult interface{}
		if err := json.Unmarshal([]byte(value), &jsonResult); err == nil {
			converted[i] = jsonResult
		} else {
			converted[i] = value
		}
	}
	return converted, nil
}

// LTrim implements cache.IDistributedCache.
func (r *RedisDistributedCache) LTrim(ctx context.Context, key string, start int64, stop int64) error {
	err := r.client.LTrim(ctx, key, start, stop).Err()
	if handleErrorRedis(err) != nil {
		return fmt.Errorf("failed to trim list %s: %w", key, err)
	}
	return nil
}

// LuaScript implements cache.IDistributedCache.
func (r *RedisDistributedCache) LuaScript(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error) {
	result, err := r.client.Eval(ctx, script, keys, args...).Result()
	if handleErrorRedis(err) != nil {
		return "", fmt.Errorf("failed to execute lua script: %w", err)
	}
	return result, nil
}

// Publish implements cache.IDistributedCache.
func (r *RedisDistributedCache) Publish(ctx context.Context, channel string, message interface{}) error {
	// Convert message to JSON string if it's not already a string
	var messageStr string
	if str, ok := message.(string); ok {
		messageStr = str
	} else {
		jsonBytes, err := json.Marshal(message)
		if (err) != nil {
			return fmt.Errorf("failed to marshal message for publish: %w", err)
		}
		messageStr = string(jsonBytes)
	}

	err := r.client.Publish(ctx, channel, messageStr).Err()
	if handleErrorRedis(err) != nil {
		return fmt.Errorf("failed to publish to channel %s: %w", channel, err)
	}
	return nil
}

// SAdd implements cache.IDistributedCache.
func (r *RedisDistributedCache) SAdd(ctx context.Context, key string, value interface{}) error {
	// Convert value to JSON string if it's not already a string
	var valueStr string
	if str, ok := value.(string); ok {
		valueStr = str
	} else {
		jsonBytes, err := json.Marshal(value)
		if (err) != nil {
			return fmt.Errorf("failed to marshal value for set add: %w", err)
		}
		valueStr = string(jsonBytes)
	}

	err := r.client.SAdd(ctx, key, valueStr).Err()
	if handleErrorRedis(err) != nil {
		return fmt.Errorf("failed to add to set %s: %w", key, err)
	}
	return nil
}

// SInter implements cache.IDistributedCache.
func (r *RedisDistributedCache) SInter(ctx context.Context, keys ...string) ([]interface{}, error) {
	result, err := r.client.SInter(ctx, keys...).Result()
	if handleErrorRedis(err) != nil {
		return nil, fmt.Errorf("failed to get intersection of sets: %w", err)
	}

	// Convert []string to []interface{}
	converted := make([]interface{}, len(result))
	for i, value := range result {
		// Try to parse as JSON first, if it fails keep as string
		var jsonResult interface{}
		if err := json.Unmarshal([]byte(value), &jsonResult); err == nil {
			converted[i] = jsonResult
		} else {
			converted[i] = value
		}
	}
	return converted, nil
}

// SIsMember implements cache.IDistributedCache.
func (r *RedisDistributedCache) SIsMember(ctx context.Context, key string, value interface{}) (bool, error) {
	// Convert value to JSON string if it's not already a string
	var valueStr string
	if str, ok := value.(string); ok {
		valueStr = str
	} else {
		jsonBytes, err := json.Marshal(value)
		if (err) != nil {
			return false, fmt.Errorf("failed to marshal value for set membership check: %w", err)
		}
		valueStr = string(jsonBytes)
	}

	result, err := r.client.SIsMember(ctx, key, valueStr).Result()
	if handleErrorRedis(err) != nil {
		return false, fmt.Errorf("failed to check membership in set %s: %w", key, err)
	}
	return result, nil
}

// SRem implements cache.IDistributedCache.
func (r *RedisDistributedCache) SRem(ctx context.Context, key string, value interface{}) error {
	// Convert value to JSON string if it's not already a string
	var valueStr string
	if str, ok := value.(string); ok {
		valueStr = str
	} else {
		jsonBytes, err := json.Marshal(value)
		if (err) != nil {
			return fmt.Errorf("failed to marshal value for set remove: %w", err)
		}
		valueStr = string(jsonBytes)
	}

	err := r.client.SRem(ctx, key, valueStr).Err()
	if handleErrorRedis(err) != nil {
		return fmt.Errorf("failed to remove from set %s: %w", key, err)
	}
	return nil
}

// Set implements cache.IDistributedCache.
func (r *RedisDistributedCache) Set(ctx context.Context, key string, value interface{}) error {
	// Convert value to JSON string if it's not already a string
	var valueStr string
	if str, ok := value.(string); ok {
		valueStr = str
	} else {
		jsonBytes, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value for set: %w", err)
		}
		valueStr = string(jsonBytes)
	}

	err := r.client.Set(ctx, key, valueStr, 0).Err()
	if handleErrorRedis(err) != nil {
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}
	return nil
}

// SetTTL implements cache.IDistributedCache.
func (r *RedisDistributedCache) SetTTL(ctx context.Context, key string, value interface{}, ttl int64) error {
	// Convert value to JSON string if it's not already a string
	var valueStr string
	if str, ok := value.(string); ok {
		valueStr = str
	} else {
		jsonBytes, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value for set with TTL: %w", err)
		}
		valueStr = string(jsonBytes)
	}

	err := r.client.Set(ctx, key, valueStr, time.Duration(ttl)*time.Second).Err()
	if handleErrorRedis(err) != nil {
		return fmt.Errorf("failed to set key %s with TTL: %w", key, err)
	}
	return nil
}

// Subscribe implements cache.IDistributedCache.
func (r *RedisDistributedCache) Subscribe(ctx context.Context, channel string) (<-chan interface{}, error) {
	pubsub := r.client.Subscribe(ctx, channel)

	// Create a channel to return messages
	messages := make(chan interface{}, 100) // Buffer to prevent blocking

	// Start a goroutine to listen for messages
	go func() {
		defer close(messages)
		defer pubsub.Close()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				msg, err := pubsub.ReceiveMessage(ctx)
				if err != nil {
					// If context is cancelled, return gracefully
					if ctx.Err() != nil {
						return
					}
					// For other errors, you might want to log them
					continue
				}

				// Try to parse as JSON first, if it fails send as string
				var jsonResult interface{}
				if err := json.Unmarshal([]byte(msg.Payload), &jsonResult); err == nil {
					select {
					case messages <- jsonResult:
					case <-ctx.Done():
						return
					}
				} else {
					select {
					case messages <- msg.Payload:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()

	return messages, nil
}

/**
 * NewRedisDistributedCache creates a new RedisDistributedCache instance
 */
func NewRedisDistributedCache(config *config.RedisSetting) (domainCache.IDistributedCache, error) {
	er := infraConn.InitRedisClient(config)
	if er != nil {
		return nil, er
	}
	client, er := infraConn.GetRedisClient()
	if er != nil {
		return nil, er
	}
	return &RedisDistributedCache{
		client: client,
	}, nil
}

// ===========================================================================
//
//	Redis Helper Functions
//
// ===========================================================================
func handleErrorRedis(err error) error {
	if err == redis.Nil {
		return nil // Key does not exist
	}
	return err
}
