package model

import "github.com/google/uuid"

// ========================
// Action model received
// ========================
type (
	// For action typing client send to server
	ActionTypingStatusReceived struct {
		SenderId       uuid.UUID `json:"sender_id"`
		ConversationId uuid.UUID `json:"conversation_id"`
		Status         bool      `json:"status"`
		Timestamp      int64     `json:"timestamp"`
	}

	// For action read receipt client send to server
	ActionReadReceiptReceived struct {
		SenderId       uuid.UUID `json:"sender_id"`
		ConversationId uuid.UUID `json:"conversation_id"`
		MessageId      uuid.UUID `json:"message_id"`
		Timestamp      int64     `json:"timestamp"`
	}

	// For action edit message client send to server
	ActionEditMessageReceived struct {
		SenderId       uuid.UUID `json:"sender_id"`
		ConversationId uuid.UUID `json:"conversation_id"`
		MessageId      uuid.UUID `json:"message_id"`
		NewContent     string    `json:"new_content"`
		Timestamp      int64     `json:"timestamp"`
	}

	// For action delete message client send to server
	ActionDeleteMessageReceived struct {
		SenderId       uuid.UUID         `json:"sender_id"`
		ConversationId uuid.UUID         `json:"conversation_id"`
		MessageId      uuid.UUID         `json:"message_id"`
		DeleteType     DeleteMessageType `json:"delete_type"` // 0 - for me | 1 - for everyone
		Timestamp      int64             `json:"timestamp"`
	}

	// For action react message client send to server
	ActionReactMessageReceived struct {
		SenderId       uuid.UUID        `json:"sender_id"`
		ConversationId uuid.UUID        `json:"conversation_id"`
		MessageId      uuid.UUID        `json:"message_id"`
		ReactionType   ReactMessageType `json:"reaction_type"` // 0 - love | 1 - like | 2 - dislike | .................
		Timestamp      int64            `json:"timestamp"`
	}

	// For action create group client send to server
	ActionCreateGroupReceived struct {
		SenderId           uuid.UUID   `json:"sender_id"`
		ConversationIdTemp string      `json:"conversation_id"`
		GroupName          string      `json:"group_name"`
		MemberIds          []uuid.UUID `json:"member_ids"`
		Timestamp          int64       `json:"timestamp"`
	}

	// For action add members to group - client send to server
	ActionAddMembersToGroupReceived struct {
		ConversationId uuid.UUID   `json:"conversation_id"`
		MemberIds      []uuid.UUID `json:"member_ids"`
		Timestamp      int64       `json:"timestamp"`
	}

	// For action leave group - client send to server
	ActionLeaveGroupReceived struct {
		ConversationId uuid.UUID `json:"conversation_id"`
		UserId         uuid.UUID `json:"user_id"`
		Timestamp      int64     `json:"timestamp"`
	}
)

// ========================
// Action model send
// ========================
type (
	// For action add members to group - server send to clients
	ActionAddMembersToGroupSend struct {
		ConversationId   uuid.UUID `json:"conversation_id"`
		ConversationName string    `json:"conversation_name"`
		Timestamp        int64     `json:"timestamp"`
	}

	// For action typing server send to clients
	ActionTypingStatusSend struct {
		SenderId       uuid.UUID `json:"sender_id"`
		ConversationId uuid.UUID `json:"conversation_id"`
		Status         bool      `json:"status"`
		Timestamp      int64     `json:"timestamp"`
	}

	// For action read receipt server send to clients
	ActionReadReceiptSend struct {
		SenderId       uuid.UUID `json:"sender_id"`
		ConversationId uuid.UUID `json:"conversation_id"`
		MessageId      uuid.UUID `json:"message_id"`
		Timestamp      int64     `json:"timestamp"`
	}

	// For action edit message server send to clients
	ActionEditMessageSend struct {
		SenderId       uuid.UUID `json:"sender_id"`
		ConversationId uuid.UUID `json:"conversation_id"`
		MessageId      uuid.UUID `json:"message_id"`
		NewContent     string    `json:"new_content"`
		Timestamp      int64     `json:"timestamp"`
	}

	// For action delete message server send to clients
	ActionDeleteMessageSend struct {
		ConversationId uuid.UUID `json:"conversation_id"`
		MessageId      uuid.UUID `json:"message_id"`
		Timestamp      int64     `json:"timestamp"`
	}

	// For action react message server send to clients
	ActionReactMessageSend struct {
		SenderId       uuid.UUID        `json:"sender_id"`
		ConversationId uuid.UUID        `json:"conversation_id"`
		MessageId      uuid.UUID        `json:"message_id"`
		ReactionType   ReactMessageType `json:"reaction_type"` // 0 - love | 1 - like | 2 - dislike | .................
		IconUrl        string           `json:"icon_url"`
		Timestamp      int64            `json:"timestamp"`
	}

	// For action create group server send to clients
	ActionCreateGroupSend struct {
		SenderId       uuid.UUID `json:"sender_id"`
		ConversationId uuid.UUID `json:"conversation_id"`
		GroupName      string    `json:"group_name"`
		Timestamp      int64     `json:"timestamp"`
	}
)
