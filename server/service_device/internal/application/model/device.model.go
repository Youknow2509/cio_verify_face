package model

import (
	"github.com/google/uuid"
)

// =================================================
// Device model
// =================================================

// Update status device
type UpdateStatusDeviceInput struct {
	// Info req
	DeviceId uuid.UUID `json:"device_id"`
	Status   int       `json:"status"` // 0: OFFLINE, 1: ONLINE, 2: MAINTENANCE, 3: ERROR
	// Info client req
	UserId      uuid.UUID `json:"user_id"`
	Role        int       `json:"role"` // 0: ADMIN, 1: Admin company, 2: STAFF
	SessionId   uuid.UUID `json:"session_id"`
	ClientIp    string    `json:"client_ip"`
	ClientAgent string    `json:"client_agent"`
}

// Refresh device token
type RefreshDeviceTokenInput struct {
	// Info req
	DeviceId uuid.UUID `json:"device_id"`
	// Info client req
	UserId      uuid.UUID `json:"user_id"`
	Role        int       `json:"role"` // 0: ADMIN, 1: Admin company, 2: STAFF
	SessionId   uuid.UUID `json:"session_id"`
	ClientIp    string    `json:"client_ip"`
	ClientAgent string    `json:"client_agent"`
}
type RefreshDeviceTokenOutput struct {
	DeviceId    string `json:"device_id"`
	DeviceToken string `json:"device_token"`
}

// Get device token
type GetDeviceTokenInput struct {
	// Info req
	DeviceId uuid.UUID `json:"device_id"`
	// Info client req
	UserId      uuid.UUID `json:"user_id"`
	Role        int       `json:"role"` // 0: ADMIN, 1: Admin company, 2: STAFF
	SessionId   uuid.UUID `json:"session_id"`
	ClientIp    string    `json:"client_ip"`
	ClientAgent string    `json:"client_agent"`
}
type GetDeviceTokenOutput struct {
	DeviceId    string `json:"device_id"`
	DeviceToken string `json:"device_token"`
}

// CreateNewDevice
type CreateNewDeviceInput struct {
	// Info req
	DeviceName   string    `json:"device_name"`
	Address      string    `json:"address"`
	DeviceType   int       `json:"device_type"` // 0: FACE_TERMINAL, 1: MOBILE_APP, 2: WEB_CAMERA, 3: IOT_SENSOR
	SerialNumber string    `json:"serial_number"`
	MacAddress   string    `json:"mac_address"`
	// Info client req
	UserId      uuid.UUID `json:"user_id"`
	Role        int       `json:"role"` // 0: ADMIN, 1: Admin company, 2: STAFF
	SessionId   uuid.UUID `json:"session_id"`
	ClientIp    string    `json:"client_ip"`
	ClientAgent string    `json:"client_agent"`
}
type CreateNewDeviceOutput struct {
	DeviceId     string `json:"device_id"`
	CompanyId    string `json:"company_id"`
	Name         string `json:"name"`
	Address      string `json:"address"`
	SerialNumber string `json:"serial_number"`
	MacAddress   string `json:"mac_address"`
}

// ListDevices
type ListDevicesInput struct {
	// Info req
	CompanyId uuid.UUID `json:"company_id"`
	Size      int       `json:"size"`
	Page      int       `json:"page"`
	// Info client req
	UserId      uuid.UUID `json:"user_id"`
	Role        int       `json:"role"` // 0: ADMIN, 1: Admin company, 2: STAFF
	SessionId   uuid.UUID `json:"session_id"`
	ClientIp    string    `json:"client_ip"`
	ClientAgent string    `json:"client_agent"`
}
type ListDevicesOutput struct {
	Devices []*GetDeviceByIdOutput `json:"devices"`
}

// GetDeviceById
type GetDeviceByIdInput struct {
	// Info req
	DeviceId uuid.UUID `json:"device_id"`
	// Info client req
	UserId      uuid.UUID `json:"user_id"`
	Role        int       `json:"role"` // 0: ADMIN, 1: Admin company, 2: STAFF
	SessionId   uuid.UUID `json:"session_id"`
	ClientIp    string    `json:"client_ip"`
	ClientAgent string    `json:"client_agent"`
}

