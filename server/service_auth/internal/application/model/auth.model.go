package model

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
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	LogoutAllInput struct {
		// TODO: Add validation for logout all input
	}

	RefreshTokenInput struct {
		ClientIp     string `json:"client_ip"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	GetMyInfoInput struct {
		ClientIp  string `json:"client_ip"`
		UserId    string `json:"user_id"`
		SessionId string `json:"session_id"`
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

	CreateDeviceSessionInput struct {
		// TODO: Add fields for CreateDeviceSessionInput
	}

	DeleteDeviceSessionInput struct {
		// TODO: Add fields
	}

	RefreshDeviceSessionInput struct {
		// TODO: Add fields

	}
)

// =======================================================
//
//	For Output Auth Model
//
// =======================================================
type (
	RefreshDeviceSessionOutput struct {
		// TODO: Add fields
	}

	CreateDeviceSessionOutput struct {
		// TODO: Add fields
	}

	GetMyInfoOutput struct {
		// TODO: Add fields for GetMyInfoOutput
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

	ValidateJwtUserOutput struct {
		// TODO: Add fields for ValidateJwtUserOutput
	}

	ValidateJwtServiceOutput struct {
		// TODO: Add fields for ValidateJwtServiceOutput
	}
)
