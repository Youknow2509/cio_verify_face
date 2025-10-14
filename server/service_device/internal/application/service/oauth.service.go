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
	// OAuth & Social Login
	IOAuthService interface {
	}
)

// =======================================================
//
//	Variables instance interfaces for Auth service
//
// =======================================================
var (
	_IOAuthService IOAuthService
)

// =======================================================
//
//	Getter, setter for Auth service interfaces
//
// =======================================================
func GetAuthOAuthService() IOAuthService {
	return _IOAuthService
}

func SetAuthOAuthService(s IOAuthService) error {
	if _IOAuthService != nil {
		return errors.New("auth OAuth service is already set")
	}
	_IOAuthService = s
	return nil
}
