package model

import (
	"time"

	"github.com/google/uuid"
)

/**
 * Token application model
 */
// CheckTokenDevice
type CheckTokenDeviceInput struct {
	DeviceId uuid.UUID `json:"device_id" validate:"required"`
	Token    string    `json:"token" validate:"required"`
}

// Create token user
type CreateTokenUserInput struct {
	UserId uuid.UUID `json:"user_id" validate:"required"`
	Role   int       `json:"role" validate:"required"`
}
type CreateTokenUserOutput struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Create device token
type CreateTokenDeviceInput struct {
	DeviceId  uuid.UUID `json:"device_id" validate:"required"`
	CompanyId uuid.UUID `json:"company_id" validate:"required"`
}
type CreateTokenDeviceOutput struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Block user access token
type BlockTokenUserInput struct {
	UserId  uuid.UUID `json:"user_id" validate:"required"`
	TokenId uuid.UUID `json:"token_id" validate:"required"`
}

// Block user refresh token
type BlockTokenUserRefreshInput struct {
	UserId  uuid.UUID `json:"user_id" validate:"required"`
	TokenId uuid.UUID `json:"token_id" validate:"required"`
}

// Block device access token
type BlockTokenDeviceInput struct {
	DeviceId uuid.UUID `json:"device_id" validate:"required"`
	TokenId  uuid.UUID `json:"token_id" validate:"required"`
}

// Block device refresh token
type BlockTokenDeviceRefreshInput struct {
	DeviceId uuid.UUID `json:"device_id" validate:"required"`
	TokenId  uuid.UUID `json:"token_id" validate:"required"`
}

// Parse user access token
type ParseTokenUserInput struct {
	Token string `json:"token" validate:"required"`
}
type ParseTokenUserOutput struct {
	UserId    string    `json:"user_id"`
	TokenId   string    `json:"token_id"`
	Expires   time.Time `json:"expires"`
	Role      int       `json:"role"`
	CompanyId string    `json:"company_id"`
}

// Parse device access token
type ParseTokenDeviceInput struct {
	Token string `json:"token" validate:"required"`
}
type ParseTokenDeviceOutput struct {
	DeviceId  string    `json:"device_id"`
	Expires   time.Time `json:"expires"`
	CompanyId string    `json:"company_id"`
	TokenId   string    `json:"token_id"`
}

// Refresh user token
type RefreshTokenUserInput struct {
	AccessToken  string `json:"access_token" validate:"required"`
	RefreshToken string `json:"refresh_token" validate:"required"`
}
type RefreshTokenUserOutput struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Refresh device token
type RefreshTokenDeviceInput struct {
	AccessToken  string `json:"access_token" validate:"required"`
	RefreshToken string `json:"refresh_token" validate:"required"`
}
type RefreshTokenDeviceOutput struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// ========================================
//
//	Token input model
//
// ========================================
type (
	TokenUserJwtInput struct {
		UserId  string    `json:"user_id" validate:"required"`
		TokenId string    `json:"token_id" validate:"required"`
		Expires time.Time `json:"expires" validate:"required"`
	}

	TokenUserRefreshInput struct {
	}

	TokenServiceJwtInput struct {
		ServiceName string    `json:"service_name" validate:"required"`
		TokenId     string    `json:"token_id" validate:"required"`
		Type        int       `json:"type" validate:"required"`
		Expires     time.Time `json:"expires" validate:"required"`
	}

	TokenVerificationRegisterInput struct{}
)

// ========================================
//
//	Token output model
//
// ========================================
type (
	TokenUserJwtOutput struct {
		UserId    string    `json:"user_id"`
		TokenId   string    `json:"jti,omitempty"`
		Issuer    string    `json:"iss,omitempty"`
		Subject   string    `json:"sub,omitempty"`
		Audience  []string  `json:"aud,omitempty"`
		ExpiresAt time.Time `json:"exp,omitempty"`
		NotBefore time.Time `json:"nbf,omitempty"`
		IssuedAt  time.Time `json:"iat,omitempty"`
	}

	TokenServiceJwtOutput struct {
		ServiceId   string    `json:"service_id" validate:"required"`
		ServiceName string    `json:"service_name" validate:"required"`
		TokenId     string    `json:"token_id" validate:"required"`
		Type        int       `json:"type" validate:"required"`
		Issuer      string    `json:"iss,omitempty"`
		Subject     string    `json:"sub,omitempty"`
		Audience    []string  `json:"aud,omitempty"`
		ExpiresAt   time.Time `json:"exp,omitempty"`
		NotBefore   time.Time `json:"nbf,omitempty"`
		IssuedAt    time.Time `json:"iat,omitempty"`
	}
)
