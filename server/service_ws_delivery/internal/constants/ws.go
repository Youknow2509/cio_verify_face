package constants

// ==================================================
// 			WebSocket Constants
// ==================================================

// WS endpoint
const WS_ENDPOINT = "/ws"

// WS route event
const (
	ROUTES_HANDLE_EVENT_TO_KAFKA = iota
	ROUTES_HANDLE_EVENT_TO_CLIENT
)

// WS event types Client to Server - C2S
const (
	WS_C2S_TYPE_MESSAGE = iota
	WS_C2S_TYPE_TYPING
	// v.v
)

