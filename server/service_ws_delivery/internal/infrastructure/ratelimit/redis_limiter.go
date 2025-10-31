package ratelimit

import (
	"context"
	"strconv"

	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/constants"
	domainCache "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/cache"
	domainModel "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/model"
	domainLimiter "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/ratelimit"
)

// ===============================================
// Redis cached limiter
// ===============================================
type RedisLimiter struct {
	distributed domainCache.IDistributedCache
	policy      domainLimiter.Policy
	prefix      string
}

// Create implements ratelimit.ILimiter.
func (r *RedisLimiter) Create(ctx context.Context, key string) error {
	script := getLuaScriptCreat()
	keySave := r.getKey(key)
	ttl := int(r.policy.Window.Milliseconds())
	_, err := r.distributed.LuaScript(ctx, script, []string{keySave}, ttl)
	if err != nil {
		return err
	}
	return nil
}

// Allow implements ratelimit.ILimiter.
func (r *RedisLimiter) Allow(ctx context.Context, key string) (domainModel.Verdict, error) {
	script := getLuaScriptAllow()
	keySave := r.getKey(key)
	limit := r.policy.Limit
	ttl := int(r.policy.Window.Milliseconds())

	result, err := r.distributed.LuaScript(ctx, script, []string{keySave}, limit, ttl)
	if err != nil {
		return domainModel.Denied, err
	}

	resultInt, ok := result.(int64)
	if !ok {
		return domainModel.Denied, nil
	}

	if resultInt == 0 {
		return domainModel.Denied, nil
	}
	return domainModel.Allowed, nil
}

// Check implements ratelimit.ILimiter.
func (r *RedisLimiter) Check(ctx context.Context, key string) (int, bool, error) {
	val, err := r.distributed.Get(ctx, r.getKey(key))
	if err != nil {
		return 0, false, err
	}
	if val == "" {
		return 0, false, nil
	}
	valInt, err := strconv.Atoi(val)
	if err != nil {
		return 0, false, err
	}
	if valInt > r.policy.Limit {
		return valInt, false, nil
	}
	return valInt, true, nil
}

// Upgrade implements ratelimit.ILimiter.
func (r *RedisLimiter) Upgrade(ctx context.Context, key string, val int) error {
	script := getLuaScriptUpgrade()
	keySave := r.getKey(key)
	_, err := r.distributed.LuaScript(ctx, script, []string{keySave}, val)
	if err != nil {
		return err
	}
	return nil
}

/**
 * New and impl redis limiter
 */
func NewRedisLimiter(
	distributed domainCache.IDistributedCache,
	policy domainLimiter.Policy,
	prefix string,
) domainLimiter.ILimiter {
	if prefix == "" {
		prefix = constants.RedisPrefixRateLimiter
	}
	return &RedisLimiter{
		distributed: distributed,
		policy:      policy,
		prefix:      prefix,
	}
}

// ===============================================
// Helper functions
// ===============================================
func (store *RedisLimiter) getKey(key string) string {
	return store.prefix + key
}

// Lua allow - Fixed version
func getLuaScriptAllow() string {
	return `
        local key = KEYS[1]
        local limit = tonumber(ARGV[1])
        local ttl_ms = tonumber(ARGV[2])
        
        local ttl = redis.call("PTTL", key)
        
        -- Key doesn't exist (-2) or expired (-1)
        if ttl == -2 or ttl == -1 then
            -- Create new key with value 1
            redis.call("SET", key, "1", "PX", ttl_ms)
            return 1
        end
        
        -- Key exists, check current value
        local current = tonumber(redis.call("GET", key) or "0")
        if current >= limit then
            return 0  -- Rate limit exceeded
        else
            redis.call("INCR", key)
            return 1  -- Request allowed
        end
    `
}

// Lua upgrade
func getLuaScriptUpgrade() string {
	return `
		local key = KEYS[1]
		local new_val = tonumber(ARGV[1])
		local ttl = redis.call("PTTL", key)
		if ttl <= 0 then
			return {-1, -1}
		end
		local r = redis.call("SET", key, new_val, "PX", ttl)
		if r ~= "OK" then
			return {-1, -1}
		end
		return {new_val, ttl}
	`
}

// Lua create
func getLuaScriptCreat() string {
	return `
		local key = KEYS[1]
		local val = 1
		local ttl = ARGV[1]
		if redis.call("EXISTS", key) == 1 then
			redis.call("INCR", key)
			return -1
		end
		local r = redis.call("SET", key, val, "PX", ttl)
		if r ~= "OK" then
			return -1
		end
		return val
	`
}
