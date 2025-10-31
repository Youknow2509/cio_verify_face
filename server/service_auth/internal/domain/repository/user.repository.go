package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/youknow2509/cio_verify_face/server/service_auth/internal/domain/model"
)

/**
 * Interface for user repository
 */
type IUserRepository interface {
	// ======================================
	// 			Use for core auth
	// ======================================
	// Get user info with id
	GetUserInfoByID(ctx context.Context, userID uuid.UUID) (*model.UserInfoOutput, error)
	// Get user base info by ID
	GetUserBaseByID(ctx context.Context, userID uuid.UUID) (*model.UserBaseInfoOutput, error)
	// Get user base by phone
	GetUserBaseByEmail(ctx context.Context, email string) (*model.UserBaseInfoOutput, error)
	// Create user session
	CreateUserSession(ctx context.Context, data *model.CreateUserSessionInput) error
	// Remove user session
	RemoveUserSession(ctx context.Context, data *model.RemoveUserSessionInput) error
	// Get user session by ID
	GetUserSessionByID(ctx context.Context, sessionID uuid.UUID) (*model.UserSessionOutput, error)
	// Refresh user session
	RefreshSession(ctx context.Context, data *model.RefreshSessionInput) error
	// v.v

	// ======================================================
	//
	// ======================================================

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
