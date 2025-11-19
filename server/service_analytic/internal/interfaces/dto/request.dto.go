package dto

import "time"

// DailyReportQuery represents query parameters for daily report
type DailyReportQuery struct {
	Date      string `form:"date" binding:"required"` // YYYY-MM-DD
	CompanyID string `form:"company_id" binding:"required,uuid"`
	DeviceID  string `form:"device_id,omitempty" binding:"omitempty,uuid"`
}

// SummaryReportQuery represents query parameters for summary report
type SummaryReportQuery struct {
	Month     string `form:"month" binding:"required"` // YYYY-MM
	CompanyID string `form:"company_id" binding:"required,uuid"`
}

// ExportReportRequest represents request body for export report
type ExportReportRequest struct {
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
	Format    string `json:"format" binding:"required,oneof=excel pdf csv"`
	CompanyID string `json:"company_id" binding:"required,uuid"`
	Email     string `json:"email,omitempty" binding:"omitempty,email"`
}

// ParseDailyReportQuery parses and validates daily report query
func (q *DailyReportQuery) Parse() (time.Time, error) {
	date, err := time.Parse("2006-01-02", q.Date)
	if err != nil {
		return time.Time{}, err
	}
	return date, nil
}
