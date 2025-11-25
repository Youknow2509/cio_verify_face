package service

import (
	"context"
	"errors"

	appErrors "github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/application/errors"
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/application/model"
)

// =================================
// Face Profile Update Service Interface:
// =================================
type IFaceProfileUpdateService interface {
	// Employee: Create a new face profile update request
	CreateUpdateRequest(ctx context.Context, input *model.CreateFaceProfileUpdateRequestInput) (*model.CreateFaceProfileUpdateRequestOutput, *appErrors.Error)

	// Employee: Get their own update requests
	GetMyUpdateRequests(ctx context.Context, input *model.GetMyUpdateRequestsInput) (*model.GetMyUpdateRequestsOutput, *appErrors.Error)

	// Manager: Get pending requests for their company
	GetPendingRequests(ctx context.Context, input *model.GetPendingRequestsInput) (*model.GetPendingRequestsOutput, *appErrors.Error)

	// Manager: Approve a face profile update request
	ApproveRequest(ctx context.Context, input *model.ApproveRequestInput) (*model.ApproveRequestOutput, *appErrors.Error)

	// Manager: Reject a face profile update request
	RejectRequest(ctx context.Context, input *model.RejectRequestInput) (*model.RejectRequestOutput, *appErrors.Error)

	// Public: Validate update token
	ValidateUpdateToken(ctx context.Context, input *model.ValidateUpdateTokenInput) (*model.ValidateUpdateTokenOutput, *appErrors.Error)

	// Employee: Update face profile using valid token
	UpdateFaceProfile(ctx context.Context, input *model.UpdateFaceProfileInput) (*model.UpdateFaceProfileOutput, *appErrors.Error)
}

// =================================
// Password Reset Service Interface:
// =================================
type IPasswordResetService interface {
	// Manager: Reset an employee's password
	ResetEmployeePassword(ctx context.Context, input *model.ResetEmployeePasswordInput) (*model.ResetEmployeePasswordOutput, *appErrors.Error)
}

// =================================
// Service Variables:
// =================================
var (
	_faceProfileUpdateService IFaceProfileUpdateService
	_passwordResetService     IPasswordResetService
)

// =================================
// Setters and Getters:
// =================================
func SetFaceProfileUpdateService(svc IFaceProfileUpdateService) error {
	if _faceProfileUpdateService != nil {
		return errors.New("face profile update service already initialized")
	}
	_faceProfileUpdateService = svc
	return nil
}

func GetFaceProfileUpdateService() IFaceProfileUpdateService {
	return _faceProfileUpdateService
}

func SetPasswordResetService(svc IPasswordResetService) error {
	if _passwordResetService != nil {
		return errors.New("password reset service already initialized")
	}
	_passwordResetService = svc
	return nil
}

func GetPasswordResetService() IPasswordResetService {
	return _passwordResetService
}
