package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/youknow2509/cio_verify_face/server/service_auth/internal/constants"
	domainCache "github.com/youknow2509/cio_verify_face/server/service_auth/internal/domain/cache"
	domainModel "github.com/youknow2509/cio_verify_face/server/service_auth/internal/domain/model"
	"github.com/youknow2509/cio_verify_face/server/service_auth/internal/global"
	utilsCache "github.com/youknow2509/cio_verify_face/server/service_auth/internal/shared/utils/cache"
	utilsCrypto "github.com/youknow2509/cio_verify_face/server/service_auth/internal/shared/utils/crypto"
)

// CacheStrategy handles multi-level caching (local + distributed)
type CacheStrategy struct {
	localCache       domainCache.ILocalCache
	distributedCache domainCache.IDistributedCache
}

// NewCacheStrategy creates a new cache strategy instance
func NewCacheStrategy() (*CacheStrategy, error) {
	localCache, err := domainCache.GetLocalCache()
	if err != nil {
		return nil, err
	}

	distributedCache, err := domainCache.GetDistributedCache()
	if err != nil {
		return nil, err
	}

	return &CacheStrategy{
		localCache:       localCache,
		distributedCache: distributedCache,
	}, nil
}

// GetWithFallback tries local cache first, then distributed cache, then executes fallback function
func (cs *CacheStrategy) GetWithFallback(ctx context.Context, key string, localTTL, distributedTTL int64, fallback func() (interface{}, error)) (interface{}, error) {
	// Try local cache first
	if val, err := cs.localCache.Get(ctx, key); err == nil && val != "" {
		global.Logger.Info("Cache hit: local cache", "key", key)
		return val, nil
	}

	// Try distributed cache
	if val, err := cs.distributedCache.Get(ctx, key); err == nil && val != "" {
		global.Logger.Info("Cache hit: distributed cache", "key", key)

		// Store in local cache for faster access next time
		if localTTL > 0 {
			if err := cs.localCache.SetTTL(ctx, key, val, localTTL); err != nil {
				global.Logger.Warn("Failed to set local cache", "key", key, "error", err)
			}
		}
		return val, nil
	}

	// Execute fallback function
	global.Logger.Info("Cache miss: executing fallback", "key", key)
	result, err := fallback()
	if err != nil {
		return nil, err
	}

	// Store in both caches
	cs.SetMultiLevel(ctx, key, result, localTTL, distributedTTL)

	return result, nil
}

// SetMultiLevel sets value in both local and distributed cache
func (cs *CacheStrategy) SetMultiLevel(ctx context.Context, key string, value interface{}, localTTL, distributedTTL int64) {
	// Convert value to string for storage
	var strValue string
	switch v := value.(type) {
	case string:
		strValue = v
	case map[string]string:
		if jsonData, err := json.Marshal(v); err == nil {
			strValue = string(jsonData)
		}
	case map[string]interface{}:
		if jsonData, err := json.Marshal(v); err == nil {
			strValue = string(jsonData)
		}
	default:
		if jsonData, err := json.Marshal(v); err == nil {
			strValue = string(jsonData)
		}
	}

	// Set in local cache
	if localTTL > 0 {
		if err := cs.localCache.SetTTL(ctx, key, strValue, localTTL); err != nil {
			global.Logger.Warn("Failed to set local cache", "key", key, "error", err)
		}
	}

	// Set in distributed cache
	if distributedTTL > 0 {
		if err := cs.distributedCache.SetTTL(ctx, key, value, distributedTTL); err != nil {
			global.Logger.Warn("Failed to set distributed cache", "key", key, "error", err)
		}
	}
}

// Delete removes value from both caches
func (cs *CacheStrategy) Delete(ctx context.Context, key string) {
	// Delete from local cache
	if err := cs.localCache.Delete(ctx, key); err != nil {
		global.Logger.Warn("Failed to delete from local cache", "key", key, "error", err)
	}

	// Delete from distributed cache
	if err := cs.distributedCache.Delete(ctx, key); err != nil {
		global.Logger.Warn("Failed to delete from distributed cache", "key", key, "error", err)
	}
}

// User-specific cache methods

// GetUserInfo caches user information
func (cs *CacheStrategy) GetUserInfo(ctx context.Context, userID string, fallback func() (*domainModel.UserInfoOutput, error)) (*domainModel.UserInfoOutput, error) {
	userIDHash := utilsCrypto.GetHash(userID)
	key := utilsCache.GetKeyCacheUserInfoView(userIDHash)

	result, err := cs.GetWithFallback(ctx, key,
		constants.TTL_LOCAL_USER_INFO_VIEW,
		constants.TTL_USER_INFO_VIEW,
		func() (interface{}, error) {
			return fallback()
		})

	if err != nil {
		return nil, err
	}

	// Parse result
	if str, ok := result.(string); ok {
		var userInfo domainModel.UserInfoOutput
		if err := json.Unmarshal([]byte(str), &userInfo); err == nil {
			return &userInfo, nil
		}
	}

	if userInfo, ok := result.(*domainModel.UserInfoOutput); ok {
		return userInfo, nil
	}

	return nil, fmt.Errorf("invalid user info format in cache")
}

