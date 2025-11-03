package router

// ============================================
// 			Http router manager
// ============================================
type RouterManager struct{
	AttendanceRouter
}

// NewRouterManager creates a new instance of RouterManager
func NewRouterManager() *RouterManager {
	return &RouterManager{}
}