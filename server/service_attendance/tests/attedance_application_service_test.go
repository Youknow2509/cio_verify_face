package tests

import (
	"context"
	"testing"
	"time"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	applicationModel "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/application/model"
	applicationServiceImpl "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/application/service/impl"
	domainCache "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/cache"
	domainConfig "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/config"
	domainLogger "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/logger"
	domainRepository "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/domain/repository"
	infraCache "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/infrastructure/cache"
	infraConn "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/infrastructure/conn"
	infraLogger "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/infrastructure/logger"
	infraRepository "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/infrastructure/repository"
)

// ============================================
// Test application service - attendance
// ============================================

// Test add attendance record
func TestAddAttendanceRecordInApplicationService(t *testing.T) {
	if err := initAttendanceApplicationServiceTest(); err != nil {
		t.Fatalf("failed to init service test: %v", err)
	}
	service := applicationServiceImpl.NewAttendanceService()
	ctx := context.Background()
	checkOutTime := time.Date(2025, 11, 24, 17, 35, 0, 0, time.UTC) // Sử dụng time.UTC hoặc location của công ty
	// checkInTime := time.Date(2025, 11, 24, 8, 25, 0, 0, time.UTC) // Sử dụng time.UTC hoặc location của công ty
	req := &applicationModel.AddAttendanceModel{
		Session: &applicationModel.SessionReq{
			SessionId:   uuid.MustParse("11111111-1111-1111-1111-111111111111"),
			UserId:      uuid.MustParse("16584457-022a-4929-818f-96d36e2c4678"),
			CompanyId:   uuid.MustParse("7f204e50-4628-406d-bcd2-40ceb1351256"),
			Role:        1,
			ClientIp:    "127.0.0.1",
			ClientAgent: "Mozilla/5.0",
		},
		CompanyID:           uuid.MustParse("7f204e50-4628-406d-bcd2-40ceb1351256"),
		EmployeeID:          uuid.MustParse("459ce207-d091-4dd4-a994-c87fdde47a83"),
		DeviceID:            uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		RecordTime:          checkOutTime,
		VerificationMethod:  "face",
		VerificationScore:   0.95,
		FaceImageURL:        "http://example.com/face.jpg",
		LocationCoordinates: "37.7749,-122.4194",
	}
	if err := service.AddAttendance(ctx, req); err != nil {
		t.Fatalf("failed to add attendance record: %+v", err)
	}
	t.Log("add attendance record successfully")
}

// init service use for application service test - attendance
func initAttendanceApplicationServiceTest() error {
	// init logger
	dataInitLogger := &infraLogger.ZapLoggerInitializer{
		FolderStore:    "./logs",
		FileMaxSize:    10,
		FileMaxBackups: 1,
		FileMaxAge:     1,
		Compress:       false,
	}
	loggerServiceImpl, er := infraLogger.NewZapLogger(dataInitLogger)
	if er != nil {
		return er
	}
	err := domainLogger.SetLogger(loggerServiceImpl)
	if err != nil {
		return err
	}
	// init db
	posgresClient, err := initConnectionPostgreSQL(
		&domainConfig.PostgresSetting{
			Address:               []string{"localhost:5433"},
			Username:              "postgres",
			Password:              "root1234",
			SSLMode:               "disable",
			Database:              "cio_verify_face",
			ConnectionTimeout:     5,
			Timezone:              "UTC",
			MaxConns:              10,
			MinConns:              2,
			MinIdleConns:          2,
			HealthCheckPeriod:     60,
			MaxConnIdleTime:       300,
			MaxConnLifetimeJitter: 60,
		},
	)
	if err != nil {
		return err
	}
	scylladbClient, err := initConnectionScyllaDB(
		&domainConfig.ScyllaDbSetting{
			Address:  []string{"localhost:9042"},
			Keyspace: "cio_verify_face",
			Authentication: struct {
				Username string `mapstructure:"username"`
				Password string `mapstructure:"password"`
			}{
				Username: "cassandra",
				Password: "root1234",
			},
		},
	)
	if err != nil {
		return err
	}
	if err := domainRepository.SetAttendanceRepository(
		infraRepository.NewAttendanceRepository(scylladbClient),
	); err != nil {
		return err
	}
	if err := domainRepository.SetUserRepository(
		infraRepository.NewUserRepository(posgresClient),
	); err != nil {
		return err
	}
	// init cache
	if err := initRedisDistributedCache(
		&domainConfig.RedisSetting{
			Host:     "127.0.0.1",
			Port:     6379,
			DB:       0,
			Password: "root1234",
			Type:     1,
			UseTLS:   false,
		},
	); err != nil {
		return err
	}
	localCacheImpl, err := infraCache.NewRistrettoLocalCache()
	if err != nil {
		return err
	}
	if err := domainCache.SetLocalCache(localCacheImpl); err != nil {
		return err
	}
	return nil
}

func initConnectionScyllaDB(setting *domainConfig.ScyllaDbSetting) (*gocql.Session, error) {
	if err := infraConn.InitScylladbClient(setting); err != nil {
		return nil, err
	}
	return infraConn.GetScylladbClient()
}

func initRedisDistributedCache(setting *domainConfig.RedisSetting) error {
	distributedCacheImpl, err := infraCache.NewRedisDistributedCache(setting)
	if err != nil {
		return err
	}
	if err := domainCache.SetDistributedCache(distributedCacheImpl); err != nil {
		return err
	}
	return nil
}

func initConnectionPostgreSQL(setting *domainConfig.PostgresSetting) (*pgxpool.Pool, error) {
	if err := infraConn.InitPostgresqlClient(setting); err != nil {
		return nil, err
	}
	return infraConn.GetPostgresqlClient()
}
