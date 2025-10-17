package dto

// =======================================
// Device DTO
// device_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
// company_id UUID NOT NULL REFERENCES companies(company_id) ON DELETE CASCADE,
// location_id UUID, -- For future location management
// name VARCHAR(255) NOT NULL,
// address TEXT,
// device_type int2 DEFAULT 0 CHECK (device_type IN (0, 1, 2, 3)), -- 0: FACE_TERMINAL, 1: MOBILE_APP, 2: WEB_CAMERA, 3: IOT_SENSOR
// serial_number VARCHAR(100),
// mac_address VARCHAR(17),
// ip_address INET,
// firmware_version VARCHAR(20),
// status int2 DEFAULT 1 CHECK (status IN (0, 1, 2, 3)), -- 0: OFFLINE, 1: ONLINE, 2: MAINTENANCE, 3: ERROR
// token VARCHAR(512) NOT NULL, -- Device authentication token
// last_heartbeat TIMESTAMP WITH TIME ZONE,
// settings JSONB DEFAULT '{}', -- Device-specific settings
// created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
// updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
// =======================================

// Create device request
type CreateDeviceRequest struct {
	CompanyId    string `json:"company_id" validate:"required"`
	LocationId   string `json:"location_id" validate:"omitempty"`
	DeviceName   string `json:"device_name" validate:"required,min=3,max=100"`
	Address      string `json:"address" validate:"omitempty,max=255"`
	DeviceType   int    `json:"device_type" validate:"omitempty,oneof=0 1 2 3"` // 0: FACE_TERMINAL, 1: MOBILE_APP, 2: WEB_CAMERA, 3: IOT_SENSOR
	SerialNumber string `json:"serial_number" validate:"omitempty,max=100"`
	MacAddress   string `json:"mac_address" validate:"omitempty,mac"`
}

// Delete device request
type DeleteDeviceRequest struct {
	DeviceId string `json:"device_id" validate:"required"`
}

// Update device request
type UpdateDeviceRequest struct {
	LocationId   string `json:"location_id" validate:"omitempty"`
	DeviceName   string `json:"device_name" validate:"omitempty,min=3,max=100"`
	Address      string `json:"address" validate:"omitempty,max=255"`
	DeviceType   int    `json:"device_type" validate:"omitempty,oneof=0 1 2 3"` // 0: FACE_TERMINAL, 1: MOBILE_APP, 2: WEB_CAMERA, 3: IOT_SENSOR
	SerialNumber string `json:"serial_number" validate:"omitempty,max=100"`
	MacAddress   string `json:"mac_address" validate:"omitempty,mac"`
	Status       int    `json:"status" validate:"omitempty,oneof=0 1 2 3"` // 0: OFFLINE, 1: ONLINE, 2: MAINTENANCE, 3: ERROR
}

// Get device by ID request
type GetDeviceByIdRequest struct {
	DeviceId string `json:"device_id" validate:"required"`
}

// Get list devices request
type ListDevicesRequest struct {
	CompanyId string `json:"company_id" validate:"required"`
	Size      int    `json:"size" validate:"omitempty,min=1,max=100"`
	Page      int    `json:"page" validate:"omitempty,min=1"`
}
