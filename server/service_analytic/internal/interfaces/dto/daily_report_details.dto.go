package dto

import (
	"time"

	"github.com/google/uuid"
)

// DailyReportDetailsResponse represents the response for detailed daily report
// Used in: GET /api/v1/daily-summaries/details
// Contains total, items (list of employees' daily report details), and next_page for pagination
type DailyReportDetailsResponse struct {
	Total    int                         `json:"total"`
	Items    []DailyReportDetailEmployee `json:"items"`
	NextPage string                      `json:"next_page,omitempty"`
}

type DailyReportDetailEmployee struct {
	CompanyID            uuid.UUID  `json:"company_id"`
	SummaryMonth         string     `json:"summary_month"` // YYYY-MM format
	WorkDate             time.Time  `json:"work_date"`
	EmployeeID           uuid.UUID  `json:"employee_id"`
	ShiftID              uuid.UUID  `json:"shift_id"`
	ActualCheckIn        *time.Time `json:"actual_check_in"`
	ActualCheckOut       *time.Time `json:"actual_check_out"`
	AttendanceStatus     int        `json:"attendance_status"`
	LateMinutes          int        `json:"late_minutes"`
	EarlyLeaveMinutes    int        `json:"early_leave_minutes"`
	TotalWorkMinutes     int        `json:"total_work_minutes"`
	Notes                string     `json:"notes"`
	UpdatedAt            time.Time  `json:"updated_at"`
	OvertimeMinutes      int        `json:"overtime_minutes"`      // Calculated field
	AttendancePercentage float64    `json:"attendance_percentage"` // Calculated field
}
