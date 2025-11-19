package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	constants "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/constants"
	"github.com/youknow2509/cio_verify_face/server/service_analytic/internal/global"
)

// Helpers to access global caches safely
func local() *ristretto.Cache    { return global.LocalCache }
func redisClient() *redis.Client { return global.RedisClient }

// BuildDailyByDateKey builds cache key for daily summaries by date
func BuildDailyByDateKey(companyID uuid.UUID, date time.Time, deviceID *uuid.UUID) string {
	base := fmt.Sprintf(constants.CacheKeyDailyByDate, companyID.String(), date.Format("2006-01-02"))
	if deviceID != nil && *deviceID != uuid.Nil {
		return base + ":device:" + deviceID.String()
	}
	return base
}

// BuildDailyByMonthKey builds cache key for daily summaries by month
func BuildDailyByMonthKey(companyID uuid.UUID, month string) string {
	return fmt.Sprintf(constants.CacheKeyDailyByMonth, companyID.String(), month)
}

// BuildTotalEmployeesKey builds cache key for total employees per company
func BuildTotalEmployeesKey(companyID uuid.UUID) string {
	return fmt.Sprintf(constants.CacheKeyTotalEmployees, companyID.String())
}

// BuildExportKey builds cache key for export reports by company and date range
func BuildExportKey(companyID uuid.UUID, startDate, endDate, format string) string {
	return fmt.Sprintf(constants.CacheKeyExportReport, companyID.String(), startDate, endDate, format)
}

// SetLocal sets a value in local cache with TTL
func SetLocal(key string, value interface{}, ttl time.Duration) bool {
	lc := local()
	if lc == nil {
		return false
	}
	return lc.SetWithTTL(key, value, 1, ttl)
}

// GetLocal retrieves a value from local cache
func GetLocal(key string) (interface{}, bool) {
	lc := local()
	if lc == nil {
		return nil, false
	}
	return lc.Get(key)
}

// SetDistributed marshals value to JSON and stores in Redis with TTL
func SetDistributed(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	rc := redisClient()
	if rc == nil {
		return nil
	}
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return rc.Set(ctx, key, b, ttl).Err()
}

// GetDistributed fetches JSON from Redis and unmarshals into dst (must be pointer)
func GetDistributed(ctx context.Context, key string, dst interface{}) (bool, error) {
	rc := redisClient()
	if rc == nil {
		return false, nil
	}
	str, err := rc.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, json.Unmarshal([]byte(str), dst)
}
