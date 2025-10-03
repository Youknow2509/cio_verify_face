package impl

import (
	"context"
	"net/netip"
	"strconv"
	"time"

	"github.com/youknow2509/cio_verify_face/server/service_auth/internal/application/errors"
	applicationModel "github.com/youknow2509/cio_verify_face/server/service_auth/internal/application/model"
	"github.com/youknow2509/cio_verify_face/server/service_auth/internal/application/service"
	constants "github.com/youknow2509/cio_verify_face/server/service_auth/internal/constants"
	domainCache "github.com/youknow2509/cio_verify_face/server/service_auth/internal/domain/cache"
	domainError "github.com/youknow2509/cio_verify_face/server/service_auth/internal/domain/errors"
	domainModel "github.com/youknow2509/cio_verify_face/server/service_auth/internal/domain/model"
	domainRepository "github.com/youknow2509/cio_verify_face/server/service_auth/internal/domain/repository"
	domainToken "github.com/youknow2509/cio_verify_face/server/service_auth/internal/domain/token"
	"github.com/youknow2509/cio_verify_face/server/service_auth/internal/global"
	utilsCache "github.com/youknow2509/cio_verify_face/server/service_auth/internal/shared/utils/cache"
	utilsCrypto "github.com/youknow2509/cio_verify_face/server/service_auth/internal/shared/utils/crypto"
	utilsRandom "github.com/youknow2509/cio_verify_face/server/service_auth/internal/shared/utils/random"
	utilsUuid "github.com/youknow2509/cio_verify_face/server/service_auth/internal/shared/utils/uuid"
)

/**
 * Define CoreAuthService struct implementing
 */
type CoreAuthService struct{}

// CreateDeviceSession implements service.ICoreAuthService.
func (c *CoreAuthService) UpdateDeviceSession(ctx context.Context, input *applicationModel.UpdateDeviceSessionInput) (*applicationModel.UpdateDeviceSessionOutput, *errors.Error) {
	// Check user exists in company
	companyRepo, err := domainRepository.GetCompanyRepository()
	if err != nil {
		global.Logger.Error("Error getting company repository: ", err)
		return nil, errors.GetError(errors.SystemTemporaryUnavailableErrorCode)
	}
	ok, err := companyRepo.CheckUserIsManagementInCompany(
		ctx,
		&domainModel.CheckCompanyIsManagementInCompanyInput{
			CompanyID: input.CompanyId,
			UserID:    input.UserId,
		},
	)
	if err != nil {
		global.Logger.Error("Error checking user is management in company: ", err)
		return nil, errors.GetError(errors.SystemTemporaryUnavailableErrorCode)
	}
	if !ok {
		// User not found or not management in company
		return nil, errors.GetError(errors.AuthDontHavePermissionErrorCode)
	}
	// Check device exists in company
	ok, err = companyRepo.CheckDeviceExistsInCompany(
		ctx,
		&domainModel.CheckDeviceExistsInCompanyInput{
			CompanyID: input.CompanyId,
			DeviceID:  input.DeviceId,
		},
	)
	if err != nil {
		global.Logger.Error("Error checking device exists in company: ", err)
		return nil, errors.GetError(errors.SystemTemporaryUnavailableErrorCode)
	}
	if !ok {
		// Device not found in company
		return nil, errors.GetError(errors.DeviceNotFoundErrorCode)
	}
	// Create token
	tokenService := domainToken.GetTokenService()
	if tokenService == nil {
		global.Logger.Error("Error getting token service: ", err)
		return nil, errors.GetError(errors.SystemTemporaryUnavailableErrorCode)
	}
	tokenId := utilsRandom.GenerateUUID()
	tokenTtl := time.Duration(constants.TTL_TOKEN_DEVICE) * time.Second
	token, err := tokenService.CreateDeviceRefreshToken(
		ctx,
		&domainModel.TokenDeviceRefreshInput{
			TokenId:   tokenId.String(),
			CompanyId: input.CompanyId.String(),
			DeviceId:  input.DeviceId.String(),
			Expires:   time.Now().Add(tokenTtl),
		},
	)
	if err != nil {
		global.Logger.Error("Error creating device token: ", err)
		return nil, errors.GetError(errors.SystemTemporaryUnavailableErrorCode)
	}
	// Update device session to database
	if err := companyRepo.UpdateDeviceSession(
		ctx,
		&domainModel.UpdateDeviceSessionInput{
			DeviceId: input.DeviceId,
			Token:    token,
		},
	); err != nil {
		global.Logger.Error("Error updating device session: ", err)
		return nil, errors.GetError(errors.SystemTemporaryUnavailableErrorCode)
	}
	// Save audit log
	auditRepo, err := domainRepository.GetAuditRepository()
	if err != nil {
		global.Logger.Error("Error getting audit repository: ", err)
		return nil, errors.GetError(errors.SystemTemporaryUnavailableErrorCode)
	}
	if err := auditRepo.AddAuditLog(
		ctx,
		&domainModel.AuditLog{
			UserId:       input.UserId,
			Action:       constants.AuditActionUpdateSessionDevice,
			ResourceType: constants.AuditResourceTypeDevice,
			ResourceId:   input.DeviceId,
			OldValues:    nil,
			NewValues: map[string]interface{}{
				"token": token,
			},
			IpAddress: input.ClientIp,
			UserAgent: input.UserAgent,
			Timestamp: time.Now().Unix(),
		},
	); err != nil {
		global.Logger.Error("Error logging audit log: ", err)
		// Not return error
	}
	return &applicationModel.UpdateDeviceSessionOutput{
		Token:    token,
		ExpireAt: time.Now().Add(tokenTtl).Unix(),
	}, nil
}

