package start

import (
	"fmt"

	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/domain/repository"
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/global"
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/infrastructure/conn"
	infraRepo "github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/infrastructure/repository"
)

func initDomain() error {
	// Get PostgreSQL client
	pgClient, err := conn.GetPostgresqlClient()
	if err != nil {
		return fmt.Errorf("failed to get PostgreSQL client: %w", err)
	}

	// Initialize User Repository
	userRepo := infraRepo.NewUserRepository(pgClient)
	if err := repository.SetUserRepository(userRepo); err != nil {
		return fmt.Errorf("failed to set user repository: %w", err)
	}
	global.Logger.Info("User repository initialized")

	// Initialize Face Profile Update Request Repository
	faceProfileRepo := infraRepo.NewFaceProfileUpdateRequestRepository(pgClient)
	if err := repository.SetFaceProfileUpdateRequestRepository(faceProfileRepo); err != nil {
		return fmt.Errorf("failed to set face profile update request repository: %w", err)
	}
	global.Logger.Info("Face profile update request repository initialized")

	// Initialize Password Reset Request Repository
	passwordResetRepo := infraRepo.NewPasswordResetRequestRepository(pgClient)
	if err := repository.SetPasswordResetRequestRepository(passwordResetRepo); err != nil {
		return fmt.Errorf("failed to set password reset request repository: %w", err)
	}
	global.Logger.Info("Password reset request repository initialized")

	global.Logger.Info("Domain layer initialized successfully")
	return nil
}
