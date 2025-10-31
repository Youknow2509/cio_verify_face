package start

import (
	"errors"

	libsInfraRateLimit "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/infrastructure/ratelimit"
	libsDomainCache "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/cache"
	libsConfig "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/config"
	libsDomainRateLimit "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/ratelimit"

	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/constants"
	global "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/global"
)

// init rate limit
func initRateLimit(setting []libsConfig.RateLimitPolicy) error {
	// convert config data
	var policies []libsDomainRateLimit.PolicyInput
	for _, item := range setting {
		policies = append(policies, libsDomainRateLimit.PolicyInput{
			Name:   item.Name,
			Limit:  item.Limit,
			Window: item.Window,
		})
	}
	// init manager rate limit
	policyManager, err := libsDomainRateLimit.NewPolicyManager(policies)
	if err != nil {
		return err
	}
	global.RateLimitPolicyManager = policyManager
	// init ws read limiter
	wsReadPolicy, ok := policyManager.GetPolicy(constants.POLICY_RATE_LIMIT_WS_READ)
	if !ok {
		return errors.New("ws_read policy not found")
	}
	localCacheInstance, err := libsDomainCache.GetLocalCache()
	if err != nil {
		return err
	}
	wsReadLimiter := libsInfraRateLimit.NewMemoryLimiter(
		localCacheInstance,
		wsReadPolicy,
		constants.WS_PREFIX_RATE_LIMIT_READ,
	)
	global.RateLimitWsRead = wsReadLimiter
	return nil
}
