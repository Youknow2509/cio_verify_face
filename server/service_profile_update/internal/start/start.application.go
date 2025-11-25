package start

import (
	"fmt"

	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/application/service"
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/application/service/impl"
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/global"
)

func initApplication() error {
	// Initialize face profile update service
	// Construct base URL from server configuration
	baseURL := fmt.Sprintf("http://localhost:%d", global.SettingServer.Server.Port)
	faceProfileUpdateService := impl.NewFaceProfileUpdateService(baseURL)
	if err := service.SetFaceProfileUpdateService(faceProfileUpdateService); err != nil {
		return err
	}

	// Initialize password reset service
	passwordResetService := impl.NewPasswordResetService()
	if err := service.SetPasswordResetService(passwordResetService); err != nil {
		return err
	}

	global.Logger.Info("Application layer initialized successfully")
	return nil
}
