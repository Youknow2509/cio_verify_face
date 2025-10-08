package model

type KafkaAttendanceVerifyReceived struct {
	ServiceId string `json:"service_id"`
	DeviceId  string `json:"device_id"`
	DataUrl   string `json:"data_url"`
	Metadata  string `json:"metadata"`
	Timestamp int64  `json:"timestamp"`
}