package dto

import (
	"github.com/google/uuid"
	"time"
)

// ============================================
// Daily Summary DTOs
// ============================================

// DailySummaryRequest represents request to create/update a daily summary
type DailySummaryRequest struct {
	CompanyID         string     `json:"company_id" binding:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	EmployeeID        string     `json:"employee_id" binding:"required,uuid" example:"660e8400-e29b-41d4-a716-446655440001"`
	WorkDate          time.Time  `json:"work_date" binding:"required" example:"2024-01-15T00:00:00Z"`
	CheckInTime       *time.Time `json:"check_in_time,omitempty" example:"2024-01-15T08:30:00Z"`
	CheckOutTime      *time.Time `json:"check_out_time,omitempty" example:"2024-01-15T17:30:00Z"`
	TotalWorkMinutes  int        `json:"total_work_minutes" example:"480"`
	LateMinutes       int        `json:"late_minutes" example:"0"`
	EarlyLeaveMinutes int        `json:"early_leave_minutes" example:"0"`
	IsPresent         bool       `json:"is_present" example:"true"`
	IsLate            bool       `json:"is_late" example:"false"`
	IsEarlyLeave      bool       `json:"is_early_leave" example:"false"`
	IsAbsent          bool       `json:"is_absent" example:"false"`
	AttendanceStatus  string     `json:"attendance_status" example:"present"`
	ShiftID           *string    `json:"shift_id,omitempty" example:"shift_001"`
	Notes             string     `json:"notes,omitempty" example:"Normal working day"`
}

// DailySummaryResponse represents a daily summary response
type DailySummaryResponse struct {
	CompanyID         uuid.UUID  `json:"company_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	EmployeeID        uuid.UUID  `json:"employee_id" example:"660e8400-e29b-41d4-a716-446655440001"`
	WorkDate          time.Time  `json:"work_date" example:"2024-01-15T00:00:00Z"`
	SummaryMonth      string     `json:"summary_month" example:"2024-01"`
	CheckInTime       *time.Time `json:"check_in_time,omitempty" example:"2024-01-15T08:30:00Z"`
	CheckOutTime      *time.Time `json:"check_out_time,omitempty" example:"2024-01-15T17:30:00Z"`
	TotalWorkMinutes  int        `json:"total_work_minutes" example:"480"`
	LateMinutes       int        `json:"late_minutes" example:"0"`
	EarlyLeaveMinutes int        `json:"early_leave_minutes" example:"0"`
	IsPresent         bool       `json:"is_present" example:"true"`
	IsLate            bool       `json:"is_late" example:"false"`
	IsEarlyLeave      bool       `json:"is_early_leave" example:"false"`
	IsAbsent          bool       `json:"is_absent" example:"false"`
	AttendanceStatus  string     `json:"attendance_status" example:"present"`
	ShiftID           *uuid.UUID `json:"shift_id,omitempty" example:"880e8400-e29b-41d4-a716-446655440003"`
	Notes             string     `json:"notes,omitempty" example:"Normal working day"`
	CreatedAt         time.Time  `json:"created_at" example:"2024-01-15T00:00:00Z"`
	UpdatedAt         time.Time  `json:"updated_at" example:"2024-01-15T18:00:00Z"`
}

// DailySummariesResponse represents a list of daily summaries
type DailySummariesResponse struct {
	Summaries    []DailySummaryResponse `json:"summaries"`
	TotalDays    int                    `json:"total_days"`
	SummaryMonth string                 `json:"summary_month,omitempty" example:"2024-01"`
}

type DailyReportDetailsRequest struct {
	CompanyId string `json:"company_id" binding:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Date      string `json:"date" binding:"required,datetime=2006-01-02" example:"2024-01-15"`
	PageSize  int    `json:"page_size" binding:"omitempty,min=1,max=100" example:"20"`
	PageState  string `json:"page_state,omitempty" example:"eyJQYWdlTnV4dCI6IjIifQ=="`
}
