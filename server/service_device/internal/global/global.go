package global

import (
	"sync"

	domainConfig "github.com/youknow2509/cio_verify_face/server/service_device/internal/domain/config"
	domainLogger "github.com/youknow2509/cio_verify_face/server/service_device/internal/domain/logger"
)

var (
	WaitGroup     *sync.WaitGroup
	Logger        domainLogger.ILogger
	SettingServer domainConfig.Setting
)
