package impl

import (
	"context"
	"time"

	"github.com/youknow2509/cio_verify_face/server/service_auth/internal/application/errors"
	applicationModel "github.com/youknow2509/cio_verify_face/server/service_auth/internal/application/model"
	"github.com/youknow2509/cio_verify_face/server/service_auth/internal/application/service"
	constants "github.com/youknow2509/cio_verify_face/server/service_auth/internal/constants"
	domainModel "github.com/youknow2509/cio_verify_face/server/service_auth/internal/domain/model"
	domainRepository "github.com/youknow2509/cio_verify_face/server/service_auth/internal/domain/repository"
	domainToken "github.com/youknow2509/cio_verify_face/server/service_auth/internal/domain/token"
	utilsCrypto "github.com/youknow2509/cio_verify_face/server/service_auth/internal/shared/utils/crypto"
	utilsRandom "github.com/youknow2509/cio_verify_face/server/service_auth/internal/shared/utils/random"
)

/**
 * Define CoreAuthService struct implementing
 */
type CoreAuthService struct{}

// CreateDeviceSession implements service.ICoreAuthService.
func (c *CoreAuthService) CreateDeviceSession(ctx context.Context, input *applicationModel.CreateDeviceSessionInput) (*applicationModel.CreateDeviceSessionOutput, *errors.Error) {
	return &applicationModel.CreateDeviceSessionOutput{
		// TODO: Add fields
	}, nil
}

// DeleteDeviceSession implements service.ICoreAuthService.
func (c *CoreAuthService) DeleteDeviceSession(ctx context.Context, input *applicationModel.DeleteDeviceSessionInput) *errors.Error {
	// TODO: Implement the delete device session logic here
	return nil
}

// GetMyInfo implements service.ICoreAuthService.
func (c *CoreAuthService) GetMyInfo(ctx context.Context, input *applicationModel.GetMyInfoInput) (*applicationModel.GetMyInfoOutput, *errors.Error) {
	// TODO: Implement the get my info logic here
	return &applicationModel.GetMyInfoOutput{
		// TODO: Add fields
	}, nil
}

// Login implements service.ICoreAuthService.
func (c *CoreAuthService) Login(ctx context.Context, input *applicationModel.LoginInput) (*applicationModel.LoginOutput, *errors.Error) {
	// Get info user with mail
	domainRepo, err := domainRepository.GetUserRepository()
	if err != nil {
		return nil, errors.GetError(errors.SystemTemporaryUnavailableErrorCode)
	}
	response, err := domainRepo.GetUserBaseByEmail(ctx, input.UserName)
	if err != nil {
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
	tokenId := utilsRandom.GenerateUUID().String()
	// Create access token 
	timeTtlAccessToken := time.Duration(constants.TTL_ACCESS_TOKEN) * time.Hour
	accessToken, err := tokenService.CreateUserToken(
		ctx,
		&domainModel.TokenUserJwtInput{
			UserId: response.UserID,
			TokenId: tokenId,
			Role: domainModel.RoleUser,
			Expires: time.Now().Add(timeTtlAccessToken),
		},
	)
	if err != nil {
		return nil, errors.GetError(errors.SystemTemporaryUnavailableErrorCode)
	}
	// Create refresh token
	timeTtlRefreshToken := time.Duration(constants.TTL_REFRESH_TOKEN) * time.Hour
	refreshToken, err := tokenService.CreateUserRefreshToken(
		ctx,
		&domainModel.TokenUserRefreshInput{
			UserId: response.UserID,
			TokenId: tokenId,
			Expires: time.Now().Add(timeTtlRefreshToken),
		},
	)
	if err != nil {
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
		return nil, errors.GetError(errors.SystemTemporaryUnavailableErrorCode)
	}
	response, err := domainRepo.GetUserBaseByEmail(ctx, input.UserName)
	if err != nil {
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
	tokenId := utilsRandom.GenerateUUID().String()
	// Create access token 
	timeTtlAccessToken := time.Duration(constants.TTL_ACCESS_TOKEN) * time.Hour
	accessToken, err := tokenService.CreateUserToken(
		ctx,
		&domainModel.TokenUserJwtInput{
			UserId: response.UserID,
			TokenId: tokenId,
			Role: domainModel.RoleAdmin,
			Expires: time.Now().Add(timeTtlAccessToken),
		},
	)
	if err != nil {
		return nil, errors.GetError(errors.SystemTemporaryUnavailableErrorCode)
	}
	// Create refresh token
	timeTtlRefreshToken := time.Duration(constants.TTL_REFRESH_TOKEN) * time.Hour
	refreshToken, err := tokenService.CreateUserRefreshToken(
		ctx,
		&domainModel.TokenUserRefreshInput{
			UserId: response.UserID,
			TokenId: tokenId,
			Expires: time.Now().Add(timeTtlRefreshToken),
		},
	)
	if err != nil {
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
	// TODO: Implement
	return nil
}

// RefreshDeviceSession implements service.ICoreAuthService.
func (c *CoreAuthService) RefreshDeviceSession(ctx context.Context, input *applicationModel.RefreshDeviceSessionInput) (*applicationModel.RefreshDeviceSessionOutput, *errors.Error) {
	// TODO: Implement
	return &applicationModel.RefreshDeviceSessionOutput{
		// TODO: Add fields
	}, nil
}

// RefreshToken implements service.ICoreAuthService.
func (c *CoreAuthService) RefreshToken(ctx context.Context, input *applicationModel.RefreshTokenInput) (*applicationModel.RefreshTokenOutput, *errors.Error) {
	// TODO: Implement
	return &applicationModel.RefreshTokenOutput{
		// TODO: Add fields
	}, nil
}

/**
 * NewCoreAuthService creates a new instance of CoreAuthService
 */
func NewCoreAuthService() service.ICoreAuthService {
	return &CoreAuthService{}
}

