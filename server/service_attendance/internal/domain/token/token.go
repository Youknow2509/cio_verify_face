package token

import (
	"context"
	"errors"

	domainErrors "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/errors"
	"github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/model"
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
	CreateUserToken(ctx context.Context, input *model.TokenUserJwtInput) (*model.UserTokenOutput, error)

	/**
	 * Create a device token - Device token use for device authentication
	 * @param ctx context.Context
	 * @param input *model.TokenDeviceJwtInput
	 * @return string, error - jwt token, error
	 */
	CreateDeviceToken(ctx context.Context, input *model.TokenDeviceJwtInput) (string, error)
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
	 * Check device token
	 * @param ctx context.Context
	 * @param token string - jwt token
	 * @return *model.TokenDeviceJwtOutput, *domainErrors.TokenValidationError
	 */
	CheckDeviceToken(ctx context.Context, token string) (bool, *domainErrors.TokenValidationError)
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
