package conn

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	domainConfig "github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/domain/config"
)

var vPostgresqlClient *pgxpool.Pool

func InitPostgresqlClient(postgresSetting *domainConfig.PostgresSetting) error {
	if postgresSetting == nil {
		return errors.New("PostgreSQL settings are nil")
	}

	address := strings.Join(postgresSetting.Address, ",")
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=%s",
		postgresSetting.Username,
		postgresSetting.Password,
		address,
		postgresSetting.Database,
		postgresSetting.SSLMode,
	)

	connConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return fmt.Errorf("failed to parse PostgreSQL connection string: %w", err)
	}

	connConfig.MaxConns = int32(postgresSetting.MaxConns)
	connConfig.MinConns = int32(postgresSetting.MinConns)
	connConfig.MaxConnIdleTime = time.Second * time.Duration(postgresSetting.MaxConnIdleTime)
	connConfig.HealthCheckPeriod = time.Second * time.Duration(postgresSetting.HealthCheckPeriod)

	if postgresSetting.AppName != "" {
		connConfig.ConnConfig.RuntimeParams["application_name"] = postgresSetting.AppName
	}
	if postgresSetting.Timezone != "" {
		connConfig.ConnConfig.RuntimeParams["timezone"] = postgresSetting.Timezone
	}
	if postgresSetting.ConnectionTimeout > 0 {
		connConfig.ConnConfig.ConnectTimeout = time.Second * time.Duration(postgresSetting.ConnectionTimeout)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, connConfig)
	if err != nil {
		return fmt.Errorf("failed to create pgxpool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return fmt.Errorf("failed to ping database: %w", err)
	}

	vPostgresqlClient = pool
	return nil
}

func GetPostgresqlClient() (*pgxpool.Pool, error) {
	if vPostgresqlClient == nil {
		return nil, errors.New("PostgreSQL client is not initialized")
	}
	return vPostgresqlClient, nil
}

func ClosePostgresqlClient() {
	if vPostgresqlClient != nil {
		vPostgresqlClient.Close()
		vPostgresqlClient = nil
	}
}
