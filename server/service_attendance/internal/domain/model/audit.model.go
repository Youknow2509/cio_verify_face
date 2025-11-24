package model

import "github.com/google/uuid"

type AuditLog struct {
	CompanyID      uuid.UUID
	YearMonth      string
	CreatedAt      int64
	ActorID        uuid.UUID
	ActionCategory string
	ActionName     string
	ResourceType   string
	ResourceID     uuid.UUID
	Details        string
	IP_Address     string
	UserAgent      string
	Status         string
}
