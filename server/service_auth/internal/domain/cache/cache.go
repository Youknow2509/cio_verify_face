package cache

import (
	"context"
	"errors"
	"time"
)

type (
	// Interfaces distributed cache
	IDistributedCache interface {
		// For get
		Get(ctx context.Context, key string) (string, error)
		// For set
		Set(ctx context.Context, key string, value interface{}) error
		SetTTL(ctx context.Context, key string, value interface{}, ttl int64) error
		TTL(ctx context.Context, key string) (time.Duration, error)
		// For delete
		Delete(ctx context.Context, key string) error
		// For checking existence
		Exists(ctx context.Context, key string) (bool, error)
		// For hash
		HSet(ctx context.Context, key string, field string, value interface{}) error
		HGet(ctx context.Context, key string, field string) (string, error)
		HDel(ctx context.Context, key string, field string) error
		HDelAll(ctx context.Context, key string) error
		HSETExpire(ctx context.Context, key string, field string, value interface{}, ttl int64) error
		HGetAll(ctx context.Context, key string) (map[string]interface{}, error)
		// For list
		LPush(ctx context.Context, key string, value interface{}) error
		LPop(ctx context.Context, key string) (string, error)
		LLength(ctx context.Context, key string) (int64, error)
		LMove(ctx context.Context, sourceKey string, destinationKey string, value interface{}) error
		LRange(ctx context.Context, key string, start int64, stop int64) ([]interface{}, error)
		LTrim(ctx context.Context, key string, start int64, stop int64) error
		// For set
		SAdd(ctx context.Context, key string, value interface{}) error
		SRem(ctx context.Context, key string, value interface{}) error
		SIsMember(ctx context.Context, key string, value interface{}) (bool, error)
		SInter(ctx context.Context, keys ...string) ([]interface{}, error)
		// Publish and subscribe
		Publish(ctx context.Context, channel string, message interface{}) error
		Subscribe(ctx context.Context, channel string) (<-chan interface{}, error)
		// Option
		Increment(ctx context.Context, key string, delta int64) (int64, error)
		Decrement(ctx context.Context, key string, delta int64) (int64, error)
		// Lua script execution
		LuaScript(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error)
		// v.v
	}

	// Interfaces local cache in memory
	ILocalCache interface {
		Get(ctx context.Context, key string) (string, error)
		Set(ctx context.Context, key string, value string) error
		SetTTL(ctx context.Context, key string, value string, ttl int64) error
		Delete(ctx context.Context, key string) error
		Exists(ctx context.Context, key string) (bool, error)
		//
		DecreHaveTTL(ctx context.Context, key string, val int) error
		Decre(ctx context.Context, key string, val int) error
		IncreHaveTTL(ctx context.Context, key string, val int) error
		Incre(ctx context.Context, key string, val int) error
		//
		All(ctx context.Context) (map[string]string, error)
		Keys(ctx context.Context) ([]string, error)
		Values(ctx context.Context) ([]string, error)
		//
		ClearUp(ctx context.Context) error
	}
)

/**
 * Variable for cache
 */
var (
	vIDistributedCache IDistributedCache
	vILocalCache       ILocalCache
)

/**
 * Set distributed cache
 */
func SetDistributedCache(cache IDistributedCache) error {
	if vIDistributedCache != nil {
		return errors.New("distributed cache is already initialized")
	}
	vIDistributedCache = cache
	return nil
}

/**
 * Get distributed cache
 */
func GetDistributedCache() (IDistributedCache, error) {
	if vIDistributedCache == nil {
		return nil, errors.New("distributed cache is not initialized, please call SetDistributedCache first")
	}
	return vIDistributedCache, nil
}

/**
 * Set local cache
 */
func SetLocalCache(cache ILocalCache) error {
	if vILocalCache != nil {
		return errors.New("local cache is already initialized")
	}
	vILocalCache = cache
	return nil
}

/**
 * Get local cache
 */
func GetLocalCache() (ILocalCache, error) {
	if vILocalCache == nil {
		return nil, errors.New("local cache is not initialized, please call SetLocalCache first")
	}
	return vILocalCache, nil
}
