package cache

import (
	"context"
	"errors"
	"time"
)

// =================================
// Distributed Cache Interface (Redis):
// =================================
type IDistributedCache interface {
	// Basic operations
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}) error
	SetTTL(ctx context.Context, key string, value interface{}, ttl int64) error
	TTL(ctx context.Context, key string) (time.Duration, error)
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)

	// Atomic increment/decrement
	Increment(ctx context.Context, key string, delta int64) (int64, error)
	Decrement(ctx context.Context, key string, delta int64) (int64, error)

	// Hash operations
	HSet(ctx context.Context, key string, field string, value interface{}) error
	HGet(ctx context.Context, key string, field string) (string, error)
	HDel(ctx context.Context, key string, field string) error
	HGetAll(ctx context.Context, key string) (map[string]interface{}, error)
	HSETExpire(ctx context.Context, key string, field string, value interface{}, ttl int64) error

	// Set operations
	SAdd(ctx context.Context, key string, value interface{}) error
	SRem(ctx context.Context, key string, value interface{}) error
	SIsMember(ctx context.Context, key string, value interface{}) (bool, error)

	// Lua script for atomic operations
	LuaScript(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error)
}

// =================================
// Local Cache Interface (In-memory):
// =================================
type ILocalCache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string) error
	SetTTL(ctx context.Context, key string, value string, ttl int64) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)

	// Atomic operations
	DecreHaveTTL(ctx context.Context, key string, val int) error
	Decre(ctx context.Context, key string, val int) error
	IncreHaveTTL(ctx context.Context, key string, val int) error
	Incre(ctx context.Context, key string, val int) error

	// Utility
	ClearUp(ctx context.Context) error
}

// =================================
// Cache Variables:
// =================================
var (
	_distributedCache IDistributedCache
	_localCache       ILocalCache
)

// =================================
// Setters and Getters:
// =================================
func SetDistributedCache(cache IDistributedCache) error {
	if _distributedCache != nil {
		return errors.New("distributed cache already initialized")
	}
	_distributedCache = cache
	return nil
}

func GetDistributedCache() (IDistributedCache, error) {
	if _distributedCache == nil {
		return nil, errors.New("distributed cache not initialized")
	}
	return _distributedCache, nil
}

func SetLocalCache(cache ILocalCache) error {
	if _localCache != nil {
		return errors.New("local cache already initialized")
	}
	_localCache = cache
	return nil
}

func GetLocalCache() (ILocalCache, error) {
	if _localCache == nil {
		return nil, errors.New("local cache not initialized")
	}
	return _localCache, nil
}
