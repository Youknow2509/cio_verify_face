package repository

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/model"
	domainRepo "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/repository"
)

// ============================================
// Attendance repository
// ============================================
type AttendanceRepository struct {
	dbSession *gocql.Session
}

// DeleteDailySummariesCompany implements repository.IAttendanceRepository.
func (a *AttendanceRepository) DeleteDailySummariesCompany(ctx context.Context, input *model.DeleteDailySummariesInput) error {
	// 1. Get All Records
	summaryMonth := input.WorkDate.Format("2006-01")

	fetchSql := `
		SELECT work_date, employee_id FROM daily_summaries
		WHERE company_id = ? AND summary_month = ? AND work_date = ?;
	`
	iter := a.dbSession.Query(fetchSql,
		marshalUuid(input.CompanyID),
		summaryMonth,
		input.WorkDate,
	).WithContext(ctx).Iter()
	var workDate time.Time
	var gocqlUUIDEmployeeID gocql.UUID
	var employeeIDs []gocql.UUID
	for iter.Scan(&workDate, &gocqlUUIDEmployeeID) {
		employeeIDs = append(employeeIDs, gocqlUUIDEmployeeID)
	}
	if err := iter.Close(); err != nil {
		return err
	}
	if len(employeeIDs) == 0 {
		return nil // không có gì để xoá
	}
	// 2. Delete Records by Batch
	const batchSize = 50
	deleteCompanySQL := `
		DELETE FROM daily_summaries
		WHERE company_id = ? AND summary_month = ? AND work_date = ? AND employee_id = ?;`
	for i := 0; i < len(employeeIDs); i += batchSize {
		end := i + batchSize
		if end > len(employeeIDs) {
			end = len(employeeIDs)
		}
		batch := a.dbSession.NewBatch(gocql.LoggedBatch)
		for _, empID := range employeeIDs[i:end] {
			batch.Query(deleteCompanySQL,
				marshalUuid(input.CompanyID),
				summaryMonth,
				input.WorkDate,
				empID,
			)
		}
		if err := a.dbSession.ExecuteBatch(batch); err != nil {
			return err
		}
	}

	return nil
}

// DeleteDailySummariesCompanyBeforeDate implements repository.IAttendanceRepository.
func (a *AttendanceRepository) DeleteDailySummariesCompanyBeforeDate(
	ctx context.Context,
	input *model.DeleteDailySummariesInput,
) error {

	summaryMonth := input.WorkDate.Format("2006-01")

	// 1. Lấy danh sách employee_id + work_date
	fetchSql := `
        SELECT work_date, employee_id
        FROM daily_summaries
        WHERE company_id = ? AND summary_month = ?
          AND work_date <= ?
    `

	iter := a.dbSession.Query(fetchSql,
		marshalUuid(input.CompanyID),
		summaryMonth,
		input.WorkDate,
	).WithContext(ctx).Iter()

	var wd time.Time
	var emp gocql.UUID
	type Record struct {
		WorkDate   time.Time
		EmployeeID gocql.UUID
	}
	var items []Record

	for iter.Scan(&wd, &emp) {
		items = append(items, Record{WorkDate: wd, EmployeeID: emp})
	}
	if err := iter.Close(); err != nil {
		return err
	}
	if len(items) == 0 {
		return nil
	}

	// 2. Delete theo batch
	const batchSize = 50
	deleteSQL := `
        DELETE FROM daily_summaries
        WHERE company_id = ? AND summary_month = ? 
          AND work_date = ? AND employee_id = ?;
    `

	for i := 0; i < len(items); i += batchSize {
		end := i + batchSize
		if end > len(items) {
			end = len(items)
		}

		batch := a.dbSession.NewBatch(gocql.LoggedBatch)

		for _, r := range items[i:end] {
			batch.Query(deleteSQL,
				marshalUuid(input.CompanyID),
				summaryMonth,
				r.WorkDate,
				r.EmployeeID,
			)
		}

		if err := a.dbSession.ExecuteBatch(batch.WithContext(ctx)); err != nil {
			return err
		}
	}

	return nil
}

