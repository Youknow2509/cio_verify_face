package cache

import "fmt"

// ===================================
// 	Utils for cache
// ===================================

// ===================================
// 	Keys cache
// ===================================

// Get user connection ws key
func GetUserConnectionWsKey(userId string) string {
	return fmt.Sprintf("ws:connection:%s", userId)
}

// Get service ws connection key
func GetServiceWsConnectionKey(serviceId string) string {
	return fmt.Sprintf("ws:connection:%s", serviceId)
}

// Get ws connection info key
func GetWsConnectionInfoKey(serviceId string, connectionId string) string {
	return fmt.Sprintf("ws:connection:%s:%s", serviceId, connectionId)
}
