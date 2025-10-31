package model

// =====================================
// Model for domain
// =====================================
type (
	WsDataReceived struct {
		Type    WSEventType `json:"type"`
		Payload interface{} `json:"payload"`
	}

	WsDataSend struct {
		Type    WSEventType `json:"type"`
		Payload interface{} `json:"payload"`
	}
)
