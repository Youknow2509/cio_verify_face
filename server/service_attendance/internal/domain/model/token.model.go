package model

import (
	"time"
)

// Create user token
type UserTokenOutput struct {
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
		UserId string `json:"user_id" validate:"required"`
		Role   int    `json:"role" validate:"required"`
	}

	TokenDeviceJwtInput struct {
		DeviceId  string    `json:"device_id" validate:"required"`
		CompanyId string    `json:"company_id" validate:"required"`
		TokenId   string    `json:"token_id" validate:"required"`
		Expires   time.Time `json:"expires" validate:"required"`
	}

	ServiceRefreshTokenInput struct {
		ServiceId string    `json:"service_id" validate:"required"`
		TokenId   string    `json:"token_id" validate:"required"`
		Expires   time.Time `json:"expires" validate:"required"`
	}
)

// ========================================
//
//	Token output model
//
// ========================================
type (
	TokenUserJwtOutput struct {
		UserId    string    `json:"user_id"`
		Role      int       `json:"role"`
		TokenId   string    `json:"jti,omitempty"`
		CompanyId string    `json:"company_id,omitempty"`
		Issuer    string    `json:"iss,omitempty"`
		Subject   string    `json:"sub,omitempty"`
		Audience  []string  `json:"aud,omitempty"`
		ExpiresAt time.Time `json:"exp,omitempty"`
		NotBefore time.Time `json:"nbf,omitempty"`
		IssuedAt  time.Time `json:"iat,omitempty"`
	}

	TokenUserRefreshOutput struct {
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

	TokenDeviceJwtOutput struct {
		DeviceId  string    `json:"device_id" validate:"required"`
		CompanyId string    `json:"company_id" validate:"required"`
		TokenId   string    `json:"token_id" validate:"required"`
		Issuer    string    `json:"iss,omitempty"`
		Subject   string    `json:"sub,omitempty"`
		Audience  []string  `json:"aud,omitempty"`
		ExpiresAt time.Time `json:"exp,omitempty"`
		NotBefore time.Time `json:"nbf,omitempty"`
		IssuedAt  time.Time `json:"iat,omitempty"`
	}
)
