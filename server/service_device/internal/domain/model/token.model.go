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
		UserId string `json:"user_id"`
		Role   int    `json:"role"`
	}

	TokenDeviceJwtInput struct {
		DeviceId  string `json:"device_id"`
		CompanyId string `json:"company_id"`
	}

	ServiceRefreshTokenInput struct {
		ServiceId string    `json:"service_id"`
		TokenId   string    `json:"token_id"`
		Expires   time.Time `json:"expires"`
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
		ServiceId   string    `json:"service_id"`
		ServiceName string    `json:"service_name"`
		TokenId     string    `json:"token_id"`
		Type        int       `json:"type"`
		Issuer      string    `json:"iss,omitempty"`
		Subject     string    `json:"sub,omitempty"`
		Audience    []string  `json:"aud,omitempty"`
		ExpiresAt   time.Time `json:"exp,omitempty"`
		NotBefore   time.Time `json:"nbf,omitempty"`
		IssuedAt    time.Time `json:"iat,omitempty"`
	}

	TokenDeviceJwtOutput struct {
		DeviceId  string    `json:"device_id"`
		CompanyId string    `json:"company_id"`
		TokenId   string    `json:"token_id"`
		Issuer    string    `json:"iss,omitempty"`
		Subject   string    `json:"sub,omitempty"`
		Audience  []string  `json:"aud,omitempty"`
		ExpiresAt time.Time `json:"exp,omitempty"`
		NotBefore time.Time `json:"nbf,omitempty"`
		IssuedAt  time.Time `json:"iat,omitempty"`
	}
)
