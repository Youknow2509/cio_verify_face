package dto

// ============================================
// Export Response DTOs
// ============================================

// ExportJobResponse represents response after creating an export job
type ExportJobResponse struct {
	JobID   string `json:"job_id" example:"job_660e8400-e29b-41d4-a716-446655440001"`
	Status  string `json:"status" example:"processing"`
	Message string `json:"message" example:"Export job queued successfully"`
}

// ExportJobStatusResponse represents the status of an export job
type ExportJobStatusResponse struct {
	JobID        string `json:"job_id" example:"job_660e8400-e29b-41d4-a716-446655440001"`
	Status       string `json:"status" example:"completed"`
	Progress     int    `json:"progress" example:"100"`
	DownloadURL  string `json:"download_url,omitempty" example:"https://storage.example.com/exports/report_123.xlsx"`
	ErrorMessage string `json:"error_message,omitempty" example:""`
	CreatedAt    string `json:"created_at" example:"2024-01-15T09:00:00Z"`
	CompletedAt  string `json:"completed_at,omitempty" example:"2024-01-15T09:05:00Z"`
}
