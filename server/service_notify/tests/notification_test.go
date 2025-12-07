package tests

import (
	"encoding/json"
	"testing"

	applicationModel "github.com/youknow2509/cio_verify_face/server/service_notify/internal/application/model"
	"github.com/youknow2509/cio_verify_face/server/service_notify/internal/interfaces/dto"
)

/**
 * Test unmarshaling password reset notification event
 */
func TestPasswordResetNotificationUnmarshal(t *testing.T) {
	// Sample JSON from problem statement
	jsonData := `{
		"event_type": 11,
		"metadata": {
			"created_at": "2025-12-07T07:43:22Z",
			"message_id": "1af454f9-97f6-4629-9e30-b253376221b9",
			"request_id": "d40365dc-f850-4623-84a5-c141c364ce52",
			"user_id": "e84b21ae-4d3e-4c2b-b638-d177e55a32f7"
		},
		"payload": {
			"expires_in": 24,
			"full_name": "Tran Thi B",
			"reset_url": "https://your-domain.com/api/v1/password/reset/confirm?token=0336413f00169fad8c5bda3245c043dd73cdd573d9f11771d727115066cb2e3d",
			"to": "employee2.fpt@example.com"
		}
	}`

	var event dto.PasswordResetNotificationEvent
	err := json.Unmarshal([]byte(jsonData), &event)
	if err != nil {
		t.Fatalf("Failed to unmarshal password reset notification: %v", err)
	}

	// Verify event type
	if event.EventType != 11 {
		t.Errorf("Expected event_type 11, got %d", event.EventType)
	}

	// Verify metadata
	if event.Metadata.UserID != "e84b21ae-4d3e-4c2b-b638-d177e55a32f7" {
		t.Errorf("Expected user_id e84b21ae-4d3e-4c2b-b638-d177e55a32f7, got %s", event.Metadata.UserID)
	}

	// Verify payload
	if event.Payload.FullName != "Tran Thi B" {
		t.Errorf("Expected full_name 'Tran Thi B', got %s", event.Payload.FullName)
	}
	if event.Payload.To != "employee2.fpt@example.com" {
		t.Errorf("Expected to 'employee2.fpt@example.com', got %s", event.Payload.To)
	}
	if event.Payload.ExpiresIn != 24 {
		t.Errorf("Expected expires_in 24, got %d", event.Payload.ExpiresIn)
	}

	// Convert to application model
	model := applicationModel.PasswordResetNotification{
		To:        event.Payload.To,
		FullName:  event.Payload.FullName,
		ResetURL:  event.Payload.ResetURL,
		ExpiresIn: event.Payload.ExpiresIn,
	}

	if model.To != "employee2.fpt@example.com" {
		t.Errorf("Model conversion failed for To field")
	}

	t.Log("Password reset notification unmarshal test passed")
}

/**
 * Test unmarshaling report attention notification event
 */
func TestReportAttentionNotificationUnmarshal(t *testing.T) {
	// Sample JSON from problem statement
	jsonData := `{
		"company_id": "72f1eb04-3909-45ff-800d-63a9a9832d54",
		"created_at": "2025-12-07T19:16:08Z",
		"download_url": "http://127.0.0.1:9000/reports/reports/2025/12/08/72f1eb04-3909-45ff-800d-63a9a9832d54_daily_detail_2024-01-15.csv?X-Amz-Algorithm=AWS4-HMAC-SHA256",
		"email": "admin@example.com",
		"end_date": "2024-01-15",
		"format": "csv",
		"start_date": "2024-01-15",
		"type": "export_report"
	}`

	var event dto.ReportAttentionNotificationEvent
	err := json.Unmarshal([]byte(jsonData), &event)
	if err != nil {
		t.Fatalf("Failed to unmarshal report attention notification: %v", err)
	}

	// Verify fields
	if event.CompanyID != "72f1eb04-3909-45ff-800d-63a9a9832d54" {
		t.Errorf("Expected company_id 72f1eb04-3909-45ff-800d-63a9a9832d54, got %s", event.CompanyID)
	}
	if event.Email != "admin@example.com" {
		t.Errorf("Expected email 'admin@example.com', got %s", event.Email)
	}
	if event.Type != "export_report" {
		t.Errorf("Expected type 'export_report', got %s", event.Type)
	}
	if event.Format != "csv" {
		t.Errorf("Expected format 'csv', got %s", event.Format)
	}
	if event.StartDate != "2024-01-15" {
		t.Errorf("Expected start_date '2024-01-15', got %s", event.StartDate)
	}

	// Convert to application model
	model := applicationModel.ReportAttentionNotification{
		Email:       event.Email,
		CompanyID:   event.CompanyID,
		DownloadURL: event.DownloadURL,
		Type:        event.Type,
		Format:      event.Format,
		StartDate:   event.StartDate,
		EndDate:     event.EndDate,
		CreatedAt:   event.CreatedAt,
	}

	if model.Email != "admin@example.com" {
		t.Errorf("Model conversion failed for Email field")
	}

	t.Log("Report attention notification unmarshal test passed")
}
