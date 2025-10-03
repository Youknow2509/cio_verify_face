package uuid

import (
	"time"

	ggUuid "github.com/google/uuid"
)

// ========================
//
//	UUID utils
//
// ========================
// Create time uuid
func GetTimeUUID() (ggUuid.UUID, error) {
	return ggUuid.NewV7()
}

// ParseUUID
func ParseUUID(uuidStr string) (ggUuid.UUID, error) {
	return ggUuid.Parse(uuidStr)
}

// Parse uuid time
func ParseTimeUUID(uuidStr string) (ggUuid.UUID, error) {
	return ggUuid.Parse(uuidStr)
}

// Get time in uuid
func GetTimeInUUID(uuid ggUuid.UUID) time.Time {
	sec, nsec := uuid.Time().UnixTime()
	return time.Unix(sec, nsec)
}