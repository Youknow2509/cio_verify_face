package errors

// ========================================
//
//	Token when parse error
//
// ========================================
type TokenValidationError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Define constants for error codes
const (
	TokenValidationErrorCode     = 1001
	TokenMalformedErrorCode      = 1002
	TokenSignatureInvalidErrCode = 1003
	TokenExpiredErrorCode        = 1004
	TokenErrorNotFoundCode       = 1005
)

var (
	mapMessageError = map[int]string{
		TokenValidationErrorCode:     "Token validation error",
		TokenMalformedErrorCode:      "Token is malformed",
		TokenSignatureInvalidErrCode: "Token signature is invalid",
		TokenExpiredErrorCode:        "Token has expired",
		TokenErrorNotFoundCode:       "Token error not found",
	}
)

/**
 * Get error token parse with code
 */
func GetTokenValidationError(code int) *TokenValidationError {
	if message, ok := mapMessageError[code]; ok {
		return &TokenValidationError{
			Code:    code,
			Message: message,
		}
	}
	return &TokenValidationError{
		Code:    code,
		Message: "Unknown token error",
	}
}

/**
 * New error with error
 */
func NewTokenServiceValidationError(err error) *TokenValidationError {
	return &TokenValidationError{
		Code:    TokenValidationErrorCode,
		Message: err.Error(),
	}
}
