package model

import "github.com/google/uuid"

// ======================================
// Model for user presence
// ======================================

// For client send to server
type (
	// For updating user presence
	UserPresenceUpdate struct {
		UserId    uuid.UUID `json:"user_id"`
		Status    int       `json:"status"`
		Timestamp int64     `json:"timestamp"`
	}
)

// For server send to clients
type (
	// For user presence updates - server send to clients
	UserPresenceUpdateSend struct {
		UserId    uuid.UUID `json:"user_id"`
		Status    int       `json:"status"`
		Timestamp int64     `json:"timestamp"`
	}
)
