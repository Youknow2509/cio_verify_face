package start

import (
	"fmt"

	domainCache "github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/domain/cache"
	infraCache "github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/infrastructure/cache"
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/global"
)

func initLocalCache() error {
	localCache, err := infraCache.NewRistrettoLocalCache()
	if err != nil {
		return fmt.Errorf("failed to initialize local cache: %w", err)
	}

	if err := domainCache.SetLocalCache(localCache); err != nil {
		return fmt.Errorf("failed to set local cache: %w", err)
	}

	global.Logger.Info("Local cache initialized successfully")
	return nil
}
