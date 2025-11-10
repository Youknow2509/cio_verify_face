package dto

// =======================================
// Device DTO
// =======================================

// Update status device request
type UpdateStatusDeviceRequest struct {
	DeviceId string `json:"device_id" validate:"required"`
	Status   int    `json:"status" validate:"oneof=0 1 2 3"` // 0: OFFLINE, 1: ONLINE, 2: MAINTENANCE, 3: ERROR
}

// Create device request
type CreateDeviceRequest struct {
	CompanyId    string `json:"company_id" validate:"omitempty"`
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

// UpdateLocationDeviceRequest
type UpdateLocationDeviceRequest struct {
	DeviceId      string `json:"device_id" validate:"required"`
	NewLocationId string `json:"location_id" validate:"required"`
	NewAddress    string `json:"address" validate:"required,max=255"`
}

// Update name device request
type UpdateNameDeviceRequest struct {
	DeviceId string `json:"device_id" validate:"required"`
	NewName  string `json:"device_name" validate:"required,min=3,max=100"`
}

// Update info device request
type UpdateInfoDeviceRequest struct {
	DeviceId        string `json:"device_id" validate:"required"`
	NewDeviceType   int    `json:"device_type" validate:"omitempty,oneof=0 1 2 3"` // 0: FACE_TERMINAL, 1: MOBILE_APP, 2: WEB_CAMERA, 3: IOT_SENSOR
	NewSerialNumber string `json:"serial_number" validate:"omitempty,max=100"`
	NewMacAddress   string `json:"mac_address" validate:"omitempty,mac"`
}
