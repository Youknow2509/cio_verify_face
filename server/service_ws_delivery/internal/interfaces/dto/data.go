package dto

// ==========================================
//
//	Details payload
//
// ==========================================
type (
	SendDataVerifyFace struct {
		DataUrl   string `json:"data_url" validate:"required,url"`
		Metadata  string `json:"metadata" validate:"omitempty"`
	}
)
