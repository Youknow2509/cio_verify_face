package routes

// ======================================
// Routes manager system
// ======================================
type RouteManager struct {
	HealthRoute
}

// GetRouteManager creates a new RouteManager
func GetRouteManager() *RouteManager {
	return &RouteManager{}
}
