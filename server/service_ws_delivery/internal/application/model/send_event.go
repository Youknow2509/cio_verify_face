package model

import "github.com/google/uuid"

// ==========================================
// Model from send event
// ==========================================
type CallEnd struct {
	CallId   string `json:"call_id"`
	Reason   int    `json:"reason"`
	Duration int64  `json:"duration"`
	// Info user
	UserId       uuid.UUID `json:"user_id"`
	SessionId    uuid.UUID `json:"session_id"`
	ConnectionId uuid.UUID `json:"connection_id"`
	IpAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
}

type CallIceCandidate struct {
	Candidate     string `json:"candidate"`
	SdpMid        string `json:"sdp_mid"`
	SdpMLineIndex int    `json:"sdp_m_line_index"`
	//
	ReceiverId uuid.UUID `json:"receiver_id"`
	CallId     string    `json:"call_id"`
	// Info user
	UserId       uuid.UUID `json:"user_id"`
	SessionId    uuid.UUID `json:"session_id"`
	ConnectionId uuid.UUID `json:"connection_id"`
	IpAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
}

type CallAnswer struct {
	CallerId   string    `json:"caller_id"`
	ReceiverId uuid.UUID `json:"receiver_id"`
	CallType   int       `json:"call_type"`
	SdpAnswer  string    `json:"sdp_answer"`
	// Info user
	UserId       uuid.UUID `json:"user_id"`
	SessionId    uuid.UUID `json:"session_id"`
	ConnectionId uuid.UUID `json:"connection_id"`
	IpAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
}

type CallOffer struct {
	ReceiverId uuid.UUID `json:"receiver_id"`
	CallType   int       `json:"call_type"`
	SdpOffer   string    `json:"sdp_offer"`
	// Info user
	SenderId     uuid.UUID `json:"sender_id"`
	SessionId    uuid.UUID `json:"session_id"`
	ConnectionId uuid.UUID `json:"connection_id"`
	IpAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
}

type CreateGroup struct {
	UserCreateId       uuid.UUID   `json:"user_create_id"`
	GroupName          string      `json:"group_name"`
	ConversationIdTemp string      `json:"conversation_id_temp"`
	MemberIds          []uuid.UUID `json:"member_ids"`
	Timestamp          int64       `json:"timestamp"`
}

type ReactMessage struct {
	Status         bool      `json:"status"`
	ConversationId uuid.UUID `json:"conversation_id"`
	MessageId      uuid.UUID `json:"message_id"`
	Reaction       int       `json:"reaction"`
	Timestamp      int64     `json:"timestamp"`
	// Info user
	UserId       uuid.UUID `json:"user_id"`
	SessionId    uuid.UUID `json:"session_id"`
	ConnectionId uuid.UUID `json:"connection_id"`
	IpAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
}

type NewMessageText struct {
	SenderId       uuid.UUID `json:"sender_id"`
	ConversationId uuid.UUID `json:"conversation_id"`
	TempId         string    `json:"temp_id,omitempty"`
	Message        string    `json:"message"`
	ReplyToId      uuid.UUID `json:"reply_to_id,omitempty"`
	Timestamp      int64     `json:"timestamp"`
	// Info user
	SessionId    uuid.UUID `json:"session_id"`
	ConnectionId uuid.UUID `json:"connection_id"`
	IpAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
}

type UserTypingStatus struct {
	UserId         uuid.UUID `json:"user_id"`
	ConversationId uuid.UUID `json:"conversation_id"`
	Status         bool      `json:"status"`
	Timestamp      int64     `json:"timestamp"`
	// Info user
	SessionId    uuid.UUID `json:"session_id"`
	ConnectionId uuid.UUID `json:"connection_id"`
	IpAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
}

type MessageReadStatus struct {
	UserId         uuid.UUID `json:"user_id"`
	ConversationId uuid.UUID `json:"conversation_id"`
	MessageId      uuid.UUID `json:"message_id"`
	Timestamp      int64     `json:"timestamp"`
	// Info user
	SessionId    uuid.UUID `json:"session_id"`
	ConnectionId uuid.UUID `json:"connection_id"`
	IpAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
}

type EditMessage struct {
	UserId         uuid.UUID `json:"user_id"`
	ConversationId uuid.UUID `json:"conversation_id"`
	MessageId      uuid.UUID `json:"message_id"`
	NewContent     string    `json:"new_content"`
	Timestamp      int64     `json:"timestamp"`
	// Info user
	SessionId    uuid.UUID `json:"session_id"`
	ConnectionId uuid.UUID `json:"connection_id"`
	IpAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
}

type DeleteMessage struct {
	Type           int       `json:"type"`
	ConversationId uuid.UUID `json:"conversation_id"`
	MessageId      uuid.UUID `json:"message_id"`
	Timestamp      int64     `json:"timestamp"`
	// Info user
	UserId       uuid.UUID `json:"user_id"`
	SessionId    uuid.UUID `json:"session_id"`
	ConnectionId uuid.UUID `json:"connection_id"`
	IpAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
}
