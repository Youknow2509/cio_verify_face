package config

// ==========================================================
//
//	Configuration Sections
//
// ==========================================================
type (
	Setting struct {
		GrpcServer        GrpcSetting          `mapstructure:"grpc"`
		Server            ServerSetting        `mapstructure:"server"`
		WsServer          WsSetting            `mapstructure:"ws"`
		Cassandra         CassandraSetting     `mapstructure:"cassandra"`
		Elasticsearch     ElasticsearchSetting `mapstructure:"elasticsearch"`
		Jaeger            JaegerSetting        `mapstructure:"jaeger"`
		Kafka             KafkaSetting         `mapstructure:"kafka"`
		Memcached         MemcachedSetting     `mapstructure:"memcached"`
		Minio             MinioSetting         `mapstructure:"minio"`
		Postgres          PostgresSetting      `mapstructure:"postgres"`
		Redis             RedisSetting         `mapstructure:"redis"`
		ScyllaDb          ScyllaDbSetting      `mapstructure:"scylladb"`
		Logstash          LogstashSetting      `mapstructure:"logstash"`
		SMTP              SMTPSetting          `mapstructure:"smtp"`
		JWT               JWTSetting           `mapstructure:"jwt"`
		Logger            LoggerSetting        `mapstructure:"logger"`
		RateLimitPolicies []RateLimitPolicy    `mapstructure:"policy_rate_limit"`
	}
)

// ==========================================================
//
//	List of configuration sections
//
// ==========================================================
// RateLimitPolicySetting
type RateLimitPolicySetting struct {
	Policies []RateLimitPolicy 
}
type RateLimitPolicy struct {
	Name   string `mapstructure:"name"`
	Limit  int    `mapstructure:"limit"`
	Window string `mapstructure:"window"`
}

// server
type ServerSetting struct {
	Name                  string  `mapstructure:"name"`
	Id                    string  `mapstructure:"id"`
	Region                string  `mapstructure:"region"`
	ShardId               string  `mapstructure:"shardId"`
	Port                  int     `mapstructure:"port"`
	Mode                  string  `mapstructure:"mode"`
	DegradedThreshold     float64 `mapstructure:"degraded_threshold"`
	OutOfServiceThreshold float64 `mapstructure:"out_of_service_threshold"`
}

// grpc server
type GrpcSetting struct {
	Network string `mapstructure:"network"`
	Port    int    `mapstructure:"port"`
	Tls     struct {
		Enabled  bool   `mapstructure:"enabled"`
		CertFile string `mapstructure:"cert_file"`
		KeyFile  string `mapstructure:"key_file"`
	}
}

// WsSetting
type WsSetting struct {
	MaxConnectionSystem int   `mapstructure:"max_conn_system"`
	NumWorkers          int   `mapstructure:"num_workers"`
	SizeBufferChan      int   `mapstructure:"size_buffer_chan"`
	MaxMessageSize      int64 `mapstructure:"max_message_size"`
	MaxSendQueueSize    int   `mapstructure:"max_send_queue_size"`
	ReadWait            int   `mapstructure:"read_wait"`
	WriteWait           int   `mapstructure:"write_wait"`
	PingPeriod          int   `mapstructure:"ping_period"`
	ReadBufferSize      int   `mapstructure:"read_buffer_size"`
	WriteBufferSize     int   `mapstructure:"write_buffer_size"`
	HandshakeTimeout    int   `mapstructure:"handshake_timeout"`
	EnableCompression   bool  `mapstructure:"enable_compression"`
	MaxConnsPerUser     int   `mapstructure:"max_conn_per_user"`
}

// cassandra
type CassandraSetting struct {
	Username        string   `mapstructure:"username"`
	Password        string   `mapstructure:"password"`
	Keyspace        string   `mapstructure:"keyspace"`
	MaxIdleConns    int      `mapstructure:"maxIdleConns"`
	MaxOpenConns    int      `mapstructure:"maxOpenConns"`
	ConnMaxLifetime int      `mapstructure:"connMaxLifetime"`
	NodeAddrs       []string `mapstructure:"node_addrs"`
}

// elasticsearch
type ElasticsearchSetting struct {
	AuthType     int      `mapstructure:"auth_type"`     // 1: username/password, 2: service_token/http bear, 3: cloud_id
	Address      []string `mapstructure:"address"`       // e.g. ["http://127.0.1:9200"]
	HaveCert     bool     `mapstructure:"have_cert"`     // true if you have a
	Username     string   `mapstructure:"username"`      // for auth_type 1
	Password     string   `mapstructure:"password"`      // for auth_type 1
	CertPath     string   `mapstructure:"cert_path"`     // path to your certificate file, only if have_cert is true
	ServiceToken string   `mapstructure:"service_token"` // for auth_type 2
	CloudID      string   `mapstructure:"cloud_id"`      // for auth_type 3
	APIKey       string   `mapstructure:"api_key"`       // for auth_type 3F
}

