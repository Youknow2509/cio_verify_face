package start

import (
	"github.com/youknow2509/cio_verify_face/server/service_attendance/internal/global"
	domainConfig "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/config"
	domainLogger "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/logger"
	infraLogger "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/infrastructure/logger"
)

func initLogger(setting *domainConfig.LoggerSetting) error {
	dataInitLogger := &infraLogger.ZapLoggerInitializer{
		FolderStore:    setting.FolderStore,
		FileMaxSize:    setting.FileMaxSize,
		FileMaxBackups: setting.FileMaxBackups,
		FileMaxAge:     setting.FileMaxAge,
		Compress:       setting.Compress,
	}
	loggerServiceImpl, er := infraLogger.NewZapLogger(dataInitLogger)
	if er != nil {
		return er
	}
	err := domainLogger.SetLogger(loggerServiceImpl)
	if err != nil {
		return err
	}
	global.Logger = loggerServiceImpl
	return nil
}
