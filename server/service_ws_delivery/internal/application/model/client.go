package model

import "github.com/google/uuid"

// ==========================================
// Model from client
// ==========================================
type SendMessageToClientInput struct {
	ConnectionId uuid.UUID `json:"connection_id"`
	Message      []byte    `json:"message"`
}
