package model

// ===============================================================
//
//	Define model for manager connection
//
// ===============================================================
type CreateConnectionInput struct {
	// KEYS:
	// 1: device_conns_key (ví dụ: device_conns:device123)
	// 2: service_conns_key (ví dụ: service_conns:notif-A)

	// ARGV:
	// 1: connection_id
	// 2: device_id
	// 3: service_id
	// 4: ip_address
	// 5: connected_at (timestamp)
	// 6: user_agent
	DeviceConnectionsKey  string `json:"device_connections_key"`
	ServiceConnectionsKey string `json:"service_connections_key"`
	ConnectionId          string `json:"connection_id"`
	DeviceId              string `json:"device_id"`
	ServiceId             string `json:"service_id"`
	IpAddress             string `json:"ip_address"`
	ConnectedAt           string `json:"connected_at"`
	UserAgent             string `json:"user_agent"`
}

type RemoveConnectionInput struct {
	// KEYS:
	// 1: device_conns_key
	// 2: service_conns_key
	// ARGV:
	// 1: device_id
	DeviceId              string `json:"device_id"`
	DeviceConnectionsKey  string `json:"device_connections_key"`
	ServiceConnectionsKey string `json:"service_connections_key"`
}
