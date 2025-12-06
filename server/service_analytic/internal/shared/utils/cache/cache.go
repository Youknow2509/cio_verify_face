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

// BuildDailyDetailByDateKey builds cache key for detailed daily report by date
func BuildDailyDetailByDateKey(companyID uuid.UUID, date time.Time) string {
	return fmt.Sprintf("%s:daily_detail:%s", companyID.String(), date.Format("2006-01-02"))
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

// ============================================
// Additional Cache Key Builders for Attendance Records
// ============================================

// BuildAttendanceRecordsKey builds cache key for attendance records by company and month
func BuildAttendanceRecordsKey(companyID uuid.UUID, yearMonth string, limit int) string {
	return fmt.Sprintf("attendance_records:%s:%s:limit:%d", companyID.String(), yearMonth, limit)
}

// BuildAttendanceRecordsByTimeRangeKey builds cache key for attendance records by time range
func BuildAttendanceRecordsByTimeRangeKey(companyID uuid.UUID, yearMonth string, startTime, endTime time.Time) string {
	return fmt.Sprintf("attendance_records_time_range:%s:%s:%s:%s", companyID.String(), yearMonth, startTime.Format("2006-01-02"), endTime.Format("2006-01-02"))
}

// BuildAttendanceRecordsByEmployeeKey builds cache key for attendance records by employee
func BuildAttendanceRecordsByEmployeeKey(companyID uuid.UUID, yearMonth string, employeeID uuid.UUID) string {
	return fmt.Sprintf("attendance_records_employee:%s:%s:%s", companyID.String(), yearMonth, employeeID.String())
}

// BuildAttendanceRecordsByUserKey builds cache key for attendance records by user
func BuildAttendanceRecordsByUserKey(companyID, employeeID uuid.UUID, yearMonth string, limit int) string {
	return fmt.Sprintf("attendance_records_user:%s:%s:%s:limit:%d", companyID.String(), employeeID.String(), yearMonth, limit)
}

// BuildAttendanceRecordsByUserTimeRangeKey builds cache key for attendance records by user time range
func BuildAttendanceRecordsByUserTimeRangeKey(companyID, employeeID uuid.UUID, yearMonth string, startTime, endTime time.Time) string {
	return fmt.Sprintf("attendance_records_user_time_range:%s:%s:%s:%s:%s", companyID.String(), employeeID.String(), yearMonth, startTime.Format("2006-01-02"), endTime.Format("2006-01-02"))
}

// BuildAttendanceRecordsNoShiftKey builds cache key for attendance records without shift
func BuildAttendanceRecordsNoShiftKey(companyID uuid.UUID, yearMonth string, limit int) string {
	return fmt.Sprintf("attendance_records_no_shift:%s:%s:limit:%d", companyID.String(), yearMonth, limit)
}

// ============================================
// Additional Cache Key Builders for Daily Summaries
// ============================================

// BuildDailySummariesByMonthKey builds cache key for daily summaries by month
func BuildDailySummariesByMonthKey(companyID uuid.UUID, month string) string {
	return fmt.Sprintf("daily_summaries:%s:%s", companyID.String(), month)
}

// BuildDailySummaryByEmployeeDateKey builds cache key for daily summary by employee and date
func BuildDailySummaryByEmployeeDateKey(companyID uuid.UUID, month string, workDate time.Time, employeeID uuid.UUID) string {
	return fmt.Sprintf("daily_summary_employee_date:%s:%s:%s:%s", companyID.String(), month, workDate.Format("2006-01-02"), employeeID.String())
}

// BuildDailySummariesByUserKey builds cache key for daily summaries by user
func BuildDailySummariesByUserKey(companyID, employeeID uuid.UUID, month string) string {
	return fmt.Sprintf("daily_summaries_user:%s:%s:%s", companyID.String(), employeeID.String(), month)
}

// BuildDailySummaryByUserDateKey builds cache key for daily summary by user and date
func BuildDailySummaryByUserDateKey(companyID, employeeID uuid.UUID, month string, workDate time.Time) string {
	return fmt.Sprintf("daily_summary_user_date:%s:%s:%s:%s", companyID.String(), employeeID.String(), month, workDate.Format("2006-01-02"))
}

// ============================================
// Additional Cache Key Builders for Audit Logs
// ============================================

// BuildAuditLogsKey builds cache key for audit logs
func BuildAuditLogsKey(companyID uuid.UUID, yearMonth string, limit int) string {
	return fmt.Sprintf("audit_logs:%s:%s:limit:%d", companyID.String(), yearMonth, limit)
}

// BuildAuditLogsByTimeRangeKey builds cache key for audit logs by time range
func BuildAuditLogsByTimeRangeKey(companyID uuid.UUID, yearMonth string, startTime, endTime time.Time) string {
	return fmt.Sprintf("audit_logs_time_range:%s:%s:%s:%s", companyID.String(), yearMonth, startTime.Format("2006-01-02"), endTime.Format("2006-01-02"))
}

// ============================================
// Additional Cache Key Builders for Face Enrollment Logs
// ============================================

// BuildFaceEnrollmentLogsKey builds cache key for face enrollment logs
func BuildFaceEnrollmentLogsKey(companyID uuid.UUID, yearMonth string, limit int) string {
	return fmt.Sprintf("face_enrollment_logs:%s:%s:limit:%d", companyID.String(), yearMonth, limit)
}

// BuildFaceEnrollmentLogsByEmployeeKey builds cache key for face enrollment logs by employee
func BuildFaceEnrollmentLogsByEmployeeKey(companyID uuid.UUID, yearMonth string, employeeID uuid.UUID) string {
	return fmt.Sprintf("face_enrollment_logs_employee:%s:%s:%s", companyID.String(), yearMonth, employeeID.String())
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

// ============================================
// Distributed Lock Functions (for scale-out)
// ============================================

// AcquireLock attempts to acquire a distributed lock using Redis SETNX
// Returns true if lock was acquired, false if already locked by another instance
func AcquireLock(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	rc := redisClient()
	if rc == nil {
		return false, fmt.Errorf("redis client not available")
	}
	lockKey := "lock:" + key
	return rc.SetNX(ctx, lockKey, time.Now().Unix(), ttl).Result()
}

// ReleaseLock releases a distributed lock
func ReleaseLock(ctx context.Context, key string) error {
	rc := redisClient()
	if rc == nil {
		return nil
	}
	lockKey := "lock:" + key
	return rc.Del(ctx, lockKey).Err()
}

// ExtendLock extends the TTL of an existing lock (for long-running operations)
func ExtendLock(ctx context.Context, key string, ttl time.Duration) error {
	rc := redisClient()
	if rc == nil {
		return nil
	}
	lockKey := "lock:" + key
	return rc.Expire(ctx, lockKey, ttl).Err()
}

// ============================================
// Redis-Only Cache Functions (Skip Local Cache)
// For mutable data that changes frequently across instances
// ============================================

// SetDistributedOnly sets value ONLY in Redis, skipping local cache
// Use this for mutable data like export job status that changes across instances
func SetDistributedOnly(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return SetDistributed(ctx, key, value, ttl)
}

// GetDistributedOnly gets value ONLY from Redis, skipping local cache
// Use this for mutable data to avoid stale local cache in scale-out scenarios
func GetDistributedOnly(ctx context.Context, key string, dst interface{}) (bool, error) {
	return GetDistributed(ctx, key, dst)
}

// DeleteDistributed removes a key from Redis
func DeleteDistributed(ctx context.Context, key string) error {
	rc := redisClient()
	if rc == nil {
		return nil
	}
	return rc.Del(ctx, key).Err()
}

// DeleteLocal removes a key from local cache
func DeleteLocal(key string) {
	lc := local()
	if lc != nil {
		lc.Del(key)
	}
}
