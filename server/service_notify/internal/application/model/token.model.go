package model

import (
	"time"
)

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
