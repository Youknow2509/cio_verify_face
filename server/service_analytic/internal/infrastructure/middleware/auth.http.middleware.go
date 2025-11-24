package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	applicationModel "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/application/model"
	"github.com/youknow2509/cio_verify_face/server/service_analytic/internal/global"
	authClient "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/infrastructure/grpc"
)

// HTTPAuthMiddleware validates JWT tokens via service_auth and sets session info in context
type HTTPAuthMiddleware struct {
	authClient *authClient.AuthServiceClient
}

// NewHTTPAuthMiddleware creates a new HTTP authentication middleware
func NewHTTPAuthMiddleware() *HTTPAuthMiddleware {
	client, err := authClient.NewAuthServiceClient()
	if err != nil {
		global.Logger.Error("Failed to create auth service client", "error", err)
		return nil
	}
	return &HTTPAuthMiddleware{
		authClient: client,
	}
}

// Apply returns a gin middleware handler that validates JWT tokens
func (m *HTTPAuthMiddleware) Apply() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		token, ok := extractBearerToken(c)
		if !ok || token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "UNAUTHORIZED",
				"message": "Missing or invalid authorization token",
			})
			c.Abort()
			return
		}
		// Validate token via service_auth gRPC
		sessionInfo, err := m.authClient.ValidateToken(c.Request.Context(), token)
		if err != nil {
			global.Logger.Warn("Token validation failed", "error", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "UNAUTHORIZED",
				"message": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Enrich session info with request metadata
		sessionInfo.ClientIP = c.ClientIP()
		sessionInfo.ClientAgent = c.Request.UserAgent()

		// Set session info in context (as a single object)
		c.Set("session", sessionInfo)

		// Also set individual fields for backward compatibility
		c.Set("user_id", sessionInfo.UserID)
		c.Set("role", sessionInfo.Role)
		c.Set("session_id", sessionInfo.SessionID)
		c.Set("company_id", sessionInfo.CompanyID)
		// Continue to next handler
		c.Next()
	}
}

// extractBearerToken extracts JWT token from Authorization header
func extractBearerToken(c *gin.Context) (string, bool) {
	authHeader := c.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer "), true
	}
	return "", false
}

// GetSessionFromContext retrieves session info from gin context
func GetSessionFromContext(c *gin.Context) (*applicationModel.SessionInfo, bool) {
	if sessionData, exists := c.Get("session"); exists {
		if session, ok := sessionData.(*applicationModel.SessionInfo); ok {
			return session, true
		}
	}
	return nil, false
}
