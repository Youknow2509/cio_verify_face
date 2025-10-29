package ratelimit

import (
	"context"
	"hash/fnv"
	"strconv"
	"sync"

	domainCache "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/cache"
	domainModel "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/model"
	domainLimiter "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/ratelimit"
)

// ===============================================
// Memory cached limiter
// ===============================================
type MemoryLimiter struct {
	localCache domainCache.ILocalCache
	policy     domainLimiter.Policy
	prefix     string
	mutexes    [32]sync.Mutex
}

// Create implements ratelimit.ILimiter.
func (m *MemoryLimiter) Create(ctx context.Context, key string) error {
	return m.localCache.SetTTL(
		ctx,
		m.getKey(key),
		"1",
		int64(m.policy.Window.Seconds()),
	)
}

// Check implements ratelimit.ILimiter.
func (m *MemoryLimiter) Check(ctx context.Context, key string) (int, bool, error) {
	val, err := m.localCache.Get(ctx, m.getKey(key))
	if err != nil {
		return 0, false, err
	}
	if val == "" {
		return 0, true, nil
	}
	valInt, err := strconv.Atoi(val)
	if err != nil {
		return 0, false, err
	}
	if valInt+1 >= m.policy.Limit {
		return valInt, false, nil
	}
	return valInt, true, nil
}

// Upgrade implements ratelimit.ILimiter.
func (m *MemoryLimiter) Upgrade(ctx context.Context, key string, val int) error {
	err := m.localCache.SetTTL(
		ctx,
		m.getKey(key),
		strconv.Itoa(val),
		int64(m.policy.Window.Seconds()),
	)
	return err
}

// Allow implements ratelimit.ILimiter.
func (m *MemoryLimiter) Allow(ctx context.Context, key string) (domainModel.Verdict, error) {
	mutexIndex := m.getMutexIndex(m.getKey(key))
	m.mutexes[mutexIndex].Lock()
	defer m.mutexes[mutexIndex].Unlock()

	v, ok, err := m.Check(ctx, key)
	if err != nil {
		return domainModel.Denied, err
	}
	if !ok {
		return domainModel.Denied, err
	}
	if v == 0 {
		err := m.Create(ctx, key)
		if err != nil {
			return domainModel.Denied, err
		}
		return domainModel.Allowed, nil
	}
	err = m.localCache.IncreHaveTTL(ctx, m.getKey(key), 1)
	if err != nil {
		return domainModel.Denied, err
	}
	return domainModel.Allowed, nil
}

func (m *MemoryLimiter) getKey(key string) string {
	return m.prefix + key
}

// getMutexIndex returns the shard index for a given key
func (m *MemoryLimiter) getMutexIndex(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32()) % len(m.mutexes)
}

/**
 * New and impl memory limiter
 */
func NewMemoryLimiter(
	localCache domainCache.ILocalCache,
	policy domainLimiter.Policy,
	prefix string,
) domainLimiter.ILimiter {
	return &MemoryLimiter{
		localCache: localCache,
		policy:     policy,
		prefix:     prefix,
	}
}
