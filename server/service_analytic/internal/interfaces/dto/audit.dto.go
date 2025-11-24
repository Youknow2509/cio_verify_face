package dto

import (
	"time"
	"github.com/google/uuid"
)

// ============================================
// Audit Log DTOs
// ============================================

// AuditLogRequest represents request to create an audit log
type AuditLogRequest struct {
	CompanyID     string    `json:"company_id" binding:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	ActorID       string    `json:"actor_id" binding:"required,uuid" example:"660e8400-e29b-41d4-a716-446655440001"`
	Action        string    `json:"action" binding:"required" example:"create_employee"`
	EntityType    string    `json:"entity_type" binding:"required" example:"employee"`
	EntityID      string    `json:"entity_id,omitempty" example:"770e8400-e29b-41d4-a716-446655440002"`
	Description   string    `json:"description,omitempty" example:"Created new employee record"`
	IPAddress     string    `json:"ip_address,omitempty" example:"192.168.1.100"`
	UserAgent     string    `json:"user_agent,omitempty" example:"Mozilla/5.0"`
	Changes       string    `json:"changes,omitempty" example:"{\"name\": \"old_value -> new_value\"}"`
}

// AuditLogResponse represents an audit log response
type AuditLogResponse struct {
	CompanyID     uuid.UUID  `json:"company_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	ActorID       uuid.UUID  `json:"actor_id" example:"660e8400-e29b-41d4-a716-446655440001"`
	CreatedAt     time.Time  `json:"created_at" example:"2024-01-15T09:00:00Z"`
	YearMonth     string     `json:"year_month" example:"2024-01"`
	Action        string     `json:"action" example:"create_employee"`
	EntityType    string     `json:"entity_type" example:"employee"`
	EntityID      *uuid.UUID `json:"entity_id,omitempty" example:"770e8400-e29b-41d4-a716-446655440002"`
	Description   string     `json:"description,omitempty" example:"Created new employee record"`
	IPAddress     string     `json:"ip_address,omitempty" example:"192.168.1.100"`
	UserAgent     string     `json:"user_agent,omitempty" example:"Mozilla/5.0"`
	Changes       string     `json:"changes,omitempty" example:"{\"name\": \"old_value -> new_value\"}"`
}

// AuditLogsResponse represents a list of audit logs
type AuditLogsResponse struct {
	Logs       []AuditLogResponse `json:"logs"`
	TotalLogs  int                `json:"total_logs"`
	YearMonth  string             `json:"year_month,omitempty" example:"2024-01"`
}

// ============================================
// Face Enrollment Log DTOs
// ============================================

// FaceEnrollmentLogRequest represents request to create a face enrollment log
type FaceEnrollmentLogRequest struct {
	CompanyID         string    `json:"company_id" binding:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	EmployeeID        string    `json:"employee_id" binding:"required,uuid" example:"660e8400-e29b-41d4-a716-446655440001"`
	EnrollmentType    string    `json:"enrollment_type" binding:"required,oneof=initial update removal" example:"initial"`
	FaceImageURLs     []string  `json:"face_image_urls" binding:"required" example:"https://storage.example.com/faces/img1.jpg,https://storage.example.com/faces/img2.jpg"`
	FaceEmbeddings    string    `json:"face_embeddings,omitempty" example:"[0.123, 0.456, ...]"`
	Quality           *float64  `json:"quality,omitempty" example:"0.95"`
	PerformedBy       string    `json:"performed_by,omitempty" example:"admin_user_id"`
	Status            string    `json:"status" example:"success"`
	ErrorMessage      string    `json:"error_message,omitempty" example:""`
	Notes             string    `json:"notes,omitempty" example:"Initial face enrollment for new employee"`
}

// FaceEnrollmentLogResponse represents a face enrollment log response
type FaceEnrollmentLogResponse struct {
	CompanyID         uuid.UUID  `json:"company_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	EmployeeID        uuid.UUID  `json:"employee_id" example:"660e8400-e29b-41d4-a716-446655440001"`
	CreatedAt         time.Time  `json:"created_at" example:"2024-01-15T09:00:00Z"`
	YearMonth         string     `json:"year_month" example:"2024-01"`
	EnrollmentType    string     `json:"enrollment_type" example:"initial"`
	FaceImageURLs     []string   `json:"face_image_urls" example:"https://storage.example.com/faces/img1.jpg,https://storage.example.com/faces/img2.jpg"`
	FaceEmbeddings    string     `json:"face_embeddings,omitempty" example:"[0.123, 0.456, ...]"`
	Quality           *float64   `json:"quality,omitempty" example:"0.95"`
	PerformedBy       *uuid.UUID `json:"performed_by,omitempty" example:"990e8400-e29b-41d4-a716-446655440004"`
	Status            string     `json:"status" example:"success"`
	ErrorMessage      string     `json:"error_message,omitempty" example:""`
	Notes             string     `json:"notes,omitempty" example:"Initial face enrollment for new employee"`
}

// FaceEnrollmentLogsResponse represents a list of face enrollment logs
type FaceEnrollmentLogsResponse struct {
	Logs       []FaceEnrollmentLogResponse `json:"logs"`
	TotalLogs  int                         `json:"total_logs"`
	YearMonth  string                      `json:"year_month,omitempty" example:"2024-01"`
}
