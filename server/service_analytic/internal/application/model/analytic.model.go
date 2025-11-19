package model

import "time"

// SessionInfo represents authenticated session information
type SessionInfo struct {
	UserID      string `json:"user_id"`
	Role        int32  `json:"role"`
	SessionID   string `json:"session_id"`
	ClientIP    string `json:"client_ip"`
	ClientAgent string `json:"client_agent"`
	CompanyID   string `json:"company_id"`
}

// DailyReportInput represents input for daily report
type DailyReportInput struct {
	Session   *SessionInfo `json:"-"` // Session info for authorization
	Date      time.Time    `json:"date" binding:"required"`
	CompanyID *string      `json:"company_id,omitempty"`
	DeviceID  *string      `json:"device_id,omitempty"`
}

// DailyReportOutput represents output for daily report
type DailyReportOutput struct {
	Date                string             `json:"date"`
	TotalEmployees      int                `json:"total_employees"`
	PresentEmployees    int                `json:"present_employees"`
	LateEmployees       int                `json:"late_employees"`
	EarlyLeaveEmployees int                `json:"early_leave_employees"`
	AbsentEmployees     int                `json:"absent_employees"`
	AttendanceRate      float64            `json:"attendance_rate"`
	Departments         []DepartmentReport `json:"departments"`
	Shifts              []ShiftReport      `json:"shifts"`
}

// DepartmentReport represents department-wise attendance
type DepartmentReport struct {
	DepartmentName   string  `json:"department_name"`
	TotalEmployees   int     `json:"total_employees"`
	PresentEmployees int     `json:"present_employees"`
	AttendanceRate   float64 `json:"attendance_rate"`
}

// ShiftReport represents shift-wise attendance
type ShiftReport struct {
	ShiftName        string  `json:"shift_name"`
	StartTime        string  `json:"start_time"`
	EndTime          string  `json:"end_time"`
	TotalEmployees   int     `json:"total_employees"`
	PresentEmployees int     `json:"present_employees"`
	AttendanceRate   float64 `json:"attendance_rate"`
}

// SummaryReportInput represents input for summary report
type SummaryReportInput struct {
	Session   *SessionInfo `json:"-"` // Session info for authorization
	Month     string       `json:"month" binding:"required"` // Format: YYYY-MM
	CompanyID *string      `json:"company_id,omitempty"`
}

// SummaryReportOutput represents output for summary report
type SummaryReportOutput struct {
	Month                  string                   `json:"month"`
	TotalWorkingDays       int                      `json:"total_working_days"`
	TotalEmployees         int                      `json:"total_employees"`
	AverageAttendanceRate  float64                  `json:"average_attendance_rate"`
	TotalWorkingHours      int                      `json:"total_working_hours"`
	TotalOvertimeHours     int                      `json:"total_overtime_hours"`
	WeeklySummary          []WeeklySummary          `json:"weekly_summary"`
	TopAttendanceEmployees []EmployeeAttendanceStat `json:"top_attendance_employees"`
	LowAttendanceEmployees []EmployeeAttendanceStat `json:"low_attendance_employees"`
}

// WeeklySummary represents weekly attendance summary
type WeeklySummary struct {
	Week           int     `json:"week"`
	StartDate      string  `json:"start_date"`
	EndDate        string  `json:"end_date"`
	AttendanceRate float64 `json:"attendance_rate"`
	TotalHours     int     `json:"total_hours"`
}

// EmployeeAttendanceStat represents employee attendance statistics
type EmployeeAttendanceStat struct {
	EmployeeCode   string  `json:"employee_code"`
	FullName       string  `json:"full_name"`
	PresentDays    int     `json:"present_days"`
	AttendanceRate float64 `json:"attendance_rate"`
	TotalHours     int     `json:"total_hours"`
}

// ExportReportInput represents input for export report
type ExportReportInput struct {
	Session   *SessionInfo `json:"-"` // Session info for authorization
	StartDate string       `json:"start_date" binding:"required"`
	EndDate   string       `json:"end_date" binding:"required"`
	Format    string       `json:"format" binding:"required,oneof=excel pdf csv"`
	CompanyID *string      `json:"company_id,omitempty"`
	Email     *string      `json:"email,omitempty"`
}

// ExportReportOutput represents output for export report
type ExportReportOutput struct {
	JobID       string  `json:"job_id"`
	Status      string  `json:"status"`
	Message     string  `json:"message"`
	DownloadURL *string `json:"download_url,omitempty"`
}

// HealthCheckOutput represents health check response
type HealthCheckOutput struct {
	Status   string                 `json:"status"`
	Version  string                 `json:"version"`
	Services map[string]interface{} `json:"services"`
}
