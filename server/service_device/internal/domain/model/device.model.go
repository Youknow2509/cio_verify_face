package model

import "github.com/google/uuid"

// UpdateTokenDevice
type UpdateTokenDeviceInput struct {
	DeviceId uuid.UUID `json:"device_id"`
	NewToken string    `json:"new_token"`
}

// GetDeviceToken
type GetDeviceTokenInput struct {
	DeviceId uuid.UUID `json:"device_id"`
}
type GetDeviceTokenOutput struct {
	DeviceId uuid.UUID `json:"device_id"`
	Token    string    `json:"token"`
}

// DeviceExist
type DeviceExistInput struct {
	DeviceId uuid.UUID `json:"device_id"`
}

// UpdateDeviceInfo
type UpdateDeviceInfoInput struct {
	DeviceId        uuid.UUID `json:"device_id"`
	SerialNumber    string    `json:"serial_number"`
	MacAddress      string    `json:"mac_address"`
	FirmwareVersion string    `json:"firmware_version"`
}

// UpdateDeviceLocation
type UpdateDeviceLocationInput struct {
	DeviceId   uuid.UUID `json:"device_id"`
	LocationId uuid.UUID `json:"location_id"`
	Address    string    `json:"address"`
}

// UpdateDeviceName
type UpdateDeviceNameInput struct {
	DeviceId uuid.UUID `json:"device_id"`
	Name     string    `json:"name"`
}

// DisableDevice
type DisableDeviceInput struct {
	DeviceId uuid.UUID `json:"device_id"`
}

// EnableDevice
type EnableDeviceInput struct {
	DeviceId uuid.UUID `json:"device_id"`
}

// New device
type NewDevice struct {
	DeviceId     uuid.UUID `json:"device_id"`
	CompanyId    uuid.UUID `json:"company_id"`
	Name         string    `json:"name"`
	Address      string    `json:"address"`
	SerialNumber string    `json:"serial_number"`
	MacAddress   string    `json:"mac_address"`
}

// DeviceInfoBase
type DeviceInfoBaseInput struct {
	DeviceId uuid.UUID `json:"device_id"`
}
type DeviceInfoBaseOutput struct {
	DeviceId     uuid.UUID `json:"device_id"`
	CompanyId    uuid.UUID `json:"company_id"`
	Name         string    `json:"name"`
	Address      string    `json:"address"`
	SerialNumber string    `json:"serial_number"`
	MacAddress   string    `json:"mac_address"`
	Status       int       `json:"status"`
	CreateAt     string    `json:"create_at"`
	UpdateAt     string    `json:"update_at"`
	Token        string    `json:"token"`
}

// DeviceInfo
type DeviceInfoInput struct {
	DeviceId uuid.UUID `json:"device_id"`
}
type DeviceInfoOutput struct {
	DeviceId        uuid.UUID   `json:"device_id"`
	CompanyId       uuid.UUID   `json:"company_id"`
	Name            string      `json:"name"`
	Address         string      `json:"address"`
	SerialNumber    string      `json:"serial_number"`
	MacAddress      string      `json:"mac_address"`
	IpAddress       string      `json:"ip_address"`
	FirmwareVersion string      `json:"firmware_version"`
	LastHeartbeat   string      `json:"last_heartbeat"`
	Settings        interface{} `json:"settings"`
	CreateAt        string      `json:"create_at"`
	UpdateAt        string      `json:"update_at"`
	Token           string      `json:"token"`
}

// ListDeviceInCompany
type ListDeviceInCompanyInput struct {
	CompanyId uuid.UUID `json:"company_id"`
	Limit     int       `json:"limit"`
	Offset    int       `json:"offset"`
}

// List of DeviceInfoBaseOutput
type ListDeviceInCompanyOutput struct {
	Devices []*DeviceInfoBaseOutput `json:"devices"`
}

// DeleteDevice
type DeleteDeviceInput struct {
	DeviceId uuid.UUID `json:"device_id"`
}
