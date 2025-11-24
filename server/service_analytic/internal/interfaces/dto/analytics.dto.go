package dto

// ============================================
// Analytics DTOs - Company Admin
// ============================================

// DailyAttendanceStatusResponse represents comprehensive daily attendance status
type DailyAttendanceStatusResponse struct {
	Date               string                  `json:"date" example:"2024-01-15"`
	TotalRecords       int                     `json:"total_records" example:"150"`
	TotalEmployees     int                     `json:"total_employees" example:"75"`
	AttendanceRecords  []AttendanceRecordResponse `json:"attendance_records"`
	DailySummaries     []DailySummaryResponse     `json:"daily_summaries"`
}

// AttendanceStatusRangeResponse represents attendance status for a time range
type AttendanceStatusRangeResponse struct {
	StartDate          string                  `json:"start_date" example:"2024-01-01"`
	EndDate            string                  `json:"end_date" example:"2024-01-31"`
	TotalRecords       int                     `json:"total_records" example:"2500"`
	TotalDays          int                     `json:"total_days" example:"31"`
	AttendanceRecords  []AttendanceRecordResponse `json:"attendance_records"`
	DailySummaries     []DailySummaryResponse     `json:"daily_summaries"`
}

// EmployeeStatistics represents statistics for a single employee
type EmployeeStatistics struct {
	EmployeeID          string  `json:"employee_id" example:"660e8400-e29b-41d4-a716-446655440001"`
	EmployeeName        string  `json:"employee_name,omitempty" example:"Nguyen Van A"`
	TotalDays           int     `json:"total_days" example:"22"`
	PresentDays         int     `json:"present_days" example:"20"`
	AbsentDays          int     `json:"absent_days" example:"2"`
	LateDays            int     `json:"late_days" example:"3"`
	TotalWorkMinutes    int     `json:"total_work_minutes" example:"9600"`
	TotalLateMinutes    int     `json:"total_late_minutes" example:"45"`
	TotalEarlyLeaveMinutes int  `json:"total_early_leave_minutes" example:"30"`
	AttendanceRate      float64 `json:"attendance_rate" example:"90.91"`
	PunctualityRate     float64 `json:"punctuality_rate" example:"85.00"`
	AvgWorkHoursPerDay  float64 `json:"avg_work_hours_per_day" example:"8.0"`
}

// MonthlySummaryResponse represents comprehensive monthly summary
type MonthlySummaryResponse struct {
	Month                   string               `json:"month" example:"2024-01"`
	TotalEmployees          int                  `json:"total_employees" example:"75"`
	TotalWorkingDays        int                  `json:"total_working_days" example:"22"`
	TotalAttendanceRecords  int                  `json:"total_attendance_records" example:"2500"`
	TotalNoShiftRecords     int                  `json:"total_no_shift_records" example:"15"`
	CompanyStatistics       CompanyStatistics    `json:"company_statistics"`
	EmployeeStatistics      []EmployeeStatistics `json:"employee_statistics"`
}

// CompanyStatistics represents company-level statistics
type CompanyStatistics struct {
	TotalPresentDays       int     `json:"total_present_days" example:"1500"`
	TotalAbsentDays        int     `json:"total_absent_days" example:"150"`
	TotalLateDays          int     `json:"total_late_days" example:"225"`
	TotalWorkMinutes       int     `json:"total_work_minutes" example:"720000"`
	TotalLateMinutes       int     `json:"total_late_minutes" example:"3375"`
	TotalEarlyLeaveMinutes int     `json:"total_early_leave_minutes" example:"2250"`
	AvgAttendanceRate      float64 `json:"avg_attendance_rate" example:"90.91"`
	AvgPunctualityRate     float64 `json:"avg_punctuality_rate" example:"85.00"`
	AvgWorkHoursPerEmployee float64 `json:"avg_work_hours_per_employee" example:"160.0"`
}

// ============================================
// Analytics DTOs - Employee Self-Service
// ============================================

// EmployeeDailyStatusResponse represents employee's daily attendance status
type EmployeeDailyStatusResponse struct {
	Date               string                     `json:"date" example:"2024-01-15"`
	TotalRecords       int                        `json:"total_records" example:"2"`
	AttendanceRecords  []AttendanceRecordResponse `json:"attendance_records"`
	DailySummary       *DailySummaryResponse      `json:"daily_summary,omitempty"`
}

// EmployeeStatusRangeResponse represents employee's attendance status for a time range
type EmployeeStatusRangeResponse struct {
	EmployeeID         string                     `json:"employee_id" example:"660e8400-e29b-41d4-a716-446655440001"`
	StartDate          string                     `json:"start_date" example:"2024-01-01"`
	EndDate            string                     `json:"end_date" example:"2024-01-31"`
	TotalRecords       int                        `json:"total_records" example:"44"`
	TotalDays          int                        `json:"total_days" example:"31"`
	AttendanceRecords  []AttendanceRecordResponse `json:"attendance_records"`
	DailySummaries     []DailySummaryResponse     `json:"daily_summaries"`
}

// EmployeeMonthlySummaryResponse represents employee's detailed monthly summary
type EmployeeMonthlySummaryResponse struct {
	EmployeeID             string                 `json:"employee_id" example:"660e8400-e29b-41d4-a716-446655440001"`
	Month                  string                 `json:"month" example:"2024-01"`
	TotalDays              int                    `json:"total_days" example:"22"`
	PresentDays            int                    `json:"present_days" example:"20"`
	AbsentDays             int                    `json:"absent_days" example:"2"`
	LateDays               int                    `json:"late_days" example:"3"`
	TotalWorkMinutes       int                    `json:"total_work_minutes" example:"9600"`
	TotalLateMinutes       int                    `json:"total_late_minutes" example:"45"`
	TotalEarlyLeaveMinutes int                    `json:"total_early_leave_minutes" example:"30"`
	AttendanceRate         float64                `json:"attendance_rate" example:"90.91"`
	PunctualityRate        float64                `json:"punctuality_rate" example:"85.00"`
	AvgWorkHoursPerDay     float64                `json:"avg_work_hours_per_day" example:"8.0"`
	DailySummaries         []DailySummaryResponse `json:"daily_summaries"`
}

// EmployeeStatsResponse represents employee's overall statistics
type EmployeeStatsResponse struct {
	EmployeeID             string  `json:"employee_id" example:"660e8400-e29b-41d4-a716-446655440001"`
	Month                  string  `json:"month" example:"2024-01"`
	TotalDays              int     `json:"total_days" example:"22"`
	PresentDays            int     `json:"present_days" example:"20"`
	AbsentDays             int     `json:"absent_days" example:"2"`
	LateDays               int     `json:"late_days" example:"3"`
	EarlyLeaveDays         int     `json:"early_leave_days" example:"1"`
	TotalWorkMinutes       int     `json:"total_work_minutes" example:"9600"`
	TotalLateMinutes       int     `json:"total_late_minutes" example:"45"`
	TotalEarlyLeaveMinutes int     `json:"total_early_leave_minutes" example:"30"`
	AttendanceRate         float64 `json:"attendance_rate" example:"90.91"`
	PunctualityRate        float64 `json:"punctuality_rate" example:"85.00"`
	AvgWorkHoursPerDay     float64 `json:"avg_work_hours_per_day" example:"8.0"`
	TotalWorkHours         float64 `json:"total_work_hours" example:"160.0"`
}
