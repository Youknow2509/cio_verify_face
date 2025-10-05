package model

// ================================================
//
//	Map connection model
//
// ================================================
type RegisterConnection struct {
	ConnectionId string `json:"connection_id"`
	UserId       string `json:"user_id"`
	IPAddress    string `json:"ip_address"`
	ConnectedAt  string `json:"connected_at"`
	UserAgent    string `json:"user_agent"`
}

type UnregisterConnection struct {
	ConnectionId string `json:"connection_id"`
	UserId       string `json:"user_id"`
}
