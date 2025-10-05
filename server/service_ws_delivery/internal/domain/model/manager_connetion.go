package model

// ===============================================================
//
//	Define model for manager connection
//
// ===============================================================
type CreateConnectionInput struct {
	// Keys for Redis cache
	UserConnectionsKey    string `json:"user_connections_key"`
	ServiceConnectionsKey string `json:"service_connections_key"`
	ConnectionKey         string `json:"connection_key"`
	// Values for Redis cache
	ConnectionId string `json:"connection_id"`
	UserId       string `json:"user_id"`
	ServiceId    string `json:"service_id"`
	IPAddress    string `json:"ip_address"`
	ConnectedAt  string `json:"connected_at"`
	UserAgent    string `json:"user_agent"`
	// More fields can be added as needed
	MaxConnsPerUser int `json:"max_conns_per_user"`
}

type RemoveConnectionInput struct {
	UserConnectionsKey    string `json:"user_connections_key"`
	ServiceConnectionsKey string `json:"service_connections_key"`
	ConnectionKey         string `json:"connection_key"`
	ConnectionId          string `json:"connection_id"`
}
