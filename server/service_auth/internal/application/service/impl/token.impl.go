package impl

import (
	"context"
	"encoding/json"
	"errors"
	"net/netip"
	"time"

	"github.com/google/uuid"
	"github.com/youknow2509/cio_verify_face/server/service_auth/internal/application/model"
	appService "github.com/youknow2509/cio_verify_face/server/service_auth/internal/application/service"
	constants "github.com/youknow2509/cio_verify_face/server/service_auth/internal/constants"
	domainCache "github.com/youknow2509/cio_verify_face/server/service_auth/internal/domain/cache"
	domainError "github.com/youknow2509/cio_verify_face/server/service_auth/internal/domain/errors"
	domainModel "github.com/youknow2509/cio_verify_face/server/service_auth/internal/domain/model"
	domainRepo "github.com/youknow2509/cio_verify_face/server/service_auth/internal/domain/repository"
	domainToken "github.com/youknow2509/cio_verify_face/server/service_auth/internal/domain/token"
	"github.com/youknow2509/cio_verify_face/server/service_auth/internal/global"
	sharedCache "github.com/youknow2509/cio_verify_face/server/service_auth/internal/shared/utils/cache"
	sharedCrypto "github.com/youknow2509/cio_verify_face/server/service_auth/internal/shared/utils/crypto"
)

/**
 * Token service implementations
 */
type TokenService struct {
}

func (t *TokenService) ParseTokenDevice(ctx context.Context, input model.ParseTokenDeviceInput) (*model.ParseTokenDeviceOutput, error) {
	// 1. Check token status token in cache
	cacheLocal, _ := domainCache.GetLocalCache()
	if cacheLocal == nil {
		global.Logger.Error("local cache not initialized")
		return nil, errors.New("internal error")
	}
	key := sharedCache.GetKeyStatusTokenUser(sharedCrypto.GetHash(input.Token))
	cached, err := cacheLocal.Get(ctx, key)
	if err != nil {
		global.Logger.Error("failed to get token status from cache", "error", err.Error(), "key", key)
		return nil, err
	}
	if cached != "" {
		var cachedRes model.ParseTokenDeviceOutput
		if err := json.Unmarshal([]byte(cached), &cachedRes); err != nil {
			global.Logger.Error("failed to unmarshal token status from cache", "error", err.Error(), "key", key)
			return nil, err
		}
		return &cachedRes, nil
	}
	// 2. Validate token validate
	tokenService := domainToken.GetTokenService()
	tokenData, tokenErr := tokenService.ParseDeviceToken(ctx, input.Token)
	if tokenErr != nil {
		switch tokenErr.Code {
		case domainError.TokenErrorNotFoundCode:
			global.Logger.Warn("device token not found", "token", input.Token)
			return nil, nil
		case domainError.TokenMalformedErrorCode:
			global.Logger.Warn("device token malformed", "token", input.Token)
			return nil, nil
		case domainError.TokenSignatureInvalidErrCode:
			global.Logger.Warn("device token signature invalid", "token", input.Token)
			return nil, nil
		case domainError.TokenExpiredErrorCode:
			global.Logger.Warn("device token expired", "token", input.Token)
			return nil, nil
		case domainError.TokenValidationErrorCode:
			global.Logger.Warn("device token validation error", "token", input.Token)
			return nil, nil
		}
	}
	// 3. Cache token status
	output := &model.ParseTokenDeviceOutput{
		TokenId:   tokenData.TokenId,
		DeviceId:  tokenData.DeviceId,
		CompanyId: tokenData.CompanyId,
		Expires:   tokenData.ExpiresAt,
	}
	valBytes, err := json.Marshal(output)
	if err != nil {
		global.Logger.Error("failed to marshal token status for cache", "error", err.Error(), "key", key)
		return nil, err
	}
	if err := cacheLocal.SetTTL(ctx, key, string(valBytes), constants.TTL_CACHE_TOKEN); err != nil {
		global.Logger.Warn("failed to set token status in cache", "error", err.Error(), "key", key)
	}
	return output, nil
}

