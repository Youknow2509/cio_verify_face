package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	appErrors "github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/application/errors"
)

// Response is the standard API response structure
type Response struct {
	Success bool        `json:"success"`
	Code    int         `json:"code,omitempty"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success sends a success response
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Success",
		Data:    data,
	})
}

// SuccessWithMessage sends a success response with custom message
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// BadRequest sends a 400 error response
func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Code:    appErrors.ErrCodeInvalidInput,
		Message: message,
	})
}

// ValidationError sends a 400 error response for validation failures
func ValidationError(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Code:    appErrors.ErrCodeInvalidInput,
		Message: "Validation error: " + message,
	})
}

// Unauthorized sends a 401 error response
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, Response{
		Success: false,
		Code:    appErrors.ErrCodeUnauthorized,
		Message: message,
	})
}

// Forbidden sends a 403 error response
func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, Response{
		Success: false,
		Code:    appErrors.ErrCodeForbidden,
		Message: message,
	})
}

// NotFound sends a 404 error response
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, Response{
		Success: false,
		Code:    appErrors.ErrCodeNotFound,
		Message: message,
	})
}

// Conflict sends a 409 error response
func Conflict(c *gin.Context, message string) {
	c.JSON(http.StatusConflict, Response{
		Success: false,
		Code:    appErrors.ErrCodeConflict,
		Message: message,
	})
}

// TooManyRequests sends a 429 error response
func TooManyRequests(c *gin.Context, message string) {
	c.JSON(http.StatusTooManyRequests, Response{
		Success: false,
		Code:    appErrors.ErrCodeRateLimitExceeded,
		Message: message,
	})
}

// InternalError sends a 500 error response
func InternalError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, Response{
		Success: false,
		Code:    appErrors.ErrCodeServiceUnavailable,
		Message: message,
	})
}

// FromAppError converts an application error to HTTP response
func FromAppError(c *gin.Context, err *appErrors.Error) {
	var httpStatus int

	switch {
	case err.Code >= 1000 && err.Code < 2000:
		// General errors
		switch err.Code {
		case appErrors.ErrCodeUnauthorized:
			httpStatus = http.StatusUnauthorized
		case appErrors.ErrCodeForbidden:
			httpStatus = http.StatusForbidden
		case appErrors.ErrCodeNotFound:
			httpStatus = http.StatusNotFound
		case appErrors.ErrCodeConflict:
			httpStatus = http.StatusConflict
		case appErrors.ErrCodeServiceUnavailable:
			httpStatus = http.StatusServiceUnavailable
		default:
			httpStatus = http.StatusBadRequest
		}
	case err.Code >= 2000 && err.Code < 3000:
		// Rate limiting errors
		httpStatus = http.StatusTooManyRequests
	case err.Code >= 3000 && err.Code < 4000:
		// Face profile update errors
		switch err.Code {
		case appErrors.ErrCodeRequestNotFound:
			httpStatus = http.StatusNotFound
		case appErrors.ErrCodeRequestAlreadyPending, appErrors.ErrCodeRequestAlreadyProcessed:
			httpStatus = http.StatusConflict
		case appErrors.ErrCodeInvalidUpdateToken, appErrors.ErrCodeUpdateTokenExpired:
			httpStatus = http.StatusUnauthorized
		default:
			httpStatus = http.StatusBadRequest
		}
	case err.Code >= 4000 && err.Code < 5000:
		// Password reset errors
		switch err.Code {
		case appErrors.ErrCodePasswordResetSpam:
			httpStatus = http.StatusTooManyRequests
		case appErrors.ErrCodeEmployeeNotFound:
			httpStatus = http.StatusNotFound
		default:
			httpStatus = http.StatusBadRequest
		}
	default:
		httpStatus = http.StatusInternalServerError
	}

	message := err.Message
	if err.Details != "" {
		message = message + ": " + err.Details
	}

	c.JSON(httpStatus, Response{
		Success: false,
		Code:    err.Code,
		Message: message,
	})
}