// jaeger
type JaegerSetting struct {
	ServiceName                         string  `mapstructure:"service_name"`
	AgentHost                           string  `mapstructure:"agent_host"`
	AgentPort                           int     `mapstructure:"agent_port"`
	Endpoint                            string  `mapstructure:"endpoint"`
	Username                            string  `mapstructure:"username"`
	Password                            string  `mapstructure:"password"`
	ReporterLogSpans                    bool    `mapstructure:"reporter_log_spans"`
	ReporterMaxQueueSize                int     `mapstructure:"reporter_max_queue_size"`
	ReporterFlushInterval               int     `mapstructure:"reporter_flush_interval"`
	ReporterAttemptReconnectingDisabled bool    `mapstructure:"reporter_attempt_reconnecting_disabled"`
	ReporterAttemptReconnectInterval    int     `mapstructure:"reporter_attempt_reconnect_interval"`
	SamplerType                         string  `mapstructure:"sampler_type"`
	SamplerParam                        float64 `mapstructure:"sampler_param"`
	SamplerManagerHostPort              string  `mapstructure:"sampler_manager_host_port"`
	SamplingEndpoint                    string  `mapstructure:"sampling_endpoint"`
	SamplerMaxOperations                int     `mapstructure:"sampler_max_operations"`
	SamplerRefreshInterval              string  `mapstructure:"sampler_refresh_interval"`
	Tags                                string  `mapstructure:"tags"`
	TraceID128Bit                       bool    `mapstructure:"traceid_128bit"`
	Disabled                            bool    `mapstructure:"disabled"`
	RPCMetrics                          bool    `mapstructure:"rpc_metrics"`
}

// kafka
type KafkaSetting struct {
	Brokers  []string         `mapstructure:"brokers"`   // Danh sách broker Kafka (ví dụ: ["127.0.0.1:9092"])
	SASL     KafkaSASLSetting `mapstructure:"sasl"`      // Cấu hình xác thực SASL
	TLS      KafkaTLSSetting  `mapstructure:"tls"`       // Cấu hình bảo mật TLS
	Producer KafkaProducer    `mapstructure:"producer"`  // Cấu hình producer
	Consumer KafkaConsumer    `mapstructure:"consumers"` // Danh sách consumer (hỗ trợ nhiều consumer group/topic)
}

// Cấu hình xác thực SASL cho Kafka
type KafkaSASLSetting struct {
	Enabled   bool   `mapstructure:"enabled"`   // true: bật xác thực SASL, false: tắt
	Mechanism int    `mapstructure:"mechanism"` // 0: plain, 1: scram-sha-256, 2: scram-sha-512
	Username  string `mapstructure:"username"`
	Password  string `mapstructure:"password"`
}

// Cấu hình bảo mật TLS cho Kafka
type KafkaTLSSetting struct {
	Enabled    bool   `mapstructure:"enabled"`     // true: bật TLS, false: tắt
	SkipVerify bool   `mapstructure:"skip_verify"` // true: bỏ qua xác thực cert (không khuyến nghị cho production)
	CAFile     string `mapstructure:"ca_file"`     // Đường dẫn file CA nếu cần
}

// Cấu hình producer Kafka
type KafkaProducer struct {
	CompressionType int  `mapstructure:"compression_type"` // 0: none, 1: gzip, 2: snappy, 3: lz4, 4: zstd
	Retries         int  `mapstructure:"retries"`          // Số lần retry khi lỗi
	RetryBackoffMs  int  `mapstructure:"retry_backoff_ms"` // Thời gian giữa các lần retry (ms)
	LingerMs        int  `mapstructure:"linger_ms"`        // Delay gửi batch nếu batch chưa đầy (ms)
	BatchSize       int  `mapstructure:"batch_size"`       // Số lượng message tối đa trong
	BatchBytes      int  `mapstructure:"batch_bytes"`      // Giới hạn kích thước batch (1MB)
	MaxAttempts     int  `mapstructure:"max_attempts"`     // Tối đa số lần thử gửi (tổng thể, bao gồm retry logic)
	Async           bool `mapstructure:"async"`            // true: gửi không chờ phản hồi (mất message nếu lỗi)
	WriteTimeoutMs  int  `mapstructure:"write_timeout_ms"` // Timeout khi ghi message (ms)
	ReadTimeoutMs   int  `mapstructure:"read_timeout_ms"`  // Timeout khi nhận phản hồi từ Kafka (ms)
	Balancer        int  `mapstructure:"balancer"`         // 0: custom config, 1: RoundRobin, 2: Least
}