// BlockTokenDevice implements service.ITokenService.
func (t *TokenService) BlockTokenDevice(ctx context.Context, input model.BlockTokenDeviceInput) error {
	//
	db := domainRepo.GetDeviceRepository()
	exist, err := db.DeviceExist(ctx, input.DeviceId)
	if err != nil {
		return err
	}
	if !exist {
		global.Logger.Warn("device not found", "device_id", input.DeviceId)
		return errors.New("device not found")
	}
	if err := db.BlockDeviceToken(ctx, input.DeviceId); err != nil {
		return err
	}
	return nil
}

// BlockTokenUser implements service.ITokenService.
func (t *TokenService) BlockTokenUser(ctx context.Context, input model.BlockTokenUserInput) error {
	//
	db, _ := domainRepo.GetUserRepository()
	resp, err := db.GetUserSessionByID(ctx, input.TokenId)
	if err != nil {
		return err
	}
	if resp == nil {
		global.Logger.Warn("user session not found", "token_id", input.TokenId.String())
		return errors.New("user session not found")
	}
	// Rm session in db
	if err := db.RemoveUserSession(
		ctx,
		&domainModel.RemoveUserSessionInput{
			SessionID: input.TokenId,
		},
	); err != nil {
		return err
	}
	// Block access token in cache
	cache, _ := domainCache.GetDistributedCache()
	key := sharedCache.GetKeyUserAccessTokenIsActive(sharedCrypto.GetHash(input.TokenId.String()))
	if err := cache.Delete(ctx, key); err != nil {
		// Log cache deletion error
		global.Logger.Error("failed to delete user access token from cache", "error", err.Error(), "key", key)
		return err
	}
	return nil
}

// BlockTokenUserRefresh implements service.ITokenService.
func (t *TokenService) BlockTokenUserRefresh(ctx context.Context, input model.BlockTokenUserRefreshInput) error {
	//
	db, _ := domainRepo.GetUserRepository()
	resp, err := db.GetUserSessionByID(ctx, input.TokenId)
	if err != nil {
		return err
	}
	if resp == nil {
		global.Logger.Warn("user session not found", "token_id", input.TokenId.String())
		return errors.New("user session not found")
	}
	// Rm session in db
	if err := db.RemoveUserSession(
		ctx,
		&domainModel.RemoveUserSessionInput{
			SessionID: input.TokenId,
		},
	); err != nil {
		return err
	}
	return nil
}

// CheckTokenDevice implements service.ITokenService.
func (t *TokenService) CheckTokenDevice(ctx context.Context, input model.CheckTokenDeviceInput) (bool, string, error) {
	db := domainRepo.GetDeviceRepository()
	ok, _, err := db.CheckTokenDevice(
		ctx,
		&domainModel.CheckTokenDeviceInput{
			DeviceId: input.DeviceId,
			Token:    input.Token,
		},
	)
	if err != nil {
		global.Logger.Error("failed to check device token", "error", err.Error(), "device_id", input.DeviceId)
		return false, "", err
	}
	if !ok {
		return false, "", nil
	}
	return true, "", nil
}

// CreateTokenDevice implements service.ITokenService.
func (t *TokenService) CreateTokenDevice(ctx context.Context, input model.CreateTokenDeviceInput) (string, error) {
	// Check device exist
	db := domainRepo.GetDeviceRepository()
	exist, err := db.DeviceExist(ctx, input.DeviceId)
	if err != nil {
		return "", err
	}
	if !exist {
		global.Logger.Warn("device not found", "device_id", input.DeviceId)
		return "", errors.New("device not found")
	}
	// Generate token
	tokenService := domainToken.GetTokenService()
	tokenId := uuid.New()
	token, err := tokenService.CreateDeviceToken(ctx, &domainModel.TokenDeviceJwtInput{
		DeviceId:  input.DeviceId.String(),
		CompanyId: input.CompanyId.String(),
		TokenId:   tokenId.String(),
		Expires:   time.Now().Add(constants.TTL_DEVICE_TOKEN_LONG * time.Second),
	})
	if err != nil {
		global.Logger.Error("failed to create device token", "error", err.Error(), "device_id", input.DeviceId)
		return "", err
	}
	// Store token in db
	if err := db.CreateDeviceToken(
		ctx,
		&domainModel.CreateDeviceTokenInput{
			DeviceId: input.DeviceId,
			NewToken: token,
		},
	); err != nil {
		global.Logger.Error("failed to create device token", "error", err.Error(), "device_id", input.DeviceId)
		return "", err
	}
	return token, nil
}

