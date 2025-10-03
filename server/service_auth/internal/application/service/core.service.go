package service

import (
	"context"
	"errors"

	errorService "github.com/youknow2509/cio_verify_face/server/service_auth/internal/application/errors"
	"github.com/youknow2509/cio_verify_face/server/service_auth/internal/application/model"
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
		// Create device token
		CreateDeviceSession(ctx context.Context, input *model.CreateDeviceSessionInput) (*model.CreateDeviceSessionOutput, *errorService.Error)
		// Delete device token
		DeleteDeviceSession(ctx context.Context, input *model.DeleteDeviceSessionInput) *errorService.Error	
		// Refresh device token
		RefreshDeviceSession(ctx context.Context, input *model.RefreshDeviceSessionInput) (*model.RefreshDeviceSessionOutput, *errorService.Error)
	}
)

// =======================================================
//
//	Variables instance interfaces for Auth service
//
// =======================================================
var (
	_ICoreAuthService ICoreAuthService
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
