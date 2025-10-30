package start

import (
	domainConfig "github.com/youknow2509/cio_verify_face/server/service_notify/internal/domain/config"
	domainLogger "github.com/youknow2509/cio_verify_face/server/service_notify/internal/domain/logger"
	"github.com/youknow2509/cio_verify_face/server/service_notify/internal/global"
	infraLogger "github.com/youknow2509/cio_verify_face/server/service_notify/internal/infrastructure/logger"
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
