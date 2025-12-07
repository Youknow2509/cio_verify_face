package dto

// ======================
//
//	Event DTOs
//
// ======================
type KafkaEvent struct {
	EventType int         `json:"event_type" validate:"gte=0"`
	Payload   interface{} `json:"payload" validate:"required"`
}

// PasswordResetNotificationEvent represents the full structure for password reset notifications
type PasswordResetNotificationEvent struct {
	EventType int                            `json:"event_type" validate:"gte=0"`
	Metadata  NotificationMetadata           `json:"metadata" validate:"required"`
	Payload   PasswordResetNotificationPayload `json:"payload" validate:"required"`
}

// NotificationMetadata represents metadata for notification events
type NotificationMetadata struct {
	CreatedAt string `json:"created_at" validate:"required"`
	MessageID string `json:"message_id" validate:"required"`
	RequestID string `json:"request_id" validate:"required"`
	UserID    string `json:"user_id" validate:"required"`
}

// PasswordResetNotificationPayload represents the payload for password reset notifications
type PasswordResetNotificationPayload struct {
	ExpiresIn int    `json:"expires_in" validate:"required,gt=0"`
	FullName  string `json:"full_name" validate:"required"`
	ResetURL  string `json:"reset_url" validate:"required,url"`
	To        string `json:"to" validate:"required,email"`
}

// ReportAttentionNotificationEvent represents the full structure for report attention notifications
type ReportAttentionNotificationEvent struct {
	CompanyID   string `json:"company_id" validate:"required"`
	CreatedAt   string `json:"created_at" validate:"required"`
	DownloadURL string `json:"download_url" validate:"required,url"`
	Email       string `json:"email" validate:"required,email"`
	EndDate     string `json:"end_date" validate:"required"`
	Format      string `json:"format" validate:"required"`
	StartDate   string `json:"start_date" validate:"required"`
	Type        string `json:"type" validate:"required"`
}
