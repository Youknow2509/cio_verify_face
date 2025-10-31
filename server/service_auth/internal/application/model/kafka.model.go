package model

// =================================
//
//	Model For Kafka Writer
//
// =================================
type KafkaWriterModel struct {
	Type int         `json:"type"`
	Data interface{} `json:"data"`
}

type KafkaWriterEventFriendRequestType int

const (
	KafkaWriterEventFriendRequestTypeSend KafkaWriterEventFriendRequestType = iota
	KafkaWriterEventFriendRequestTypeAccept
	KafkaWriterEventFriendRequestTypeReject
)

type KafkaWriterEventFriendRequest struct {
	Type            KafkaWriterEventFriendRequestType `json:"type"` // 0: send, 1: accept, 2: reject
	SenderId        string                            `json:"sender_id"`
	SenderFristName string                            `json:"sender_first_name"`
	SenderLastName  string                            `json:"sender_last_name"`
	Avatar          string                            `json:"avatar"`
}

type KafkaWriterEventSendOTPRegister struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

type KafkaWriterEventResetPassword struct {
	Email       string `json:"email"`
	PasswordNew string `json:"password_new"`
	Token       string `json:"token"`
	TokenTTL    int64  `json:"token_ttl"`
}

type KafkaWriterEventPushNotification struct {
	Title   string `json:"title"`
	Message string `json:"message"`
	To      string `json:"to"`
}

// v.v
