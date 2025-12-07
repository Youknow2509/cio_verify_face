package constants

// =================================
// Configuration File Paths:
// =================================
const (
	DEFAULT_CONFIG_FILE_PATH = "./config/config.yaml"
	FilePathConfigDev        = "./config/config.dev.yaml"
	FilePathConfigProd       = "./config/config.yaml"
)

// =================================
// Default Certificate File Paths:
// =================================
const (
	DEFAULT_ELASTIC_CERT_FILE_PATH     = "./config/elasticsearch/http_ca.crt"
	DEFAULT_KAFKA_CERT_FILE_PATH       = "./config/kafka/kafka.crt"
	DEFAULT_MEMCACHED_CONFIG_FILE_PATH = "./config/memcached/memcached.conf"
	//
	DEFAULT_POSTGRESQL_SSL_CERT_FILE_PATH = "./config/postgresql/ssl/server.crt"
	DEFAULT_POSTGRESQL_SSL_KEY_FILE_PATH  = "./config/postgresql/ssl/server.key"
	DEFAULT_POSTGRESQL_SSL_CA_FILE_PATH   = "./config/postgresql/ssl/root.crt"
	DEFAULT_POSTGRESQL_SSL_PASSWORD       = "./config/postgresql/ssl/server.key.password"
	//
	DEFAULT_REDIS_SSL_CERT_FILE_PATH = "./config/redis/cert.pem"
	DEFAULT_REDIS_SSL_KEY_FILE_PATH  = "./config/redis/key.pem"
)

// =================================
// Redis Type Constants:
// =================================
const (
	RedisTypeStandalone = 1
	RedisTypeSentinel   = 2
	RedisTypeCluster    = 3
)

// =================================
// Rate Limiting Constants:
// =================================
const (
	// Maximum face profile update requests per user per month
	MaxFaceProfileRequestsPerMonth = 3

	// Update link TTL (24 hours)
	UpdateLinkTTLSeconds = 60 * 60 * 24

	// Password reset link TTL (24 hours)
	PasswordResetLinkTTLSeconds = 60 * 60 * 24

	// Minimum interval between password resets for the same user (1 hour)
	PasswordResetCooldownSeconds = 60 * 60

	// Maximum password resets per manager per hour
	MaxPasswordResetsPerManagerPerHour = 10

	// Spam prevention window (seconds)
	SpamPreventionWindowSeconds = 60
)

// =================================
// Cache TTL Constants:
// =================================
const (
	// Local cache TTLs
	TTLLocalPendingCheck         = 30 // 30 seconds
	TTLLocalMonthlyCount         = 60 // 1 minute
	TTLLocalUpdateToken          = 60 // 1 minute
	TTLLocalSpamCheck            = 5  // 5 seconds
	TTLLocalPasswordResetToken   = 60 // 1 minute

	// Distributed cache TTLs
	TTLDistributedPendingCheck      = 120          // 2 minutes
	TTLDistributedMonthlyCount      = 300          // 5 minutes
	TTLDistributedUpdateToken       = 60 * 60 * 24 // 24 hours
	TTLDistributedSpamCheck         = 60           // 1 minute
	TTLDistributedPasswordResetToken = 60 * 60 * 24 // 24 hours
)

// =================================
// Kafka Topics:
// =================================
const (
	// Topic for face profile update requests (for high concurrency processing)
	KafkaTopicFaceProfileUpdateRequests = "face_profile_update_requests"

	// Topic for password reset notifications
	KafkaTopicPasswordResetNotifications = "password_reset_notifications"

	// Topic for email notifications
	KafkaTopicNotifications = "notification_requests"
)

// =================================
// Kafka Event Types:
// =================================
const (
	KafkaEventTypeFaceProfileUpdateRequest = 10
	KafkaEventTypePasswordReset            = 11
)

// =================================
// Cache Key Prefixes:
// =================================
const (
	CacheKeyPrefixPendingRequest      = "fpr:pending:"
	CacheKeyPrefixMonthlyCount        = "fpr:monthly:"
	CacheKeyPrefixUpdateToken         = "fpr:token:"
	CacheKeyPrefixPasswordResetSpam   = "prr:spam:"
	CacheKeyPrefixPasswordResetToken  = "prr:token:"
	CacheKeyPrefixRequestLock         = "fpr:lock:"
	CacheKeyPrefixApprovalLock        = "fpr:approval_lock:"
)

// =================================
// Audit Action Constants:
// =================================
const (
	AuditActionFaceProfileUpdateRequest = "face_profile_update_request"
	AuditActionFaceProfileUpdateApprove = "face_profile_update_approve"
	AuditActionFaceProfileUpdateReject  = "face_profile_update_reject"
	AuditActionFaceProfileUpdate        = "face_profile_update"
	AuditActionPasswordReset            = "password_reset"
)

// =================================
// Resource Types:
// =================================
const (
	AuditResourceTypeFaceProfileRequest = "face_profile_request"
	AuditResourceTypeUser               = "user"
)

// =================================
// Kafka Constants:
// =================================
const (
	KAFKA_SASL_MECHANISM_PLAIN        = 0
	KAFKA_SASL_MECHANISM_SCRAM_SHA256 = 1
	KAFKA_SASL_MECHANISM_SCRAM_SHA512 = 2

	KAFKA_ACKS_NONE   = 0
	KAFKA_ACKS_LEADER = 1
	KAFKA_ACKS_ALL    = 2

	KAFKA_COMPRESSION_NONE   = 0
	KAFKA_COMPRESSION_GZIP   = 1
	KAFKA_COMPRESSION_SNAPPY = 2
	KAFKA_COMPRESSION_LZ4    = 3
	KAFKA_COMPRESSION_ZSTD   = 4

	KAFKA_BALANCER_ROUND_ROBIN = 1
	KAFKA_BALANCER_HASH        = 3
)
