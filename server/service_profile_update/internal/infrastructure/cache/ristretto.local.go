package cache

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/dgraph-io/ristretto/v2"
	domainCache "github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/domain/cache"
)

// RistrettoLocalCache implements ILocalCache using Ristretto
type RistrettoLocalCache struct {
	cache *ristretto.Cache[string, string]
}

// NewRistrettoLocalCache creates a new RistrettoLocalCache instance
func NewRistrettoLocalCache() (domainCache.ILocalCache, error) {
	cache, err := ristretto.NewCache(&ristretto.Config[string, string]{
		NumCounters: 1e7,     // number of keys to track frequency of (10M)
		MaxCost:     1 << 30, // maximum cost of cache (1GB)
		BufferItems: 64,      // number of keys per Get buffer
	})
	if err != nil {
		return nil, err
	}
	return &RistrettoLocalCache{cache: cache}, nil
}

func (r *RistrettoLocalCache) Get(ctx context.Context, key string) (string, error) {
	value, found := r.cache.Get(key)
	if !found {
		return "", nil
	}
	return value, nil
}

func (r *RistrettoLocalCache) Set(ctx context.Context, key string, value string) error {
	r.cache.Set(key, value, 1)
	r.cache.Wait()
	return nil
}

func (r *RistrettoLocalCache) SetTTL(ctx context.Context, key string, value string, ttl int64) error {
	r.cache.SetWithTTL(key, value, 1, time.Duration(ttl)*time.Second)
	r.cache.Wait()
	return nil
}

func (r *RistrettoLocalCache) Delete(ctx context.Context, key string) error {
	r.cache.Del(key)
	return nil
}

func (r *RistrettoLocalCache) Exists(ctx context.Context, key string) (bool, error) {
	_, found := r.cache.Get(key)
	return found, nil
}

func (r *RistrettoLocalCache) DecreHaveTTL(ctx context.Context, key string, val int) error {
	v, ok := r.cache.Get(key)
	if !ok {
		return fmt.Errorf("key not found")
	}
	vInt, err := strconv.Atoi(v)
	if err != nil {
		return fmt.Errorf("value is not an integer")
	}
	newVal := vInt - val
	ttl, ok := r.cache.GetTTL(key)
	if !ok {
		return fmt.Errorf("key not found")
	}
	r.cache.SetWithTTL(key, strconv.Itoa(newVal), 1, ttl)
	r.cache.Wait()
	return nil
}

func (r *RistrettoLocalCache) Decre(ctx context.Context, key string, val int) error {
	v, ok := r.cache.Get(key)
	if !ok {
		return fmt.Errorf("key not found")
	}
	vInt, err := strconv.Atoi(v)
	if err != nil {
		return fmt.Errorf("value is not an integer")
	}
	newVal := vInt - val
	r.cache.Set(key, strconv.Itoa(newVal), 1)
	r.cache.Wait()
	return nil
}

func (r *RistrettoLocalCache) IncreHaveTTL(ctx context.Context, key string, val int) error {
	v, ok := r.cache.Get(key)
	if !ok {
		return fmt.Errorf("key not found")
	}
	vInt, err := strconv.Atoi(v)
	if err != nil {
		return fmt.Errorf("value is not an integer")
	}
	newVal := vInt + val
	ttl, ok := r.cache.GetTTL(key)
	if !ok {
		return fmt.Errorf("key not found")
	}
	r.cache.SetWithTTL(key, strconv.Itoa(newVal), 1, ttl)
	r.cache.Wait()
	return nil
}

func (r *RistrettoLocalCache) Incre(ctx context.Context, key string, val int) error {
	v, ok := r.cache.Get(key)
	if !ok {
		return fmt.Errorf("key not found")
	}
	vInt, err := strconv.Atoi(v)
	if err != nil {
		return fmt.Errorf("value is not an integer")
	}
	newVal := vInt + val
	r.cache.Set(key, strconv.Itoa(newVal), 1)
	r.cache.Wait()
	return nil
}

func (r *RistrettoLocalCache) ClearUp(ctx context.Context) error {
	r.cache.Clear()
	return nil
}
