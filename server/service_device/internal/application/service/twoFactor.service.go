package service

import (
	"errors"
	// errorService "github.com/youknow2509/cio_verify_face/server/service_device/internal/application/errors"
)

// =======================================================
//
//	Define interfaces for Auth service
//
// =======================================================
type (
	// Two-Factor Authentication
	ITwoFactorAuthService interface {
	}
)

// =======================================================
//
//	Variables instance interfaces for Auth service
//
// =======================================================
var (
	_ITwoFactorAuthService ITwoFactorAuthService
)

// =======================================================
//
//	Getter, setter for Auth service interfaces
//
// =======================================================
func GetAuthTwoFactorAuthService() ITwoFactorAuthService {
	return _ITwoFactorAuthService
}

func SetAuthTwoFactorAuthService(s ITwoFactorAuthService) error {
	if _ITwoFactorAuthService != nil {
		return errors.New("auth Two Factor Authentication service is already set")
	}
	_ITwoFactorAuthService = s
	return nil
}
