package attendance

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/google/uuid"
	domainConfig "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/config"
	domainLogger "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/logger"
	domainModel "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/model"
	domainRepo "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/repository"
	"github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/worker"
	"github.com/youknow2509/cio_verify_face/server/service_attendance/internal/global"
)

// ============================================
// Worker for Attendance service
// ============================================

type AttendanceServiceWorker struct {
	logger         domainLogger.ILogger
	attendanceRepo domainRepo.IAttendanceRepository
	summaryJobChan chan *domainModel.AddDailySummariesInput
	config         domainConfig.WorkerAttendanceSetting
}

// Running worker to process daily summary jobs
func (w *AttendanceServiceWorker) RunDailySummaryWorker() error {
	if w.config.NumWorkers <= 0 {
		w.logger.Warn("number of daily summary workers is set to 0 or less, skipping worker startup")
		return errors.New("number worker less than 0")
	}
	for i := 0; i < w.config.NumWorkers; i++ {
		global.WaitGroup.Add(1)
		go func(workerID int) {
			defer global.WaitGroup.Done()
			ctx := context.Background()
			for job := range w.summaryJobChan {
				w.logger.Info("processing daily summary job", "worker number", workerID, "employeeID", job.EmployeeID, "workDate", job.WorkDate) // TODO: remove log in production
				if err := w.attendanceRepo.AddDailySummaries(ctx, job); err != nil {
					w.logger.Error("add daily simmaries", "worker number", workerID, "error", err)
				}
			}
		}(i)
	}
	return nil
}

// add job to worker channel v2
func (w *AttendanceServiceWorker) AddJobToDailySummaryWorkerV2(
	companyID uuid.UUID,
	employeeID uuid.UUID,
	recordTime time.Time,
	matchedShift domainModel.ShiftTimeEmployee,
) {
	ctx := context.Background()
	workDate := time.Date(recordTime.Year(), recordTime.Month(), recordTime.Day(), 0, 0, 0, 0, recordTime.Location())
	w.logger.Info("recordTime -> work date", "recordTime", recordTime, "workDate", workDate) // TODO: remove log in production
	dailySummary, err := w.calculateDailySummary(ctx, companyID, employeeID, workDate, recordTime, matchedShift)
	if err != nil {
		w.logger.Error("calculate daily summary", "error", err)
		return
	}
	if dailySummary != nil {
		select {
		case w.summaryJobChan <- dailySummary:
			w.logger.Info("added daily summary job to channel", "employeeID", employeeID, "workDate", workDate) // TODO: remove log in production
		default:
			w.logger.Warn("daily summary job channel is full, dropping job")
			ctx := context.Background()
			if err := w.attendanceRepo.AddDailySummaries(ctx, dailySummary); err != nil {
				w.logger.Error("failed to process synchronous fallback", "error", err)
			}
		}
	} else {
		w.logger.Warn("no daily summary to add, skipping", "employeeID", employeeID, "workDate", workDate) // TODO: remove log in production
	}
}

// add job to worker channel
func (w *AttendanceServiceWorker) AddJobToDailySummaryWorker(job *domainModel.AddDailySummariesInput) {
	select {
	case w.summaryJobChan <- job:
	default:
		w.logger.Warn("daily summary job channel is full, dropping job")
		ctx := context.Background()
		if err := w.attendanceRepo.AddDailySummaries(ctx, job); err != nil {
			w.logger.Error("failed to process synchronous fallback", "error", err)
		}
	}
}

func NewAttendanceServiceWorker(
	config domainConfig.WorkerAttendanceSetting,
	logger domainLogger.ILogger,
	attendanceRepo domainRepo.IAttendanceRepository,
) worker.IWorkerAttendanceServiceWorker {
	// initialize worker
	return &AttendanceServiceWorker{
		logger:         logger,
		attendanceRepo: attendanceRepo,
		summaryJobChan: make(chan *domainModel.AddDailySummariesInput, config.SizeBufferChan),
		config:         config,
	}
}

// ============================================
// Helper functions
// ============================================

