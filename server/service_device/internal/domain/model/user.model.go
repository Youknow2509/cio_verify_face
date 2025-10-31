package model

import "github.com/google/uuid"

// user permission device - user and device same company
type UserPermissionDeviceInput struct {
	DeviceID uuid.UUID `json:"device_id"`
	UserID   uuid.UUID `json:"user_id"`
}

// user exists in company
type UserExistsInCompanyInput struct {
	CompanyID uuid.UUID `json:"company_id"`
	UserID    uuid.UUID `json:"user_id"`
}

// for GetCompanyIdOfUser
type GetCompanyIdOfUserInput struct {
	UserID uuid.UUID `json:"user_id"`
}

type GetCompanyIdOfUserOutput struct {
	CompanyID uuid.UUID `json:"company_id"`
}
