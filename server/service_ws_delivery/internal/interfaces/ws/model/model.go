package model

import (
	"encoding/json"

	"github.com/google/uuid"
)

/**
 * Client info
 */
type ClientInfo struct {
	UserId       uuid.UUID `json:"user_id"`
	SessionId    uuid.UUID `json:"session_id"`
	ConnectionId uuid.UUID `json:"connection_id"`
	IpAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
}

/**
 * Client writer data to worker
 */
type ClientWriterData struct {
	ClientInfo ClientInfo `json:"client_info"`
	Data       []byte     `json:"data"`
}

/**
 * Data client receive
 */
type DataClientReceive struct {
	TypeEvent int             `json:"type_event"`
	Payload   json.RawMessage `json:"payload"`
}

/**
 * Data send
 */
type DataSend struct {
	ConnectionId uuid.UUID `json:"connection_id"`
	Data         []byte    `json:"data"`
}