type GetDeviceByIdOutput struct {
	DeviceId     string `json:"device_id"`
	CompanyId    string `json:"company_id"`
	Name         string `json:"name"`
	Address      string `json:"address"`
	SerialNumber string `json:"serial_number"`
	MacAddress   string `json:"mac_address"`
	Status       int    `json:"status"`
	CreateAt     string `json:"create_at"`
	UpdateAt     string `json:"update_at"`
	Token        string `json:"token"`
}

// UpdateDevice
type UpdateDeviceInput struct {
	// Info req
	LocationId   uuid.UUID `json:"location_id" validate:"omitempty"`
	DeviceName   string    `json:"device_name" validate:"omitempty,min=3,max=100"`
	Address      string    `json:"address" validate:"omitempty,max=255"`
	DeviceType   int       `json:"device_type" validate:"omitempty,oneof=0 1 2 3"` // 0: FACE_TERMINAL, 1: MOBILE_APP, 2: WEB_CAMERA, 3: IOT_SENSOR
	SerialNumber string    `json:"serial_number" validate:"omitempty,max=100"`
	MacAddress   string    `json:"mac_address" validate:"omitempty,mac"`
	Status       int       `json:"status" validate:"omitempty,oneof=0 1 2 3"` // 0: OFFLINE, 1: ONLINE, 2: MAINTENANCE, 3: ERROR
	// Info client req
	UserId      uuid.UUID `json:"user_id"`
	Role        int       `json:"role"` // 0: ADMIN, 1: Admin company, 2: STAFF
	SessionId   uuid.UUID `json:"session_id"`
	ClientIp    string    `json:"client_ip"`
	ClientAgent string    `json:"client_agent"`
}
type UpdateDeviceOutput struct{}

// DeleteDevice
type DeleteDeviceInput struct {
	// Info req
	DeviceId uuid.UUID `json:"device_id"`
	// Info client req
	UserId      uuid.UUID `json:"user_id"`
	Role        int       `json:"role"` // 0: ADMIN, 1: Admin company, 2: STAFF
	SessionId   uuid.UUID `json:"session_id"`
	ClientIp    string    `json:"client_ip"`
	ClientAgent string    `json:"client_agent"`
}

// UpdateLocationDevice
type UpdateLocationDeviceInput struct {
	// Info req
	DeviceId      uuid.UUID `json:"device_id"`
	NewLocationId uuid.UUID `json:"location_id"`
	NewAddress    string    `json:"address"`
	// Info client req
	UserId      uuid.UUID `json:"user_id"`
	Role        int       `json:"role"` // 0: ADMIN, 1: Admin company, 2: STAFF
	SessionId   uuid.UUID `json:"session_id"`
	ClientIp    string    `json:"client_ip"`
	ClientAgent string    `json:"client_agent"`
}

// UpdateNameDeviceInput
type UpdateNameDeviceInput struct {
	// Info req
	DeviceId uuid.UUID `json:"device_id"`
	NewName  string    `json:"device_name"`
	// Info client req
	UserId      uuid.UUID `json:"user_id"`
	Role        int       `json:"role"` // 0: ADMIN, 1: Admin company, 2: STAFF
	SessionId   uuid.UUID `json:"session_id"`
	ClientIp    string    `json:"client_ip"`
	ClientAgent string    `json:"client_agent"`
}

// UpdateInfoDeviceInput
type UpdateInfoDeviceInput struct {
	// Info req
	DeviceId        uuid.UUID `json:"device_id"`
	NewDeviceType   int       `json:"device_type"` // 0: FACE_TERMINAL, 1: MOBILE_APP, 2: WEB_CAMERA, 3: IOT_SENSOR
	NewSerialNumber string    `json:"serial_number"`
	NewMacAddress   string    `json:"mac_address"`
	// Info client req
	UserId      uuid.UUID `json:"user_id"`
	Role        int       `json:"role"` // 0: ADMIN, 1: Admin company, 2: STAFF
	SessionId   uuid.UUID `json:"session_id"`
	ClientIp    string    `json:"client_ip"`
	ClientAgent string    `json:"client_agent"`
}