// CreateUserToken implements service.ITokenService.
func (t *TokenService) CreateUserToken(ctx context.Context, input model.CreateTokenUserInput) (*model.CreateTokenUserOutput, error) {
	// Check user is exist
	userRepo, _ := domainRepo.GetUserRepository()
	userExist, err := userRepo.GetUserBaseByID(
		ctx,
		input.UserId,
	)
	if err != nil {
		return nil, err
	}
	if userExist == nil {
		global.Logger.Warn("user not found", "user_id", input.UserId.String())
		return nil, errors.New("user not found")
	}
	// Get company ID
	companyRepo, err := domainRepo.GetCompanyRepository()
	if err != nil {
		global.Logger.Error("Error getting company repository: ", err)
		return nil, err
	}
	companyReps, err := companyRepo.GetCompanyUser(ctx, &domainModel.GetCompanyUserInput{UserID: input.UserId})
	if err != nil {
		global.Logger.Error("Error getting company user: ", err)
		return nil, err
	}
	companyId := ""
	if companyReps != nil {
		companyId = companyReps.CompanyID.String()
	}
	// Create access token
	tokenService := domainToken.GetTokenService()
	tokenUuid := uuid.New()
	accessTokenTTl := constants.TTL_ACCESS_TOKEN * time.Second
	accessToken, err := tokenService.CreateUserToken(
		ctx,
		&domainModel.TokenUserJwtInput{
			UserId:    input.UserId.String(),
			TokenId:   tokenUuid.String(),
			Expires:   time.Now().Add(accessTokenTTl),
			Role:      userExist.Role,
			CompanyId: companyId,
		},
	)
	if err != nil {
		global.Logger.Error("failed to create user access token", "error", err.Error(), "user_id", input.UserId.String())
		return nil, err
	}
	// Create refresh token
	refreshTokenTtl := constants.TTL_REFRESH_TOKEN * time.Second
	refreshToken, err := tokenService.CreateUserRefreshToken(
		ctx,
		&domainModel.TokenUserRefreshInput{
			UserId:  input.UserId.String(),
			TokenId: tokenUuid.String(),
			Expires: time.Now().Add(refreshTokenTtl),
		},
	)
	if err != nil {
		global.Logger.Error("failed to create user refresh token", "error", err.Error(), "user_id", input.UserId.String())
		return nil, err
	}
	// Store session in db
	if err := userRepo.CreateUserSession(
		ctx,
		&domainModel.CreateUserSessionInput{
			SessionID:    tokenUuid,
			UserID:       input.UserId,
			IPAddress:    netip.Addr{},
			UserAgent:    "",
			RefreshToken: refreshToken,
			ExpiredAt:    time.Now().Add(refreshTokenTtl),
		},
	); err != nil {
		global.Logger.Error("failed to create user session", "error", err.Error(), "user_id", input.UserId.String())
		return nil, err
	}
	// Save access token in cache
	cache, _ := domainCache.GetDistributedCache()
	key := sharedCache.GetKeyUserAccessTokenIsActive(sharedCrypto.GetHash(tokenUuid.String()))
	val := "1"
	if err := cache.SetTTL(
		ctx,
		key,
		val,
		constants.TTL_ACCESS_TOKEN,
	); err != nil {
		global.Logger.Error("failed to set user access token in cache", "error", err.Error(), "key", key)
		return nil, err
	}
	return &model.CreateTokenUserOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// ParseTokenUser implements service.ITokenService.
func (t *TokenService) ParseTokenUser(ctx context.Context, input model.ParseTokenUserInput) (*model.ParseTokenUserOutput, error) {
	// 1. Check token status token in cache
	cacheLocal, _ := domainCache.GetLocalCache()
	if cacheLocal == nil {
		global.Logger.Error("local cache not initialized")
		return nil, errors.New("internal error")
	}
	key := sharedCache.GetKeyStatusTokenUser(sharedCrypto.GetHash(input.Token))
	cached, err := cacheLocal.Get(ctx, key)
	if err != nil {
		global.Logger.Error("failed to get token status from cache", "error", err.Error(), "key", key)
		return nil, err
	}
	if cached != "" {
		var cachedRes model.ParseTokenUserOutput
		if err := json.Unmarshal([]byte(cached), &cachedRes); err != nil {
			global.Logger.Error("failed to unmarshal token status from cache", "error", err.Error(), "key", key)
			return nil, err
		}
		return &cachedRes, nil
	}
	// 2. Validate token validate
	tokenService := domainToken.GetTokenService()
	tokenResp, tokenErr := tokenService.ParseUserToken(ctx, input.Token)
	if tokenErr != nil {
		switch tokenErr.Code {
		case domainError.TokenErrorNotFoundCode:
			global.Logger.Warn("token not found", "token", input.Token)
			return nil, errors.New("token not found")
		case domainError.TokenMalformedErrorCode:
			global.Logger.Warn("token malformed", "token", input.Token)
			return nil, errors.New("token malformed")
		case domainError.TokenSignatureInvalidErrCode:
			global.Logger.Warn("token signature invalid", "token", input.Token)
			return nil, errors.New("token signature invalid")
		case domainError.TokenExpiredErrorCode:
			global.Logger.Warn("token expired", "token", input.Token)
			return nil, errors.New("token expired")
		case domainError.TokenValidationErrorCode:
			global.Logger.Warn("token validation error", "token", input.Token)
			return nil, errors.New("token validation error")
		default:
			global.Logger.Error("unknown token error", "error", tokenErr.Message, "token", input.Token)
			return nil, errors.New("unknown token error")
		}
	}
	// 3. Cache token status
	output := &model.ParseTokenUserOutput{
		UserId:    tokenResp.UserId,
		TokenId:   tokenResp.TokenId,
		Role:      tokenResp.Role,
		CompanyId: tokenResp.CompanyId,
	}
	valBytes, err := json.Marshal(output)
	if err != nil {
		global.Logger.Error("failed to marshal token status for cache", "error", err.Error(), "key", key)
		return nil, err
	}
	if err := cacheLocal.SetTTL(ctx, key, string(valBytes), constants.TTL_CACHE_TOKEN); err != nil {
		global.Logger.Error("failed to set token status in cache", "error", err.Error(), "key", key)
		return nil, err
	}
	return output, nil
}

// RefreshTokenUser implements service.ITokenService.
func (t *TokenService) RefreshTokenUser(ctx context.Context, input model.RefreshTokenUserInput) (*model.RefreshTokenUserOutput, error) {
	// Validate access token
	tokenService := domainToken.GetTokenService()
	accessTokenResp, errToken := tokenService.ParseUserToken(ctx, input.AccessToken)
	if errToken != nil {
		switch errToken.Code {
		case domainError.TokenErrorNotFoundCode:
			global.Logger.Warn("token not found", "token", input.AccessToken)
			return nil, errors.New("token not found")
		case domainError.TokenMalformedErrorCode:
			global.Logger.Warn("token malformed", "token", input.AccessToken)
			return nil, errors.New("token malformed")
		case domainError.TokenSignatureInvalidErrCode:
			global.Logger.Warn("token signature invalid", "token", input.AccessToken)
			return nil, errors.New("token signature invalid")
		case domainError.TokenValidationErrorCode:
			global.Logger.Warn("token validation error", "token", input.AccessToken)
			return nil, errors.New("token validation error")
		}
	}
	// Validate refresh token
	_, errToken = tokenService.ParseUserRefreshToken(ctx, input.RefreshToken)
	if errToken != nil {
		switch errToken.Code {
		case domainError.TokenErrorNotFoundCode:
			global.Logger.Warn("refresh token not found", "token", input.RefreshToken)
			return nil, errors.New("refresh token not found")
		case domainError.TokenMalformedErrorCode:
			global.Logger.Warn("refresh token malformed", "token", input.RefreshToken)
			return nil, errors.New("refresh token malformed")
		case domainError.TokenSignatureInvalidErrCode:
			global.Logger.Warn("refresh token signature invalid", "token", input.RefreshToken)
			return nil, errors.New("refresh token signature invalid")
		case domainError.TokenExpiredErrorCode:
			global.Logger.Warn("refresh token expired", "token", input.RefreshToken)
			return nil, errors.New("refresh token expired")
		case domainError.TokenValidationErrorCode:
			global.Logger.Warn("refresh token validation error", "token", input.RefreshToken)
			return nil, errors.New("refresh token validation error")
		default:
			global.Logger.Error("unknown refresh token error", "error", errToken.Message, "token", input.RefreshToken)
			return nil, errors.New("unknown refresh token error")
		}
	}
	// Check token
	domainRepo, _ := domainRepo.GetUserRepository()
	sessionID, _ := uuid.Parse(accessTokenResp.TokenId)
	userSession, err := domainRepo.GetUserSessionByID(
		ctx,
		sessionID,
	)
	if err != nil {
		return nil, err
	}
	if userSession == nil {
		global.Logger.Warn("user session not found", "session_id", sessionID.String())
		return nil, errors.New("user session not found")
	}
	if userSession.RefreshToken != input.RefreshToken {
		global.Logger.Warn("refresh token does not match", "session_id", sessionID.String())
		return nil, errors.New("refresh token does not match")
	}
	// Create new access token and refresh token
	accessTokenTTl := constants.TTL_ACCESS_TOKEN * time.Second
	accessToken, err := tokenService.CreateUserToken(
		ctx,
		&domainModel.TokenUserJwtInput{
			UserId:  accessTokenResp.UserId,
			TokenId: accessTokenResp.TokenId,
			Expires: time.Now().Add(accessTokenTTl),
			Role:    accessTokenResp.Role,
		},
	)
	if err != nil {
		global.Logger.Error("failed to create user access token", "error", err.Error(), "user_id", accessTokenResp.UserId)
		return nil, err
	}
	// Create refresh token
	refreshTokenTtl := constants.TTL_REFRESH_TOKEN * time.Second
	refreshToken, err := tokenService.CreateUserRefreshToken(
		ctx,
		&domainModel.TokenUserRefreshInput{
			UserId:  accessTokenResp.UserId,
			TokenId: accessTokenResp.TokenId,
			Expires: time.Now().Add(refreshTokenTtl),
		},
	)
	if err != nil {
		global.Logger.Error("failed to create user refresh token", "error", err.Error(), "user_id", accessTokenResp.UserId)
		return nil, err
	}
	// Update session in db
	if err := domainRepo.RefreshSession(
		ctx,
		&domainModel.RefreshSessionInput{
			SessionID:    sessionID,
			RefreshToken: refreshToken,
			ExpiredAt:    time.Now().Add(refreshTokenTtl),
		},
	); err != nil {
		global.Logger.Error("failed to refresh user session", "error", err.Error(), "user_id", accessTokenResp.UserId)
		return nil, err
	}
	// Update access token in cache
	cache, _ := domainCache.GetDistributedCache()
	key := sharedCache.GetKeyUserAccessTokenIsActive(sharedCrypto.GetHash(sessionID.String()))
	if err := cache.SetTTL(
		ctx,
		key,
		"1",
		constants.TTL_ACCESS_TOKEN,
	); err != nil {
		global.Logger.Error("failed to update user access token in cache", "error", err.Error(), "key", key)
		return nil, err
	}
	return &model.RefreshTokenUserOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// New token service and impl
func NewTokenService() appService.ITokenService {
	return &TokenService{}
}
