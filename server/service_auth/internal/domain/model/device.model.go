package model

import "github.com/google/uuid"

// CreateDeviceToken
type CreateDeviceTokenInput struct {
	DeviceId uuid.UUID `json:"device_id"`
	NewToken string    `json:"new_token"`
}

// CheckTokenDevice
type CheckTokenDeviceInput struct {
	DeviceId uuid.UUID `json:"device_id"`
	Token    string    `json:"token"`
}
