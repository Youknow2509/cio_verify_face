package start

import (
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/global"
)

// StartService initializes and starts the profile update service
func StartService() error {
	// Load configuration
	setting, err := loadConfig()
	if err != nil {
		return err
	}
	global.SettingServer = *setting

	// Initialize logger
	if err := initLogger(&setting.Logger); err != nil {
		return err
	}

	// Initialize local cache (Ristretto)
	if err := initLocalCache(); err != nil {
		return err
	}

	// Initialize auth service client
	if err := initAuthClient(); err != nil {
		return err
	}

	// Initialize infrastructure connections
	if err := initConnectionToInfrastructure(setting); err != nil {
		return err
	}

	// Initialize domain layer
	if err := initDomain(); err != nil {
		return err
	}

	// Initialize application layer
	if err := initApplication(); err != nil {
		return err
	}

	// Initialize HTTP router
	if err := initGinRouter(&setting.Server); err != nil {
		return err
	}

	return nil
}
