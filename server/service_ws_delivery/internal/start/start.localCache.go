package start

import (
	domainCache "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/cache"
	infraCache "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/infrastructure/cache"
)

// initLocalCache initializes the local cache for the service.
func initLocalCache() error {
	localCacheImpl, err := infraCache.NewRistrettoLocalCache()
	if err != nil {
		return err
	}
	if err := domainCache.SetLocalCache(localCacheImpl); err != nil {
		return err
	}
	return nil
}
