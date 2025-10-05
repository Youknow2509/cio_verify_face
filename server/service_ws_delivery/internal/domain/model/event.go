package model

import "github.com/google/uuid"

// ===============================
// WebSocket events - server send to clients
// ===============================

type (
	MessageAck struct {
		MessageIdTemp  string    `json:"message_id_temp"`
		MessageIdAck   string    `json:"message_id_ack"`
		SenderId       uuid.UUID `json:"sender_id"`
		ConversationId uuid.UUID `json:"conversation_id"`
		Timestamp      int64     `json:"timestamp"`
	}

	CreateGroupAck struct {
		ConversationIdTemp string    `json:"conversation_id_temp"`
		ConversationId     uuid.UUID `json:"conversation_id"`
	}
)
