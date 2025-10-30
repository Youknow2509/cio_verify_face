package model

// ======================================
// Enums use in model
// ======================================

// Ws event types
type WSEventType int

const (
	WSEventReceivedAttendance WSEventType = iota
	WSEventSendAttendance
	WSEventDeviceStatus
	WSEventAdminAlert
)
