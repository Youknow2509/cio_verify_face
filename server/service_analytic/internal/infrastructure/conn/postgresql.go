package conn

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	domainConfig "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/domain/config"
)

// postgres client variables
var (
	vPostgresPool *pgxpool.Pool
)

// InitPostgresClient initializes the PostgreSQL connection pool
func InitPostgresClient(pgConfig *domainConfig.PostgresConfig) error {
	if vPostgresPool != nil {
		return errors.New("PostgreSQL client is already initialized")
	}

	// Build connection string
	connStr := fmt.Sprintf(
		"postgresql://%s:%s@%s/%s?sslmode=%s&application_name=%s&connect_timeout=%d&timezone=%s",
		pgConfig.Username,
		pgConfig.Password,
		pgConfig.Address[0], // Use first address
		pgConfig.Database,
		pgConfig.SSLMode,
		pgConfig.AppName,
		pgConfig.ConnectionTimeout,
		pgConfig.TZ,
	)

	// Parse pool config
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return fmt.Errorf("failed to parse postgres config: %w", err)
	}

	// Set pool configuration
	config.MaxConns = pgConfig.MaxConns
	config.MinConns = pgConfig.MinConns
	config.MaxConnIdleTime = time.Duration(pgConfig.MaxConnIdleTime) * time.Second
	config.MaxConnLifetime = time.Duration(pgConfig.MaxConnLifetimeJitter) * time.Second
	config.HealthCheckPeriod = time.Duration(pgConfig.HealthCheckPeriod) * time.Second

	// Create connection pool
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return fmt.Errorf("failed to create postgres pool: %w", err)
	}

	// Ping to verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return fmt.Errorf("failed to ping postgres: %w", err)
	}

	vPostgresPool = pool
	return nil
}

// GetPostgresPool returns the PostgreSQL connection pool
func GetPostgresPool() (*pgxpool.Pool, error) {
	if vPostgresPool == nil {
		return nil, errors.New("PostgreSQL client is not initialized, please call InitPostgresClient first")
	}
	return vPostgresPool, nil
}

// ClosePostgresClient closes the PostgreSQL connection pool
func ClosePostgresClient() {
	if vPostgresPool != nil {
		vPostgresPool.Close()
		vPostgresPool = nil
	}
}
