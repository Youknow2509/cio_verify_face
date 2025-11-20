package global

import (
	domainConfig "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/config"
	domainLogger "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/logger"
	domainWorker "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/worker"
	"sync"
)

var (
	WaitGroup               *sync.WaitGroup
	Logger                  domainLogger.ILogger
	SettingServer           domainConfig.Setting
	AttendanceServiceWorker domainWorker.IWorkerAttendanceServiceWorker
)
