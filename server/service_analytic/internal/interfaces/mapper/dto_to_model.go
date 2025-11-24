package mapper

import (
	"time"

	"github.com/google/uuid"
	domainModel "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/domain/model"
	"github.com/youknow2509/cio_verify_face/server/service_analytic/internal/interfaces/dto"
)

// ============================================
// Audit Log Mappers
// ============================================

// DailySummaryRequestToModel converts DailySummaryRequest DTO to domain model
// ============================================
// Audit Log Mappers
// ============================================

// AuditLogRequestToModel converts AuditLogRequest DTO to domain model
func AuditLogRequestToModel(req *dto.AuditLogRequest) (*domainModel.AuditLog, error) {
	companyID, err := uuid.Parse(req.CompanyID)
	if err != nil {
		return nil, err
	}

	actorID, err := uuid.Parse(req.ActorID)
	if err != nil {
		return nil, err
	}

	createdAt := time.Now()
	yearMonth := createdAt.Format("2006-01")

	// Parse changes JSON into map if provided
	details := make(map[string]string)
	if req.Changes != "" {
		details["changes"] = req.Changes
	}
	if req.Description != "" {
		details["description"] = req.Description
	}

	return &domainModel.AuditLog{
		CompanyID:      companyID,
		ActorID:        actorID,
		CreatedAt:      createdAt,
		YearMonth:      yearMonth,
		ActionCategory: "general",  // Default category
		ActionName:     req.Action,
		ResourceType:   req.EntityType,
		ResourceID:     req.EntityID,
		Details:        details,
		IPAddress:      req.IPAddress,
		UserAgent:      req.UserAgent,
		Status:         "success",  // Default status
	}, nil
}

// FaceEnrollmentLogRequestToModel converts FaceEnrollmentLogRequest DTO to domain model