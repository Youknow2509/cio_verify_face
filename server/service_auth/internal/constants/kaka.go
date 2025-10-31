package constants

// ================================================
// 		Constants for Kafka client
// ================================================

// Base
const (
	KAFKA_DIALER_TIMEOUT = 10   // seconds
	KAFKA_DUAL_STACK     = true // Dual stack support
)

// Topic name
const (
	KAFKA_TOPIC_NOTIFICATION             = "notification_requests"
	KAFKA_TOPIC_NEW_MESSAGE              = "new_messages"
	KAFKA_TOPIC_MESSAGE_ACKNOWLEDGEMENTS = "message_acknowledgements"
	KAFKA_TOPIC_MESSAGE_DELIVERIES       = "message_deliveries"
	KAFKA_TOPIC_MESSAGE_READ_STATUS      = "message_read_status"
	KAFKA_TOPIC_MESSAGE_REACT            = "message_react"
	KAFKA_TOPIC_MESSAGE_EDIT             = "message_edit"
	KAFKA_TOPIC_MESSAGE_DELETE           = "message_delete"
	KAFKA_TOPIC_MESSAGE_ACTION           = "message_action"
	KAFKA_TOPIC_USER_PRESENCE_EVENTS     = "user_presence_events"
	KAFKA_TOPIC_USER_TYPING_EVENTS       = "user_typing_events"
	KAFKA_TOPIC_USER_CALL_EVENTS         = "user_call_events"
	KAFKA_TOPIC_MESSAGE_MODIFICATIONS    = "message_modifications"
	KAFKA_TOPIC_MEDIA_PROCESSING_JOBS    = "media_processing_jobs"
	KAFKA_TOPIC_AUDIT_EVENTS             = "audit_events"
	KAFKA_TOPIC_DB_REPLICATION_EVENTS    = "db_replication_events"
	// v.v
)

// type event notification
const (
	KAFKA_EVENT_TYPE_SEND_OTP_REGISTER                = 1 // Send OTP for registration
	KAFKA_EVENT_TYPE_PUSH_NOTIFICATION                = 2 // Push notification event
	KAFKA_EVENT_TYPE_SEND_TOKEN_RESET_PASSWORD        = 3 // Send reset password event
	KAFKA_EVENT_TYPE_PUSH_NOTIFICATION_FRIEND_REQUEST = 4 // Push notification for friend request
	// v.v
)

// SASL Mechanism
const (
	KAFKA_SASL_MECHANISM_PLAIN        = 0 // PLAIN
	KAFKA_SASL_MECHANISM_SCRAM_SHA256 = 1 // SCRAM-SHA-256
	KAFKA_SASL_MECHANISM_SCRAM_SHA512 = 2 // SCRAM-SHA-512
)

// Balancer Type
const (
	// # 0: custom config, 1: RoundRobin, 2: LeastBytes, 3: Hash, 4: ReferenceHash, 5: CRC32Balancer, 6: Murmur2Balancer, ...
	KAFKA_BALANCER_CUSTOM         = 0 // Custom balancer
	KAFKA_BALANCER_ROUND_ROBIN    = 1 // Round Robin balancer
	KAFKA_BALANCER_LEAST_BYTES    = 2 // Least Bytes balancer
	KAFKA_BALANCER_HASH           = 3 // Hash balancer
	KAFKA_BALANCER_REFERENCE_HASH = 4 // Reference Hash balancer
	KAFKA_BALANCER_CRC32          = 5 // CRC32 balancer
	KAFKA_BALANCER_MURMUR2        = 6 // Murmur2 balancer
)

// Acks
const (
	KAFKA_ACKS_NONE   = 0 // No response required
	KAFKA_ACKS_LEADER = 1 // Only leader ack
	KAFKA_ACKS_ALL    = 2 // All in-sync replicas must ack
)

// Compression
const (
	KAFKA_COMPRESSION_NONE   = 0 // No compression (default)
	KAFKA_COMPRESSION_GZIP   = 1 // GZIP compression (tốt cho tỉ lệ nén, nhưng CPU cao)
	KAFKA_COMPRESSION_SNAPPY = 2 // Snappy compression (nhanh, hiệu quả, thường dùng)
	KAFKA_COMPRESSION_LZ4    = 3 // LZ4 compression (rất nhanh, tốt cho throughput)
	KAFKA_COMPRESSION_ZSTD   = 4 // ZSTD compression (hiệu suất tốt cả về tốc độ và nén)
)

// Auto Offset Reset
const (
	KAFKA_AUTO_OFFSET_RESET_EARLIEST = 0 // Start from earliest offset
	KAFKA_AUTO_OFFSET_RESET_LATEST   = 1 // Start from latest offset
)

// Enable Auto Commit
const (
	KAFKA_ENABLE_AUTO_COMMIT_FALSE = 0 // Disable auto commit
	KAFKA_ENABLE_AUTO_COMMIT_TRUE  = 1 // Enable auto commit
)
