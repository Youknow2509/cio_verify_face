package dto

import (
	"time"
	"github.com/google/uuid"
)

// ============================================
// Attendance Record DTOs
// ============================================

// AttendanceRecordRequest represents request to create an attendance record
type AttendanceRecordRequest struct {
	CompanyID      string    `json:"company_id" binding:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	EmployeeID     string    `json:"employee_id" binding:"required,uuid" example:"660e8400-e29b-41d4-a716-446655440001"`
	DeviceID       string    `json:"device_id" binding:"required,uuid" example:"770e8400-e29b-41d4-a716-446655440002"`
	RecordTime     time.Time `json:"record_time" binding:"required" example:"2024-01-15T09:00:00Z"`
	RecordType     string    `json:"record_type" binding:"required,oneof=check_in check_out" example:"check_in"`
	Temperature    *float64  `json:"temperature,omitempty" example:"36.5"`
	FaceImageURL   string    `json:"face_image_url,omitempty" example:"https://storage.example.com/faces/img123.jpg"`
	VerifyMethod   string    `json:"verify_method,omitempty" example:"face_recognition"`
	Confidence     *float64  `json:"confidence,omitempty" example:"0.98"`
	LocationLat    *float64  `json:"location_lat,omitempty" example:"21.0285"`
	LocationLng    *float64  `json:"location_lng,omitempty" example:"105.8542"`
	DeviceLocation string    `json:"device_location,omitempty" example:"Office Main Entrance"`
}

// AttendanceRecordResponse represents an attendance record response
type AttendanceRecordResponse struct {
	CompanyID      uuid.UUID  `json:"company_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	EmployeeID     uuid.UUID  `json:"employee_id" example:"660e8400-e29b-41d4-a716-446655440001"`
	DeviceID       uuid.UUID  `json:"device_id" example:"770e8400-e29b-41d4-a716-446655440002"`
	RecordTime     time.Time  `json:"record_time" example:"2024-01-15T09:00:00Z"`
	YearMonth      string     `json:"year_month" example:"2024-01"`
	RecordType     string     `json:"record_type" example:"check_in"`
	Temperature    *float64   `json:"temperature,omitempty" example:"36.5"`
	FaceImageURL   string     `json:"face_image_url,omitempty" example:"https://storage.example.com/faces/img123.jpg"`
	VerifyMethod   string     `json:"verify_method,omitempty" example:"face_recognition"`
	Confidence     *float64   `json:"confidence,omitempty" example:"0.98"`
	LocationLat    *float64   `json:"location_lat,omitempty" example:"21.0285"`
	LocationLng    *float64   `json:"location_lng,omitempty" example:"105.8542"`
	DeviceLocation string     `json:"device_location,omitempty" example:"Office Main Entrance"`
	CreatedAt      time.Time  `json:"created_at" example:"2024-01-15T09:00:00Z"`
}

// AttendanceRecordsResponse represents a list of attendance records
type AttendanceRecordsResponse struct {
	Records      []AttendanceRecordResponse `json:"records"`
	TotalRecords int                        `json:"total_records"`
	YearMonth    string                     `json:"year_month,omitempty" example:"2024-01"`
}

// AttendanceRecordNoShiftRequest represents request for attendance without shift
type AttendanceRecordNoShiftRequest struct {
	CompanyID    string    `json:"company_id" binding:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	EmployeeID   string    `json:"employee_id" binding:"required,uuid" example:"660e8400-e29b-41d4-a716-446655440001"`
	RecordTime   time.Time `json:"record_time" binding:"required" example:"2024-01-15T09:00:00Z"`
	RecordType   string    `json:"record_type" binding:"required,oneof=check_in check_out" example:"check_in"`
	EmployeeName string    `json:"employee_name,omitempty" example:"Nguyen Van A"`
	Department   string    `json:"department,omitempty" example:"IT Department"`
}

// AttendanceRecordNoShiftResponse represents response for attendance without shift
type AttendanceRecordNoShiftResponse struct {
	CompanyID    uuid.UUID `json:"company_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	EmployeeID   uuid.UUID `json:"employee_id" example:"660e8400-e29b-41d4-a716-446655440001"`
	RecordTime   time.Time `json:"record_time" example:"2024-01-15T09:00:00Z"`
	YearMonth    string    `json:"year_month" example:"2024-01"`
	RecordType   string    `json:"record_type" example:"check_in"`
	EmployeeName string    `json:"employee_name,omitempty" example:"Nguyen Van A"`
	Department   string    `json:"department,omitempty" example:"IT Department"`
}
