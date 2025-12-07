package grpc

import (
"context"
"time"
)

// =================================
// Face Service Domain Interfaces
// =================================

// IFaceServiceClient defines the interface for interacting with the AI face service
type IFaceServiceClient interface {
// EnrollFace enrolls a new face profile for a user
EnrollFace(ctx context.Context, req *EnrollFaceRequest) (*EnrollFaceResponse, error)

// DeleteProfile deletes a face profile
DeleteProfile(ctx context.Context, req *DeleteProfileRequest) (*DeleteProfileResponse, error)

// GetUserProfiles gets all face profiles for a user
GetUserProfiles(ctx context.Context, req *GetUserProfilesRequest) (*GetUserProfilesResponse, error)

// UpdateProfile updates an existing face profile
UpdateProfile(ctx context.Context, req *UpdateProfileRequest) (*UpdateProfileResponse, error)

// Close closes the gRPC connection
Close() error
}

// =================================
// Request/Response Models
// =================================

type EnrollFaceRequest struct {
ImageData   []byte
UserID      string
CompanyID   string
DeviceID    string
MakePrimary bool
Filename    string
}

type EnrollFaceResponse struct {
Status        string
Message       string
ProfileID     string
QualityScore  float32
DuplicateIDs  []string
}

type DeleteProfileRequest struct {
ProfileID  string
CompanyID  string
HardDelete bool
}

type DeleteProfileResponse struct {
Status  string
Message string
}

type UpdateProfileRequest struct {
ProfileID   string
CompanyID   string
ImageData   []byte
MakePrimary *bool
Filename    string
}

type UpdateProfileResponse struct {
Status  string
Message string
}

type GetUserProfilesRequest struct {
UserID     string
CompanyID  string
PageNumber int32
PageSize   int32
}

type FaceProfile struct {
ProfileID         string
UserID            string
CompanyID         string
EmbeddingVersion  string
IsPrimary         bool
CreatedAt         time.Time
UpdatedAt         time.Time
DeletedAt         *time.Time
QualityScore      *float32
}

type GetUserProfilesResponse struct {
Profiles []*FaceProfile
}

// =================================
// Global Service Instance
// =================================

var _faceServiceClient IFaceServiceClient

// SetFaceServiceClient sets the face service client
func SetFaceServiceClient(client IFaceServiceClient) {
_faceServiceClient = client
}

// GetFaceServiceClient gets the face service client
func GetFaceServiceClient() IFaceServiceClient {
return _faceServiceClient
}

// =================================
// Status Constants
// =================================
const (
FaceServiceStatusSuccess = "ok"
FaceServiceStatusFailed  = "failed"
FaceServiceStatusError   = "error"
)
