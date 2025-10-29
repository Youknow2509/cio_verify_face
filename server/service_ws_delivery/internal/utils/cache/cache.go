package cache

import "fmt"

// ===================================
// 	Utils for cache
// ===================================

// ===================================
// 	Keys cache
// ===================================

// Get Device connection ws key
func GetDeviceConnectionWsKey(deviceId string) string {
	return fmt.Sprintf("ws:device:%s", deviceId)
}

// Get service ws connection key
func GetServiceWsConnectionKey(serviceId string) string {
	return fmt.Sprintf("ws:service:%s", serviceId)
}
