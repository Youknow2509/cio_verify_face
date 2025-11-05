package tests

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	domainConfig "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/config"
	model "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/model"
	domainRepo "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/repository"
	conn "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/infrastructure/conn"
	repo "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/infrastructure/repository"
)

func TestAddAndGetAttendanceRecordIntegration(t *testing.T) {
	if err := conn.InitScylladbClient(&domainConfig.ScyllaDbSetting{
		Authentication: struct {
			Username string `mapstructure:"username"`
			Password string `mapstructure:"password"`
		}{
			Username: "cassandra",
			Password: "root1234",
		},
		Address:  []string{"127.0.1:9042"},
		Keyspace: "cio_verify_face_test",
		SSL: struct {
			Enabled      bool   `mapstructure:"enabled"`       // true/false to enable/disable SSL
			CertFilePath string `mapstructure:"certfile_path"` // Path to the CA certificate file, default: ./config/scylladb/rootca.crt
			Validate     bool   `mapstructure:"validate"`      // true/false to validate the certificate
			UserKeyPath  string `mapstructure:"userkey_path"`  // Path to the user key file, default: ./config/scylladb/client_key.key
			UserCertPath string `mapstructure:"usercert_path"` // Path to the user certificate
		}{
			Enabled: false,
		},
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: 3600,
	}); err != nil {
		t.Fatalf("failed to initialize scylladb client: %v", err)
	}

	session, err := conn.GetScylladbClient()
	if err != nil {
		t.Fatalf("failed to get scylladb client: %v", err)
	}
	defer session.Close()
	// Test connection
	if err := session.Query("SELECT now() FROM system.local").Exec(); err != nil {
		t.Fatalf("failed to execute query: %v", err)
	}

	// Create infra repo
	repository := repo.NewAttendanceRepository(session)
	ctx := context.Background()
	//
	companyID := uuid.New()
	employeeID := uuid.New()
	deviceID := uuid.New()
	recordTime := time.Now().Unix()
	t.Logf("Company id :: %s", companyID.String())
	t.Logf("Employee id :: %s", employeeID.String())
	t.Logf("Device id :: %s", deviceID.String())
	t.Logf("Record time :: %d", recordTime)
	// Add check-in record
	// if err := addDataCheckIn(
	// 	ctx,
	// 	repository,
	// 	companyID,
	// 	deviceID,
	// 	employeeID,
	// 	recordTime,
	// ); err != nil {
	// 	t.Fatalf("failed to add check-in record: %v", err)
	// }
	// t.Logf("Added check-in record for employee %s at time %d", employeeID.String(), recordTime)

	// Get attendance record
	companyID = uuid.MustParse("2593aad0-a633-4723-8853-58188b3cdbf7")
	deviceID = uuid.MustParse("db93627c-47b7-46c9-a5e8-daa1b5e7edc0")
	recordTime = 1762357421
	records, err := getAttendanceRecord(
		ctx,
		repository,
		companyID,
		deviceID,
		recordTime-3600, // 1 hour before
		recordTime+3600, // 1 hour after
	)
	if err != nil {
		t.Fatalf("failed to get attendance record: %v", err)
	}
	t.Logf("Retrieved %d attendance records", len(records))
	for _, record := range records {
		t.Logf("Record: %+v", record)
	}
}

// Get attendance record
func getAttendanceRecord(
	ctx context.Context,
	repo domainRepo.IAttendanceRepository,
	companyID uuid.UUID,
	deviceID uuid.UUID,
	startTime int64,
	endTime int64,
) ([]*model.AttendanceRecord, error) {
	return repo.GetAttendanceRecordRangeTime(
		ctx,
		&model.GetAttendanceRecordRangeTimeInput{
			CompanyID: companyID,
			DeviceID:  deviceID,
			StartTime: startTime,
			EndTime:   endTime,
		},
	)
}

// Add data check in
func addDataCheckIn(
	ctx context.Context,
	repo domainRepo.IAttendanceRepository,
	companyID uuid.UUID,
	deviceID uuid.UUID,
	employeeID uuid.UUID,
	recordTime int64,
) error {
	return repo.AddCheckInRecord(
		ctx,
		&model.AddCheckInRecordInput{
			CompanyID:           companyID,
			EmployeeID:          employeeID,
			DeviceID:            deviceID,
			VerificationMethod:  "face_recognition",
			VerificationScore:   0.95,
			FaceImageURL:        "http://example.com/face.jpg",
			LocationCoordinates: "37.7749,-122.4194",
			Metadata: map[string]string{
				"ip_address": "192.168.1.1",
			},
			RecordTime: recordTime,
		},
	)
}
