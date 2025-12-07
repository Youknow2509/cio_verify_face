package face

import (
	"context"

	domainFace "github.com/youknow2509/cio_verify_face/server/service_device/internal/domain/face"
	domainModel "github.com/youknow2509/cio_verify_face/server/service_device/internal/domain/model"
	pb "github.com/youknow2509/cio_verify_face/server/service_device/proto"
)

// FaceVerificationService bridges to the AI face verification gRPC service.
type FaceVerificationService struct {
	client pb.FaceVerificationServiceClient
}

// VerifyFace calls the downstream gRPC service and maps the response to domain models.
func (s *FaceVerificationService) VerifyFace(ctx context.Context, input *domainModel.FaceVerifyInput) (*domainModel.FaceVerifyOutput, error) {
	req := &pb.VerifyRequest{
		ImageData:  input.ImageData,
		CompanyId:  input.CompanyId,
		SearchMode: input.SearchMode,
		TopK:       input.TopK,
	}
	if input.UserId != "" {
		req.UserId = &input.UserId
	}
	if input.DeviceId != "" {
		req.DeviceId = &input.DeviceId
	}
	resp, err := s.client.VerifyFace(ctx, req)
	if err != nil {
		return nil, err
	}
	matches := make([]domainModel.FaceVerifyMatch, 0, len(resp.Matches))
	for _, m := range resp.Matches {
		matches = append(matches, domainModel.FaceVerifyMatch{
			UserId:     m.UserId,
			ProfileId:  m.ProfileId,
			Similarity: m.Similarity,
			Confidence: m.Confidence,
			IsPrimary:  m.IsPrimary,
		})
	}
	var bestMatch *domainModel.FaceVerifyMatch
	if resp.BestMatch != nil {
		bestMatch = &domainModel.FaceVerifyMatch{
			UserId:     resp.BestMatch.UserId,
			ProfileId:  resp.BestMatch.ProfileId,
			Similarity: resp.BestMatch.Similarity,
			Confidence: resp.BestMatch.Confidence,
			IsPrimary:  resp.BestMatch.IsPrimary,
		}
	}
	var livenessScore *float32
	if resp.LivenessScore != nil {
		value := float32(resp.GetLivenessScore())
		livenessScore = &value
	}
	return &domainModel.FaceVerifyOutput{
		Status:        resp.Status,
		Verified:      resp.Verified,
		Matches:       matches,
		BestMatch:     bestMatch,
		Message:       resp.GetMessage(),
		LivenessScore: livenessScore,
	}, nil
}

// NewFaceVerificationService builds a face verification service implementation.
func NewFaceVerificationService(client pb.FaceVerificationServiceClient) domainFace.IFaceVerificationService {
	return &FaceVerificationService{client: client}
}
