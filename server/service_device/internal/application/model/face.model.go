package model

// VerifyFaceInput holds the data needed to call the face verification service.
type VerifyFaceInput struct {
	ImageData  []byte
	CompanyId  string
	DeviceId   string
	UserId     string
	SearchMode string
	TopK       int32
}

// VerifyFaceMatch describes an individual match.
type VerifyFaceMatch struct {
	UserId     string  `json:"user_id"`
	ProfileId  string  `json:"profile_id"`
	Similarity float32 `json:"similarity"`
	Confidence float32 `json:"confidence"`
	IsPrimary  bool    `json:"is_primary"`
}

// VerifyFaceOutput returns verification results in API-friendly shape.
type VerifyFaceOutput struct {
	Status        string            `json:"status"`
	Verified      bool              `json:"verified"`
	Matches       []VerifyFaceMatch `json:"matches"`
	BestMatch     *VerifyFaceMatch  `json:"best_match,omitempty"`
	Message       string            `json:"message,omitempty"`
	LivenessScore *float32          `json:"liveness_score,omitempty"`
}
