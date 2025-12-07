package model

import (
	"time"

	"github.com/google/uuid"
	domainModel "github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/domain/model"
)

// =================================
// Session Info (from JWT token):
// =================================
type SessionInfo struct {
	UserID      string `json:"user_id"`
	CompanyID   string `json:"company_id"`
	Role        int32  `json:"role"`
	SessionID   string `json:"session_id"`
	ClientIP    string `json:"client_ip"`
	ClientAgent string `json:"client_agent"`
}

// =================================
// Face Profile Update Request Inputs/Outputs:
// =================================

// CreateFaceProfileUpdateRequestInput - Employee creates request to update face profile
type CreateFaceProfileUpdateRequestInput struct {
	Session *SessionInfo `json:"-"`
	Reason  string       `json:"reason,omitempty"` // Optional reason for request
}

type CreateFaceProfileUpdateRequestOutput struct {
	RequestID   string `json:"request_id"`
	Status      string `json:"status"`
	Message     string `json:"message"`
	CanRetry    bool   `json:"can_retry"`
	RetryAfter  int    `json:"retry_after_days,omitempty"` // Days until can retry
	RequestsRemaining int `json:"requests_remaining_this_month"`
}

// GetMyUpdateRequestsInput - Employee gets their own update requests
type GetMyUpdateRequestsInput struct {
	Session *SessionInfo `json:"-"`
	Month   string       `json:"month,omitempty"` // Format: YYYY-MM, optional
	Limit   int          `json:"limit,omitempty"`
	Offset  int          `json:"offset,omitempty"`
}

type GetMyUpdateRequestsOutput struct {
	Requests []*FaceProfileUpdateRequestDTO `json:"requests"`
	Total    int                            `json:"total"`
}

// GetPendingRequestsInput - Manager gets pending requests for their company
type GetPendingRequestsInput struct {
	Session *SessionInfo `json:"-"`
	Limit   int          `json:"limit,omitempty"`
	Offset  int          `json:"offset,omitempty"`
}

type GetPendingRequestsOutput struct {
	Requests []*FaceProfileUpdateRequestDTO `json:"requests"`
	Total    int                            `json:"total"`
}

// ApproveRequestInput - Manager approves a face profile update request
type ApproveRequestInput struct {
	Session   *SessionInfo `json:"-"`
	RequestID string       `json:"request_id" binding:"required"`
	CompanyID string       `json:"company_id,omitempty"` // Optional, defaults to session company
}

type ApproveRequestOutput struct {
	RequestID  string    `json:"request_id"`
	UpdateLink string    `json:"update_link"`
	ExpiresAt  time.Time `json:"expires_at"`
	Message    string    `json:"message"`
}

// RejectRequestInput - Manager rejects a face profile update request
type RejectRequestInput struct {
	Session   *SessionInfo `json:"-"`
	RequestID string       `json:"request_id" binding:"required"`
	CompanyID string       `json:"company_id,omitempty"` // Optional, defaults to session company
	Reason    string       `json:"reason,omitempty"`     // Optional rejection reason
}

type RejectRequestOutput struct {
	RequestID string `json:"request_id"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

// ValidateUpdateTokenInput - Validate update token before allowing face update
type ValidateUpdateTokenInput struct {
	Token string `json:"token" binding:"required"`
}

type ValidateUpdateTokenOutput struct {
	Valid     bool      `json:"valid"`
	UserID    string    `json:"user_id,omitempty"`
	CompanyID string    `json:"company_id,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	Message   string    `json:"message"`
}

// UpdateFaceProfileInput - Employee updates their face profile using update link
type UpdateFaceProfileInput struct {
	Token     string `json:"token" binding:"required"`
	ImageData []byte `json:"image_data" binding:"required"`
	Filename  string `json:"filename" binding:"required"`
}

type UpdateFaceProfileOutput struct {
	Success      bool    `json:"success"`
	ProfileID    string  `json:"profile_id,omitempty"`
	QualityScore float64 `json:"quality_score,omitempty"`
	Message      string  `json:"message"`
}

