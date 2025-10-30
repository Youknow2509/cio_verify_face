package constants

const (
	// Rate limit spam register send otp
	RateLimitRegisterSendOTP = 5 // 5 requests per minute
	// Rate limit verify otp register
	RateLimitVerifyOTPRegisterFail = 5 // 5 requests per minute
	// Rate limit spam reset password send mail
	RateLimitResetPasswordSendMail = 5 // 5 requests per minute
)

const (
	RateDataListFriendRequestToUser = 10 // 10 requests per minute
	RateDataListFriendsOfUser       = 10 // 10 requests per minute
)
