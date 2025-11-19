package cache

import (
	"errors"

	"github.com/dgraph-io/ristretto"
)

// ristretto local cache variables
var (
	vLocalCache *ristretto.Cache
)

// InitLocalCache initializes the ristretto local cache
func InitLocalCache() error {
	if vLocalCache != nil {
		return errors.New("local cache is already initialized")
	}

	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M)
		MaxCost:     1 << 30, // maximum cost of cache (1GB)
		BufferItems: 64,      // number of keys per Get buffer
	})
	if err != nil {
		return errors.New("failed to create local cache: " + err.Error())
	}

	vLocalCache = cache
	return nil
}

// GetLocalCache returns the local cache instance
func GetLocalCache() (*ristretto.Cache, error) {
	if vLocalCache == nil {
		return nil, errors.New("local cache is not initialized, please call InitLocalCache first")
	}
	return vLocalCache, nil
}

// CloseLocalCache closes the local cache
func CloseLocalCache() {
	if vLocalCache != nil {
		vLocalCache.Close()
		vLocalCache = nil
	}
}
