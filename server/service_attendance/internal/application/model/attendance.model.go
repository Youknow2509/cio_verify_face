package model

import (
	"time"

	"github.com/google/uuid"
)

// For GetMyRecords method
type GetMyRecordsInput struct {
	StartDate time.Time `json:"start_date,omitempty"`
	EndDate   time.Time `json:"end_date,omitempty"`
	Page      int       `json:"page,omitempty"`
	Size      int       `json:"size,omitempty"`
	// Session info
	UserID      uuid.UUID `json:"user_id"`
	SessionID   uuid.UUID `json:"session_id"`
	Role        int       `json:"role"`
	ClientIp    string    `json:"client_ip"`
	ClientAgent string    `json:"client_agent"`
}
type GetMyRecordsOutput struct {
	Records []AttendanceRecordInfo `json:"records"`
	Total   int                    `json:"total"`
}

// For GetRecordByID method
type GetAttendanceRecordByIDInput struct {
	RecordID uuid.UUID `json:"record_id"`
	// Session info
	UserID      uuid.UUID `json:"user_id"`
	SessionID   uuid.UUID `json:"session_id"`
	Role        int       `json:"role"`
	ClientIp    string    `json:"client_ip"`
	ClientAgent string    `json:"client_agent"`
}

// For GetRecords method
type GetAttendanceRecordsInput struct {
	DeviceID  uuid.UUID `json:"device_id,omitempty"`
	CompanyId uuid.UUID `json:"company_id,omitempty"`
	StartDate time.Time `json:"start_date,omitempty"`
	EndDate   time.Time `json:"end_date,omitempty"`
	Page      int       `json:"page,omitempty"`
	Size      int       `json:"size,omitempty"`
	// Session info
	UserID      uuid.UUID `json:"user_id"`
	SessionID   uuid.UUID `json:"session_id"`
	Role        int       `json:"role"`
	ClientIp    string    `json:"client_ip"`
	ClientAgent string    `json:"client_agent"`
}
type AttendanceRecordOutput struct {
	Records []AttendanceRecordInfo `json:"records"`
	Total   int                    `json:"total"`
}
type AttendanceRecordInfo struct {
	RecordID uuid.UUID `json:"record_id"`
	UserID   uuid.UUID `json:"user_id"`
	CheckIn  string    `json:"check_in"`
	CheckOut string    `json:"check_out"`
	Location string    `json:"location"`
	DeviceID uuid.UUID `json:"device_id"`
}

// For CheckIn method
type CheckInInput struct {
	UserCheckInId uuid.UUID `json:"user_checkin_id"`
	Timestamp     string    `json:"timestamp"`
	Location      string    `json:"location"`
	DeviceId      uuid.UUID `json:"device_id"`
	// Session info
	UserID      uuid.UUID `json:"user_id"`
	SessionID   uuid.UUID `json:"session_id"`
	Role        int       `json:"role"`
	ClientIp    string    `json:"client_ip"`
	ClientAgent string    `json:"client_agent"`
}

// For CheckOut method
type CheckOutInput struct {
	UserCheckOutId uuid.UUID `json:"user_check_out_id"`
	Timestamp      string    `json:"timestamp"`
	Location       string    `json:"location"`
	DeviceId       uuid.UUID `json:"device_id"`
	// Session info
	UserID      uuid.UUID `json:"user_id"`
	SessionID   uuid.UUID `json:"session_id"`
	Role        int       `json:"role"`
	ClientIp    string    `json:"client_ip"`
	ClientAgent string    `json:"client_agent"`
}
