package dto

// ======================
//
//	Event DTOs
//
// ======================
type KafkaEvent struct {
	EventType int         `json:"event_type" validate:"gte=0"`
	Payload   interface{} `json:"payload" validate:"required"`
}
