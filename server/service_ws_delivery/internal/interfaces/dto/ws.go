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
		Type    int             `json:"type"`
		Payload json.RawMessage `json:"payload"`
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

// ==========================================
//
//	Details payload
//
// ==========================================
type (
	NewMessageText struct {
		ConversationId string `json:"conversation_id"`
		TempId         string `json:"temp_id,omitempty"`
		Message        string `json:"message"`
		ReplyToId      string `json:"reply_to_id,omitempty"`
	}

	EditMessage struct {
		ConversationId string `json:"conversation_id"`
		MessageId      string `json:"message_id"`
		NewContent     string `json:"new_content"`
	}

	DeleteMessage struct {
		Type           int    `json:"type"`
		ConversationId string `json:"conversation_id"`
		MessageId      string `json:"message_id"`
	}

	ReactMessage struct {
		Status         bool   `json:"status"`
		ConversationId string `json:"conversation_id"`
		MessageId      string `json:"message_id"`
		Reaction       int    `json:"reaction"`
	}

	UserTypingStatus struct {
		ConversationId string `json:"conversation_id"`
		IsTyping       bool   `json:"is_typing"`
	}

	MessageReadStatus struct {
		ConversationId string `json:"conversation_id"`
		MessageId      string `json:"message_id"`
	}

	CallOffer struct {
		ReceiverId string `json:"receiver_id"`
		CallType   int    `json:"call_type"`
		SdpOffer   string `json:"sdp_offer"`
	}

	CallAnswer struct {
		CallerId   string `json:"caller_id"`
		ReceiverId string `json:"receiver_id"`
		CallType   int    `json:"call_type"`
		SdpAnswer  string `json:"sdp_answer"`
	}

	CallIceCandidate struct {
		Candidate     string `json:"candidate"`
		SdpMid        string `json:"sdp_mid"`
		SdpMLineIndex int    `json:"sdp_m_line_index"`
		//
		ReceiverId string `json:"receiver_id"`
		CallId     string `json:"call_id"`
	}

	CallEnd struct {
		CallId   string `json:"call_id"`
		Reason   int    `json:"reason"`
		Duration int64  `json:"duration"`
	}
	// v.v
)
