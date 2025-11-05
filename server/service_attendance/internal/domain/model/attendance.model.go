package model

import "github.com/google/uuid"

// ============================================
// Attendance Record Model
// ============================================

// For GetAttendanceRecordRangeTimeWithUserId
type GetAttendanceRecordRangeTimeWithUserIdInput struct {
	CompanyID uuid.UUID
	UserID    uuid.UUID
	StartTime int64
	EndTime   int64
	Limit     int
	Offset    int
}

// For GetAttendanceRecordRangeTime
type GetAttendanceRecordRangeTimeInput struct {
	CompanyID uuid.UUID
	DeviceID  uuid.UUID
	StartTime int64
	EndTime   int64
	Limit     int
	Offset    int
}

// Attendance Record structure
type AttendanceRecord struct {
	CompanyID           uuid.UUID
	EmployeeID          uuid.UUID
	DeviceID            uuid.UUID
	RecordTime          int64
	Type                int // 0: CHECK_IN, 1: CHECK_OUT
	VerificationMethod  string
	VerificationScore   float32
	FaceImageURL        string
	LocationCoordinates string // "lat,lng" format
	Metadata            map[string]string
	CreatedAt           int64
}

// For adding a check-out record
type AddCheckOutRecordInput struct {
	CompanyID           uuid.UUID
	EmployeeID          uuid.UUID
	DeviceID            uuid.UUID
	VerificationMethod  string
	VerificationScore   float32
	FaceImageURL        string
	LocationCoordinates string // "lat,lng" format
	Metadata            map[string]string
	RecordTime          int64
}

// For adding a check-in record
type AddCheckInRecordInput struct {
	CompanyID           uuid.UUID
	EmployeeID          uuid.UUID
	DeviceID            uuid.UUID
	VerificationMethod  string
	VerificationScore   float32
	FaceImageURL        string
	LocationCoordinates string // "lat,lng" format
	Metadata            map[string]string
	RecordTime          int64
}
