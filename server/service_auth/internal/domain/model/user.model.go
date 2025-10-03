package model

import (
	"net/netip"
	"time"

	"github.com/google/uuid"
)

// =======================================================
//
//	Define input and output models for Auth Base
//
// =======================================================
type (
	/**
	 * Input create user session struct
	 */
	CreateUserSessionInput struct {
		SessionID    uuid.UUID  `json:"session_id"`
		UserID       uuid.UUID  `json:"user_id"`
		IPAddress    netip.Addr `json:"ip_address"`
		RefreshToken string     `json:"refresh_token"`
		ExpiredAt    time.Time  `json:"expired_at"`
	}

	/**
	 * Input create user session app struct
	 */
	CreateUserSessionAppInput struct {
		SessionID    uuid.UUID  `json:"session_id"`
		UserID       uuid.UUID  `json:"user_id"`
		IPAddress    netip.Addr `json:"ip_address"`
		RefreshToken string     `json:"refresh_token"`
		ExpiredAt    time.Time  `json:"expired_at"`
		DeviceID     string     `json:"device_id"`
		PushToken    string     `json:"push_token"`
	}

	/**
	 * Input remove user session struct
	 */
	RemoveUserSessionInput struct {
		SessionID uuid.UUID `json:"session_id"`
	}

	// GetSessionInfoInput
	GetSessionInfoInput struct {
		SessionID uuid.UUID `json:"session_id"`
	}

	// GetRefreshSessionInfoInput
	GetRefreshSessionInfoInput struct {
		SessionID uuid.UUID `json:"session_id"`
	}

	// RefreshSessionInput
	RefreshSessionInput struct {
		SessionID    uuid.UUID `json:"session_id"`
		RefreshToken string    `json:"refresh_token"`
		ExpiredAt    time.Time `json:"expired_at"`
	}

	// v.v
)

type (
	/**
	 * Output user base info struct
	 */
	UserBaseInfoOutput struct {
		UserID       string `json:"user_id"`
		UserEmail    string `json:"user_email"`
		UserSalt     string `json:"user_salt"`
		UserPassword string `json:"user_password"`
		IsBlocked    bool   `json:"is_blocked"`
		Role         int    `json:"role"`
	}

	// GetRefreshSessionInfoOutput
	GetRefreshSessionInfoOutput struct {
		RefreshToken string    `json:"refresh_token"`
		ExpiredAt    time.Time `json:"expired_at"`
		SessionID    uuid.UUID `json:"session_id"`
	}
)
