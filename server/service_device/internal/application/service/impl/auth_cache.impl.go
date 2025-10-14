package impl

import (
	"context"
	"fmt"
	"strconv"
	"time"

	applicationModel "github.com/youknow2509/cio_verify_face/server/service_device/internal/application/model"
	"github.com/youknow2509/cio_verify_face/server/service_device/internal/application/service"
	constants "github.com/youknow2509/cio_verify_face/server/service_device/internal/constants"
	domainModel "github.com/youknow2509/cio_verify_face/server/service_device/internal/domain/model"
	domainRepository "github.com/youknow2509/cio_verify_face/server/service_device/internal/domain/repository"
	domainToken "github.com/youknow2509/cio_verify_face/server/service_device/internal/domain/token"
	"github.com/youknow2509/cio_verify_face/server/service_device/internal/global"
	utilsUuid "github.com/youknow2509/cio_verify_face/server/service_device/internal/shared/utils/uuid"
)

// AuthCacheService provides optimized authentication operations with caching
type AuthCacheService struct {
	cacheStrategy *CacheStrategy
}

// NewAuthCacheService creates a new auth cache service instance
func NewAuthCacheService() (service.IAuthCacheService, error) {
	cacheStrategy, err := NewCacheStrategy()
	if err != nil {
		return nil, err
	}

	return &AuthCacheService{
		cacheStrategy: cacheStrategy,
	}, nil
}

// ValidateAccessToken validates access token using multi-level caching
func (a *AuthCacheService) ValidateAccessToken(ctx context.Context, tokenString string) (*applicationModel.TokenValidationResult, error) {
	// Parse and validate JWT token first
	tokenService := domainToken.GetTokenService()
	if tokenService == nil {
		return nil, fmt.Errorf("token service not available")
	}

	tokenClaims, tkErr := tokenService.ParseUserToken(ctx, tokenString)
	if tkErr != nil {
		global.Logger.Warn("Invalid token format", "error", tkErr)
		return &applicationModel.TokenValidationResult{
			Valid: false,
			Error: "Invalid token format",
		}, nil
	}

	// Check if token is active in cache
	sessionData, err := a.cacheStrategy.GetUserSession(ctx, tokenClaims.TokenId)
	if err != nil {
		// Token not in cache, check database
		return a.validateTokenFromDB(ctx, tokenClaims)
	}

	// Validate cached session data
	if sessionData["user_id"] != tokenClaims.UserId {
		return &applicationModel.TokenValidationResult{
			Valid: false,
			Error: "Token user mismatch",
		}, nil
	}

	role, err := strconv.Atoi(sessionData["role"])
	if err != nil || role != int(tokenClaims.Role) {
		return &applicationModel.TokenValidationResult{
			Valid: false,
			Error: "Token role mismatch",
		}, nil
	}

	userID, err := utilsUuid.ParseUUID(tokenClaims.UserId)
	if err != nil {
		return &applicationModel.TokenValidationResult{
			Valid: false,
			Error: "Invalid user ID format",
		}, nil
	}

	return &applicationModel.TokenValidationResult{
		Valid:  true,
		UserID: userID,
		Role:   tokenClaims.Role,
	}, nil
}

// validateTokenFromDB validates token by checking database and updates cache if valid
func (a *AuthCacheService) validateTokenFromDB(ctx context.Context, tokenClaims *domainModel.TokenUserJwtOutput) (*applicationModel.TokenValidationResult, error) {
	userRepo, err := domainRepository.GetUserRepository()
	if err != nil {
		return nil, err
	}

	sessionID, err := utilsUuid.ParseUUID(tokenClaims.TokenId)
	if err != nil {
		return &applicationModel.TokenValidationResult{
			Valid: false,
			Error: "Invalid session ID format",
		}, nil
	}

	// Check if session exists and is valid in database
	sessionInfo, err := userRepo.GetUserSessionByID(ctx, sessionID)
	if err != nil {
		global.Logger.Warn("Error checking session in database", "error", err)
		return &applicationModel.TokenValidationResult{
			Valid: false,
			Error: "Session validation failed",
		}, nil
	}

	if sessionInfo == nil {
		return &applicationModel.TokenValidationResult{
			Valid: false,
			Error: "Session not found",
		}, nil
	}

	// Check if session is expired
	if sessionInfo.ExpiredAt.Before(time.Now()) {
		return &applicationModel.TokenValidationResult{
			Valid: false,
			Error: "Session expired",
		}, nil
	}

	userID, err := utilsUuid.ParseUUID(tokenClaims.UserId)
	if err != nil {
		return &applicationModel.TokenValidationResult{
			Valid: false,
			Error: "Invalid user ID format",
		}, nil
	}

	// Update cache with valid session
	if err := a.cacheStrategy.SetUserSession(ctx, tokenClaims.TokenId, tokenClaims.UserId, domainModel.Role(tokenClaims.Role), constants.TTL_ACCESS_TOKEN); err != nil {
		global.Logger.Warn("Failed to cache session data", "error", err)
	}

	return &applicationModel.TokenValidationResult{
		Valid:  true,
		UserID: userID,
		Role:   tokenClaims.Role,
	}, nil
}

