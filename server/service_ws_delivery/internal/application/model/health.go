package model

type SystemDetailsInput struct {
	ClientIp string `json:"clientIp"`
}

// SystemDetails chứa chi tiết các kiểm tra con.
type SystemDetails struct {
	WebSocketServer    *ComponentCheck `json:"webSocketServer"`
	DownstreamServices *ComponentCheck `json:"downstreamServices"`
	SystemResources    *ComponentCheck `json:"systemResources"`
}

// ComponentCheck đại diện cho một thành phần được kiểm tra.
type ComponentCheck struct {
	Status  string                 `json:"status"`
	Details map[string]interface{} `json:"details"`
}
