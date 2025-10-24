package model

import (
	"time"
)

// gRPC specific models for inter-service communication

// TokenValidationResult for gRPC token validation
type GRPCTokenValidationResult struct {
	Valid       bool      `json:"valid"`
	UserID      string    `json:"user_id,omitempty"`
	CompanyID   string    `json:"company_id,omitempty"`
	SessionID   string    `json:"session_id,omitempty"`
	ExpiresAt   time.Time `json:"expires_at,omitempty"`
	Permissions []string  `json:"permissions,omitempty"`
	Error       string    `json:"error,omitempty"`
}

// GRPCUserInfoOutput for gRPC user info retrieval
type GRPCUserInfoOutput struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CompanyID string    `json:"company_id"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// DeviceValidationResult for device session validation
type DeviceValidationResult struct {
	Valid     bool      `json:"valid"`
	DeviceID  string    `json:"device_id"`
	CompanyID string    `json:"company_id"`
	UserID    string    `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	Error     string    `json:"error,omitempty"`
}

// CompanyPermissionResult for permission checking
type CompanyPermissionResult struct {
	HasPermission bool     `json:"has_permission"`
	UserRole      string   `json:"user_role"`
	Permissions   []string `json:"permissions"`
}

// DeviceCompanyResult for device company validation
type DeviceCompanyResult struct {
	DeviceExists bool   `json:"device_exists"`
	IsActive     bool   `json:"is_active"`
	DeviceName   string `json:"device_name"`
	DeviceType   string `json:"device_type"`
}

// UserCompanyInfo for user company information
type UserCompanyInfo struct {
	CompanyID   string   `json:"company_id"`
	CompanyName string   `json:"company_name"`
	CompanyCode string   `json:"company_code"`
	IsActive    bool     `json:"is_active"`
	UserRole    string   `json:"user_role"`
	Permissions []string `json:"permissions"`
}