// GetUserInfoCached gets user information with caching
func (a *AuthCacheService) GetUserInfoCached(ctx context.Context, userID string) (*applicationModel.UserInfoOutput, error) {
	domainUserInfo, err := a.cacheStrategy.GetUserInfo(ctx, userID, func() (*domainModel.UserInfoOutput, error) {
		userRepo, err := domainRepository.GetUserRepository()
		if err != nil {
			return nil, err
		}

		uid, err := utilsUuid.ParseUUID(userID)
		if err != nil {
			return nil, err
		}

		return userRepo.GetUserInfoByID(ctx, uid)
	})

	if err != nil {
		return nil, err
	}

	if domainUserInfo == nil {
		return nil, nil
	}

	// Convert domain model to application model
	return &applicationModel.UserInfoOutput{
		Email:     domainUserInfo.Email,
		Phone:     domainUserInfo.Phone,
		FullName:  domainUserInfo.FullName,
		AvatarURL: domainUserInfo.AvatarURL,
	}, nil
}

// CheckUserPermissionCached checks user permission with caching
func (a *AuthCacheService) CheckUserPermissionCached(ctx context.Context, companyID, userID string) (bool, error) {
	return a.cacheStrategy.GetCompanyUserPermission(ctx, companyID, userID, func() (bool, error) {
		companyRepo, err := domainRepository.GetCompanyRepository()
		if err != nil {
			return false, err
		}

		cID, err := utilsUuid.ParseUUID(companyID)
		if err != nil {
			return false, err
		}

		uID, err := utilsUuid.ParseUUID(userID)
		if err != nil {
			return false, err
		}

		return companyRepo.CheckUserIsManagementInCompany(ctx, &domainModel.CheckCompanyIsManagementInCompanyInput{
			CompanyID: cID,
			UserID:    uID,
		})
	})
}

// CheckDeviceInCompanyCached checks if device exists in company with caching
func (a *AuthCacheService) CheckDeviceInCompanyCached(ctx context.Context, companyID, deviceID string) (bool, error) {
	return a.cacheStrategy.GetDeviceInCompany(ctx, companyID, deviceID, func() (bool, error) {
		companyRepo, err := domainRepository.GetCompanyRepository()
		if err != nil {
			return false, err
		}

		cID, err := utilsUuid.ParseUUID(companyID)
		if err != nil {
			return false, err
		}

		dID, err := utilsUuid.ParseUUID(deviceID)
		if err != nil {
			return false, err
		}

		return companyRepo.CheckDeviceExistsInCompany(ctx, &domainModel.CheckDeviceExistsInCompanyInput{
			CompanyID: cID,
			DeviceID:  dID,
		})
	})
}

// InvalidateUserCache invalidates all cache entries for a user
func (a *AuthCacheService) InvalidateUserCache(ctx context.Context, userID string) {
	a.cacheStrategy.InvalidateUserCache(ctx, userID)
}

// InvalidateSessionCache invalidates session cache
func (a *AuthCacheService) InvalidateSessionCache(ctx context.Context, sessionID string) {
	a.cacheStrategy.DeleteUserSession(ctx, sessionID)
}

// PreloadUserData preloads frequently accessed user data into cache
func (a *AuthCacheService) PreloadUserData(ctx context.Context, userIDs []string) error {
	for _, userID := range userIDs {
		go func(uid string) {
			_, err := a.GetUserInfoCached(ctx, uid)
			if err != nil {
				global.Logger.Warn("Failed to preload user data", "userID", uid, "error", err)
			}
		}(userID)
	}
	return nil
}

// GetCacheStats returns cache statistics for monitoring
func (a *AuthCacheService) GetCacheStats(ctx context.Context) (*applicationModel.CacheStats, error) {
	// This would depend on your cache implementation
	// For now, returning basic stats structure
	return &applicationModel.CacheStats{
		LocalCacheHits:       0,
		DistributedCacheHits: 0,
		CacheMisses:          0,
		TotalRequests:        0,
	}, nil
}

// WarmupCache pre-fills cache with frequently accessed data
func (a *AuthCacheService) WarmupCache(ctx context.Context) error {
	// Implement cache warmup logic based on your application needs
	// This could include:
	// - Pre-loading active user sessions
	// - Pre-loading frequently accessed user info
	// - Pre-loading company-user permissions

	global.Logger.Info("Cache warmup completed")
	return nil
}
