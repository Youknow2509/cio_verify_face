package repository

import (
	"context"
	"errors"
	model "github.com/youknow2509/cio_verify_face/server/service_auth/internal/domain/model"
)

/**
 * Interface for company repository
 */
type ICompanyRepository interface {
	// Check company is management in company
	CheckUserIsManagementInCompany(ctx context.Context, data *model.CheckCompanyIsManagementInCompanyInput) (bool, error)
	// Check device exists in company
	CheckDeviceExistsInCompany(ctx context.Context, data *model.CheckDeviceExistsInCompanyInput) (bool, error)
	// Update device session
	UpdateDeviceSession(ctx context.Context, data *model.UpdateDeviceSessionInput) error
	// Delete device session
	DeleteDeviceSession(ctx context.Context, data *model.DeleteDeviceSessionInput) error
}

/**
 * Variable for Company repository instance
 */
var _vCompanyRepository ICompanyRepository

/**
 * Set the Company repository instance
 */
func SetCompanyRepository(v ICompanyRepository) error {
	if _vCompanyRepository != nil {
		return errors.New("Company repository initialization failed, not nil")
	}
	_vCompanyRepository = v
	return nil
}

/**
 * Get the Company repository instance
 */
func GetCompanyRepository() (ICompanyRepository, error) {
	if _vCompanyRepository == nil {
		return nil, errors.New("Company repository not initialized")
	}
	return _vCompanyRepository, nil
}