// DeleteDailySummariesEmployee implements repository.IAttendanceRepository.
func (a *AttendanceRepository) DeleteDailySummariesEmployee(
	ctx context.Context,
	input *model.DeleteDailySummariesEmployeeInput,
) error {
	summaryMonth := input.WorkDate.Format("2006-01")
	sql_raw := `DELETE FROM daily_summaries_by_user
		WHERE company_id = ? AND employee_id = ?
		AND summary_month = ? AND work_date = ?;`
	err := a.dbSession.Query(sql_raw,
		marshalUuid(input.CompanyID),
		marshalUuid(input.EmployeeID),
		summaryMonth,
		input.WorkDate,
	).WithContext(ctx).Exec()
	if err != nil {
		return err
	}
	return nil
}

// DeleteDailySummariesEmployeeBeforeDate implements repository.IAttendanceRepository.
func (a *AttendanceRepository) DeleteDailySummariesEmployeeBeforeDate(ctx context.Context, input *model.DeleteDailySummariesEmployeeInput) error {
	summaryMonth := input.WorkDate.Format("2006-01")
	// 1. Lấy danh sách work_date
	fetchSql := `
		SELECT work_date FROM daily_summaries_by_user
		WHERE company_id = ? AND employee_id = ?
		  AND summary_month = ? AND work_date <= ?;
	`
	iter := a.dbSession.Query(fetchSql,
		marshalUuid(input.CompanyID),
		marshalUuid(input.EmployeeID),
		summaryMonth,
		input.WorkDate,
	).WithContext(ctx).Iter()
	var wd time.Time
	var workDates []time.Time
	for iter.Scan(&wd) {
		workDates = append(workDates, wd)
	}
	if err := iter.Close(); err != nil {
		return err
	}
	if len(workDates) == 0 {
		return nil
	}
	// 2. Delete theo batch
	const batchSize = 50
	deleteSQL := `
		DELETE FROM daily_summaries_by_user
		WHERE company_id = ? AND employee_id = ?
		  AND summary_month = ? AND work_date = ?;
	`
	for i := 0; i < len(workDates); i += batchSize {
		end := i + batchSize
		if end > len(workDates) {
			end = len(workDates)
		}
		batch := a.dbSession.NewBatch(gocql.LoggedBatch)
		for _, wd := range workDates[i:end] {
			batch.Query(deleteSQL,
				marshalUuid(input.CompanyID),
				marshalUuid(input.EmployeeID),
				summaryMonth,
				wd,
			)
		}
		if err := a.dbSession.ExecuteBatch(batch.WithContext(ctx)); err != nil {
			return err
		}
	}
	return nil
}

// DeleteAttendanceRecord implements repository.IAttendanceRepository.
func (a *AttendanceRepository) DeleteAttendanceRecord(ctx context.Context, input *model.DeleteAttendanceRecordInput) error {
	// BEGIN BATCH
	//     DELETE FROM attendance_records
	//     WHERE company_id = uuid_company AND year_month = '2023-10'
	//     AND record_time = '2023-10-25 08:00:00' AND employee_id = uuid_employee;
	//     DELETE FROM attendance_records_by_user
	//     WHERE company_id = uuid_company AND employee_id = uuid_employee
	//     AND year_month = '2023-10' AND record_time = '2023-10-25 08:00:00';
	// APPLY BATCH;
	sql_raw := `BEGIN BATCH
	DELETE FROM attendance_records
		WHERE company_id = ? AND year_month = ? 
		AND record_time = ? AND employee_id = ?;
	DELETE FROM attendance_records_by_user
		WHERE company_id = ? AND employee_id = ?
		AND year_month = ? AND record_time = ?;
	APPLY BATCH;`
	yearMonth := input.RecordTime.Format("2006-01")
	err := a.dbSession.Query(sql_raw,
		// first delete
		marshalUuid(input.CompanyID),
		yearMonth,
		input.RecordTime,
		marshalUuid(input.EmployeeID),
		// second delete
		marshalUuid(input.CompanyID),
		marshalUuid(input.EmployeeID),
		yearMonth,
		input.RecordTime,
	).WithContext(ctx).Exec()
	if err != nil {
		if errors.Is(err, gocql.ErrNotFound) {
			return nil
		}
		return err
	}
	return nil
}