// calculateDailySummary tính toán thông tin tổng hợp chấm công hàng ngày
func (w *AttendanceServiceWorker) calculateDailySummary(
	ctx context.Context,
	companyID uuid.UUID,
	employeeID uuid.UUID,
	workDate time.Time,
	checkOutTime time.Time,
	shift domainModel.ShiftTimeEmployee,
) (*domainModel.AddDailySummariesInput, error) {

	// 1. Dựng lại shiftStart/shiftEnd theo cùng logic dùng cho check-in (qua đêm)
	shiftStart, shiftEnd := buildShiftBounds(workDate, checkOutTime, shift.StartTime, shift.EndTime)

	yearMonth := workDate.Format("2006-01")

	checkInRecord, err := w.attendanceRepo.GetFirstCheckIn(ctx, &domainModel.GetFirstCheckInInput{
		CompanyID:      companyID,
		EmployeeID:     employeeID,
		YearMonth:      yearMonth,
		ShiftTimeStart: shift.StartTime, // time-of-day
		ShiftTimeEnd:   shift.EndTime,   // time-of-day
		DateCheckOut:   checkOutTime,
	})
	if err != nil {
		return nil, err
	}

	// Nếu không có check-in → tạo summary vắng mặt (tuỳ yêu cầu)
	if checkInRecord == nil {
		return &domainModel.AddDailySummariesInput{
			CompanyID:         companyID,
			SummaryMonth:      yearMonth,
			WorkDate:          workDate,
			EmployeeID:        employeeID,
			ShiftID:           shift.ShiftID,
			ActualCheckIn:     time.Time{},
			ActualCheckOut:    checkOutTime,
			AttendanceStatus:  domainModel.StatusAbsent,
			LateMinutes:       0,
			EarlyLeaveMinutes: 0,
			TotalWorkMinutes:  0,
			Notes:             "Absent: no check-in",
			UpdatedAt:         time.Now().UTC(),
		}, nil
	}

	actualCheckIn := checkInRecord.RecordTime

	// 2. Lateness
	grace := time.Duration(shift.GracePeriodMinutes) * time.Minute
	lateness := actualCheckIn.Sub(shiftStart) - grace
	if lateness < 0 {
		lateness = 0
	}
	lateMinutes := int(lateness.Minutes())

	// 3. Early leave
	allowedEarly := time.Duration(shift.EarlyDepartureMinutes) * time.Minute
	earlyLeave := shiftEnd.Sub(checkOutTime) - allowedEarly
	if earlyLeave < 0 {
		earlyLeave = 0
	}
	earlyLeaveMinutes := int(earlyLeave.Minutes())

	// 4. Total work
	totalWorkMinutes := int(checkOutTime.Sub(actualCheckIn).Minutes())
	if totalWorkMinutes < 0 {
		totalWorkMinutes = 0
	}

	// 5. Overtime (nếu cần)
	overtimeMinutes := 0
	if checkOutTime.After(shiftEnd) {
		overtimeMinutes = int(checkOutTime.Sub(shiftEnd).Minutes())
	}
	_ = overtimeMinutes // currently not used

	// 6. Scheduled minutes
	scheduledMinutes := int(shiftEnd.Sub(shiftStart).Minutes())
	if scheduledMinutes < 0 {
		scheduledMinutes = 0
	}

	// 7. Attendance status
	var attendanceStatus int
	switch {
	case lateMinutes > 0 && earlyLeaveMinutes > 0:
		attendanceStatus = domainModel.StatusLateAndEarlyLeave
	case lateMinutes > 0:
		attendanceStatus = domainModel.StatusLate
	case earlyLeaveMinutes > 0:
		attendanceStatus = domainModel.StatusEarlyLeave
	default:
		attendanceStatus = domainModel.StatusPresent
	}

	// Edge: inconsistent times
	notes := buildAttendanceNotes(lateMinutes, earlyLeaveMinutes, totalWorkMinutes, shift)
	if actualCheckIn.After(checkOutTime) {
		attendanceStatus = domainModel.StatusAbsent
		notes += " | Inconsistent timestamps (check-in after check-out)"
		totalWorkMinutes = 0
	}

	// 8. Trả về
	return &domainModel.AddDailySummariesInput{
		CompanyID:         companyID,
		SummaryMonth:      yearMonth,
		WorkDate:          workDate,
		EmployeeID:        employeeID,
		ShiftID:           shift.ShiftID,
		ActualCheckIn:     actualCheckIn,
		ActualCheckOut:    checkOutTime,
		AttendanceStatus:  attendanceStatus,
		LateMinutes:       lateMinutes,
		EarlyLeaveMinutes: earlyLeaveMinutes,
		TotalWorkMinutes:  totalWorkMinutes,
		// OvertimeMinutes:   overtimeMinutes,  // (thêm field nếu model hỗ trợ)
		// ScheduledMinutes:  scheduledMinutes, // (thêm field nếu model hỗ trợ)
		Notes:             notes,
		UpdatedAt:         time.Now().UTC(),
	}, nil
}

// buildShiftBounds: chuẩn hóa xử lý ca qua đêm
func buildShiftBounds(workDate, checkOut time.Time, startTimeOfDay, endTimeOfDay time.Time) (time.Time, time.Time) {
	// Giờ/phút/giây:
	shStartH, shStartM, shStartS := startTimeOfDay.Hour(), startTimeOfDay.Minute(), startTimeOfDay.Second()
	shEndH, shEndM, shEndS := endTimeOfDay.Hour(), endTimeOfDay.Minute(), endTimeOfDay.Second()

	// Giả định workDate là ngày bắt đầu ca
	year, month, day := workDate.Date()
	loc := workDate.Location()

	start := time.Date(year, month, day, shStartH, shStartM, shStartS, 0, loc)
	end := time.Date(year, month, day, shEndH, shEndM, shEndS, 0, loc)

	// Ca qua đêm
	if end.Before(start) {
		end = end.Add(24 * time.Hour)
		// Nếu checkOut thuộc ngày tiếp theo và workDate = ngày của start thì không đổi.
	}

	return start, end
}

// buildAttendanceNotes tạo ghi chú cho bản tổng hợp chấm công
func buildAttendanceNotes(lateMinutes, earlyLeaveMinutes, totalWorkMinutes int, shift domainModel.ShiftTimeEmployee) string {
	var notes string

	if lateMinutes > 0 {
		notes += "Late: " + strconv.Itoa(lateMinutes) + " minutes. "
	}
	if earlyLeaveMinutes > 0 {
		notes += "Early leave: " + strconv.Itoa(earlyLeaveMinutes) + " minutes. "
	}

	// Tính thời gian ca làm việc tiêu chuẩn
	standardWorkMinutes := int(shift.EndTime.Sub(shift.StartTime).Minutes())
	if standardWorkMinutes < 0 {
		standardWorkMinutes += 24 * 60 // Ca đêm
	}

	notes += "Total work: " + strconv.Itoa(totalWorkMinutes) + "/" + strconv.Itoa(standardWorkMinutes) + " minutes."

	return notes
}
