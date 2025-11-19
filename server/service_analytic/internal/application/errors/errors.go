package errors

import (
	"fmt"
	"net/http"
)

// Error represents application error
type Error struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
	Details    string `json:"details,omitempty"`
}

func (e *Error) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewError creates a new application error
func NewError(code, message string, statusCode int) *Error {
	return &Error{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

// WithDetails adds details to the error
func (e *Error) WithDetails(details string) *Error {
	e.Details = details
	return e
}

// Common errors
var (
	ErrInvalidInput       = NewError("INVALID_INPUT", "Invalid input parameters", http.StatusBadRequest)
	ErrNotFound           = NewError("NOT_FOUND", "Resource not found", http.StatusNotFound)
	ErrUnauthorized       = NewError("UNAUTHORIZED", "Unauthorized access", http.StatusUnauthorized)
	ErrForbidden          = NewError("FORBIDDEN", "Access forbidden", http.StatusForbidden)
	ErrInternalServer     = NewError("INTERNAL_ERROR", "Internal server error", http.StatusInternalServerError)
	ErrDatabaseError      = NewError("DATABASE_ERROR", "Database operation failed", http.StatusInternalServerError)
	ErrInvalidDateFormat  = NewError("INVALID_DATE_FORMAT", "Invalid date format", http.StatusBadRequest)
	ErrInvalidDateRange   = NewError("INVALID_DATE_RANGE", "Invalid date range", http.StatusBadRequest)
	ErrExportFailed       = NewError("EXPORT_FAILED", "Failed to export report", http.StatusInternalServerError)
	ErrTooManyRequests    = NewError("TOO_MANY_REQUESTS", "Too many requests", http.StatusTooManyRequests)
)
