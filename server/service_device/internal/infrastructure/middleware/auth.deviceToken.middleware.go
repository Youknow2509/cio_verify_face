package middleware

import (
	"github.com/gin-gonic/gin"
	// domainCache "github.com/youknow2509/cio_verify_face/server/service_device/internal/domain/cache"
	// domainModel "github.com/youknow2509/cio_verify_face/server/service_device/internal/domain/model"
	domainToken "github.com/youknow2509/cio_verify_face/server/service_device/internal/domain/token"
	// utilsCache "github.com/youknow2509/cio_verify_face/server/service_device/internal/shared/utils/cache"
	utilsContext "github.com/youknow2509/cio_verify_face/server/service_device/internal/shared/utils/context"
	// utilsCrypto "github.com/youknow2509/cio_verify_face/server/service_device/internal/shared/utils/crypto"
)

// Define the AuthDeviceKTokenJwtMiddleware struct
type AuthDeviceTokenJwtMiddleware struct{}

/**
 * Apply method to process the JWT access token.
 */
func (m *AuthDeviceTokenJwtMiddleware) Apply() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Instance use
		var (
			tokenService domainToken.ITokenService
		)
		tokenService = domainToken.GetTokenService()
		var (
			tokenStr string
		)
		// Extract the JWT token from the Authorization header
		tokenStr, ok := utilsContext.ExtractBearerToken(c)
		if !ok || tokenStr == "" {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		// Validate the JWT token
		tokenInfo, errToken := tokenService.ParseDeviceToken(
			c,
			tokenStr,
		)
		if errToken != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		if tokenInfo == nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		// Store the token information in the context for further use
		utilsContext.SaveDeviceSessionToContext(c, tokenInfo.DeviceId, tokenInfo.CompanyId)
		// If the token is valid, proceed to the next middleware/handler
		c.Next()
	}
}
