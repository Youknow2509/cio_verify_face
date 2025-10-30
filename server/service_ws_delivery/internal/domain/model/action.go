package model

import "github.com/google/uuid"

// ========================
// Action model received
// ========================
type (
	ActionAttendanceResultReceived struct {
		DeviceId      uuid.UUID `json:"device_id"`
		DataObjectUrl string    `json:"data_object_url"`
		Timestamp     int64     `json:"timestamp"`
	}
)

// ========================
// Action model send
// ========================
type (
	// For action attendance result - server send to clients
	ActionAttendanceResultSend struct {
		DeviceId  uuid.UUID `json:"device_id"`
		UserId    uuid.UUID `json:"user_id"`
		Result    bool      `json:"result"` // true - success | false - fail
		Timestamp int64     `json:"timestamp"`
	}
)
