package constants

// ================================================
//
//	Constants for system configuration
//
// ================================================
const (
	DEFAULT_CONFIG_FILE_PATH = "./config/config.yaml"
	//
	DEFAULT_ELASTIC_CERT_FILE_PATH     = "./config/elasticsearch/http_ca.crt"
	DEFAULT_KAFKA_CERT_FILE_PATH       = "./config/kafka/kafka.crt"
	DEFAULT_MEMCACHED_CONFIG_FILE_PATH = "./config/memcached/memcached.conf"
	//
	DEFAULT_POSTGRESQL_SSL_CERT_FILE_PATH = "./config/postgresql/ssl/server.crt"
	DEFAULT_POSTGRESQL_SSL_KEY_FILE_PATH  = "./config/postgresql/ssl/server.key"
	DEFAULT_POSTGRESQL_SSL_CA_FILE_PATH   = "./config/postgresql/ssl/root.crt"
	DEFAULT_POSTGRESQL_SSL_PASSWORD       = "./config/postgresql/ssl/server.key.password" // Password for the PostgreSQL SSL key file
	//
	DEFAULT_REDIS_SSL_CERT_FILE_PATH = "./config/redis/cert.pem"
	DEFAULT_REDIS_SSL_KEY_FILE_PATH  = "./config/redis/key.pem"
	//
	MAX_CONNECTION_USER_TO_WS = 5
)

const (
	DEFAULT_AVATAR_USER = "https://avatars.githubusercontent.com/u/88392742?v=4"
)

const (
	DeviceTypeWeb     = 1
	DeviceTypeMobile  = 2
	DeviceTypeDesktop = 3
)

const (
	PageDefault = 1
	SizeDefault = 20
)