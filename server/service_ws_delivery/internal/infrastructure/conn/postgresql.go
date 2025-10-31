package clients

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/constants"
	domainConfig "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/config"
)

// Global pool variable
var (
	vPostgresqlClient *pgxpool.Pool
)

/**
 * Initialize PostgreSQL client
 * @param postgresSetting *config.PostgresSetting - The PostgreSQL settings
 * @return error - The error if any
 */
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
		return errors.New("failed to parse PostgreSQL connection string")
	}

	// Pool settings
	connConfig.MaxConns = int32(postgresSetting.MaxConns)
	connConfig.MinConns = int32(postgresSetting.MinConns)
	connConfig.MinIdleConns = int32(postgresSetting.MinIdleConns)
	connConfig.MaxConnIdleTime = time.Second * time.Duration(postgresSetting.MaxConnIdleTime)
	connConfig.MaxConnLifetimeJitter = time.Second * time.Duration(postgresSetting.MaxConnLifetimeJitter)
	connConfig.HealthCheckPeriod = time.Second * time.Duration(postgresSetting.HealthCheckPeriod)

	// Optional params
	if postgresSetting.AppName != "" {
		connConfig.ConnConfig.RuntimeParams["application_name"] = postgresSetting.AppName
	}
	if postgresSetting.Timezone != "" {
		connConfig.ConnConfig.RuntimeParams["timezone"] = postgresSetting.Timezone
	}
	if postgresSetting.ConnectionTimeout > 0 {
		connConfig.ConnConfig.ConnectTimeout = time.Second * time.Duration(postgresSetting.ConnectionTimeout)
	}

	// SSL/TLS handling
	switch postgresSetting.SSLMode {
	case constants.POSTGRESQL_SSL_MODE_DISABLE:
		// Nothing to do, sslmode=disable is already in the DSN
	case constants.POSTGRESQL_SSL_MODE_ALLOW, constants.POSTGRESQL_SSL_MODE_PREFER, constants.POSTGRESQL_SSL_MODE_REQUIRE:
		// Client will use SSL if server supports, but does not verify CA or hostname
		// Optionally, user can provide client cert/key for mutual auth even in these modes
		if postgresSetting.SSLCertPath != "" && postgresSetting.SSLCertKeyPath != "" {
			tlsConfig, tlsErr := loadClientCertificate(
				postgresSetting.SSLCertPath,
				postgresSetting.SSLCertKeyPath,
				postgresSetting.SSLRootCertPath,
			)
			if tlsErr != nil {
				return fmt.Errorf("cannot load client certificate: %v", tlsErr)
			}
			connConfig.ConnConfig.TLSConfig = tlsConfig
		}
	case constants.POSTGRESQL_SSL_MODE_VERIFY_CA, constants.POSTGRESQL_SSL_MODE_VERIFY_FULL:
		// These modes require CA at minimum, and full also checks hostname
		tlsConfig, tlsErr := loadClientCertificate(
			postgresSetting.SSLCertPath,
			postgresSetting.SSLCertKeyPath,
			postgresSetting.SSLRootCertPath,
		)
		if tlsErr != nil {
			return fmt.Errorf("cannot load certificates/CA: %v", tlsErr)
		}
		if postgresSetting.SSLMode == constants.POSTGRESQL_SSL_MODE_VERIFY_FULL {
			tlsConfig.InsecureSkipVerify = false // Full verification
		} else {
			// CA only, don't check hostname
			tlsConfig.InsecureSkipVerify = true
		}
		connConfig.ConnConfig.TLSConfig = tlsConfig
	default:
		return fmt.Errorf("unsupported SSL mode: %s", postgresSetting.SSLMode)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, connConfig)
	if err != nil {
		return fmt.Errorf("failed to create pgxpool: %v", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return fmt.Errorf("failed to ping database: %v", err)
	}
	vPostgresqlClient = pool
	return nil
}

/**
 * GetPostgresqlClient returns the PostgreSQL client connection pool.
 * @return (*pgxpool.Pool, error) - The PostgreSQL client pool and error if any.
 */
func GetPostgresqlClient() (*pgxpool.Pool, error) {
	if vPostgresqlClient == nil {
		return nil, errors.New("PostgreSQL client is not initialized, please call InitPostgresqlClient first")
	}
	return vPostgresqlClient, nil
}

// ClosePostgresqlClient closes the PostgreSQL client connection pool.
func ClosePostgresqlClient() {
	if vPostgresqlClient != nil {
		vPostgresqlClient.Close()
		vPostgresqlClient = nil
	}
}

// ===================================================
// 		Helper functions for PostgreSQL client
// ===================================================

// Helper function to load client and CA certs for TLS
func loadClientCertificate(certFile, keyFile, caFile string) (*tls.Config, error) {
	tlsConfig := &tls.Config{}
	if caFile != "" {
		caCert, err := os.ReadFile(caFile)
		if err != nil {
			return nil, fmt.Errorf("error reading CA file: %w", err)
		}
		caPool := x509.NewCertPool()
		caPool.AppendCertsFromPEM(caCert)
		tlsConfig.RootCAs = caPool
	}
	if certFile != "" && keyFile != "" {
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return nil, fmt.Errorf("error loading client cert/key: %w", err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}
	return tlsConfig, nil
}
