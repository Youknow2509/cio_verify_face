package model

// StatusType định nghĩa các loại trạng thái có thể có.
type StatusType string

const (
	StatusUp           StatusType = "UP"
	StatusDegraded     StatusType = "DEGRADED"
	StatusOutOfService StatusType = "OUT_OF_SERVICE"
	StatusDown         StatusType = "DOWN"
)

// HealthResponse là cấu trúc tổng thể của response.
type HealthResponse struct {
	Status           StatusType `json:"status"`
	UptimeSeconds    float64    `json:"uptime_seconds"`
	GracefulShutdown bool       `json:"gracefulShutdown"`
	Build            BuildInfo  `json:"build"`
	Details          Details    `json:"details"`
}

// BuildInfo chứa thông tin về phiên bản build.
type BuildInfo struct {
	Version    string `json:"version"`
	CommitHash string `json:"commit_hash"`
	BuildTime  string `json:"build_time"`
}

// Details chứa chi tiết các kiểm tra con.
type Details struct {
	WebSocketServer    *ComponentCheck `json:"webSocketServer"`
	DownstreamServices *ComponentCheck `json:"downstreamServices"`
	SystemResources    *ComponentCheck `json:"systemResources"`
}

// ComponentCheck đại diện cho một thành phần được kiểm tra.
type ComponentCheck struct {
	Status  StatusType             `json:"status"`
	Details map[string]interface{} `json:"details"`
}
