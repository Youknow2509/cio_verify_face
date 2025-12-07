package dto

// VerifyFaceRequest represents the payload from device to verify a face.
// Supports both multipart form-data (binary image) and JSON (base64 image) for backward compatibility.
type VerifyFaceRequest struct {
	ImageBase64 string `json:"image_base64" form:"image_base64" validate:"omitempty"`
	UserId      string `json:"user_id,omitempty" form:"user_id"`
	SearchMode  string `json:"search_mode,omitempty" form:"search_mode"`
	TopK        int32  `json:"top_k,omitempty" form:"top_k" validate:"omitempty,min=1,max=10"`
}

// VerifyFaceMatchResponse maps a single match.
type VerifyFaceMatchResponse struct {
	UserId     string  `json:"user_id"`
	ProfileId  string  `json:"profile_id"`
	Similarity float32 `json:"similarity"`
	Confidence float32 `json:"confidence"`
	IsPrimary  bool    `json:"is_primary"`
}

// VerifyFaceResponse is returned to device clients.
type VerifyFaceResponse struct {
	Status        string                    `json:"status"`
	Verified      bool                      `json:"verified"`
	Matches       []VerifyFaceMatchResponse `json:"matches"`
	BestMatch     *VerifyFaceMatchResponse  `json:"best_match,omitempty"`
	Message       string                    `json:"message,omitempty"`
	LivenessScore *float32                  `json:"liveness_score,omitempty"`
}
