package dto

import (
	"encoding/json"

	"github.com/google/uuid"
)

// ====================================
//
//	Main model use for system
//
// ====================================
type (
	// DataClientReceive
	DataClientReceive struct {
		Type    int             `json:"type"`
		Payload json.RawMessage `json:"payload"`
	}

	DataServerReceive struct {
		Type    int             `json:"type" validate:"required"`
		Payload json.RawMessage `json:"payload" validate:"required"`
	}

	// Client info
	ClientInfo struct {
		UserId       uuid.UUID `json:"user_id"`
		SessionId    uuid.UUID `json:"session_id"`
		ConnectionId uuid.UUID `json:"connection_id"`
		IpAddress    string    `json:"ip_address"`
		UserAgent    string    `json:"user_agent"`
	}

	// Event Job
	EventJob struct {
		Type       int        `json:"type"`
		ClientInfo ClientInfo `json:"client_info"`
		Payload    []byte     `json:"payload"`
	}
)
