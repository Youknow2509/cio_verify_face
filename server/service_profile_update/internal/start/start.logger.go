package start

import (
	"fmt"

	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/domain/config"
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/global"
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/infrastructure/logger"
)

func initLogger(cfg *config.LoggerSetting) error {
	zapLogger, err := logger.NewZapLogger(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}
	global.Logger = zapLogger
	global.Logger.Info("Logger initialized successfully")
	return nil
}
