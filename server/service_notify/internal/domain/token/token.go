package token

import (
	"context"
	"errors"

	domainErrors "github.com/youknow2509/cio_verify_face/server/service_notify/internal/domain/errors"
	"github.com/youknow2509/cio_verify_face/server/service_notify/internal/domain/model"
)

// ========================================
//
//	Token interface
//
// ========================================
type ITokenService interface {
	/**
	 * Create a user token - Access token
	 * @param ctx context.Context
	 * @param input *model.TokenUserJwtInput
	 * @return string, error - jwt token, error
	 */
	CreateUserToken(ctx context.Context, input *model.TokenUserJwtInput) (string, error)

	/**
	 * Create a service token - Service token use for internal communication
	 * @param ctx context.Context
	 * @param input *model.TokenServiceJwtInput
	 * @return string, error - jwt token, error
	 */
	CreateServiceToken(ctx context.Context, input *model.TokenServiceJwtInput) (string, error)

	/**
	 * Create a device token - Device token use for device authentication
	 * @param ctx context.Context
	 * @param input *model.TokenDeviceJwtInput
	 * @return string, error - jwt token, error
	 */
	CreateDeviceToken(ctx context.Context, input *model.TokenDeviceJwtInput) (string, error)

	/**
	 * Create a user refresh token
	 * @param ctx context.Context
	 * @param input *model.TokenUserRefreshInput
	 * @return string, error
	 */
	CreateUserRefreshToken(ctx context.Context, input *model.TokenUserRefreshInput) (string, error)

	/**
	 * Create a device refresh token
	 * @param ctx context.Context
	 * @param input *model.TokenDeviceRefreshInput
	 * @return string, error
	 */
	CreateDeviceRefreshToken(ctx context.Context, input *model.TokenDeviceRefreshInput) (string, error)

	/**
	 * Create a service refresh token
	 * @param ctx context.Context
	 * @param input *model.TokenServiceJwtInput
	 * @return string, error
	 */
	CreateServiceRefreshToken(ctx context.Context, input *model.ServiceRefreshTokenInput) (string, error)

	/**
	 * Parse a user token
	 * @param ctx context.Context
	 * @param token string - jwt token
	 * @return *model.TokenUserJwtOutput, *domainErrors.TokenValidationError
	 */
	ParseUserToken(ctx context.Context, token string) (*model.TokenUserJwtOutput, *domainErrors.TokenValidationError)

	/**
	 * Parse a user refresh token
	 * @param ctx context.Context
	 * @param token string - jwt token
	 * @return *model.TokenUserRefreshOutput, *domainErrors.TokenValidationError
	 */
	ParseUserRefreshToken(ctx context.Context, token string) (*model.TokenUserRefreshOutput, *domainErrors.TokenValidationError)

	/**
	 * Parse a service token
	 * @param ctx context.Context
	 * @param token string - jwt token
	 * @return *model.TokenServiceJwtOutput, *domainErrors.TokenValidationError
	 */
	ParseServiceToken(ctx context.Context, token string) (*model.TokenServiceJwtOutput, *domainErrors.TokenValidationError)

	/**
	 * Parse a device token
	 * @param ctx context.Context
	 * @param token string - jwt token
	 * @return *model.TokenDeviceJwtOutput, *domainErrors.TokenValidationError
	 */
	ParseDeviceToken(ctx context.Context, token string) (*model.TokenDeviceJwtOutput, *domainErrors.TokenValidationError)
	// v.v
}

/**
 * Variable save interface of ITokenService
 */
var _ITokenService ITokenService

// ================================================================================
//
//	Getter and setter for token service
//
// ================================================================================
// GetTokenService returns the current token service instance.
func GetTokenService() ITokenService {
	return _ITokenService
}

// SetTokenService sets the token service instance.
func SetTokenService(tokenService ITokenService) error {
	if tokenService == nil {
		return errors.New("token service cannot be nil")
	}
	if _ITokenService != nil {
		return errors.New("token service is already set, cannot be overwritten")
	}
	_ITokenService = tokenService
	return nil
}
