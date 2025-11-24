package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/youknow2509/cio_verify_face/server/service_analytic/internal/interfaces/dto"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateStruct validates a struct using validator/v10
func ValidateStruct(data interface{}) error {
	return validate.Struct(data)
}

// ValidateAndBind binds JSON to DTO and validates it
func ValidateAndBind(c *gin.Context, data interface{}) error {
	// Bind JSON to DTO
	if err := c.ShouldBindJSON(data); err != nil {
		return fmt.Errorf("invalid JSON format: %w", err)
	}

	// Validate the DTO
	if err := validate.Struct(data); err != nil {
		return formatValidationError(err)
	}

	return nil
}

// BindAndValidateMiddleware is a middleware that binds and validates request body
func BindAndValidateMiddleware(dtoType interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := ValidateAndBind(c, dtoType); err != nil {
			c.JSON(http.StatusBadRequest, dto.NewErrorResponse("VALIDATION_ERROR", "Invalid request data", err.Error()))
			c.Abort()
			return
		}
		c.Next()
	}
}

// formatValidationError formats validation errors into a readable message
func formatValidationError(err error) error {
	var errors []string
	
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			var errorMsg string
			switch e.Tag() {
			case "required":
				errorMsg = fmt.Sprintf("field '%s' is required", e.Field())
			case "uuid":
				errorMsg = fmt.Sprintf("field '%s' must be a valid UUID", e.Field())
			case "email":
				errorMsg = fmt.Sprintf("field '%s' must be a valid email address", e.Field())
			case "oneof":
				errorMsg = fmt.Sprintf("field '%s' must be one of: %s", e.Field(), e.Param())
			case "min":
				errorMsg = fmt.Sprintf("field '%s' must be at least %s", e.Field(), e.Param())
			case "max":
				errorMsg = fmt.Sprintf("field '%s' must be at most %s", e.Field(), e.Param())
			case "gte":
				errorMsg = fmt.Sprintf("field '%s' must be greater than or equal to %s", e.Field(), e.Param())
			case "lte":
				errorMsg = fmt.Sprintf("field '%s' must be less than or equal to %s", e.Field(), e.Param())
			case "gt":
				errorMsg = fmt.Sprintf("field '%s' must be greater than %s", e.Field(), e.Param())
			case "lt":
				errorMsg = fmt.Sprintf("field '%s' must be less than %s", e.Field(), e.Param())
			default:
				errorMsg = fmt.Sprintf("field '%s' failed validation on '%s' tag", e.Field(), e.Tag())
			}
			errors = append(errors, errorMsg)
		}
		return fmt.Errorf("%s", strings.Join(errors, "; "))
	}
	
	return err
}

// GetValidator returns the global validator instance
func GetValidator() *validator.Validate {
	return validate
}
