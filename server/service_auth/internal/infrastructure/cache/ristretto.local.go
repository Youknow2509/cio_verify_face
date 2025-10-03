package cache

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/dgraph-io/ristretto/v2"
	domainCache "github.com/youknow2509/cio_verify_face/server/service_auth/internal/domain/cache"
)

/**
 * Ristretto service interface implementation local cache
 */
type RistrettoLocalCache struct {
	cache *ristretto.Cache[string, string]
}

// All implements cache.ILocalCache.
func (o *RistrettoLocalCache) All(ctx context.Context) (map[string]string, error) {
	return nil, fmt.Errorf("ristretto does not support retrieving all items directly, use Keys and Get for individual items")
}

// ClearUp implements cache.ILocalCache.
func (o *RistrettoLocalCache) ClearUp(ctx context.Context) error {
	o.cache.Clear()
	return nil
}

// Delete implements cache.ILocalCache.
func (o *RistrettoLocalCache) Delete(ctx context.Context, key string) error {
	o.cache.Del(key)
	return nil
}

// Exists implements cache.ILocalCache.
func (o *RistrettoLocalCache) Exists(ctx context.Context, key string) (bool, error) {
	_, found := o.cache.Get(key)
	return found, nil
}

// Get implements cache.ILocalCache.
func (o *RistrettoLocalCache) Get(ctx context.Context, key string) (string, error) {
	value, found := o.cache.Get(key)
	if !found {
		return "", nil
	}
	return value, nil
}

// Keys implements cache.ILocalCache.
func (o *RistrettoLocalCache) Keys(ctx context.Context) ([]string, error) {
	return nil, fmt.Errorf("ristretto does not support retrieving all keys directly, use All for all items or Get for individual items")
}

// Set implements cache.ILocalCache.
func (o *RistrettoLocalCache) Set(ctx context.Context, key string, value string) error {
	o.cache.Set(key, value, 1)
	o.cache.Wait()
	return nil
}

// SetTTL implements cache.ILocalCache.
func (o *RistrettoLocalCache) SetTTL(ctx context.Context, key string, value string, ttl int64) error {
	// Cost = 1, TTL được convert từ seconds sang Duration
	o.cache.SetWithTTL(key, value, 1, time.Duration(ttl)*time.Second)

	// Đợi cho đến khi value được set
	o.cache.Wait()
	return nil
}

// Values implements cache.ILocalCache.
func (o *RistrettoLocalCache) Values(ctx context.Context) ([]string, error) {
	return nil, fmt.Errorf("ristretto does not support retrieving all values directly, use All for all items or Get for individual items")
}

// DecreHaveTTL implements cache.ILocalCache.
func (o *RistrettoLocalCache) DecreHaveTTL(ctx context.Context, key string, val int) error {
	v, ok := o.cache.Get(key)
	if !ok {
		return fmt.Errorf("key not found")
	}
	vInt, err := strconv.Atoi(v)
	if err != nil {
		return fmt.Errorf("value is not an integer")
	}
	newVal := vInt - val
	ttl, ok := o.cache.GetTTL(key)
	if !ok {
		return fmt.Errorf("key not found")
	}
	o.cache.SetWithTTL(key, strconv.Itoa(newVal), 1, ttl)
	o.cache.Wait()
	return nil
}

// IncreHaveTTL implements cache.ILocalCache.
func (o *RistrettoLocalCache) IncreHaveTTL(ctx context.Context, key string, val int) error {
	v, ok := o.cache.Get(key)
	if !ok {
		return fmt.Errorf("key not found")
	}
	vInt, err := strconv.Atoi(v)
	if err != nil {
		return fmt.Errorf("value is not an integer")
	}
	newVal := vInt + val
	ttl, ok := o.cache.GetTTL(key)
	if !ok {
		return fmt.Errorf("key not found")
	}
	o.cache.SetWithTTL(key, strconv.Itoa(newVal), 1, ttl)
	o.cache.Wait()
	return nil
}

// Decre implements cache.ILocalCache.
func (o *RistrettoLocalCache) Decre(ctx context.Context, key string, val int) error {
	v, ok := o.cache.Get(key)
	if !ok {
		return fmt.Errorf("key not found")
	}
	vInt, err := strconv.Atoi(v)
	if err != nil {
		return fmt.Errorf("value is not an integer")
	}
	newVal := vInt - val
	o.cache.Set(key, strconv.Itoa(newVal), 1)
	o.cache.Wait()
	return nil
}

// Incre implements cache.ILocalCache.
func (o *RistrettoLocalCache) Incre(ctx context.Context, key string, val int) error {
	v, ok := o.cache.Get(key)
	if !ok {
		return fmt.Errorf("key not found")
	}
	vInt, err := strconv.Atoi(v)
	if err != nil {
		return fmt.Errorf("value is not an integer")
	}
	newVal := vInt + val
	o.cache.Set(key, strconv.Itoa(newVal), 1)
	o.cache.Wait()
	return nil
}

// NewRistrettoLocalCache creates a new instance of RistrettoLocalCache
func NewRistrettoLocalCache() (domainCache.ILocalCache, error) {
	cache, err := ristretto.NewCache(&ristretto.Config[string, string]{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
		return nil, err
	}
	return &RistrettoLocalCache{
		cache: cache,
	}, nil
}