// Cấu hình consumer Kafka
type KafkaConsumer struct {
	GroupID             string `mapstructure:"group_id"`              // ID của consumer group
	CommitIntervalMs    int    `mapstructure:"commit_interval_ms"`    // 0 = sync commit, >0 = auto commit theo thời gian (ms)
	MinBytes            int    `mapstructure:"min_bytes"`             // 10KB: Minimum data per fetch
	MaxBytes            int    `mapstructure:"max_bytes"`             // 1MB: Maximum fetch size
	MaxWaitMs           int    `mapstructure:"max_wait_ms"`           // Maximum time to wait for batch fill (ms)
	ReadBatchTimeoutMs  int    `mapstructure:"read_batch_timeout_ms"` // Timeout khi đọc batch (ms)
	HeartbeatIntervalMs int    `mapstructure:"heartbeat_interval_ms"` // Interval heartbeat (ms)
	SessionTimeoutMs    int    `mapstructure:"session_timeout_ms"`    // Thời gian timeout phiên (ms)
	RebalanceTimeoutMs  int    `mapstructure:"rebalance_timeout_ms"`  // Thời gian timeout khi tái cân bằng (ms)
	JoinGroupBackoffMs  int    `mapstructure:"join_group_backoff_ms"` // Thời gian chờ khi tham gia nhóm (ms)
	ReadBackoffMinMs    int    `mapstructure:"read_backoff_min_ms"`   // Min delay between poll retries (ms)
	ReadBackoffMaxMs    int    `mapstructure:"read_backoff_max_ms"`   // Max delay between poll retries (ms)
	ReadLagIntervalMs   int    `mapstructure:"read_lag_interval_ms"`  // -1 = disable, interval để kiểm tra lag (ms)
	MaxAttempts         int    `mapstructure:"max_attempts"`          // Số lần tối đa để thử gửi (bao gồm retry logic)
	QueueCapacity       int    `mapstructure:"queue_capacity"`        // Số lượng tối đa của message trong hàng đợi
	RetentionTimeMs     int    `mapstructure:"retention_time_ms"`     // -1 = sử dụng giá trị mặc định của broker, thời gian giữ message
}

// memcached
type MemcachedSetting struct {
	Address      []string          `mapstructure:"address"`
	Username     string            `mapstructure:"username"`
	Password     string            `mapstructure:"password"`
	Security     MemcachedSecurity `mapstructure:"security"`
	MaxIdleConns int               `mapstructure:"maxIdleConns"`
	Timeout      int               `mapstructure:"timeout"` // seconds
}

type MemcachedSecurity struct {
	EnabledSasl bool   `mapstructure:"enabled_sasl"` // true/false
	FilePath    string `mapstructure:"file_path"`    // path to your sasl file, only if enabled_sasl is true
}

// minio
type MinioSetting struct {
	Type         int    `mapstructure:"type"`          // 1: 'minio', 2: 's3', ...
	Endpoint     string `mapstructure:"endpoint"`      // Ví dụ: 'localhost:9000'
	AccessKey    string `mapstructure:"access_key"`    // Ví dụ: 'minio'
	SecretKey    string `mapstructure:"secret_key"`    // Ví dụ: 'miniosecret'
	Token        string `mapstructure:"token"`         // Nếu cần, để trống nếu không sử dụng token
	EnableSSL    bool   `mapstructure:"enable_ssl"`    // true/false
	Region       string `mapstructure:"region"`        // Ví dụ: 'us-east-1', có thể để trống nếu không cần
	BucketLookup string `mapstructure:"bucket_lookup"` // Giá trị: 'dns', 'path', hoặc 'auto'
}

