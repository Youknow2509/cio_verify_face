package config

import "time"

// Config represents the main configuration structure
type Config struct {
	Server        ServerConfig        `mapstructure:"server"`
	Grpc          GrpcConfig          `mapstructure:"grpc"`
	AuthService   AuthServiceConfig   `mapstructure:"auth_service"`
	Postgres      PostgresConfig      `mapstructure:"postgres"`
	Scylla        ScyllaConfig        `mapstructure:"scylladb"`
	Redis         RedisConfig         `mapstructure:"redis"`
	JWT           JWTConfig           `mapstructure:"jwt"`
	Logger        LoggerConfig        `mapstructure:"logger"`
	Export        ExportConfig        `mapstructure:"export"`
	ObjectStorage ObjectStorageConfig `mapstructure:"object_storage"`
	Kafka         KafkaConfig         `mapstructure:"kafka"`
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Name                  string  `mapstructure:"name"`
	ID                    string  `mapstructure:"id"`
	Region                string  `mapstructure:"region"`
	ShardID               string  `mapstructure:"shardId"`
	Port                  int     `mapstructure:"port"`
	Mode                  string  `mapstructure:"mode"`
	DegradedThreshold     float64 `mapstructure:"degraded_threshold"`
	OutOfServiceThreshold float64 `mapstructure:"out_of_service_threshold"`
}

// GrpcConfig represents gRPC configuration
type GrpcConfig struct {
	Network string    `mapstructure:"network"`
	Host    string    `mapstructure:"host"`
	Port    int       `mapstructure:"port"`
	TLS     TLSConfig `mapstructure:"tls"`
}

// AuthServiceConfig represents auth service connection config
type AuthServiceConfig struct {
	GrpcAddr string `mapstructure:"grpc_addr"`
	Enabled  bool   `mapstructure:"enabled"`
}

// TLSConfig represents TLS configuration
type TLSConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	CertFile string `mapstructure:"cert_file"`
	KeyFile  string `mapstructure:"key_file"`
}

// PostgresConfig represents PostgreSQL configuration
type PostgresConfig struct {
	Address               []string      `mapstructure:"address"`
	Database              string        `mapstructure:"database"`
	Username              string        `mapstructure:"username"`
	Password              string        `mapstructure:"password"`
	SSLMode               string        `mapstructure:"sslmode"`
	SSLCert               string        `mapstructure:"sslcert"`
	SSLKey                string        `mapstructure:"sslkey"`
	SSLRootCert           string        `mapstructure:"sslrootcert"`
	SSLPassword           string        `mapstructure:"sslpassword"`
	AppName               string        `mapstructure:"appname"`
	ConnectionTimeout     int           `mapstructure:"connectionTimeout"`
	TZ                    string        `mapstructure:"tz"`
	MaxConns              int32         `mapstructure:"maxConns"`
	MinConns              int32         `mapstructure:"minConns"`
	MinIdleConns          int32         `mapstructure:"minIdleConns"`
	MaxConnIdleTime       time.Duration `mapstructure:"maxConnIdleTime"`
	MaxConnLifetimeJitter time.Duration `mapstructure:"maxConnLifetimeJitter"`
	HealthCheckPeriod     time.Duration `mapstructure:"healthCheckPeriod"`
}

// ScyllaConfig represents ScyllaDB configuration
type ScyllaConfig struct {
	Address         []string         `mapstructure:"address"`
	Keyspace        string           `mapstructure:"keyspace"`
	Authentication  ScyllaAuthConfig `mapstructure:"authentication"`
	SSL             ScyllaSSLConfig  `mapstructure:"ssl"`
	MaxIdleConns    int              `mapstructure:"maxIdleConns"`
	MaxOpenConns    int              `mapstructure:"maxOpenConns"`
	ConnMaxLifetime int              `mapstructure:"connMaxLifetime"`
}

// ScyllaAuthConfig represents ScyllaDB authentication configuration
type ScyllaAuthConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// ScyllaSSLConfig represents ScyllaDB SSL configuration
type ScyllaSSLConfig struct {
	Enabled      bool   `mapstructure:"enabled"`
	CertFilePath string `mapstructure:"certfile_path"`
	Validate     bool   `mapstructure:"validate"`
	UserKeyPath  string `mapstructure:"userkey_path"`
	UserCertPath string `mapstructure:"usercert_path"`
}

// RedisConfig represents Redis configuration
type RedisConfig struct {
	Type           int      `mapstructure:"type"`
	UseTLS         bool     `mapstructure:"use_tls"`
	CertPath       string   `mapstructure:"cert_path"`
	KeyPath        string   `mapstructure:"key_path"`
	Password       string   `mapstructure:"password"`
	DB             int      `mapstructure:"db"`
	Host           string   `mapstructure:"host"`
	Port           int      `mapstructure:"port"`
	MasterName     string   `mapstructure:"master_name"`
	SentinelAddrs  []string `mapstructure:"sentinel_addrs"`
	Address        []string `mapstructure:"address"`
	RouteByLatency bool     `mapstructure:"route_by_latency"`
	RouteRandomly  bool     `mapstructure:"route_randomly"`
	PoolSize       int      `mapstructure:"pool_size"`
	MinIdleConns   int      `mapstructure:"min_idle_conns"`
	MaxRetries     int      `mapstructure:"max_retries"`
}

// JWTConfig represents JWT configuration
type JWTConfig struct {
	Secret   string   `mapstructure:"secret"`
	Issuer   string   `mapstructure:"issuer"`
	Subject  string   `mapstructure:"subject"`
	Audience []string `mapstructure:"audience"`
}

// LoggerConfig represents logger configuration
type LoggerConfig struct {
	FolderStore    string `mapstructure:"folder_store"`
	FileMaxSize    int    `mapstructure:"file_max_size"`
	FileMaxBackups int    `mapstructure:"file_max_backups"`
	FileMaxAge     int    `mapstructure:"file_max_age"`
	Compress       bool   `mapstructure:"compress"`
}

// ExportConfig represents export configuration
type ExportConfig struct {
	TempFolder        string `mapstructure:"temp_folder"`
	MaxConcurrentJobs int    `mapstructure:"max_concurrent_jobs"`
	JobTimeoutMinutes int    `mapstructure:"job_timeout_minutes"`
}

// ObjectStorageConfig represents S3/MinIO configuration
type ObjectStorageConfig struct {
	Endpoint             string `mapstructure:"endpoint"`
	AccessKey            string `mapstructure:"access_key"`
	SecretKey            string `mapstructure:"secret_key"`
	Bucket               string `mapstructure:"bucket"`
	UseSSL               bool   `mapstructure:"use_ssl"`
	Region               string `mapstructure:"region"`
	PresignExpireMinutes int    `mapstructure:"presign_expire_minutes"`
}

// KafkaConfig represents Kafka producer configuration
type KafkaConfig struct {
	Brokers      []string `mapstructure:"brokers"`
	NotifyTopic  string   `mapstructure:"notify_topic"`
	ClientID     string   `mapstructure:"client_id"`
	SASLEnabled  bool     `mapstructure:"sasl_enabled"`
	SASLUser     string   `mapstructure:"sasl_user"`
	SASLPassword string   `mapstructure:"sasl_password"`
	TLSEnabled   bool     `mapstructure:"tls_enabled"`
}
