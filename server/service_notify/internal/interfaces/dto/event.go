package dto

// ======================
//
//	Event DTOs
//
// ======================
type KafkaEvent struct {
	EventType int         `json:"event_type" validate:"required"`
	Payload   interface{} `json:"payload" validate:"required"`
}
