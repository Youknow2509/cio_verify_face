package model

import "github.com/google/uuid"

type CheckCompanyIsManagementInCompanyInput struct {
	CompanyID uuid.UUID `json:"company_id"`
	UserID    uuid.UUID `json:"user_id"`
}

type CheckDeviceExistsInCompanyInput struct {
	CompanyID uuid.UUID `json:"company_id"`
	DeviceID  uuid.UUID `json:"device_id"`
}

type UpdateDeviceSessionInput struct {
	DeviceId uuid.UUID `json:"device_id"`
	Token    string    `json:"token"`
}

type DeleteDeviceSessionInput struct {
	DeviceId uuid.UUID `json:"device_id"`
}
