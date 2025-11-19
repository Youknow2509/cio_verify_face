package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	domainConfig "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/config"
	domainModel "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/model"
	domainRepo "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/repository"
	conn "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/infrastructure/conn"
	infraRepo "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/infrastructure/repository"
)

func getService(ctx context.Context) (domainRepo.IAttendanceRepository, error) {
	if err := conn.InitScylladbClient(&domainConfig.ScyllaDbSetting{
		Authentication: struct {
			Username string `mapstructure:"username"`
			Password string `mapstructure:"password"`
		}{
			Username: "cassandra",
			Password: "root1234",
		},
		Address:  []string{"127.0.1:9042"},
		Keyspace: "cio_verify_face",
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
		return nil, errors.New("failed to initialize ScyllaDB client: " + err.Error())
	}

	session, err := conn.GetScylladbClient()
	if err != nil {
		return nil, errors.New("failed to get ScyllaDB client: " + err.Error())
	}
	// Test connection
	if err := session.Query("SELECT now() FROM system.local").Exec(); err != nil {
		return nil, errors.New("failed to execute query: " + err.Error())
	}
	// Creaet repo service
	repo := infraRepo.NewAttendanceRepository(session)
	return repo, nil
}

// =======================
// For add data
// =======================
func TestAddAttendanceRecord(t *testing.T) {
	serviceSession, err := getService(context.Background())
	if err != nil {
		t.Fatalf("Failed to get service session: %v", err)
	}
	ctx := context.Background()
	err = serviceSession.AddAttendanceRecord(ctx, &domainModel.AddAttendanceRecordInput{
		CompanyID:           uuid.New(),
		EmployeeID:          uuid.New(),
		DeviceID:            uuid.New(),
		RecordTime:          time.Now(),
		RecordType:          1,
		VerificationMethod:  "FACE",
		VerificationScore:   0.9,
		FaceImageURL:        "http://test.jpg",
		LocationCoordinates: "10,106",
		Metadata: map[string]string{
			"mode": "test",
			"ip":   "1.2.3.4",
			"type": "office",
		},
	})

	if err != nil {
		t.Fatalf("AddAttendanceRecord failed: %v", err)
	}
}

func TestAddDailySummaries(t *testing.T) {
	serviceSession, err := getService(context.Background())
	if err != nil {
		t.Fatalf("Failed to get service session: %v", err)
	}
	ctx := context.Background()
	err = serviceSession.AddDailySummaries(ctx, &domainModel.AddDailySummariesInput{
		CompanyID:         uuid.New(),
		SummaryMonth:      "2025-01",
		WorkDate:          time.Now(),
		EmployeeID:        uuid.New(),
		ShiftID:           uuid.New(),
		ActualCheckIn:     time.Now(),
		ActualCheckOut:    time.Now().Add(8 * time.Hour),
		AttendanceStatus:  0,
		LateMinutes:       0,
		EarlyLeaveMinutes: 0,
		TotalWorkMinutes:  480,
		Notes:             "All good",
		UpdatedAt:         time.Now(),
	})

	if err != nil {
		t.Fatalf("AddDailySummaries failed: %v", err)
	}
}

// =======================
// For get data
// =======================
func TestGetAttendanceRecordCompany(t *testing.T) {
	serviceSession, err := getService(context.Background())
	if err != nil {
		t.Fatalf("Failed to get service session: %v", err)
	}
	ctx := context.Background()
	company_id := uuid.MustParse("2ef95831-bce0-4e0f-ada4-4099a1544c4b")
	records, err := serviceSession.GetAttendanceRecordCompany(ctx, &domainModel.GetAttendanceRecordCompanyInput{
		CompanyID: company_id,
		YearMonth: "2025-11",
	})

	if err != nil {
		t.Fatalf("GetAttendanceRecordCompany failed: %v", err)
	}
	if records == nil {
		t.Fatalf("GetAttendanceRecordCompany returned nil records")
	}
	// Show records
	t.Logf("Total Records: %d", len(records.Records))
	t.Logf("Page Size: %d", records.PageSize)
	t.Logf("Next Page Token: %s", records.PageStageNext)
	for _, record := range records.Records {
		t.Logf("Record: %+v", record)
	}
}

func TestGetAttendanceRecordEmployee(t *testing.T) {
	serviceSession, err := getService(context.Background())
	if err != nil {
		t.Fatalf("Failed to get service session: %v", err)
	}
	ctx := context.Background()
	company_id := uuid.MustParse("2ef95831-bce0-4e0f-ada4-4099a1544c4b")
	employee_id := uuid.MustParse("3a91cfbf-985f-41dd-b86f-c323eb4a6a27")
	records, err := serviceSession.GetAttendanceRecordCompanyForEmployee(ctx, &domainModel.GetAttendanceRecordCompanyForEmployeeInput{
		CompanyID:  company_id,
		YearMonth:  "2025-11",
		EmployeeID: employee_id,
	})

	if err != nil {
		t.Fatalf("GetAttendanceRecordEmployee failed: %v", err)
	}
	if records == nil {
		t.Fatalf("GetAttendanceRecordEmployee returned nil records")
	}
	// Show records
	t.Logf("Total Records: %d", len(records.Records))
	t.Logf("Page Size: %d", records.PageSize)
	t.Logf("Next Page Token: %s", records.PageStageNext)
	for _, record := range records.Records {
		t.Logf("Record: %+v", record)
	}
}

func TestGetGetDailySummarieCompany(t *testing.T) {
	serviceSession, err := getService(context.Background())
	if err != nil {
		t.Fatalf("Failed to get service session: %v", err)
	}
	ctx := context.Background()
	workDate, _ := time.Parse("2006-01-02", "2025-11-19")
	company_id := uuid.MustParse("cd185bee-a7fb-4f84-8466-9425896d696b")
	records, err := serviceSession.GetDailySummarieCompany(ctx, &domainModel.GetDailySummariesCompanyInput{
		CompanyID:    company_id,
		SummaryMonth: "2025-01",
		WorkDate:     workDate,
	})

	if err != nil {
		t.Fatalf("GetGetDailySummarieCompany failed: %v", err)
	}
	if records == nil {
		t.Fatalf("GetGetDailySummarieCompany returned nil records")
	}
	// Show records
	t.Logf("Total Records: %d", len(records.Records))
	t.Logf("Page Size: %d", records.PageSize)
	t.Logf("Next Page Token: %s", records.PageStageNext)
	for _, record := range records.Records {
		t.Logf("Record: %+v", record)
	}
}

func TestGetDailySummarieEmployee(t *testing.T) {
	serviceSession, err := getService(context.Background())
	if err != nil {
		t.Fatalf("Failed to get service session: %v", err)
	}
	ctx := context.Background()
	company_id := uuid.MustParse("cd185bee-a7fb-4f84-8466-9425896d696b")
	employee_id := uuid.MustParse("114c41b3-9bcd-46f9-b626-6d2177d736e7")
	records, err := serviceSession.GetDailySummarieCompanyForEmployee(ctx, &domainModel.GetDailySummariesCompanyForEmployeeInput{
		// PageSize     int       `json:"page_size" omitempty`
		// PageStage    []byte    `json:"page_stage" omitempty`
		CompanyID:    company_id,
		EmployeeID:   employee_id,
		SummaryMonth: "2025-01",
	})

	if err != nil {
		t.Fatalf("GetDailySummarieEmployee failed: %v", err)
	}
	if records == nil {
		t.Fatalf("GetDailySummarieEmployee returned nil records")
	}
	// Show records
	t.Logf("Total Records: %d", len(records.Records))
	t.Logf("Page Size: %d", records.PageSize)
	t.Logf("Next Page Token: %s", records.PageStageNext)
	for _, record := range records.Records {
		t.Logf("Record: %+v", record)
	}
}

// =======================
// For update data
// =======================
func TestUpdateDailySummaries(t *testing.T) {
	serviceSession, err := getService(context.Background())
	if err != nil {
		t.Fatalf("Failed to get service session: %v", err)
	}
	ctx := context.Background()
	company_id := uuid.MustParse("cd185bee-a7fb-4f84-8466-9425896d696b")
	employee_id := uuid.MustParse("114c41b3-9bcd-46f9-b626-6d2177d736e7")
	work_date, _ := time.Parse("2006-01-02", "2025-11-19")
	err = serviceSession.UpdateDailySummariesEmployee(ctx, &domainModel.UpdateDailySummariesEmployeeInput{
		CompanyID:        company_id,
		SummaryMonth:     "2025-01",
		WorkDate:         work_date,
		EmployeeID:       employee_id,
		Notes:            "Updated note test case",
		AttendanceStatus: 1,
	})

	if err != nil {
		t.Fatalf("UpdateDailySummaries failed: %v", err)
	}
}

// =======================
// For delete data
// =======================

func TestDeleteAttendanceRecordBeforeTimestamp(t *testing.T) {
	serviceSession, err := getService(context.Background())
	if err != nil {
		t.Fatalf("Failed to get service session: %v", err)
	}
	ctx := context.Background()
	company_id := uuid.MustParse("2ef95831-bce0-4e0f-ada4-4099a1544c4b")
	employee_id := uuid.MustParse("3a91cfbf-985f-41dd-b86f-c323eb4a6a27")
	record_time, err := time.Parse(time.RFC3339, "2025-11-19T16:13:59.075Z")
	if err != nil {
		t.Fatalf("Failed to parse record time: %v", err)
	}
	record_time = record_time.Add(1 * time.Hour) // Delete all records before 1 hour ago
	err = serviceSession.DeleteAttendanceRecordBeforeTimestamp(ctx, &domainModel.DeleteAttendanceRecordInput{
		CompanyID:  company_id,
		YearMonth:  "2025-01",
		RecordTime: record_time,
		EmployeeID: employee_id,
	})

	if err != nil {
		t.Fatalf("DeleteAttendanceRecordBeforeTimestamp failed: %v", err)
	}
}

func TestDeleteDailySummariesCompanyBeforDate(t *testing.T) {
	serviceSession, err := getService(context.Background())
	if err != nil {
		t.Fatalf("Failed to get service session: %v", err)
	}
	ctx := context.Background()
	company_id := uuid.MustParse("cd185bee-a7fb-4f84-8466-9425896d696b")
	employee_id := uuid.MustParse("114c41b3-9bcd-46f9-b626-6d2177d736e7")
	work_date, _ := time.Parse("2006-01-02", "2025-12-19")
	work_date = work_date.Add(1 * 24 * time.Hour) // Delete all summaries before this date
	err = serviceSession.DeleteDailySummariesCompanyBeforeDate(ctx, &domainModel.DeleteDailySummariesInput{
		CompanyID:    company_id,
		SummaryMonth: "2025-01",
		WorkDate:     work_date,
		EmployeeID:   employee_id,
	})
	if err != nil {
		t.Fatalf("DeleteDailySummariesBeforDate failed: %v", err)
	}
}

func TestDeleteDailySummariesEmployee(t *testing.T) {
	serviceSession, err := getService(context.Background())
	if err != nil {
		t.Fatalf("Failed to get service session: %v", err)
	}
	ctx := context.Background()
	company_id := uuid.MustParse("cd185bee-a7fb-4f84-8466-9425896d696b")
	employee_id := uuid.MustParse("114c41b3-9bcd-46f9-b626-6d2177d736e7")
	work_date, _ := time.Parse("2006-01-02", "2025-11-19")
	err = serviceSession.DeleteDailySummariesEmployee(ctx, &domainModel.DeleteDailySummariesEmployeeInput{
		CompanyID:    company_id,
		SummaryMonth: "2025-01",
		WorkDate:     work_date,
		EmployeeID:   employee_id,
	})
	if err != nil {
		t.Fatalf("DeleteDailySummariesEmployee failed: %v", err)
	}
}
