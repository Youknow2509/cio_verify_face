package errors

// =======================================================
//
//	Define user errors returned by the service
//
// =======================================================
const (
	UserAlreadyExistsErrorCode                   = 20001
	InvalidEmailFormatErrorCode                  = 20002
	UserRegistrationFailedErrorCode              = 20003
	UserNotFoundErrorCode                        = 20004
	UserAuthenticationFailedErrorCode            = 20005
	UserSpamRegisterErrorCode                    = 20006
	UserInvalidOTPErrorCode                      = 20007
	UserInvalidTokenErrorCode                    = 20008
	UserBlockedErrorCode                         = 20009
	UserPasswordIncorrectErrorCode               = 20010
	UserFriendRequestNotFoundErrorCode           = 20011
	UserDontHaveRoleCancelFriendRequestErrorCode = 20012
)

var mapUserErrors = map[int]string{
	UserDontHaveRoleCancelFriendRequestErrorCode: "User don't have role cancel friend request",
	UserFriendRequestNotFoundErrorCode:           "Not found friend request",
	UserPasswordIncorrectErrorCode:               "Incorrect password",
	UserBlockedErrorCode:                         "User is blocked",
	UserInvalidTokenErrorCode:                    "Invalid token",
	UserInvalidOTPErrorCode:                      "Invalid OTP",
	UserAlreadyExistsErrorCode:                   "User already exists",
	InvalidEmailFormatErrorCode:                  "Invalid email format",
	UserRegistrationFailedErrorCode:              "User registration failed",
	UserNotFoundErrorCode:                        "User not found",
	UserAuthenticationFailedErrorCode:            "User authentication failed",
	UserSpamRegisterErrorCode:                    "User spam register detected",
}
