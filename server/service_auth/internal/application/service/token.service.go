package service

import (
	"context"
	"errors"

	"github.com/youknow2509/cio_verify_face/server/service_auth/internal/application/model"
)

/**
 * Token application service
 */
type ITokenService interface {
	// User token operations
	CreateUserToken(ctx context.Context, input model.CreateTokenUserInput) (*model.CreateTokenUserOutput, error)
	BlockTokenUser(ctx context.Context, input model.BlockTokenUserInput) error
	BlockTokenUserRefresh(ctx context.Context, input model.BlockTokenUserRefreshInput) error
	ParseTokenUser(ctx context.Context, input model.ParseTokenUserInput) (*model.ParseTokenUserOutput, error)
	RefreshTokenUser(ctx context.Context, input model.RefreshTokenUserInput) (*model.RefreshTokenUserOutput, error)

	// Device token operations
	CreateTokenDevice(ctx context.Context, input model.CreateTokenDeviceInput) (string, error)
	BlockTokenDevice(ctx context.Context, input model.BlockTokenDeviceInput) error
	CheckTokenDevice(ctx context.Context, input model.CheckTokenDeviceInput) (bool, string, error)
}

/**
 * Manager token instance
 */
var _vITokenService ITokenService

// Getter token service instance
func GetTokenService() ITokenService {
	return _vITokenService
}

// Setter token service instance
func SetTokenService(svc ITokenService) error {
	if svc == nil {
		return errors.New("cannot set nil token service")
	}
	if _vITokenService != nil {
		return errors.New("token service is already set")
	}
	_vITokenService = svc
	return nil
}
