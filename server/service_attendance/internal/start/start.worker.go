package start

import (
	domainConfig "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/config"
	domainLogger "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/logger"
	domainRepo "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/repository"
	domainWorker "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/worker"
	domainWorkerAttendance "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/worker/attendance"
	"github.com/youknow2509/cio_verify_face/server/service_attendance/internal/global"
)

// ============================================
// Start Attendance service worker
// ============================================
func InitAttendanceServiceWorker(config *domainConfig.WorkerAttendanceSetting) error {
	worker := domainWorkerAttendance.NewAttendanceServiceWorker(
		*config,
		domainLogger.GetLogger(),
		domainRepo.GetAttendanceRepository(),
	)
	_ = domainWorker.SetWorkerAttendanceServiceWorker(worker)
	global.AttendanceServiceWorker = worker
	worker.RunDailySummaryWorker()
	return nil
}