// DeleteDeviceSession implements service.ICoreAuthService.
func (c *CoreAuthService) DeleteDeviceSession(ctx context.Context, input *applicationModel.DeleteDeviceSessionInput) *errors.Error {
	// TODO: Implement the delete device session logic here
	return nil
}

// GetMyInfo implements service.ICoreAuthService.
func (c *CoreAuthService) GetMyInfo(ctx context.Context, input *applicationModel.GetMyInfoInput) (*applicationModel.GetMyInfoOutput, *errors.Error) {
	// Get data from data base
	domainRepo, err := domainRepository.GetUserRepository()
	if err != nil {
		global.Logger.Error("Error getting user repository: ", err)
		return nil, errors.GetError(errors.SystemTemporaryUnavailableErrorCode)
	}
	response, err := domainRepo.GetUserInfoByID(ctx, input.UserId)
	if err != nil {
		global.Logger.Warn("Error getting user info by ID: ", err)
		return nil, errors.GetError(errors.SystemTemporaryUnavailableErrorCode)
	}
	if response == nil {
		// User not found
		return nil, errors.GetError(errors.UserNotFoundErrorCode)
	}
	return &applicationModel.GetMyInfoOutput{
		UserId:    input.UserId.String(),
		Email:     response.Email,
		Phone:     response.Phone,
		FullName:  response.FullName,
		AvatarURL: response.AvatarURL,
		Role:      input.Role,
	}, nil
}

