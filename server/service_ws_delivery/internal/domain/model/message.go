package model

import "github.com/google/uuid"

// ============================================================
// Message model Received
// ============================================================
type (
	// MessageTextReceived represents a text message received by the user from ws
	MessageTextReceived struct {
		SenderId       uuid.UUID `json:"sender_id"`
		ConversationId uuid.UUID `json:"conversation_id"`
		TempId         string    `json:"temp_id,omitempty"`
		Message        string    `json:"message"`
		ReplyToId      uuid.UUID `json:"reply_to_id,omitempty"`
		Timestamp      int64     `json:"timestamp"`
	}
)

// ============================================================
// Message model Send
// ============================================================
// For media send
type (
	// MediaImageData represents image data in a media message
	MediaImageData struct {
		OriginalUrl   string `json:"original_url"`
		ThumbnailUrl  string `json:"thumbnail_url"`
		Width         int    `json:"width"`
		Height        int    `json:"height"`
		FileSizeBytes int64  `json:"file_size_bytes"`
	}

	// MediaVideoData represents video data in a media message
	MediaVideoData struct {
		ThumbnailUrl    string            `json:"thumbnail_url"`
		DurationSeconds int               `json:"duration_seconds"`
		Resolutions     map[string]string `json:"resolutions"`
		FileSizeBytes   int64             `json:"file_size_bytes"`
	}

	// MediaFileData represents file data in a media message
	MediaFileData struct {
		DownloadUrl   string `json:"download_url"`
		FileName      string `json:"file_name"`
		FileSizeBytes int64  `json:"file_size_bytes"`
		MimeType      string `json:"mime_type"` // Ex: application/pdf, image/jpeg, ...
	}

	// MediaAudioData represents audio data in a media message
	MediaAudioData struct {
		DownloadUrl     string    `json:"download_url"`
		DurationSeconds int       `json:"duration_seconds"`
		MimeType        string    `json:"mime_type"`
		Waveform        []float64 `json:"waveform"`
	}
)
type (
	// MessageTextSend represents a text message sent by the user to ws
	MessageTextSend struct {
		SenderId       uuid.UUID `json:"sender_id"`
		ConversationId uuid.UUID `json:"conversation_id"`
		MessageId      string    `json:"message_id,omitempty"`
		Message        string    `json:"message"`
		ReplyToId      uuid.UUID `json:"reply_to_id,omitempty"`
		Timestamp      int64     `json:"timestamp"`
	}

	// MessageMediaImageSend represents a media image message sent by the user to ws
	MessageMediaImageSend struct {
		SenderId       uuid.UUID        `json:"sender_id"`
		ConversationId uuid.UUID        `json:"conversation_id"`
		MessageId      string           `json:"message_id,omitempty"`
		ReplyToId      uuid.UUID        `json:"reply_to_id,omitempty"`
		MediaData      []MediaImageData `json:"media_data"`
		Timestamp      int64            `json:"timestamp"`
	}

	// MessageMediaVideoSend represents a media video message sent by the user to ws
	MessageMediaVideoSend struct {
		SenderId       uuid.UUID        `json:"sender_id"`
		ConversationId uuid.UUID        `json:"conversation_id"`
		MessageId      string           `json:"message_id,omitempty"`
		ReplyToId      uuid.UUID        `json:"reply_to_id,omitempty"`
		MediaData      []MediaVideoData `json:"media_data"`
		Timestamp      int64            `json:"timestamp"`
	}
)
