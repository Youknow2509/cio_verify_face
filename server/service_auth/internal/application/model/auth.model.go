package model

import "github.com/google/uuid"

// =======================================================
//
//	For Input Auth Model, Validation
//
// =======================================================
type (
	AuthRegisterInput struct {
		Email string `json:"email"`
	}

	VerifyOTPRegisterInput struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	AuthCreateBaseRegisterInput struct {
		Token     string `json:"token"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Birthday  string `json:"birthday"`
		Gender    int    `json:"gender"`
	}

	LoginInput struct {
		UserName  string `json:"username"`
		Password  string `json:"password"`
		ClientIp  string `json:"client_ip"`
		UserAgent string `json:"user_agent"`
	}

	LoginInputAdmin struct {
		UserName  string `json:"username"`
		Password  string `json:"password"`
		ClientIp  string `json:"client_ip"`
		UserAgent string `json:"user_agent"`
	}

	LogoutInput struct {
		UserId    uuid.UUID `json:"user_id"`
		SessionId uuid.UUID `json:"session_id"`
	}

	LogoutAllInput struct {
		// TODO: Add validation for logout all input
	}

	RefreshTokenInput struct {
		ClientIp     string    `json:"client_ip"`
		RefreshToken string    `json:"refresh_token"`
		UserId       uuid.UUID `json:"user_id"`
		SessionId    uuid.UUID `json:"session_id"`
		UserRole     int       `json:"user_role"`
	}

	GetMyInfoInput struct {
		ClientIp string    `json:"client_ip"`
		UserId   uuid.UUID `json:"user_id"`
		Role     int       `json:"role"`
	}

	// =======================================================

	ChangePasswordInput struct {
		ClientIp    string `json:"client_ip"`
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
		UserId      string `json:"user_id"`
		SessionId   string `json:"session_id"`
	}

	ResetPasswordInput struct {
		ClientIp string `json:"client_ip"`
		Email    string `json:"email"`
	}

	ValidateResetPasswordTokenInput struct {
		ClientIp string `json:"client_ip"`
		UserId   string `json:"user_id"`
		Token    string `json:"token"`
	}

	UpdateDeviceSessionInput struct {
		UserId     uuid.UUID `json:"user_id"`
		SessionId  uuid.UUID `json:"session_id"`
		ClientIp   string    `json:"client_ip"`
		UserAgent  string    `json:"user_agent"`
		Role       int       `json:"role"`
		CompanyId  uuid.UUID `json:"company_id"`
		DeviceName string    `json:"device_name"`
		DeviceType int       `json:"device_type"`
		DeviceId   uuid.UUID `json:"device_id"`
	}

	DeleteDeviceSessionInput struct {
		UserId    uuid.UUID `json:"user_id"`
		SessionId uuid.UUID `json:"session_id"`
		ClientIp  string    `json:"client_ip"`
		UserAgent string    `json:"user_agent"`
		Role      int       `json:"role"`
		CompanyId uuid.UUID `json:"company_id"`
		DeviceId  uuid.UUID `json:"device_id"`
	}
)

// =======================================================
//
//	For Output Auth Model
//
// =======================================================
type (
	UpdateDeviceSessionOutput struct {
		Token    string `json:"token"`
		ExpireAt int64  `json:"expire_at"`
	}

	GetMyInfoOutput struct {
		UserId    string `json:"user_id"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
		FullName  string `json:"full_name"`
		AvatarURL string `json:"avatar_url"`
		Role      int    `json:"role"`
	}

	VerifyOTPRegisterOutput struct {
		Token    string `json:"token"`
		ExpireAt int64  `json:"expire_at"`
	}

	LoginOutput struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	RefreshTokenOutput struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	// Token validation result for cached operations
	TokenValidationResult struct {
		Valid  bool      `json:"valid"`
		UserID uuid.UUID `json:"user_id,omitempty"`
		Role   int       `json:"role,omitempty"`
		Error  string    `json:"error,omitempty"`
	}

	// User info output for cached operations (re-export from domain)
	UserInfoOutput struct {
		Email     string `json:"email"`
		Phone     string `json:"phone"`
		FullName  string `json:"full_name"`
		AvatarURL string `json:"avatar_url"`
	}

	// Cache statistics for monitoring
	CacheStats struct {
		LocalCacheHits       int64   `json:"local_cache_hits"`
		DistributedCacheHits int64   `json:"distributed_cache_hits"`
		CacheMisses          int64   `json:"cache_misses"`
		TotalRequests        int64   `json:"total_requests"`
		HitRatio             float64 `json:"hit_ratio"`
	}
)