// GetDailySummarieCompany implements repository.IAttendanceRepository.
func (a *AttendanceRepository) GetDailySummarieCompany(ctx context.Context, input *model.GetDailySummariesCompanyInput) (*model.DailySummariesCompanyOutput, error) {
	// SELECT * FROM daily_summaries
	// WHERE company_id = uuid_company AND summary_month = '2023-10' AND work_date = '2023-10-25';
	sql_raw := `SELECT company_id, summary_month, work_date, employee_id,
		shift_id, actual_check_in, actual_check_out, attendance_status,
		late_minutes, early_leave_minutes, total_work_minutes, notes, updated_at
		FROM daily_summaries
		WHERE company_id = ? AND summary_month = ? AND work_date = ?;`
	var iter *gocql.Iter
	output := &model.DailySummariesCompanyOutput{}
	if input.PageSize > 0 && input.PageStage != nil {
		query := a.dbSession.Query(sql_raw,
			marshalUuid(input.CompanyID),
			input.SummaryMonth,
			input.WorkDate,
		).WithContext(ctx).PageSize(input.PageSize).PageState(input.PageStage)
		iter = query.Iter()
		output.PageStageNext = iter.PageState()
	} else {
		iter = a.dbSession.Query(sql_raw,
			marshalUuid(input.CompanyID),
			input.SummaryMonth,
			input.WorkDate,
		).WithContext(ctx).Iter()
		output.PageStageNext = nil
	}
	var listDailySummaries []model.DailySummariesCompanyInfo
	for {
		var r model.DailySummariesCompanyInfo
		gocqlUUIDCompanyID := gocql.UUID{}
		gocqlUUIDEmployeeID := gocql.UUID{}
		gocqlUUIDShiftID := gocql.UUID{}
		if !iter.Scan(
			&gocqlUUIDCompanyID,
			&r.SummaryMonth,
			&r.WorkDate,
			&gocqlUUIDEmployeeID,
			&gocqlUUIDShiftID,
			&r.ActualCheckIn,
			&r.ActualCheckOut,
			&r.AttendanceStatus,
			&r.LateMinutes,
			&r.EarlyLeaveMinutes,
			&r.TotalWorkMinutes,
			&r.Notes,
			&r.UpdatedAt,
		) {
			break
		}
		// Unmarshal
		r.CompanyId = uuid.UUID(gocqlUUIDCompanyID)
		r.EmployeeId = uuid.UUID(gocqlUUIDEmployeeID)
		r.ShiftId = uuid.UUID(gocqlUUIDShiftID)
		//
		listDailySummaries = append(listDailySummaries, r)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	output.Records = listDailySummaries
	output.PageSize = len(listDailySummaries)
	return output, nil
}

// GetDailySummarieCompanyForEmployee implements repository.IAttendanceRepository.
func (a *AttendanceRepository) GetDailySummarieCompanyForEmployee(ctx context.Context, input *model.GetDailySummariesCompanyForEmployeeInput) (*model.DailySummariesEmployeeOutput, error) {
	// SELECT * FROM daily_summaries_by_user
	// WHERE company_id = uuid_company AND employee_id = uuid_employee AND summary_month = '2023-10';
	sql_raw := `SELECT company_id, summary_month, work_date, employee_id,
		shift_id, actual_check_in, actual_check_out, attendance_status,
		late_minutes, early_leave_minutes, total_work_minutes, notes, updated_at
		FROM daily_summaries_by_user
		WHERE company_id = ? AND summary_month = ? AND employee_id = ?;`
	var iter *gocql.Iter
	output := &model.DailySummariesEmployeeOutput{}
	if input.PageSize > 0 && input.PageStage != nil {
		query := a.dbSession.Query(sql_raw,
			marshalUuid(input.CompanyID),
			input.SummaryMonth,
			marshalUuid(input.EmployeeID),
		).WithContext(ctx).PageSize(input.PageSize).PageState(input.PageStage)
		iter = query.Iter()
		output.PageStageNext = iter.PageState()
	} else {
		iter = a.dbSession.Query(sql_raw,
			marshalUuid(input.CompanyID),
			input.SummaryMonth,
			marshalUuid(input.EmployeeID),
		).WithContext(ctx).Iter()
		output.PageStageNext = nil
	}
	var listDailySummaries []model.DailySummariesEmployeeInfo
	for {
		var r model.DailySummariesEmployeeInfo
		gocqlUUIDCompanyID := gocql.UUID{}
		gocqlUUIDEmployeeID := gocql.UUID{}
		gocqlUUIDShiftID := gocql.UUID{}
		if !iter.Scan(
			&gocqlUUIDCompanyID,
			&r.SummaryMonth,
			&r.WorkDate,
			&gocqlUUIDEmployeeID,
			&gocqlUUIDShiftID,
			&r.ActualCheckIn,
			&r.ActualCheckOut,
			&r.AttendanceStatus,
			&r.LateMinutes,
			&r.EarlyLeaveMinutes,
			&r.TotalWorkMinutes,
			&r.Notes,
			&r.UpdatedAt,
		) {
			break
		}
		// Unmarshal
		r.CompanyId = uuid.UUID(gocqlUUIDCompanyID)
		r.EmployeeId = uuid.UUID(gocqlUUIDEmployeeID)
		r.ShiftId = uuid.UUID(gocqlUUIDShiftID)
		//
		listDailySummaries = append(listDailySummaries, r)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	output.Records = listDailySummaries
	output.PageSize = len(listDailySummaries)
	return output, nil
}

// UpdateDailySummariesEmployee implements repository.IAttendanceRepository.
func (a *AttendanceRepository) UpdateDailySummariesEmployee(ctx context.Context, input *model.UpdateDailySummariesEmployeeInput) error {
	// BEGIN BATCH
	//     -- Update note bảng Company
	//     UPDATE daily_summaries SET notes = 'Đã duyệt phép', attendance_status = 0
	//     WHERE company_id = uuid_company AND summary_month = '2023-10'
	//     AND work_date = '2023-10-25' AND employee_id = uuid_employee;
	//	-- Update note bảng User
	//	UPDATE daily_summaries_by_user SET notes = 'Đã duyệt phép', attendance_status = 0
	//	WHERE company_id = uuid_company AND employee_id = uuid_employee
	//	AND summary_month = '2023-10' AND work_date = '2023-10-25';
	//
	// APPLY BATCH;
	sql_raw := `BEGIN BATCH
	UPDATE daily_summaries SET notes = ?, attendance_status = ?
		WHERE company_id = ? AND summary_month = ?
		AND work_date = ? AND employee_id = ?;
	UPDATE daily_summaries_by_user SET notes = ?, attendance_status = ?
		WHERE company_id = ? AND employee_id = ?
		AND summary_month = ? AND work_date = ?;
	APPLY BATCH;`
	err := a.dbSession.Query(sql_raw,
		input.Notes,
		input.AttendanceStatus,
		marshalUuid(input.CompanyID),
		input.SummaryMonth,
		input.WorkDate,
		marshalUuid(input.EmployeeID),
		input.Notes,
		input.AttendanceStatus,
		marshalUuid(input.CompanyID),
		marshalUuid(input.EmployeeID),
		input.SummaryMonth,
		input.WorkDate,
	).WithContext(ctx).Exec()
	return err
}

// AddAttendanceRecord implements repository.IAttendanceRepository.
func (a *AttendanceRepository) AddAttendanceRecord(ctx context.Context, input *model.AddAttendanceRecordInput) error {
	// 	BEGIN BATCH
	//     -- 1. Ghi vào bảng Company
	//     INSERT INTO attendance_records (
	//         company_id, year_month, record_time, employee_id,
	//         device_id, record_type, verification_method, verification_score,
	//         face_image_url, location_coordinates, metadata, sync_status, created_at
	//     ) VALUES (
	//         uuid_company, '2023-10', '2023-10-25 08:00:00', uuid_employee,
	//         uuid_device, 0, 'FACE', 0.98,
	//         'http://minio/img.jpg', '10.7,106.6', {'ip': '1.2.3.4'}, 'synced', toTimestamp(now())
	//     );
	//     -- 2. Ghi vào bảng User
	//     INSERT INTO attendance_records_by_user (
	//         company_id, employee_id, year_month, record_time,
	//         device_id, record_type, verification_method, verification_score,
	//         face_image_url, location_coordinates, metadata, sync_status, created_at
	//     ) VALUES (
	//         uuid_company, uuid_employee, '2023-10', '2023-10-25 08:00:00',
	//         uuid_device, 0, 'FACE', 0.98,
	//         'http://minio/img.jpg', '10.7,106.6', {'ip': '1.2.3.4'}, 'synced', toTimestamp(now())
	//     );
	// APPLY BATCH;
	sql_raw := `BEGIN BATCH
	INSERT INTO attendance_records (
		company_id, year_month, record_time, employee_id,
		device_id, record_type, verification_method, verification_score,
		face_image_url, location_coordinates, metadata, sync_status, created_at
	) VALUES (
		?, ?, ?, ?,
		?, ?, ?, ?,
		?, ?, ?, 'synced', ?
	);
	INSERT INTO attendance_records_by_user (
		company_id, employee_id, year_month, record_time,
		device_id, record_type, verification_method, verification_score,
		face_image_url, location_coordinates, metadata, sync_status, created_at
	) VALUES (
		?, ?, ?, ?,
		?, ?, ?, ?,
		?, ?, ?, 'synced', ?
	);
	APPLY BATCH;`
	yearMonth := input.RecordTime.Format("2006-01")
	createdAt := time.Now()
	//
	err := a.dbSession.Query(sql_raw,
		marshalUuid(input.CompanyID),
		yearMonth,
		input.RecordTime,
		marshalUuid(input.EmployeeID),
		marshalUuid(input.DeviceID),
		input.RecordType,
		input.VerificationMethod,
		// marshalFloat64ToFloat32()
		marshalFloat64ToFloat32(input.VerificationScore),
		input.FaceImageURL,
		input.LocationCoordinates,
		input.Metadata,
		createdAt,
		// second insert
		marshalUuid(input.CompanyID),
		marshalUuid(input.EmployeeID),
		yearMonth,
		input.RecordTime,
		marshalUuid(input.DeviceID),
		input.RecordType,
		input.VerificationMethod,
		marshalFloat64ToFloat32(input.VerificationScore),
		input.FaceImageURL,
		input.LocationCoordinates,
		input.Metadata,
		createdAt,
	).WithContext(ctx).Exec()
	if err != nil {
		return err
	}
	return nil
}

// AddDailySummaries implements repository.IAttendanceRepository.
func (a *AttendanceRepository) AddDailySummaries(ctx context.Context, input *model.AddDailySummariesInput) error {
	// BEGIN BATCH
	//     -- 1. Cập nhật bảng Company
	//     INSERT INTO daily_summaries (
	//         company_id, summary_month, work_date, employee_id,
	//         shift_id, actual_check_in, actual_check_out, attendance_status,
	//         late_minutes, early_leave_minutes, total_work_minutes, notes, updated_at
	//     ) VALUES (
	//         uuid_company, '2023-10', '2023-10-25', uuid_employee,
	//         uuid_shift, '2023-10-25 08:00:00', '2023-10-25 17:30:00', 1,
	//         15, 0, 480, 'Đi muộn do kẹt xe', toTimestamp(now())
	//     );
	//	-- 2. Cập nhật bảng User
	//	INSERT INTO daily_summaries_by_user (
	//	    company_id, employee_id, summary_month, work_date,
	//	    shift_id, actual_check_in, actual_check_out, attendance_status,
	//	    late_minutes, early_leave_minutes, total_work_minutes, notes, updated_at
	//	) VALUES (
	//	    uuid_company, uuid_employee, '2023-10', '2023-10-25',
	//	    uuid_shift, '2023-10-25 08:00:00', '2023-10-25 17:30:00', 1,
	//	    15, 0, 480, 'Đi muộn do kẹt xe', toTimestamp(now())
	//	);
	//
	// APPLY BATCH;
	sql_raw := `BEGIN BATCH
	INSERT INTO daily_summaries (
		company_id, summary_month, work_date, employee_id,
		shift_id, actual_check_in, actual_check_out, attendance_status,
		late_minutes, early_leave_minutes, total_work_minutes, notes, updated_at
	) VALUES (
		?, ?, ?, ?,
		?, ?, ?, ?,
		?, ?, ?, ?, ?
	);
	INSERT INTO daily_summaries_by_user (
		company_id, employee_id, summary_month, work_date,
		shift_id, actual_check_in, actual_check_out, attendance_status,
		late_minutes, early_leave_minutes, total_work_minutes, notes, updated_at
	) VALUES (
		?, ?, ?, ?,
		?, ?, ?, ?,
		?, ?, ?, ?, ?
	);
	APPLY BATCH;`
	updatedAt := time.Now()
	err := a.dbSession.Query(sql_raw,
		marshalUuid(input.CompanyID),
		input.SummaryMonth,
		input.WorkDate,
		marshalUuid(input.EmployeeID),
		marshalUuid(input.ShiftID),
		input.ActualCheckIn,
		input.ActualCheckOut,
		input.AttendanceStatus,
		input.LateMinutes,
		input.EarlyLeaveMinutes,
		input.TotalWorkMinutes,
		input.Notes,
		updatedAt,
		// second insert
		marshalUuid(input.CompanyID),
		marshalUuid(input.EmployeeID),
		input.SummaryMonth,
		input.WorkDate,
		marshalUuid(input.ShiftID),
		input.ActualCheckIn,
		input.ActualCheckOut,
		input.AttendanceStatus,
		input.LateMinutes,
		input.EarlyLeaveMinutes,
		input.TotalWorkMinutes,
		input.Notes,
		updatedAt,
	).WithContext(ctx).Exec()
	if err != nil {
		return err
	}
	return nil
}

// DeleteAttendanceRecordBeforeTimestamp implements repository.IAttendanceRepository.
func (a *AttendanceRepository) DeleteAttendanceRecordBeforeTimestamp(
	ctx context.Context,
	input *model.DeleteAttendanceRecordInput,
) error {

	yearMonth := input.RecordTime.Format("2006-01")

	// ============================
	// Lấy danh sách record_time từ bảng by_user
	// ============================
	fetchSql := `
        SELECT record_time FROM attendance_records_by_user
        WHERE company_id = ? AND employee_id = ? AND year_month = ?
          AND record_time <= ?
    `
	iter := a.dbSession.Query(fetchSql,
		marshalUuid(input.CompanyID),
		marshalUuid(input.EmployeeID),
		yearMonth,
		input.RecordTime,
	).WithContext(ctx).Iter()

	var rt time.Time
	var recordTimes []time.Time

	for iter.Scan(&rt) {
		recordTimes = append(recordTimes, rt)
	}
	if err := iter.Close(); err != nil {
		return err
	}

	if len(recordTimes) == 0 {
		return nil // không có gì để xoá
	}

	// ============================
	// Xoá từng record theo PK đầy đủ bằng batch (tối đa 50/lần)
	// ============================
	const batchSize = 50
	deleteMainSQL := `
        DELETE FROM attendance_records
        WHERE company_id = ? AND year_month = ?
          AND record_time = ? AND employee_id = ?
    `
	deleteUserSQL := `
        DELETE FROM attendance_records_by_user
        WHERE company_id = ? AND employee_id = ?
          AND year_month = ? AND record_time = ?
    `

	for i := 0; i < len(recordTimes); i += batchSize {
		end := i + batchSize
		if end > len(recordTimes) {
			end = len(recordTimes)
		}

		batch := a.dbSession.NewBatch(gocql.LoggedBatch)

		for _, rt := range recordTimes[i:end] {

			// xóa bảng chính
			batch.Query(deleteMainSQL,
				marshalUuid(input.CompanyID),
				yearMonth,
				rt,
				marshalUuid(input.EmployeeID),
			)

			// xóa bảng phụ
			batch.Query(deleteUserSQL,
				marshalUuid(input.CompanyID),
				marshalUuid(input.EmployeeID),
				yearMonth,
				rt,
			)
		}

		if err := a.dbSession.ExecuteBatch(batch.WithContext(ctx)); err != nil {
			return err
		}
	}

	return nil
}

// GetAttendanceRecordCompany implements repository.IAttendanceRepository.
func (a *AttendanceRepository) GetAttendanceRecordCompany(ctx context.Context, input *model.GetAttendanceRecordCompanyInput) (*model.AttendanceRecordOutput, error) {
	// SELECT * FROM attendance_records
	// WHERE company_id = uuid_company AND year_month = '2023-10';
	sql_raw := `SELECT company_id, year_month, record_time, employee_id,
		device_id, record_type, verification_method, verification_score,
		face_image_url, location_coordinates, metadata, sync_status, created_at
		FROM attendance_records
		WHERE company_id = ? AND year_month = ?;`
	// Check have page
	var iter *gocql.Iter
	output := &model.AttendanceRecordOutput{}
	if input.PageSize > 0 && input.PageStage != nil {
		query := a.dbSession.Query(sql_raw,
			marshalUuid(input.CompanyID),
			input.YearMonth,
		).WithContext(ctx).PageSize(input.PageSize).PageState(input.PageStage)
		iter = query.Iter()
		output.PageStageNext = iter.PageState()
	} else {
		iter = a.dbSession.Query(sql_raw,
			marshalUuid(input.CompanyID),
			input.YearMonth,
		).WithContext(ctx).Iter()
		output.PageStageNext = nil
	}
	var listAttendanceRecords []model.AttendanceRecordInfo
	for {
		var r model.AttendanceRecordInfo
		gocqlUUIDCompanyID := gocql.UUID{}
		gocqlUUIDEmployeeID := gocql.UUID{}
		gocqlUUIDDeviceID := gocql.UUID{}
		gocqlFloatVerificationScore := float32(0)
		if !iter.Scan(
			&gocqlUUIDCompanyID,
			&r.YearMonth,
			&r.RecordTime,
			&gocqlUUIDEmployeeID,
			&gocqlUUIDDeviceID,
			&r.RecordType,
			&r.VerificationMethod,
			&gocqlFloatVerificationScore,
			&r.FaceImageURL,
			&r.LocationCoordinates,
			&r.Metadata,
			&r.SyncStatus,
			&r.CreatedAt,
		) {
			break
		}
		// Unmarshal
		r.CompanyID = uuid.UUID(gocqlUUIDCompanyID)
		r.EmployeeID = uuid.UUID(gocqlUUIDEmployeeID)
		r.DeviceID = uuid.UUID(gocqlUUIDDeviceID)
		r.VerificationScore = float64(gocqlFloatVerificationScore)
		//
		listAttendanceRecords = append(listAttendanceRecords, r)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	output.Records = listAttendanceRecords
	output.PageSize = len(listAttendanceRecords)
	return output, nil
}

// GetAttendanceRecordCompanyForEmployee implements repository.IAttendanceRepository.
func (a *AttendanceRepository) GetAttendanceRecordCompanyForEmployee(ctx context.Context, input *model.GetAttendanceRecordCompanyForEmployeeInput) (*model.AttendanceRecordOutput, error) {
	// SELECT * FROM attendance_records_by_user
	// WHERE company_id = uuid_company AND employee_id = uuid_employee AND year_month = '2023-10';
	sql_raw := `SELECT company_id, year_month, record_time, employee_id,
		device_id, record_type, verification_method, verification_score,
		face_image_url, location_coordinates, metadata, sync_status, created_at
		FROM attendance_records_by_user
		WHERE company_id = ? AND year_month = ? AND employee_id = ?;`
	// Check have page
	var iter *gocql.Iter
	output := &model.AttendanceRecordOutput{}
	if input.PageSize > 0 && input.PageStage != nil {
		query := a.dbSession.Query(sql_raw,
			marshalUuid(input.CompanyID),
			input.YearMonth,
			marshalUuid(input.EmployeeID),
		).WithContext(ctx).PageSize(input.PageSize).PageState(input.PageStage)
		iter = query.Iter()
		output.PageStageNext = iter.PageState()
	} else {
		iter = a.dbSession.Query(sql_raw,
			marshalUuid(input.CompanyID),
			input.YearMonth,
			marshalUuid(input.EmployeeID),
		).WithContext(ctx).Iter()
		output.PageStageNext = nil
	}
	var listAttendanceRecords []model.AttendanceRecordInfo
	for {
		var r model.AttendanceRecordInfo
		// Parsing UUID and other types
		gocqlUUIDCompanyID := gocql.UUID{}
		gocqlUUIDEmployeeID := gocql.UUID{}
		gocqlUUIDDeviceID := gocql.UUID{}
		gocqlFloatVerificationScore := float32(0)
		//
		if !iter.Scan(
			&gocqlUUIDCompanyID,
			&r.YearMonth,
			&r.RecordTime,
			&gocqlUUIDEmployeeID,
			&gocqlUUIDDeviceID,
			&r.RecordType,
			&r.VerificationMethod,
			&gocqlFloatVerificationScore,
			&r.FaceImageURL,
			&r.LocationCoordinates,
			&r.Metadata,
			&r.SyncStatus,
			&r.CreatedAt,
		) {
			break
		}
		// Unmarshal
		r.CompanyID = uuid.UUID(gocqlUUIDCompanyID)
		r.EmployeeID = uuid.UUID(gocqlUUIDEmployeeID)
		r.DeviceID = uuid.UUID(gocqlUUIDDeviceID)
		r.VerificationScore = float64(gocqlFloatVerificationScore)
		//
		listAttendanceRecords = append(listAttendanceRecords, r)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	output.Records = listAttendanceRecords
	output.PageSize = len(listAttendanceRecords)
	return output, nil
}

func NewAttendanceRepository(dbSession *gocql.Session) domainRepo.IAttendanceRepository {
	return &AttendanceRepository{
		dbSession: dbSession,
	}
}

// ===========================================
// Helper parser data
// ===========================================

// UUID marshaling/unmarshaling
func marshalUuid(id uuid.UUID) gocql.UUID {
	return gocql.UUID(id)
}

func unmarshalUuid(id gocql.UUID) uuid.UUID {
	return uuid.UUID(id)
}

func marshalUuidFromBytes(b [16]byte) gocql.UUID {
	return gocql.UUID(b)
}

func marshalUuidFromSlice(b []byte) (gocql.UUID, error) {
	if len(b) != 16 {
		return gocql.UUID{}, errors.New("UUID byte slice must be 16 bytes")
	}
	var arr [16]byte
	copy(arr[:], b)
	return gocql.UUID(arr), nil
}

func unmarshalUuidToBytes(id gocql.UUID) []byte {
	return id[:]
}

// Numeric types marshaling
func marshalFloat64ToFloat32(input float64) float32 {
	return float32(input)
}

func marshalFloat32ToFloat64(input float32) float64 {
	return float64(input)
}

func marshalIntToInt64(input int) int64 {
	return int64(input)
}

func marshalInt32ToInt64(input int32) int64 {
	return int64(input)
}

func marshalInt16ToInt64(input int16) int64 {
	return int64(input)
}

func marshalInt8ToInt64(input int8) int64 {
	return int64(input)
}

func marshalInt64ToInt(input int64) int {
	return int(input)
}

func marshalInt64ToInt32(input int64) int32 {
	return int32(input)
}

// Time marshaling/unmarshaling
func marshalTimeToTimestamp(t time.Time) int64 {
	return t.UnixMilli() // milliseconds since Unix epoch
}

func unmarshalTimestampToTime(ts int64) time.Time {
	return time.UnixMilli(ts)
}

func marshalDurationToNanos(d time.Duration) int64 {
	return int64(d) // nanoseconds
}

func unmarshalNanosToDuration(nanos int64) time.Duration {
	return time.Duration(nanos)
}

func marshalTimeOfDayToNanos(t time.Time) int64 {
	midnight := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return int64(t.Sub(midnight)) // nanoseconds since start of day
}

func unmarshalNanosToTimeOfDay(nanos int64) time.Duration {
	return time.Duration(nanos)
}

// Date marshaling (milliseconds since Unix epoch to start of day in UTC)
func marshalDateToMillis(t time.Time) int64 {
	startOfDay := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	return startOfDay.UnixMilli()
}

func unmarshalMillisToDate(millis int64) time.Time {
	return time.UnixMilli(millis).UTC()
}

// String marshaling for various types
func marshalIPToString(ip net.IP) string {
	return ip.String()
}

func unmarshalStringToIP(s string) net.IP {
	return net.ParseIP(s)
}

func marshalDateToString(t time.Time) string {
	return t.Format("2006-01-02")
}

func unmarshalStringToDate(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}

// Slice/Array marshaling
func marshalSliceToSet[T comparable](slice []T) map[T]struct{} {
	result := make(map[T]struct{}, len(slice))
	for _, v := range slice {
		result[v] = struct{}{}
	}
	return result
}

func marshalSetToSlice[T comparable](set map[T]struct{}) []T {
	result := make([]T, 0, len(set))
	for k := range set {
		result = append(result, k)
	}
	return result
}

// Map marshaling helpers
func marshalMapStringInterface(m map[string]interface{}) map[string]interface{} {
	return m // Direct pass-through for UDT
}

func unmarshalMapStringInterface(m map[string]interface{}) map[string]interface{} {
	return m
}
