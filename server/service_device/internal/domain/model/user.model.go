package model

import "github.com/google/uuid"

type UserExistsInCompanyInput struct {
	CompanyID uuid.UUID `json:"company_id"`
	UserID    uuid.UUID `json:"user_id"`
}