// SetUserSession stores user session information in cache
func (cs *CacheStrategy) SetUserSession(ctx context.Context, sessionID string, userID string, role domainModel.Role, ttl int64) error {
	sessionHash := utilsCrypto.GetHash(sessionID)
	key := utilsCache.GetKeyUserAccessTokenIsActive(sessionHash)

	value := map[string]string{
		"user_id": userID,
		"role":    strconv.Itoa(int(role)),
	}

	cs.SetMultiLevel(ctx, key, value, constants.TTL_LOCAL_ACCESS_TOKEN, ttl)
	return nil
}

// GetUserSession retrieves user session from cache
func (cs *CacheStrategy) GetUserSession(ctx context.Context, sessionID string) (map[string]string, error) {
	sessionHash := utilsCrypto.GetHash(sessionID)
	key := utilsCache.GetKeyUserAccessTokenIsActive(sessionHash)

	result, err := cs.GetWithFallback(ctx, key,
		constants.TTL_LOCAL_ACCESS_TOKEN,
		0, // Don't set in distributed if not exists
		func() (interface{}, error) {
			return nil, fmt.Errorf("session not found")
		})

	if err != nil {
		return nil, err
	}

	if str, ok := result.(string); ok {
		var sessionData map[string]string
		if err := json.Unmarshal([]byte(str), &sessionData); err == nil {
			return sessionData, nil
		}
	}

	if sessionData, ok := result.(map[string]string); ok {
		return sessionData, nil
	}

	return nil, fmt.Errorf("invalid session format in cache")
}

// DeleteUserSession removes user session from cache
func (cs *CacheStrategy) DeleteUserSession(ctx context.Context, sessionID string) {
	sessionHash := utilsCrypto.GetHash(sessionID)
	key := utilsCache.GetKeyUserAccessTokenIsActive(sessionHash)
	cs.Delete(ctx, key)
}

// CheckAndIncrementSpamCounter checks spam counter and increments it
func (cs *CacheStrategy) CheckAndIncrementSpamCounter(ctx context.Context, keyFunc func(string) string, identifier string, maxAttempts int, blockTTL, countTTL int64) (bool, error) {
	identifierHash := utilsCrypto.GetHash(identifier)

	// Check if blocked
	blockKey := keyFunc(identifierHash) + ":block"
	if exists, _ := cs.distributedCache.Exists(ctx, blockKey); exists {
		return false, fmt.Errorf("temporarily blocked due to too many attempts")
	}

	// Increment counter
	countKey := keyFunc(identifierHash) + ":count"
	count, err := cs.distributedCache.Increment(ctx, countKey, 1)
	if err != nil {
		return false, err
	}

	// Set TTL for counter if it's the first increment
	if count == 1 {
		cs.distributedCache.SetTTL(ctx, countKey, count, countTTL)
	}

	// Check if exceeded max attempts
	if int(count) >= maxAttempts {
		// Block for specified duration
		cs.distributedCache.SetTTL(ctx, blockKey, "blocked", blockTTL)
		return false, fmt.Errorf("too many attempts, temporarily blocked")
	}

	return true, nil
}

// Company and device cache methods

// GetCompanyUserPermission caches company user permission checks
func (cs *CacheStrategy) GetCompanyUserPermission(ctx context.Context, companyID, userID string, fallback func() (bool, error)) (bool, error) {
	key := fmt.Sprintf("company:%s:user:%s:permission",
		utilsCrypto.GetHash(companyID),
		utilsCrypto.GetHash(userID))

	result, err := cs.GetWithFallback(ctx, key,
		constants.TTL_LOCAL_USER_PERMISSION,
		constants.TTL_USER_PERMISSION,
		func() (interface{}, error) {
			return fallback()
		})

	if err != nil {
		return false, err
	}

	if str, ok := result.(string); ok {
		return strconv.ParseBool(str)
	}

	if val, ok := result.(bool); ok {
		return val, nil
	}

	return false, fmt.Errorf("invalid permission format in cache")
}

// GetDeviceInCompany caches device existence checks
func (cs *CacheStrategy) GetDeviceInCompany(ctx context.Context, companyID, deviceID string, fallback func() (bool, error)) (bool, error) {
	key := fmt.Sprintf("company:%s:device:%s:exists",
		utilsCrypto.GetHash(companyID),
		utilsCrypto.GetHash(deviceID))

	result, err := cs.GetWithFallback(ctx, key,
		constants.TTL_LOCAL_DEVICE_CHECK,
		constants.TTL_DEVICE_CHECK,
		func() (interface{}, error) {
			return fallback()
		})

	if err != nil {
		return false, err
	}

	if str, ok := result.(string); ok {
		return strconv.ParseBool(str)
	}

	if val, ok := result.(bool); ok {
		return val, nil
	}

	return false, fmt.Errorf("invalid device check format in cache")
}

// InvalidateUserCache invalidates all user-related cache entries
func (cs *CacheStrategy) InvalidateUserCache(ctx context.Context, userID string) {
	userIDHash := utilsCrypto.GetHash(userID)

	// Invalidate user info cache
	userInfoKey := utilsCache.GetKeyCacheUserInfoView(userIDHash)
	cs.Delete(ctx, userInfoKey)

	// Invalidate permission caches (we need to be more specific here based on your use case)
	// This is a simplified version - in practice, you might want to maintain a list of cache keys per user
}

// InvalidateCompanyCache invalidates company-related cache entries
func (cs *CacheStrategy) InvalidateCompanyCache(ctx context.Context, companyID string) {
	// Similar to user cache invalidation, implement based on your specific needs
}
