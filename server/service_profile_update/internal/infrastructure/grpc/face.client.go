package grpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	domainConfig "github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/domain/config"
	domainGrpc "github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/domain/grpc"
	facepb "github.com/youknow2509/cio_verify_face/server/service_profile_update/proto"
)

// FaceServiceClient implements domainGrpc.IFaceServiceClient over gRPC.
type FaceServiceClient struct {
	conn   *grpc.ClientConn
	client facepb.FaceVerificationServiceClient
}

// NewFaceServiceClient creates a new face service gRPC client.
func NewFaceServiceClient(cfg *domainConfig.ServiceFaceSetting) (*FaceServiceClient, error) {
	if cfg == nil {
		return nil, fmt.Errorf("face service config is nil")
	}
	if !cfg.Enabled {
		return nil, fmt.Errorf("face service is disabled")
	}

	var creds credentials.TransportCredentials
	if cfg.TLS.Enabled {
		tlsConfig := &tls.Config{}
		if cfg.TLS.CertFile != "" && cfg.TLS.KeyFile != "" {
			cert, err := tls.LoadX509KeyPair(cfg.TLS.CertFile, cfg.TLS.KeyFile)
			if err != nil {
				return nil, fmt.Errorf("failed to load face service TLS cert: %w", err)
			}
			tlsConfig.Certificates = []tls.Certificate{cert}
		}
		creds = credentials.NewTLS(tlsConfig)
	} else {
		creds = insecure.NewCredentials()
	}

	kacp := keepalive.ClientParameters{
		Time:                time.Duration(cfg.KeepaliveTimeMs) * time.Millisecond,
		Timeout:             time.Duration(cfg.KeepaliveTimeoutMs) * time.Millisecond,
		PermitWithoutStream: cfg.KeepalivePermitWithoutCalls,
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		grpc.WithKeepaliveParams(kacp),
	}

	conn, err := grpc.Dial(cfg.GrpcAddr, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to face service: %w", err)
	}

	client := facepb.NewFaceVerificationServiceClient(conn)

	return &FaceServiceClient{
		conn:   conn,
		client: client,
	}, nil
}

// EnrollFace enrolls a new profile for a user.
func (c *FaceServiceClient) EnrollFace(ctx context.Context, req *domainGrpc.EnrollFaceRequest) (*domainGrpc.EnrollFaceResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("enroll request is nil")
	}

	pbReq := &facepb.EnrollRequest{
		ImageData:   req.ImageData,
		UserId:      req.UserID,
		CompanyId:   req.CompanyID,
		MakePrimary: req.MakePrimary,
		Filename:    req.Filename,
	}

	if req.DeviceID != "" {
		pbReq.DeviceId = &req.DeviceID
	}

	resp, err := c.client.EnrollFace(ctx, pbReq)
	if err != nil {
		return nil, err
	}

	dupIDs := make([]string, 0, len(resp.DuplicateProfiles))
	for _, dup := range resp.DuplicateProfiles {
		dupIDs = append(dupIDs, dup.GetProfileId())
	}

	return &domainGrpc.EnrollFaceResponse{
		Status:       resp.GetStatus(),
		Message:      resp.GetMessage(),
		ProfileID:    resp.GetProfileId(),
		QualityScore: resp.GetQualityScore(),
		DuplicateIDs: dupIDs,
	}, nil
}

// DeleteProfile deletes a profile.
func (c *FaceServiceClient) DeleteProfile(ctx context.Context, req *domainGrpc.DeleteProfileRequest) (*domainGrpc.DeleteProfileResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("delete profile request is nil")
	}

	pbReq := &facepb.DeleteProfileRequest{
		ProfileId:  req.ProfileID,
		CompanyId:  req.CompanyID,
		HardDelete: req.HardDelete,
	}

	resp, err := c.client.DeleteProfile(ctx, pbReq)
	if err != nil {
		return nil, err
	}

	return &domainGrpc.DeleteProfileResponse{
		Status:  resp.GetStatus(),
		Message: resp.GetMessage(),
	}, nil
}

// GetUserProfiles returns all profiles for a user.
func (c *FaceServiceClient) GetUserProfiles(ctx context.Context, req *domainGrpc.GetUserProfilesRequest) (*domainGrpc.GetUserProfilesResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("get profiles request is nil")
	}

	pbReq := &facepb.GetUserProfilesRequest{
		UserId:     req.UserID,
		CompanyId:  req.CompanyID,
		PageNumber: req.PageNumber,
		PageSize:   req.PageSize,
	}

	resp, err := c.client.GetUserProfiles(ctx, pbReq)
	if err != nil {
		return nil, err
	}

	profiles := make([]*domainGrpc.FaceProfile, 0, len(resp.Profiles))
	for _, p := range resp.Profiles {
		createdAt := time.Time{}
		if ts := p.GetCreatedAt(); ts != nil {
			createdAt = ts.AsTime()
		}

		updatedAt := time.Time{}
		if ts := p.GetUpdatedAt(); ts != nil {
			updatedAt = ts.AsTime()
		}

		var deletedAt *time.Time
		if ts := p.GetDeletedAt(); ts != nil {
			t := ts.AsTime()
			deletedAt = &t
		}

		var qualityScore *float32
		if p.QualityScore != nil {
			score := p.GetQualityScore()
			qualityScore = &score
		}

		profiles = append(profiles, &domainGrpc.FaceProfile{
			ProfileID:        p.GetProfileId(),
			UserID:           p.GetUserId(),
			CompanyID:        p.GetCompanyId(),
			EmbeddingVersion: p.GetEmbeddingVersion(),
			IsPrimary:        p.GetIsPrimary(),
			CreatedAt:        createdAt,
			UpdatedAt:        updatedAt,
			DeletedAt:        deletedAt,
			QualityScore:     qualityScore,
		})
	}

	return &domainGrpc.GetUserProfilesResponse{Profiles: profiles}, nil
}

// UpdateProfile updates an existing profile.
func (c *FaceServiceClient) UpdateProfile(ctx context.Context, req *domainGrpc.UpdateProfileRequest) (*domainGrpc.UpdateProfileResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("update profile request is nil")
	}

	pbReq := &facepb.UpdateProfileRequest{
		ProfileId:   req.ProfileID,
		CompanyId:   req.CompanyID,
		ImageData:   req.ImageData,
		MakePrimary: req.MakePrimary,
	}

	if req.Filename != "" {
		pbReq.Filename = &req.Filename
	}

	resp, err := c.client.UpdateProfile(ctx, pbReq)
	if err != nil {
		return nil, err
	}

	return &domainGrpc.UpdateProfileResponse{
		Status:  resp.GetStatus(),
		Message: resp.GetMessage(),
	}, nil
}

// Close closes the underlying gRPC connection.
func (c *FaceServiceClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
