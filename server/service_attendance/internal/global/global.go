package global

import (
	domainLogger "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/logger"
	domainConfig "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/config"
	"sync"
)

var (
	WaitGroup     *sync.WaitGroup
	Logger        domainLogger.ILogger
	SettingServer domainConfig.Setting
)
