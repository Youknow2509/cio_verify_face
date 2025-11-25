package dto

// CreateFaceProfileUpdateRequestDTO represents the request body for creating a face profile update request
type CreateFaceProfileUpdateRequestDTO struct {
	Reason string `json:"reason" binding:"omitempty,max=500" example:"Need to update my face photo for better recognition"`
}

// GetMyRequestsQueryDTO represents query parameters for getting employee's own requests
type GetMyRequestsQueryDTO struct {
	Month  string `form:"month" binding:"omitempty,datetime=2006-01" example:"2024-11"`
	Limit  int    `form:"limit" binding:"omitempty,min=1,max=100" example:"50"`
	Offset int    `form:"offset" binding:"omitempty,min=0" example:"0"`
}

// GetPendingRequestsQueryDTO represents query parameters for getting pending requests
type GetPendingRequestsQueryDTO struct {
	Limit  int `form:"limit" binding:"omitempty,min=1,max=100" example:"50"`
	Offset int `form:"offset" binding:"omitempty,min=0" example:"0"`
}

// ApproveRequestParamDTO represents the path parameter for approving a request
type ApproveRequestParamDTO struct {
	RequestID string `uri:"id" binding:"required,uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
}

// RejectRequestParamDTO represents the path parameter for rejecting a request
type RejectRequestParamDTO struct {
	RequestID string `uri:"id" binding:"required,uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
}

// RejectRequestBodyDTO represents the request body for rejecting a request
type RejectRequestBodyDTO struct {
	Reason string `json:"reason" binding:"omitempty,max=500" example:"Photos are not clear enough"`
}

// ValidateTokenQueryDTO represents query parameters for validating a token
type ValidateTokenQueryDTO struct {
	Token string `form:"token" binding:"required,min=10" example:"abc123def456ghi789"`
}

// UpdateFaceProfileFormDTO represents the form data for updating face profile
type UpdateFaceProfileFormDTO struct {
	Token string `form:"token" binding:"required,min=10" example:"abc123def456ghi789"`
}

// Default values
const (
	DefaultLimit  = 50
	DefaultOffset = 0
)

// SetDefaults sets default values for GetMyRequestsQueryDTO
func (d *GetMyRequestsQueryDTO) SetDefaults() {
	if d.Limit == 0 {
		d.Limit = DefaultLimit
	}
	if d.Offset < 0 {
		d.Offset = DefaultOffset
	}
}

// SetDefaults sets default values for GetPendingRequestsQueryDTO
func (d *GetPendingRequestsQueryDTO) SetDefaults() {
	if d.Limit == 0 {
		d.Limit = DefaultLimit
	}
	if d.Offset < 0 {
		d.Offset = DefaultOffset
	}
}
