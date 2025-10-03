package context

import (
	"strings"

	"github.com/gin-gonic/gin"
	constants "github.com/youknow2509/cio_verify_face/server/service_auth/internal/constants"
)

// Extract the token from the header
func ExtractBearerToken(c *gin.Context) (string, bool) {
	// Authorization: Bearer token
	authHeader := c.GetHeader(constants.ContextAuthHeaderKey)
	if strings.HasPrefix(authHeader, constants.ContextBearerPrefixKey) {
		return strings.TrimPrefix(authHeader, constants.ContextBearerPrefixKey), true
	}
	return "", false
}

// Save session to the context
func SaveSessionToContext(c *gin.Context, userId string, sessionId string) {
	c.Set(constants.ContextUserIDKey, userId)
	c.Set(constants.ContextSessionIDKey, sessionId)
}
