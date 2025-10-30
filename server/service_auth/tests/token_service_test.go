package tests

import (
	"context"
	"testing"

	infra "github.com/youknow2509/cio_verify_face/server/service_auth/internal/infrastructure/token"
)

// Test parse token service
func TestParseUserToken(t *testing.T) {
	// Sample config for token service - FOR TESTING PURPOSES ONLY
	// secret: 'your_jwt_secret_key'
	// issuer: 'cio_verify_face'
	// subject: 'cio_verify_face'
	// audience:
	//     - vinh
	//     - hihihi
	tokenService := infra.NewTokenService(
		"your_jwt_secret_key",
		"cio_verify_face",
		"cio_verify_face",
		[]string{"vinh", "hihihi"},
	)
	if tokenService == nil {
		t.Fatal("Failed to create token service")
	}

	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaW9fdmVyaWZ5X2ZhY2UiLCJzdWIiOiJjaW9fdmVyaWZ5X2ZhY2UiLCJhdWQiOlsidmluaCIsImhpaGloaSJdLCJleHAiOjE3NjE4NTcxNDEsIm5iZiI6MTc2MTg0OTk0MSwiaWF0IjoxNzYxODQ5OTQxLCJqdGkiOiIyNTg1NGMwZi1kNjI5LTQ4MWUtODNjOS1lOTE5OGUyN2ZkMzQiLCJ1c2VyX2lkIjoiNmMxYTdlMDEtYmYwNi00ZDVjLTliOTUtMTc0MjRiOWJkNGFjIiwicm9sZSI6MH0.XekFqlo6h-fsB7jBd7RQ0Y6JaMh1p3UVjIci7VsUeHI"
	ctx := context.Background()
	parsedToken, err := tokenService.ParseUserToken(ctx, token)
	if err != nil {
		t.Fatalf("Failed to parse user token: %v", err)
	}
	t.Logf("Parsed token: %+v", parsedToken)
}
