package model

import (
	"time"

	"github.com/google/uuid"
)

// =================================
// Role definitions (same as other services):
// =================================
type Role int

const (
	RoleSystemAdmin  Role = 0 // System admin/root with full access
	RoleCompanyAdmin Role = 1 // Company admin/manager with elevated access
	RoleEmployee     Role = 2 // Regular employee with limited access
)

// =================================
// Face Profile Update Request Status:
// =================================
type RequestStatus int

const (
	RequestStatusPending   RequestStatus = 0
	RequestStatusApproved  RequestStatus = 1
	RequestStatusRejected  RequestStatus = 2
	RequestStatusExpired   RequestStatus = 3
	RequestStatusCompleted RequestStatus = 4
)

// =================================
// Password Reset Request Status:
// =================================
type PasswordResetStatus int

const (
	PasswordResetStatusPending PasswordResetStatus = 0
	PasswordResetStatusSent    PasswordResetStatus = 1
	PasswordResetStatusFailed  PasswordResetStatus = 2
)

// =================================
// Face Profile Update Request Model:
// =================================
type FaceProfileUpdateRequest struct {
	RequestID            uuid.UUID              `json:"request_id"`
	UserID               uuid.UUID              `json:"user_id"`
	CompanyID            uuid.UUID              `json:"company_id"`
	Status               RequestStatus          `json:"status"`
	RequestMonth         string                 `json:"request_month"` // Format: YYYY-MM
	RequestCountInMonth  int                    `json:"request_count_in_month"`
	UpdateToken          *string                `json:"update_token,omitempty"`
	UpdateLinkExpiresAt  *time.Time             `json:"update_link_expires_at,omitempty"`
	ApprovedBy           *uuid.UUID             `json:"approved_by,omitempty"`
	ApprovedAt           *time.Time             `json:"approved_at,omitempty"`
	RejectionReason      *string                `json:"rejection_reason,omitempty"`
	Reason               *string                `json:"reason,omitempty"`
	MetaData             map[string]interface{} `json:"meta_data"`
	CreatedAt            time.Time              `json:"created_at"`
	UpdatedAt            time.Time              `json:"updated_at"`
}

// =================================
// Password Reset Request Model:
// =================================
type PasswordResetRequest struct {
	RequestID      uuid.UUID              `json:"request_id"`
	UserID         uuid.UUID              `json:"user_id"`
	CompanyID      *uuid.UUID             `json:"company_id,omitempty"`
	RequestedBy    uuid.UUID              `json:"requested_by"`
	Status         PasswordResetStatus    `json:"status"`
	KafkaMessageID *string                `json:"kafka_message_id,omitempty"`
	KafkaSentAt    *time.Time             `json:"kafka_sent_at,omitempty"`
	MetaData       map[string]interface{} `json:"meta_data"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

// =================================
// User Info Model (for authorization):
// =================================
type UserInfo struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	FullName  string    `json:"full_name"`
	AvatarURL string    `json:"avatar_url"`
	Role      Role      `json:"role"`
	CompanyID *uuid.UUID `json:"company_id,omitempty"`
}

// =================================
// Employee Info Model:
// =================================
type EmployeeInfo struct {
	EmployeeID   uuid.UUID `json:"employee_id"`
	CompanyID    uuid.UUID `json:"company_id"`
	EmployeeCode string    `json:"employee_code"`
	Department   string    `json:"department"`
	Position     string    `json:"position"`
}
