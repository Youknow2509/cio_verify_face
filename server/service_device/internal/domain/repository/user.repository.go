package repository

import (
	"errors"

	"context"

	"github.com/youknow2509/cio_verify_face/server/service_device/internal/domain/model"
)

/**
 * Interface for user repository
 */
type IUserRepository interface {
	UserPermissionDevice(ctx context.Context, input *model.UserPermissionDeviceInput) (bool, error)
	UserExistsInCompany(ctx context.Context, input *model.UserExistsInCompanyInput) (bool, error)
	GetCompanyIdOfUser(ctx context.Context, input *model.GetCompanyIdOfUserInput) (*model.GetCompanyIdOfUserOutput, error)
}

/**
 * Variable for user repository instance
 */
var _vUserRepository IUserRepository

/**
 * Set the user repository instance
 */
func SetUserRepository(v IUserRepository) error {
	if _vUserRepository != nil {
		return errors.New("user repository initialization failed, not nil")
	}
	_vUserRepository = v
	return nil
}

/**
 * Get the user repository instance
 */
func GetUserRepository() (IUserRepository, error) {
	if _vUserRepository == nil {
		return nil, errors.New("user repository not initialized")
	}
	return _vUserRepository, nil
}