// postgres
type PostgresSetting struct {
	Address               []string `mapstructure:"address"`
	Database              string   `mapstructure:"database"`
	Username              string   `mapstructure:"username"`
	Password              string   `mapstructure:"password"`
	SSLMode               string   `mapstructure:"sslmode"` // disable, require, verify-ca, verify-full
	SSLCertPath           string   `mapstructure:"sslcert"`
	SSLCertKeyPath        string   `mapstructure:"sslkey"`
	SSLRootCertPath       string   `mapstructure:"sslrootcert"`
	SSLPassword           string   `mapstructure:"sslpassword"`
	AppName               string   `mapstructure:"appname"`               // application name
	ConnectionTimeout     int      `mapstructure:"connection_timeout"`    // seconds
	Timezone              string   `mapstructure:"tz"`                    // timezone, e.g. 'UTC'
	MaxConns              int      `mapstructure:"maxConns"`              // max idle connections
	MinConns              int      `mapstructure:"minConns"`              // min idle connections
	MinIdleConns          int      `mapstructure:"minIdleConns"`          // minimum idle connections
	HealthCheckPeriod     int      `mapstructure:"healthCheckPeriod"`     // seconds
	MaxConnIdleTime       int      `mapstructure:"maxConnIdleTime"`       // seconds
	MaxConnLifetimeJitter int      `mapstructure:"maxConnLifetimeJitter"` // seconds
}

// redis
type RedisSetting struct {
	Type           int      `mapstructure:"type"`             // 1: standalone, 2: sentinel, 3: cluster
	UseTLS         bool     `mapstructure:"use_tls"`          // true if you want to use TLS, false if not
	CertPath       string   `mapstructure:"cert_path"`        // path to your certificate file, default if use: ./config/redis/cert.pem
	KeyPath        string   `mapstructure:"key_path"`         // path to your key file, default if use: ./config/redis/key.pem
	Password       string   `mapstructure:"password"`         // password for redis
	DB             int      `mapstructure:"db"`               // default database
	Host           string   `mapstructure:"host"`             // for standalone type
	Port           int      `mapstructure:"port"`             // for standalone type
	SentinelAddrs  []string `mapstructure:"sentinel_addrs"`   // for sentinel type
	MasterName     string   `mapstructure:"master_name"`      // for sentinel type
	Address        []string `mapstructure:"address"`          // for cluster type
	RouteByLatency bool     `mapstructure:"route_by_latency"` // true if you want to route by latency, false if you want to route by hash
	RouteRandomly  bool     `mapstructure:"route_randomly"`   // true if you want to route randomly, false if you want
	PoolSize       int      `mapstructure:"pool_size"`        // maximum number of connections in the pool
	MinIdleConns   int      `mapstructure:"min_idle_conns"`   // minimum number of idle connections
	MaxRetries     int      `mapstructure:"max_retries"`      // maximum number of retries for
}

// scylladb
type ScyllaDbSetting struct {
	Authentication struct {
		Username string `mapstructure:"username"` // Username for ScyllaDB authentication
		Password string `mapstructure:"password"` // Password for ScyllaDB authentication
	} `mapstructure:"authentication"`
	Address  []string `mapstructure:"address"`  // List of ScyllaDB node addresses (e.g., ["127.0.1:9042"])
	Keyspace string   `mapstructure:"keyspace"` // Keyspace name for ScyllaDB
	SSL      struct {
		Enabled      bool   `mapstructure:"enabled"`       // true/false to enable/disable SSL
		CertFilePath string `mapstructure:"certfile_path"` // Path to the CA certificate file, default: ./config/scylladb/rootca.crt
		Validate     bool   `mapstructure:"validate"`      // true/false to validate the certificate
		UserKeyPath  string `mapstructure:"userkey_path"`  // Path to the user key file, default: ./config/scylladb/client_key.key
		UserCertPath string `mapstructure:"usercert_path"` // Path to the user certificate
	} `mapstructure:"ssl"` // SSL configuration for ScyllaDB
	MaxIdleConns    int `mapstructure:"maxIdleConns"`    // Maximum number of idle connections
	MaxOpenConns    int `mapstructure:"maxOpenConns"`    // Maximum number
	ConnMaxLifetime int `mapstructure:"connMaxLifetime"` // Maximum connection lifetime in seconds
}

// logstash
type LogstashSetting struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Timeout  int    `mapstructure:"timeout"`
	Protocol string `mapstructure:"protocol"`
}

// smtp
type SMTPSetting struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// jwt
type JWTSetting struct {
	Secret   string   `mapstructure:"secret"`
	Issuer   string   `mapstructure:"issuer"`
	Subject  string   `mapstructure:"subject"`
	Audience []string `mapstructure:"audience"`
}

// logger
type LoggerSetting struct {
	FolderStore    string `mapstructure:"folder_store"`     // Folder to store log files
	FileMaxSize    int    `mapstructure:"file_max_size"`    // Maximum size of each log file in MB
	FileMaxBackups int    `mapstructure:"file_max_backups"` // Maximum number of old log files to keep
	FileMaxAge     int    `mapstructure:"file_max_age"`     // Maximum age of log files in days
	Compress       bool   `mapstructure:"compress"`         // Whether to compress old log files
}
