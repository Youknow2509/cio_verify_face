package dto

import "time"

// ============================================
// Query DTOs
// ============================================

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

// ============================================
// Export Request DTOs
// ============================================

// ExportReportRequest represents request body for export report
type ExportReportRequest struct {
	StartDate string `json:"start_date" binding:"required" example:"2024-01-01"`
	EndDate   string `json:"end_date" binding:"required" example:"2024-01-31"`
	Format    string `json:"format" binding:"required,oneof=excel pdf csv" example:"excel"`
	CompanyID string `json:"company_id" binding:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Email     string `json:"email,omitempty" binding:"omitempty,email" example:"admin@example.com"`
}

// ExportDailyStatusRequest represents request to export daily status
type ExportDailyStatusRequest struct {
	CompanyID string `json:"company_id" binding:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Date      string `json:"date" binding:"required" example:"2024-01-15"`
	Format    string `json:"format" binding:"required,oneof=excel pdf csv" example:"excel"`
	Email     string `json:"email,omitempty" binding:"omitempty,email" example:"admin@example.com"`
}

// ExportMonthlySummaryRequest represents request to export monthly summary
type ExportMonthlySummaryRequest struct {
	CompanyID string `json:"company_id" binding:"required,uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Month     string `json:"month" binding:"required" example:"2024-01"`
	Format    string `json:"format" binding:"required,oneof=excel pdf csv" example:"excel"`
	Email     string `json:"email,omitempty" binding:"omitempty,email" example:"admin@example.com"`
}

// ExportEmployeeDailyStatusRequest represents request to export employee daily status
type ExportEmployeeDailyStatusRequest struct {
	Date   string `json:"date" binding:"required" example:"2024-01-15"`
	Format string `json:"format" binding:"required,oneof=excel pdf csv" example:"excel"`
	Email  string `json:"email,omitempty" binding:"omitempty,email" example:"employee@example.com"`
}

// ExportEmployeeMonthlySummaryRequest represents request to export employee monthly summary
type ExportEmployeeMonthlySummaryRequest struct {
	Month  string `json:"month" binding:"required" example:"2024-01"`
	Format string `json:"format" binding:"required,oneof=excel pdf csv" example:"excel"`
	Email  string `json:"email,omitempty" binding:"omitempty,email" example:"employee@example.com"`
}

// ============================================
// Helper Methods
// ============================================

// ParseDailyReportQuery parses and validates daily report query
func (q *DailyReportQuery) Parse() (time.Time, error) {
	date, err := time.Parse("2006-01-02", q.Date)
	if err != nil {
		return time.Time{}, err
	}
	return date, nil
}
