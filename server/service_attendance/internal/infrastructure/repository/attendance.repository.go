package repository

import (
	"context"
	"errors"
	"time"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/model"
	domainRepo "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/repository"
)

// CREATE TABLE IF NOT EXISTS attendance_records_by_company_date (
//     company_id UUID,
//     record_time TIMESTAMP,
//     employee_id UUID,
//     device_id UUID,
//     record_type INT, -- 0: CHECK_IN, 1: CHECK_OUT
//     verification_method TEXT,
//     verification_score float,
//     face_image_url TEXT,
//     location_coordinates TEXT, -- "lat,lng" format
//     metadata MAP<TEXT, TEXT>,
//     sync_status TEXT,
//     created_at TIMESTAMP,
//     PRIMARY KEY ((company_id, device_id), record_time)
// ) WITH CLUSTERING ORDER BY (record_time DESC)
// AND default_time_to_live = 2592000; -- 30 days TTL

// ============================================
// Attendance repository
// ============================================
type AttendanceRepository struct {
	dbSession *gocql.Session
}

// GetAttendanceRecordRangeTimeWithUserId implements repository.IAttendanceRepository.
func (a *AttendanceRepository) GetAttendanceRecordRangeTimeWithUserId(ctx context.Context, input *model.GetAttendanceRecordRangeTimeWithUserIdInput) ([]*model.AttendanceRecord, error) {
	if input == nil {
		return nil, errors.New("GetAttendanceRecordRangeTimeWithUserIdInput is nil")
	}
	// Parse data query parameters
	companyId, err := gocql.ParseUUID(input.CompanyID.String())
	if err != nil {
		return nil, err
	}
	userId, err := gocql.ParseUUID(input.UserID.String())
	if err != nil {
		return nil, err
	}
	startTime := time.Unix(input.StartTime, 0)
	endTime := time.Unix(input.EndTime, 0)
	if endTime.Before(startTime) {
		return nil, errors.New("EndTime is before StartTime")
	}
	//
	query := `SELECT
		company_id,
		record_time,
		employee_id,
		device_id,
		record_type,
		verification_method,
		verification_score,
		face_image_url,
		location_coordinates,
		metadata,
		sync_status,
		created_at
	FROM attendance_records_by_company_date
	WHERE company_id = ? AND employee_id = ? AND record_time >= ? AND record_time <= ?
	ALLOW FILTERING
	ORDER BY record_time ASC
	LIMIT ? OFFSET ?`
	// Execute query
	iter := a.dbSession.Query(
		query,
		companyId,
		userId,
		startTime,
		endTime,
		input.Limit,
		input.Offset,
	).Iter()
	var records []*model.AttendanceRecord
	var (
		cId            gocql.UUID
		rTime          time.Time
		eId            gocql.UUID
		dId            gocql.UUID
		rType          int
		vMethod        string
		vScore         float32
		fImageURL      string
		locCoordinates string
		metadata       map[string]string
		syncStatus     string
		createdAt      time.Time
	)
	for iter.Scan(
		&cId,
		&rTime,
		&eId,
		&dId,
		&rType,
		&vMethod,
		&vScore,
		&fImageURL,
		&locCoordinates,
		&metadata,
		&syncStatus,
		&createdAt,
	) {
		record := &model.AttendanceRecord{
			CompanyID:           uuid.UUID(cId.Bytes()),
			RecordTime:          rTime.Unix(),
			EmployeeID:          uuid.UUID(eId.Bytes()),
			DeviceID:            uuid.UUID(dId.Bytes()),
			Type:                rType,
			VerificationMethod:  vMethod,
			VerificationScore:   vScore,
			FaceImageURL:        fImageURL,
			LocationCoordinates: locCoordinates,
			Metadata:            metadata,
			CreatedAt:           createdAt.Unix(),
		}
		records = append(records, record)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return records, nil
}

// AddCheckInRecord implements repository.IAttendanceRepository.
func (a *AttendanceRepository) AddCheckInRecord(ctx context.Context, input *model.AddCheckInRecordInput) error {
	if input == nil {
		return errors.New("AddCheckInRecordInput is nil")
	}
	// Parse data query parameters
	companyId, err := gocql.ParseUUID(input.CompanyID.String())
	if err != nil {
		return err
	}
	employeeId, err := gocql.ParseUUID(input.EmployeeID.String())
	if err != nil {
		return err
	}
	deviceId, err := gocql.ParseUUID(input.DeviceID.String())
	if err != nil {
		return err
	}
	recordTime := time.Unix(input.RecordTime, 0)
	//
	query := `INSERT INTO attendance_records_by_company_date (
		company_id,
		record_time,
		employee_id,
		device_id,
		record_type,
		verification_method,
		verification_score,
		face_image_url,
		location_coordinates,
		metadata,
		sync_status,
		created_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) USING TTL 2592000` // 30 days TTL
	if err := a.dbSession.Query(
		query,
		companyId,
		recordTime,
		employeeId,
		deviceId,
		0, // CHECK_IN
		input.VerificationMethod,
		input.VerificationScore,
		input.FaceImageURL,
		input.LocationCoordinates,
		input.Metadata,
		"pending",
		time.Now(),
	).Exec(); err != nil {
		return err
	}
	return nil
}

// AddCheckOutRecord implements repository.IAttendanceRepository.
func (a *AttendanceRepository) AddCheckOutRecord(ctx context.Context, input *model.AddCheckOutRecordInput) error {
	if input == nil {
		return errors.New("AddCheckOutRecordInput is nil")
	}
	// Parse data query parameters
	companyId, err := gocql.ParseUUID(input.CompanyID.String())
	if err != nil {
		return err
	}
	employeeId, err := gocql.ParseUUID(input.EmployeeID.String())
	if err != nil {
		return err
	}
	deviceId, err := gocql.ParseUUID(input.DeviceID.String())
	if err != nil {
		return err
	}
	recordTime := time.Unix(input.RecordTime, 0)
	//
	query := `INSERT INTO attendance_records_by_company_date (
		company_id,
		record_time,
		employee_id,
		device_id,
		record_type,
		verification_method,
		verification_score,
		face_image_url,
		location_coordinates,
		metadata,
		sync_status,
		created_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) USING TTL 2592000` // 30 days TTL
	if err := a.dbSession.Query(
		query,
		companyId,
		recordTime,
		employeeId,
		deviceId,
		1, // CHECK_OUT
		input.VerificationMethod,
		input.VerificationScore,
		input.FaceImageURL,
		input.LocationCoordinates,
		input.Metadata,
		"pending",
		time.Now(),
	).Exec(); err != nil {
		return err
	}
	return nil
}

// GetAttendanceRecordRangeTime implements repository.IAttendanceRepository.
func (a *AttendanceRepository) GetAttendanceRecordRangeTime(ctx context.Context, input *model.GetAttendanceRecordRangeTimeInput) ([]*model.AttendanceRecord, error) {
	if input == nil {
		return nil, errors.New("GetAttendanceRecordRangeTimeInput is nil")
	}
	// Parse data query parameters
	companyId, err := gocql.ParseUUID(input.CompanyID.String())
	if err != nil {
		return nil, err
	}
	deviceId, err := gocql.ParseUUID(input.DeviceID.String())
	if err != nil {
		return nil, err
	}
	startTime := time.Unix(input.StartTime, 0)
	endTime := time.Unix(input.EndTime, 0)
	if endTime.Before(startTime) {
		return nil, errors.New("EndTime is before StartTime")
	}
	//
	query := `SELECT
		company_id,
		record_time,
		employee_id,
		device_id,
		record_type,
		verification_method,
		verification_score,
		face_image_url,
		location_coordinates,
		metadata,
		sync_status,
		created_at
	FROM attendance_records_by_company_date
	WHERE company_id = ? AND device_id = ? AND record_time >= ? AND record_time <= ?
	ORDER BY record_time ASC
	LIMIT ? OFFSET ?`
	// Execute query
	iter := a.dbSession.Query(
		query,
		companyId,
		deviceId,
		startTime,
		endTime,
		input.Limit,
		input.Offset,
	).Iter()
	var records []*model.AttendanceRecord
	var (
		cId            gocql.UUID
		rTime          time.Time
		eId            gocql.UUID
		dId            gocql.UUID
		rType          int
		vMethod        string
		vScore         float32
		fImageURL      string
		locCoordinates string
		metadata       map[string]string
		syncStatus     string
		createdAt      time.Time
	)
	for iter.Scan(
		&cId,
		&rTime,
		&eId,
		&dId,
		&rType,
		&vMethod,
		&vScore,
		&fImageURL,
		&locCoordinates,
		&metadata,
		&syncStatus,
		&createdAt,
	) {
		record := &model.AttendanceRecord{
			CompanyID:           uuid.UUID(cId.Bytes()),
			RecordTime:          rTime.Unix(),
			EmployeeID:          uuid.UUID(eId.Bytes()),
			DeviceID:            uuid.UUID(dId.Bytes()),
			Type:                rType,
			VerificationMethod:  vMethod,
			VerificationScore:   vScore,
			FaceImageURL:        fImageURL,
			LocationCoordinates: locCoordinates,
			Metadata:            metadata,
			CreatedAt:           createdAt.Unix(),
		}
		records = append(records, record)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return records, nil
}

func NewAttendanceRepository(dbSession *gocql.Session) domainRepo.IAttendanceRepository {
	return &AttendanceRepository{
		dbSession: dbSession,
	}
}
