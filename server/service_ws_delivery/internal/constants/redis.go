package constants

// ================================================
//
//	Constants for Redis
//
// ================================================
const (
	RedisTypeStandalone = 1
	RedisTypeSentinel   = 2
	RedisTypeCluster    = 3
)

const (
	RedisPrefixRateLimiter = "ratelimiter:"
)
