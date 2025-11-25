package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/domain/model"
)

// =================================
// Face Profile Update Request Repository Interface:
// =================================
type IFaceProfileUpdateRequestRepository interface {
	// Create a new face profile update request
	CreateRequest(ctx context.Context, req *model.FaceProfileUpdateRequest) error

	// Get request by ID
	GetRequestByID(ctx context.Context, requestID, companyID uuid.UUID) (*model.FaceProfileUpdateRequest, error)

	// Get request by update token
	GetRequestByToken(ctx context.Context, token string) (*model.FaceProfileUpdateRequest, error)

	// Get pending requests for a company
	GetPendingRequestsByCompany(ctx context.Context, companyID uuid.UUID, limit, offset int) ([]*model.FaceProfileUpdateRequest, error)

	// Get requests by user and month
	GetRequestsByUserAndMonth(ctx context.Context, userID uuid.UUID, month string) ([]*model.FaceProfileUpdateRequest, error)

	// Count requests by user in a month
	CountRequestsByUserInMonth(ctx context.Context, userID uuid.UUID, month string) (int, error)

	// Update request status
	UpdateRequestStatus(ctx context.Context, requestID, companyID uuid.UUID, status model.RequestStatus) error

	// Approve request (sets approved_by, approved_at, update_token, update_link_expires_at)
	ApproveRequest(ctx context.Context, requestID, companyID, approvedBy uuid.UUID, updateToken string, expiresAt interface{}) error

	// Reject request
	RejectRequest(ctx context.Context, requestID, companyID, rejectedBy uuid.UUID, reason string) error

	// Complete request (after face profile is updated)
	CompleteRequest(ctx context.Context, requestID, companyID uuid.UUID) error

	// Mark expired requests
	MarkExpiredRequests(ctx context.Context) (int64, error)

	// Check if user has pending request
	HasPendingRequest(ctx context.Context, userID, companyID uuid.UUID) (bool, error)
}

// =================================
// Password Reset Request Repository Interface:
// =================================
type IPasswordResetRequestRepository interface {
	// Create a new password reset request
	CreateRequest(ctx context.Context, req *model.PasswordResetRequest) error

	// Get request by ID
	GetRequestByID(ctx context.Context, requestID uuid.UUID) (*model.PasswordResetRequest, error)

	// Get recent requests by manager for a user (for spam check)
	GetRecentRequestsByManagerForUser(ctx context.Context, managerID, userID uuid.UUID, since interface{}) ([]*model.PasswordResetRequest, error)

	// Update request status
	UpdateRequestStatus(ctx context.Context, requestID uuid.UUID, status model.PasswordResetStatus, kafkaMessageID string) error

	// Count requests by manager in time window (for spam prevention)
	CountRequestsByManagerInWindow(ctx context.Context, managerID uuid.UUID, since interface{}) (int, error)
}

// =================================
// User Repository Interface:
// =================================
type IUserRepository interface {
	// Get user by ID
	GetUserByID(ctx context.Context, userID uuid.UUID) (*model.UserInfo, error)

	// Get user by email
	GetUserByEmail(ctx context.Context, email string) (*model.UserInfo, error)

	// Check if user belongs to company
	UserBelongsToCompany(ctx context.Context, userID, companyID uuid.UUID) (bool, error)

	// Get employee info
	GetEmployeeInfo(ctx context.Context, employeeID uuid.UUID) (*model.EmployeeInfo, error)

	// Update user password
	UpdateUserPassword(ctx context.Context, userID uuid.UUID, salt, passwordHash string) error

	// Check if user is company admin
	IsCompanyAdmin(ctx context.Context, userID, companyID uuid.UUID) (bool, error)
}

// =================================
// Repository Variables:
// =================================
var (
	_faceProfileUpdateRequestRepository IFaceProfileUpdateRequestRepository
	_passwordResetRequestRepository     IPasswordResetRequestRepository
	_userRepository                     IUserRepository
)

// =================================
// Setters and Getters:
// =================================
func SetFaceProfileUpdateRequestRepository(repo IFaceProfileUpdateRequestRepository) error {
	if _faceProfileUpdateRequestRepository != nil {
		return errors.New("face profile update request repository already initialized")
	}
	_faceProfileUpdateRequestRepository = repo
	return nil
}

func GetFaceProfileUpdateRequestRepository() (IFaceProfileUpdateRequestRepository, error) {
	if _faceProfileUpdateRequestRepository == nil {
		return nil, errors.New("face profile update request repository not initialized")
	}
	return _faceProfileUpdateRequestRepository, nil
}

func SetPasswordResetRequestRepository(repo IPasswordResetRequestRepository) error {
	if _passwordResetRequestRepository != nil {
		return errors.New("password reset request repository already initialized")
	}
	_passwordResetRequestRepository = repo
	return nil
}

func GetPasswordResetRequestRepository() (IPasswordResetRequestRepository, error) {
	if _passwordResetRequestRepository == nil {
		return nil, errors.New("password reset request repository not initialized")
	}
	return _passwordResetRequestRepository, nil
}

func SetUserRepository(repo IUserRepository) error {
	if _userRepository != nil {
		return errors.New("user repository already initialized")
	}
	_userRepository = repo
	return nil
}

func GetUserRepository() (IUserRepository, error) {
	if _userRepository == nil {
		return nil, errors.New("user repository not initialized")
	}
	return _userRepository, nil
}
