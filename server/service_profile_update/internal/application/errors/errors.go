package errors

import (
	"fmt"
)

// Error represents a structured application error
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Error implements the error interface
func (e *Error) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%d] %s: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// WithDetails returns a copy of the error with additional details
func (e *Error) WithDetails(details string) *Error {
	return &Error{
		Code:    e.Code,
		Message: e.Message,
		Details: details,
	}
}

// =================================
// Error Code Definitions:
// =================================
const (
	// General errors (1xxx)
	ErrCodeUnknown             = 1000
	ErrCodeInvalidInput        = 1001
	ErrCodeUnauthorized        = 1002
	ErrCodeForbidden           = 1003
	ErrCodeNotFound            = 1004
	ErrCodeConflict            = 1005
	ErrCodeServiceUnavailable  = 1006

	// Rate limiting errors (2xxx)
	ErrCodeRateLimitExceeded   = 2001
	ErrCodeSpamDetected        = 2002
	ErrCodeMonthlyLimitReached = 2003

	// Face profile update errors (3xxx)
	ErrCodeRequestAlreadyPending   = 3001
	ErrCodeRequestNotFound         = 3002
	ErrCodeRequestAlreadyProcessed = 3003
	ErrCodeInvalidUpdateToken      = 3004
	ErrCodeUpdateTokenExpired      = 3005
	ErrCodeFaceEnrollmentFailed    = 3006

	// Password reset errors (4xxx)
	ErrCodePasswordResetSpam      = 4001
	ErrCodeEmployeeNotFound       = 4002
	ErrCodePasswordResetFailed    = 4003
	ErrCodeKafkaPublishFailed     = 4004
)

// =================================
// Predefined Errors:
// =================================
var (
	// General errors
	ErrUnknown            = &Error{Code: ErrCodeUnknown, Message: "An unknown error occurred"}
	ErrInvalidInput       = &Error{Code: ErrCodeInvalidInput, Message: "Invalid input provided"}
	ErrUnauthorized       = &Error{Code: ErrCodeUnauthorized, Message: "Authentication required"}
	ErrForbidden          = &Error{Code: ErrCodeForbidden, Message: "Access denied"}
	ErrNotFound           = &Error{Code: ErrCodeNotFound, Message: "Resource not found"}
	ErrConflict           = &Error{Code: ErrCodeConflict, Message: "Resource conflict"}
	ErrServiceUnavailable = &Error{Code: ErrCodeServiceUnavailable, Message: "Service temporarily unavailable"}

	// Rate limiting errors
	ErrRateLimitExceeded   = &Error{Code: ErrCodeRateLimitExceeded, Message: "Rate limit exceeded"}
	ErrSpamDetected        = &Error{Code: ErrCodeSpamDetected, Message: "Spam detected, please wait before retrying"}
	ErrMonthlyLimitReached = &Error{Code: ErrCodeMonthlyLimitReached, Message: "Monthly request limit reached"}

	// Face profile update errors
	ErrRequestAlreadyPending   = &Error{Code: ErrCodeRequestAlreadyPending, Message: "You already have a pending request"}
	ErrRequestNotFound         = &Error{Code: ErrCodeRequestNotFound, Message: "Request not found"}
	ErrRequestAlreadyProcessed = &Error{Code: ErrCodeRequestAlreadyProcessed, Message: "Request has already been processed"}
	ErrInvalidUpdateToken      = &Error{Code: ErrCodeInvalidUpdateToken, Message: "Invalid update token"}
	ErrUpdateTokenExpired      = &Error{Code: ErrCodeUpdateTokenExpired, Message: "Update token has expired"}
	ErrFaceEnrollmentFailed    = &Error{Code: ErrCodeFaceEnrollmentFailed, Message: "Face enrollment failed"}

	// Password reset errors
	ErrPasswordResetSpam   = &Error{Code: ErrCodePasswordResetSpam, Message: "Too many password reset requests"}
	ErrEmployeeNotFound    = &Error{Code: ErrCodeEmployeeNotFound, Message: "Employee not found"}
	ErrPasswordResetFailed = &Error{Code: ErrCodePasswordResetFailed, Message: "Password reset failed"}
	ErrKafkaPublishFailed  = &Error{Code: ErrCodeKafkaPublishFailed, Message: "Failed to send notification"}
)
