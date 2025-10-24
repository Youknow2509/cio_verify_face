package context

import (
	"strings"

	"github.com/gin-gonic/gin"
	constants "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/constants"
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
func SaveSessionToContext(c *gin.Context, userId string, sessionId string, userRole int) {
	c.Set(constants.ContextUserIDKey, userId)
	c.Set(constants.ContextSessionIDKey, sessionId)
	c.Set(constants.ContextUserRoleKey, userRole)
}

// Get session from the context
func GetSessionFromContext(c *gin.Context) (string, string, int, bool) {
	userId, userIdExists := c.Get(constants.ContextUserIDKey)
	sessionId, sessionIdExists := c.Get(constants.ContextSessionIDKey)
	userRole, userRoleExists := c.Get(constants.ContextUserRoleKey)
	if !userIdExists || !sessionIdExists || !userRoleExists {
		return "", "", -1, false
	}
	return userId.(string), sessionId.(string), userRole.(int), true
}
