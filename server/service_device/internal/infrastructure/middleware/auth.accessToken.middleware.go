package middleware

import (
	"github.com/gin-gonic/gin"
	domainCache "github.com/youknow2509/cio_verify_face/server/service_device/internal/domain/cache"
	domainModel "github.com/youknow2509/cio_verify_face/server/service_device/internal/domain/model"
	domainToken "github.com/youknow2509/cio_verify_face/server/service_device/internal/domain/token"
	utilsCache "github.com/youknow2509/cio_verify_face/server/service_device/internal/shared/utils/cache"
	utilsContext "github.com/youknow2509/cio_verify_face/server/service_device/internal/shared/utils/context"
	utilsCrypto "github.com/youknow2509/cio_verify_face/server/service_device/internal/shared/utils/crypto"
)

// Define the AuthAccessTokenJwtMiddleware struct
type AuthAccessTokenJwtMiddleware struct{}

/**
 * Apply method to process the JWT access token.
 */
func (m *AuthAccessTokenJwtMiddleware) Apply() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Instance use
		var (
			tokenService     domainToken.ITokenService
			distributedCache domainCache.IDistributedCache
		)
		tokenService = domainToken.GetTokenService()
		distributedCache, err := domainCache.GetDistributedCache()
		if err != nil {
			c.JSON(500, gin.H{"error": "Internal Server Error"})
			c.Abort()
			return
		}
		var (
			tokenStr                    string
			tokenObj                    *domainModel.TokenUserJwtOutput
			keyCacheAccessTokenIsActive string
			cacheData                   string
		)
		// Extract the JWT token from the Authorization header
		tokenStr, ok := utilsContext.ExtractBearerToken(c)
		if !ok || tokenStr == "" {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		// Validate the JWT token
		tokenObj, errToken := tokenService.ParseUserToken(
			c,
			tokenStr,
		)
		if errToken != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		// Check token is blocked
		keyCacheAccessTokenIsActive = utilsCache.GetKeyUserAccessTokenIsActive(
			utilsCrypto.GetHash(tokenObj.TokenId),
		)
		cacheData, err = distributedCache.Get(
			c,
			keyCacheAccessTokenIsActive,
		)
		if err != nil {
			c.JSON(500, gin.H{"error": "Internal Server Error"})
			c.Abort()
			return
		}
		if cacheData == "" || cacheData == "0" {
			c.JSON(401, gin.H{"error": "Unauthorized - Token is blocked"})
			c.Abort()
			return
		}
		// Save session to the context
		utilsContext.SaveSessionToContext(
			c,
			tokenObj.UserId,
			tokenObj.TokenId,
			tokenObj.Role,
		)
		// If the token is valid, proceed to the next middleware/handler
		c.Next()
	}
}
