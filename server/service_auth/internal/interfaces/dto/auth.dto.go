package dto

// ======================
//
//	Auth DTOs
//
// ======================
type DeleteDeviceRequest struct {
	DeviceId  string `json:"device_id" validate:"required"`
	CompanyId string `json:"company_id" validate:"required"`
}

type UpdateDeviceRequest struct {
	CompanyId  string `json:"company_id" validate:"required"`
	DeviceName string `json:"device_name" validate:"required,min=2,max=100"`
	DeviceType int    `json:"device_type" validate:"required"`
	DeviceId   string `json:"device_id" validate:"required"`
}

type RegisterRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type VerifyOtpRequest struct {
	Email string `json:"email" validate:"required,email"`
	OTP   string `json:"otp" validate:"required,len=6"`
}

type CreateBaseAccountRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Token     string `json:"token" validate:"required,min=20,max=50"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string `json:"last_name" validate:"required,min=2,max=50"`
	Birthday  string `json:"birthday" validate:"required,datetime=2006-01-02"`
	Gender    int    `json:"gender" validate:"required,min=0,max=2"`
}

type LoginRequest struct {
	UserName string `json:"username" validate:"required,min=2,max=100"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginRequestApp struct {
	UserName  string `json:"username" validate:"required,min=2,max=100"`
	Password  string `json:"password" validate:"required,min=8"`
	DeviceId  string `json:"device_id" validate:"required,min=10,max=50"`
	PushToken string `json:"push_token" validate:"omitempty,min=10,max=200"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RefreshTokenRequest struct {
	AccessToken  string `json:"access_token" validate:"required"`
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required,min=8"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type ResetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ValidateResetPasswordTokenRequest struct {
	UserId string `json:"user_id" validate:"required"`
	Token  string `json:"token" validate:"required,min=20,max=50"`
}
