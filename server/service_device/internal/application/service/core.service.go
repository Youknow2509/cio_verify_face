package service

import (
	"context"
	"errors"

	errorService "github.com/youknow2509/cio_verify_face/server/service_device/internal/application/errors"
	"github.com/youknow2509/cio_verify_face/server/service_device/internal/application/model"
)

// =======================================================
//
//	Define interfaces for Auth service
//
// =======================================================
type (
	// Core Auth
	ICoreAuthService interface {
		// Login
		Login(ctx context.Context, input *model.LoginInput) (*model.LoginOutput, *errorService.Error)
		// Login Admin
		LoginAdmin(ctx context.Context, input *model.LoginInputAdmin) (*model.LoginOutput, *errorService.Error)
		// Logout
		Logout(ctx context.Context, input *model.LogoutInput) *errorService.Error
		// Refresh token
		RefreshToken(ctx context.Context, input *model.RefreshTokenInput) (*model.RefreshTokenOutput, *errorService.Error)
		// Get my info
		GetMyInfo(ctx context.Context, input *model.GetMyInfoInput) (*model.GetMyInfoOutput, *errorService.Error)
		// Update device token
		UpdateDeviceSession(ctx context.Context, input *model.UpdateDeviceSessionInput) (*model.UpdateDeviceSessionOutput, *errorService.Error)
		// Delete device token
		DeleteDeviceSession(ctx context.Context, input *model.DeleteDeviceSessionInput) *errorService.Error
	}

	// Auth Cache Service for optimized operations
	IAuthCacheService interface {
		// Validate access token with caching
		ValidateAccessToken(ctx context.Context, tokenString string) (*model.TokenValidationResult, error)
		// Get user info with caching
		GetUserInfoCached(ctx context.Context, userID string) (*model.UserInfoOutput, error)
		// Check user permission with caching
		CheckUserPermissionCached(ctx context.Context, companyID, userID string) (bool, error)
		// Check device in company with caching
		CheckDeviceInCompanyCached(ctx context.Context, companyID, deviceID string) (bool, error)
		// Invalidate user cache
		InvalidateUserCache(ctx context.Context, userID string)
		// Invalidate session cache
		InvalidateSessionCache(ctx context.Context, sessionID string)
		// Preload user data
		PreloadUserData(ctx context.Context, userIDs []string) error
		// Get cache statistics
		GetCacheStats(ctx context.Context) (*model.CacheStats, error)
		// Warmup cache
		WarmupCache(ctx context.Context) error
	}
)

// =======================================================
//
//	Variables instance interfaces for Auth service
//
// =======================================================
var (
	_ICoreAuthService  ICoreAuthService
	_IAuthCacheService IAuthCacheService
)

// =======================================================
//
//	Getter, setter for Auth service interfaces
//
// =======================================================
func GetCoreAuthService() ICoreAuthService {
	return _ICoreAuthService
}

func SetCoreAuthService(s ICoreAuthService) error {
	if _ICoreAuthService != nil {
		return errors.New("auth service is already set")
	}
	_ICoreAuthService = s
	return nil
}

func GetAuthCacheService() IAuthCacheService {
	return _IAuthCacheService
}

func SetAuthCacheService(s IAuthCacheService) error {
	if _IAuthCacheService != nil {
		return errors.New("auth cache service is already set")
	}
	_IAuthCacheService = s
	return nil
}