// =================================
// Password Reset Inputs/Outputs:
// =================================

// ResetEmployeePasswordInput - Manager resets an employee's password
type ResetEmployeePasswordInput struct {
	Session    *SessionInfo `json:"-"`
	EmployeeID string       `json:"employee_id" binding:"required"`
}

type ResetEmployeePasswordOutput struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// ConfirmPasswordResetInput - User confirms password reset using the link
type ConfirmPasswordResetInput struct {
	Token string `json:"token" binding:"required"`
}

type ConfirmPasswordResetOutput struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// =================================
// DTO for Face Profile Update Request:
// =================================
type FaceProfileUpdateRequestDTO struct {
	RequestID           string                 `json:"request_id"`
	UserID              string                 `json:"user_id"`
	UserName            string                 `json:"user_name,omitempty"`
	UserEmail           string                 `json:"user_email,omitempty"`
	CompanyID           string                 `json:"company_id"`
	Status              string                 `json:"status"`
	RequestMonth        string                 `json:"request_month"`
	RequestCountInMonth int                    `json:"request_count_in_month"`
	UpdateLink          string                 `json:"update_link,omitempty"`
	UpdateLinkExpiresAt *time.Time             `json:"update_link_expires_at,omitempty"`
	ApprovedBy          string                 `json:"approved_by,omitempty"`
	ApprovedAt          *time.Time             `json:"approved_at,omitempty"`
	RejectionReason     string                 `json:"rejection_reason,omitempty"`
	Reason              string                 `json:"reason,omitempty"`
	MetaData            map[string]interface{} `json:"meta_data,omitempty"`
	CreatedAt           time.Time              `json:"created_at"`
	UpdatedAt           time.Time              `json:"updated_at"`
}

// Convert domain model to DTO
func ToFaceProfileUpdateRequestDTO(req *domainModel.FaceProfileUpdateRequest, baseURL string) *FaceProfileUpdateRequestDTO {
	dto := &FaceProfileUpdateRequestDTO{
		RequestID:           req.RequestID.String(),
		UserID:              req.UserID.String(),
		CompanyID:           req.CompanyID.String(),
		Status:              getStatusString(req.Status),
		RequestMonth:        req.RequestMonth,
		RequestCountInMonth: req.RequestCountInMonth,
		CreatedAt:           req.CreatedAt,
		UpdatedAt:           req.UpdatedAt,
		MetaData:            req.MetaData,
	}

	if req.UpdateToken != nil && req.Status == domainModel.RequestStatusApproved {
		dto.UpdateLink = baseURL + "/api/v1/profile-update/face?token=" + *req.UpdateToken
	}
	if req.UpdateLinkExpiresAt != nil {
		dto.UpdateLinkExpiresAt = req.UpdateLinkExpiresAt
	}
	if req.ApprovedBy != nil {
		dto.ApprovedBy = req.ApprovedBy.String()
	}
	if req.ApprovedAt != nil {
		dto.ApprovedAt = req.ApprovedAt
	}
	if req.RejectionReason != nil {
		dto.RejectionReason = *req.RejectionReason
	}
	if req.Reason != nil {
		dto.Reason = *req.Reason
	}

	return dto
}

func getStatusString(status domainModel.RequestStatus) string {
	switch status {
	case domainModel.RequestStatusPending:
		return "pending"
	case domainModel.RequestStatusApproved:
		return "approved"
	case domainModel.RequestStatusRejected:
		return "rejected"
	case domainModel.RequestStatusExpired:
		return "expired"
	case domainModel.RequestStatusCompleted:
		return "completed"
	default:
		return "unknown"
	}
}

// =================================
// Helper to convert Session to UserInfo:
// =================================
func (s *SessionInfo) ToUserInfo() *domainModel.UserInfo {
	userID, _ := uuid.Parse(s.UserID)
	companyID, _ := uuid.Parse(s.CompanyID)
	return &domainModel.UserInfo{
		UserID:    userID,
		Role:      domainModel.Role(s.Role),
		CompanyID: &companyID,
	}
}
