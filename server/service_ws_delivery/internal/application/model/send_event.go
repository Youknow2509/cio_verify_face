package model

import "github.com/google/uuid"

// ==========================================
// Model from send event
// ==========================================
type SendDataVerifyFace struct {
	DeviceId  uuid.UUID `json:"device_id"`
	DataUrl   string    `json:"data_url"`
	Metadata  string    `json:"metadata"`
	Timestamp int64     `json:"timestamp"`
}