// Login implements service.ICoreAuthService.
func (c *CoreAuthService) Login(ctx context.Context, input *applicationModel.LoginInput) (*applicationModel.LoginOutput, *errors.Error) {
	// Get info user with mail
	domainRepo, err := domainRepository.GetUserRepository()
	if err != nil {
		global.Logger.Error("Error getting user repository: ", err)
		return nil, errors.GetError(errors.SystemTemporaryUnavailableErrorCode)
	}
	response, err := domainRepo.GetUserBaseByEmail(ctx, input.UserName)
	if err != nil {
		global.Logger.Warn("Error getting user by email: ", err)
		return nil, errors.GetError(errors.SystemTemporaryUnavailableErrorCode)
	}
	if response == nil {
		// User not found
		return nil, errors.GetError(errors.UserNotFoundErrorCode)
	}
	// Check role
	if response.Role != domainModel.RoleUser {
		// Role not match
		return nil, errors.GetError(errors.UserNotFoundErrorCode)
	}
	// Check password
	if !utilsCrypto.ComparePasswordWithHash(
		input.Password,
		response.UserSalt,
		response.UserPassword,
	) {
		// Password not match
		return nil, errors.GetError(errors.UserPasswordIncorrectErrorCode)
	}
	// Create session
	tokenService := domainToken.GetTokenService()
	tokenId := utilsRandom.GenerateUUID()
	// Create access token
	timeTtlAccessToken := time.Duration(constants.TTL_ACCESS_TOKEN) * time.Second
	accessToken, err := tokenService.CreateUserToken(
		ctx,
		&domainModel.TokenUserJwtInput{
			UserId:  response.UserID,
			TokenId: tokenId.String(),
			Role:    domainModel.RoleUser,
			Expires: time.Now().Add(timeTtlAccessToken),
		},
	)
	if err != nil {
		global.Logger.Warn("Error creating user token: ", err)
		return nil, errors.GetError(errors.SystemTemporaryUnavailableErrorCode)
	}
	// Create refresh token
	timeTtlRefreshToken := time.Duration(constants.TTL_REFRESH_TOKEN) * time.Second
	refreshToken, err := tokenService.CreateUserRefreshToken(
		ctx,
		&domainModel.TokenUserRefreshInput{
			UserId:  response.UserID,
			TokenId: tokenId.String(),
			Expires: time.Now().Add(timeTtlRefreshToken),
		},
	)
	if err != nil {
		global.Logger.Warn("Error creating user refresh token: ", err)
		return nil, errors.GetError(errors.SystemTemporaryUnavailableErrorCode)
	}
	// Save session to db and cache
	uuidUser, _ := utilsUuid.ParseUUID(response.UserID)
	ipAddr, _ := netip.ParseAddr(input.ClientIp)
	if err := domainRepo.CreateUserSession(
		ctx,
		&domainModel.CreateUserSessionInput{
			SessionID:    tokenId,
			UserID:       uuidUser,
			IPAddress:    ipAddr,
			UserAgent:    input.UserAgent,
			RefreshToken: refreshToken,
			ExpiredAt:    time.Now().Add(timeTtlRefreshToken),
		},
	); err != nil {
		global.Logger.Error("Error creating user session: ", err)
		return nil, errors.GetError(errors.SystemTemporaryUnavailableErrorCode)
	}
	cacheDistributed, err := domainCache.GetDistributedCache()
	if err != nil {
		global.Logger.Error("Error getting distributed cache: ", err)
		return nil, errors.GetError(errors.SystemTemporaryUnavailableErrorCode)
	}
	sessionHash := utilsCrypto.GetHash(tokenId.String())
	keyCache := utilsCache.GetKeyUserAccessTokenIsActive(sessionHash)
	var valCache = map[string]string{
		"user_id": response.UserID,
		"role":    strconv.Itoa(int(domainModel.RoleUser)),
	}
	if err := cacheDistributed.SetTTL(
		ctx,
		keyCache,
		valCache,
		constants.TTL_ACCESS_TOKEN,
	); err != nil {
		global.Logger.Error("Error setting session in cache: ", err)
		return nil, errors.GetError(errors.SystemTemporaryUnavailableErrorCode)
	}
	// Return session
	return &applicationModel.LoginOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// LoginAdmin implements service.ICoreAuthService.
func (c *CoreAuthService) LoginAdmin(ctx context.Context, input *applicationModel.LoginInputAdmin) (*applicationModel.LoginOutput, *errors.Error) {
	// Get info user with mail
	domainRepo, err := domainRepository.GetUserRepository()
	if err != nil {
		global.Logger.Error("Error getting user repository: ", err)
		return nil, errors.GetError(errors.SystemTemporaryUnavailableErrorCode)
	}
	response, err := domainRepo.GetUserBaseByEmail(ctx, input.UserName)
	if err != nil {
		global.Logger.Error("Error getting user by email: ", err)
		return nil, errors.GetError(errors.SystemTemporaryUnavailableErrorCode)
	}
	if response == nil {
		// User not found
		return nil, errors.GetError(errors.UserNotFoundErrorCode)
	}
	// Check role
	if response.Role != domainModel.RoleAdmin {
		// Role not match
		return nil, errors.GetError(errors.UserNotFoundErrorCode)
	}
	// Check password
	if !utilsCrypto.ComparePasswordWithHash(
		input.Password,
		response.UserSalt,
		response.UserPassword,
	) {
		// Password not match
		return nil, errors.GetError(errors.UserPasswordIncorrectErrorCode)
	}
	// Create session
	tokenService := domainToken.GetTokenService()
	tokenId := utilsRandom.GenerateUUID()
	// Create access token
	timeTtlAccessToken := time.Duration(constants.TTL_ACCESS_TOKEN) * time.Second
	accessToken, err := tokenService.CreateUserToken(
		ctx,
		&domainModel.TokenUserJwtInput{
			UserId:  response.UserID,
			TokenId: tokenId.String(),
			Role:    domainModel.RoleAdmin,
			Expires: time.Now().Add(timeTtlAccessToken),
		},
	)
	if err != nil {
		global.Logger.Error("Error creating user token: ", err)
		return nil, errors.GetError(errors.SystemTemporaryUnavailableErrorCode)
	}
	// Create refresh token
	timeTtlRefreshToken := time.Duration(constants.TTL_REFRESH_TOKEN) * time.Second
	refreshToken, err := tokenService.CreateUserRefreshToken(
		ctx,
		&domainModel.TokenUserRefreshInput{
			UserId:  response.UserID,
			TokenId: tokenId.String(),
			Expires: time.Now().Add(timeTtlRefreshToken),
		},
	)
	if err != nil {
		global.Logger.Error("Error creating user refresh token: ", err)
		return nil, errors.GetError(errors.SystemTemporaryUnavailableErrorCode)
	}
	// Save session to db and cache
	uuidUser, _ := utilsUuid.ParseUUID(response.UserID)
	ipAddr, _ := netip.ParseAddr(input.ClientIp)
	if err := domainRepo.CreateUserSession(
		ctx,
		&domainModel.CreateUserSessionInput{
			SessionID:    tokenId,
			UserID:       uuidUser,
			IPAddress:    ipAddr,
			UserAgent:    input.UserAgent,
			RefreshToken: refreshToken,
			ExpiredAt:    time.Now().Add(timeTtlRefreshToken),
		},
	); err != nil {
		global.Logger.Error("Error creating user session: ", err)
		return nil, errors.GetError(errors.SystemTemporaryUnavailableErrorCode)
	}
	cacheDistributed, err := domainCache.GetDistributedCache()
	if err != nil {
		global.Logger.Error("Error getting distributed cache: ", err)
		return nil, errors.GetError(errors.SystemTemporaryUnavailableErrorCode)
	}
	sessionHash := utilsCrypto.GetHash(tokenId.String())
	var valCache = map[string]string{
		"user_id": response.UserID,
		"role":    strconv.Itoa(int(domainModel.RoleAdmin)),
	}
	keyCache := utilsCache.GetKeyUserAccessTokenIsActive(sessionHash)
	if err := cacheDistributed.SetTTL(
		ctx,
		keyCache,
		valCache,
		constants.TTL_ACCESS_TOKEN,
	); err != nil {
		global.Logger.Error("Error setting session in cache: ", err)
		return nil, errors.GetError(errors.SystemTemporaryUnavailableErrorCode)
	}
	// Return session
	return &applicationModel.LoginOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Logout implements service.ICoreAuthService.
func (c *CoreAuthService) Logout(ctx context.Context, input *applicationModel.LogoutInput) *errors.Error {
	// Get user repository
	userRepo, err := domainRepository.GetUserRepository()
	if err != nil {
		global.Logger.Error("Error getting user repository: ", err)
		return errors.GetError(errors.SystemTemporaryUnavailableErrorCode)
	}
	// Remove session from db
	if err := userRepo.RemoveUserSession(
		ctx,
		&domainModel.RemoveUserSessionInput{
			SessionID: input.SessionId,
		},
	); err != nil {
		global.Logger.Warn("Error removing user session: ", err)
		return errors.GetError(errors.SystemTemporaryUnavailableErrorCode)
	}
	// Remove session from cache
	cacheDistributed, err := domainCache.GetDistributedCache()
	if err != nil {
		global.Logger.Error("Error getting distributed cache: ", err)
		return errors.GetError(errors.SystemTemporaryUnavailableErrorCode)
	}
	keyCache := utilsCache.GetKeyUserAccessTokenIsActive(utilsCrypto.GetHash(input.SessionId.String()))
	if err := cacheDistributed.Delete(ctx, keyCache); err != nil {
		global.Logger.Error("Error removing session from cache: ", err)
		return errors.GetError(errors.SystemTemporaryUnavailableErrorCode)
	}
	return nil
}

// RefreshToken implements service.ICoreAuthService.
func (c *CoreAuthService) RefreshToken(ctx context.Context, input *applicationModel.RefreshTokenInput) (*applicationModel.RefreshTokenOutput, *errors.Error) {
	// Validate token refresh
	tokenService := domainToken.GetTokenService()
	_, tkErr := tokenService.ParseUserRefreshToken(
		ctx,
		input.RefreshToken,
	)
	if tkErr != nil {
		if tkErr.Code == domainError.TokenExpiredErrorCode {
			// Token expired
			return nil, errors.GetError(errors.TokenExpiredErrorCode)
		}
		if tkErr.Code == domainError.TokenMalformedErrorCode || tkErr.Code == domainError.TokenSignatureInvalidErrCode {
			// Token invalid
			return nil, errors.GetError(errors.AuthTokenInvalidErrorCode)
		}
		// Other error
		global.Logger.Error("Error parsing user refresh token: ", tkErr.Message)
		return nil, errors.GetError(errors.SystemTemporaryUnavailableErrorCode)
	}
	// Get data session from db
	domainRepo, err := domainRepository.GetUserRepository()
	if err != nil {
		global.Logger.Error("Error getting user repository: ", err)
		return nil, errors.GetError(errors.SystemTemporaryUnavailableErrorCode)
	}
	sessionData, err := domainRepo.GetUserSessionByID(ctx, input.SessionId)
	if err != nil {
		global.Logger.Error("Error getting user session by ID: ", err)
		return nil, errors.GetError(errors.SystemTemporaryUnavailableErrorCode)
	}
	if sessionData == nil {
		// Session not found
		return nil, errors.GetError(errors.AuthCannotRefreshTokenErrorCode)
	}
	if sessionData.ExpiredAt.Before(time.Now()) {
		// Session expired
		return nil, errors.GetError(errors.AuthCannotRefreshTokenErrorCode)
	}
	if sessionData.RefreshToken != input.RefreshToken {
		// Refresh token not match
		return nil, errors.GetError(errors.AuthCannotRefreshTokenErrorCode)
	}
	// Create new access token and refresh token
	accessTokenTimeTtl := time.Duration(constants.TTL_ACCESS_TOKEN) * time.Second
	accessToken, err := tokenService.CreateUserToken(
		ctx,
		&domainModel.TokenUserJwtInput{
			UserId:  input.UserId.String(),
			TokenId: input.SessionId.String(),
			Role:    input.UserRole,
			Expires: time.Now().Add(accessTokenTimeTtl),
		},
	)
	if err != nil {
		global.Logger.Error("Error creating user token: ", err)
		return nil, errors.GetError(errors.SystemTemporaryUnavailableErrorCode)
	}
	refreshTokenTimeTtl := time.Duration(constants.TTL_REFRESH_TOKEN) * time.Second
	refreshTokenNew, err := tokenService.CreateUserRefreshToken(
		ctx,
		&domainModel.TokenUserRefreshInput{
			UserId:  input.UserId.String(),
			TokenId: input.SessionId.String(),
			Expires: time.Now().Add(refreshTokenTimeTtl),
		},
	)
	if err != nil {
		global.Logger.Error("Error creating user token: ", err)
		return nil, errors.GetError(errors.SystemTemporaryUnavailableErrorCode)
	}
	// Save new session to db
	if err := domainRepo.RefreshSession(
		ctx,
		&domainModel.RefreshSessionInput{
			SessionID:    input.SessionId,
			RefreshToken: refreshTokenNew,
			ExpiredAt:    time.Now().Add(refreshTokenTimeTtl),
		},
	); err != nil {
		global.Logger.Error("Error refreshing user session: ", err)
		return nil, errors.GetError(errors.SystemTemporaryUnavailableErrorCode)
	}
	// Save new session to cache - Use to block access token
	cacheDistributed, err := domainCache.GetDistributedCache()
	if err != nil {
		global.Logger.Error("Error getting distributed cache: ", err)
		return nil, errors.GetError(errors.SystemTemporaryUnavailableErrorCode)
	}
	sessionHash := utilsCrypto.GetHash(input.SessionId.String())
	keyCache := utilsCache.GetKeyUserAccessTokenIsActive(sessionHash)
	var valCache = map[string]string{
		"user_id": input.UserId.String(),
		"role":    strconv.Itoa(int(input.UserRole)),
	}
	if err := cacheDistributed.SetTTL(
		ctx,
		keyCache,
		valCache,
		constants.TTL_ACCESS_TOKEN,
	); err != nil {
		global.Logger.Error("Error setting session in cache: ", err)
		return nil, errors.GetError(errors.SystemTemporaryUnavailableErrorCode)
	}
	// Return new tokens
	return &applicationModel.RefreshTokenOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenNew,
	}, nil
}

/**
 * NewCoreAuthService creates a new instance of CoreAuthService
 */
func NewCoreAuthService() service.ICoreAuthService {
	return &CoreAuthService{}
}
