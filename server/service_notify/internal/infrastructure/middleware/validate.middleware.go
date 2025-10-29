package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/youknow2509/cio_verify_face/server/service_notify/internal/constants"
)

/**
 * Validate middleware struct for handling validation logic.
 */
type ValidateMiddleware struct {
}

/**
 * Apply method to process the request and perform validation.
 */
func (vm *ValidateMiddleware) Apply() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Create a new validator instance
		validate := validator.New()
		// Set instance to the context
		ctx.Set(constants.MIDDLEWARE_VALIDATE_SERVICE_NAME, validate)
		ctx.Next()
	}
}
