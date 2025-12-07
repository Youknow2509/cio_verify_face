package model

// FaceVerifyInput holds data for verifying a face against stored profiles.
type FaceVerifyInput struct {
	ImageData  []byte
	CompanyId  string
	UserId     string
	DeviceId   string
	SearchMode string
	TopK       int32
}

// FaceVerifyMatch represents a matched profile.
type FaceVerifyMatch struct {
	UserId     string  `json:"user_id"`
	ProfileId  string  `json:"profile_id"`
	Similarity float32 `json:"similarity"`
	Confidence float32 `json:"confidence"`
	IsPrimary  bool    `json:"is_primary"`
}

// FaceVerifyOutput aggregates verification results.
type FaceVerifyOutput struct {
	Status        string            `json:"status"`
	Verified      bool              `json:"verified"`
	Matches       []FaceVerifyMatch `json:"matches"`
	BestMatch     *FaceVerifyMatch  `json:"best_match,omitempty"`
	Message       string            `json:"message,omitempty"`
	LivenessScore *float32          `json:"liveness_score,omitempty"`
}
