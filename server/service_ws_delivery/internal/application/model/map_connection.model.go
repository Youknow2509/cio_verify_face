package model

// ================================================
//
//	Map connection model
//
// ================================================
type RegisterConnection struct {
	ConnectionId string `json:"connection_id"`
	DeviceId     string `json:"device_id"`
	IPAddress    string `json:"ip_address"`
	ConnectedAt  string `json:"connected_at"`
	UserAgent    string `json:"user_agent"`
}

type UnregisterConnection struct {
	DeviceId  string `json:"device_id"`
}
