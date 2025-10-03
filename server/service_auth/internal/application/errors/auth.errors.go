package errors

// =======================================================
//
//	Define auth errors returned by the service
//
// =======================================================
const (
	AuthAccessTokenRequiredErrorCode            = 10001
	AuthTokenInvalidErrorCode                   = 10002
	AuthCannotRefreshTokenErrorCode             = 10003
	AuthResetPasswordSpamErrorCode              = 10004
	AuthUserSpamRefreshTokenErrorCode           = 10005
	AuthValidateResetPasswordTokenSpamErrorCode = 10006
	AuthResetPasswordTokenInvalidErrorCode      = 10007
	AuthUserIdInvalidErrorCode                  = 10008
	AuthUUIDParseErrorCode                      = 10009
	AuthDontHavePermissionErrorCode             = 10010
	TokenExpiredErrorCode                       = 20010
)

var mapAuthErrors = map[int]string{
	AuthDontHavePermissionErrorCode:             "Don't have permission",
	TokenExpiredErrorCode:                       "Token is expired",
	AuthUUIDParseErrorCode:                      "UUID parse error",
	AuthUserIdInvalidErrorCode:                  "User ID is invalid",
	AuthResetPasswordTokenInvalidErrorCode:      "Reset password token is invalid",
	AuthValidateResetPasswordTokenSpamErrorCode: "Spam reset password token detected",
	AuthUserSpamRefreshTokenErrorCode:           "User spam refresh token detected",
	AuthResetPasswordSpamErrorCode:              "Spam reset password detected",
	AuthCannotRefreshTokenErrorCode:             "Cannot refresh token",
	AuthTokenInvalidErrorCode:                   "Token is invalid",
	AuthAccessTokenRequiredErrorCode:            "Access token is required",
}
