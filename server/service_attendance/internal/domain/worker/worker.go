package worker

import (
	"errors"
	"time"

	"github.com/google/uuid"
	domainModel "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/model"
)

// ============================================
// Worker
// ============================================

// For worker attendance service
type IWorkerAttendanceServiceWorker interface {
	RunDailySummaryWorker() error
	AddJobToDailySummaryWorker(job *domainModel.AddDailySummariesInput)
	AddJobToDailySummaryWorkerV2(companyID uuid.UUID, employeeID uuid.UUID, recordTime time.Time, matchedShift domainModel.ShiftTimeEmployee)
}

var _vIWorkerAttendanceServiceWorker IWorkerAttendanceServiceWorker

func GetWorkerAttendanceServiceWorker() IWorkerAttendanceServiceWorker {
	return _vIWorkerAttendanceServiceWorker
}

func SetWorkerAttendanceServiceWorker(worker IWorkerAttendanceServiceWorker) error {
	if worker == nil {
		return errors.New("worker attendance service worker is nil")
	}
	if _vIWorkerAttendanceServiceWorker != nil {
		return errors.New("worker attendance service worker is already set")
	}
	_vIWorkerAttendanceServiceWorker = worker
	return nil
}
