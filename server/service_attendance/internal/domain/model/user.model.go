package model

import "github.com/google/uuid"

// ============================================
// User Model
// ============================================

// For UserIsManagerCompany
type UserIsManagerCompanyInput struct {
	UserID    uuid.UUID
	CompanyID uuid.UUID
}

// For UserIsEmployeeInCompany
type UserIsEmployeeInCompanyInput struct {
	UserID    uuid.UUID
	CompanyID uuid.UUID
}

// For GetCompanyIdUser
type GetCompanyIdUserInput struct {
	UserID uuid.UUID
}

type GetCompanyIdUserOutput struct {
	CompanyID uuid.UUID
}
