package face

import (
	"context"
	"errors"

	domainModel "github.com/youknow2509/cio_verify_face/server/service_device/internal/domain/model"
)

// IFaceVerificationService defines the contract to interact with face verification backend.
type IFaceVerificationService interface {
	VerifyFace(ctx context.Context, input *domainModel.FaceVerifyInput) (*domainModel.FaceVerifyOutput, error)
}

var faceVerificationService IFaceVerificationService

// GetFaceVerificationService returns the current face verification service implementation.
func GetFaceVerificationService() IFaceVerificationService {
	return faceVerificationService
}

// SetFaceVerificationService sets the face verification service implementation once.
func SetFaceVerificationService(s IFaceVerificationService) error {
	if s == nil {
		return errors.New("face verification service cannot be nil")
	}
	if faceVerificationService != nil {
		return errors.New("face verification service already set")
	}
	faceVerificationService = s
	return nil
}
