package conn

import (
	"errors"

	"github.com/gocql/gocql"
	domainConfig "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/domain/config"
)

// scylladb client variables
var (
	vScylladbClient *gocql.Session
)

// InitScylladbClient initializes the ScyllaDB client
func InitScylladbClient(scyllaDbSetting *domainConfig.ScyllaConfig) error {
	if vScylladbClient != nil {
		return errors.New("ScyllaDB client is already initialized")
	}
	
	cluster := gocql.NewCluster(scyllaDbSetting.Address...)
	cluster.Keyspace = scyllaDbSetting.Keyspace
	cluster.Consistency = gocql.Quorum
	
	// Set authentication
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: scyllaDbSetting.Authentication.Username,
		Password: scyllaDbSetting.Authentication.Password,
	}
	
	// Set SSL configuration if enabled
	if scyllaDbSetting.SSL.Enabled {
		cluster.SslOpts = &gocql.SslOptions{
			CaPath:                 scyllaDbSetting.SSL.CertFilePath,
			EnableHostVerification: scyllaDbSetting.SSL.Validate,
			KeyPath:                scyllaDbSetting.SSL.UserKeyPath,
			CertPath:               scyllaDbSetting.SSL.UserCertPath,
		}
	}
	
	// Set pool configuration
	cluster.PoolConfig.HostSelectionPolicy = gocql.RoundRobinHostPolicy()
	cluster.NumConns = scyllaDbSetting.MaxOpenConns
	
	// connect to the cluster
	session, err := cluster.CreateSession()
	if err != nil {
		return errors.New("failed to connect to ScyllaDB: " + err.Error())
	}
	
	vScylladbClient = session
	return nil
}

// GetScylladbClient returns the ScyllaDB client session
func GetScylladbClient() (*gocql.Session, error) {
	if vScylladbClient == nil {
		return nil, errors.New("ScyllaDB client is not initialized, please call InitScylladbClient first")
	}
	return vScylladbClient, nil
}

// CloseScylladbClient closes the ScyllaDB connection
func CloseScylladbClient() {
	if vScylladbClient != nil {
		vScylladbClient.Close()
		vScylladbClient = nil
	}
}
