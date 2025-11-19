package model

// =================================
// Role definitions:
// =================================
type Role int

const (
	RoleEmployee     Role = 2 // Regular employee with limited access
	RoleCompanyAdmin Role = 1 // Company admin/manager with elevated access
	RoleSystemAdmin  Role = 0 // System admin/root with full access
)
