package model

import "github.com/google/uuid"

// ==================================================
// Model for send event mq domain
// ==================================================
type (
	WriteNewMessage struct {
		SenderId       uuid.UUID `json:"sender_id"`
		ConversationId uuid.UUID `json:"conversation_id"`
		TempId         string    `json:"temp_id,omitempty"`
		Message        string    `json:"message"`
		ReplyToId      uuid.UUID `json:"reply_to_id,omitempty"`
		Timestamp      int64     `json:"timestamp"`
	}

	UpgradeStatusTypingUser struct {
		UserId         uuid.UUID `json:"user_id"`
		ConversationId uuid.UUID `json:"conversation_id"`
		Status         bool      `json:"status"`
		Timestamp      int64     `json:"timestamp"`
	}

	UserReadMessageStatus struct {
		UserId         uuid.UUID `json:"user_id"`
		ConversationId uuid.UUID `json:"conversation_id"`
		MessageId      uuid.UUID `json:"message_id"`
		Timestamp      int64     `json:"timestamp"`
	}

	UserEditMessage struct {
		UserId     uuid.UUID `json:"user_id"`
		MessageId  uuid.UUID `json:"message_id"`
		NewContent string    `json:"new_content"`
		Timestamp  int64     `json:"timestamp"`
	}

	UserDeleteMessage struct {
		Type           DeleteMessageType `json:"type"`
		UserId         uuid.UUID         `json:"user_id"`
		MessageId      uuid.UUID         `json:"message_id"`
		ConversationId uuid.UUID         `json:"conversation_id"`
		Timestamp      int64             `json:"timestamp"`
	}

	UserReactMessage struct {
		Status         bool             `json:"status"`
		ConversationId uuid.UUID        `json:"conversation_id"`
		MessageId      uuid.UUID        `json:"message_id"`
		Reaction       ReactMessageType `json:"reaction"`
		Timestamp      int64            `json:"timestamp"`
		UserId         uuid.UUID        `json:"user_id"`
	}

	UserCallOfferInitilize struct {
		ReceiverId uuid.UUID `json:"receiver_id"`
		CallType   CallType  `json:"call_type"`
		SdpOffer   string    `json:"sdp_offer"`
		SenderId   uuid.UUID `json:"sender_id"`
	}

	UserCallAnswer struct {
		SenderId     uuid.UUID `json:"sender_id"`
		ReceiverId   uuid.UUID `json:"receiver_id"`
		CallType     CallType  `json:"call_type"`
		SdpAnswer    string    `json:"sdp_answer"`
		CallAnswerId string    `json:"call_answer_id"`
	}

	UserCallIceCandidate struct {
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

	UserCallEnd struct {
		CallId   string        `json:"call_id"`
		Reason   CallEndReason `json:"reason"`
		Duration int64         `json:"duration"`
		// Info user
		UserId       uuid.UUID `json:"user_id"`
		SessionId    uuid.UUID `json:"session_id"`
		ConnectionId uuid.UUID `json:"connection_id"`
		IpAddress    string    `json:"ip_address"`
		UserAgent    string    `json:"user_agent"`
	}
)
