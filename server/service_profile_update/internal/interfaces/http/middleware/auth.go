package middleware

import (
	"context"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/application/model"
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/global"
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/interfaces/response"
)

// AuthMiddleware creates a middleware that validates JWT token via auth service
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if auth client is available
		if global.AuthClient == nil {
			response.InternalError(c, "Auth service unavailable")
			c.Abort()
			return
		}

		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Missing authorization header")
			c.Abort()
			return
		}

		// Token should be in format: "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			response.Unauthorized(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		token := parts[1]
		if token == "" {
			response.Unauthorized(c, "Missing token")
			c.Abort()
			return
		}

		// Parse token via auth service with timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		tokenInfo, err := global.AuthClient.ParseUserToken(ctx, token)
		if err != nil {
			global.Logger.Error("Failed to parse user token: " + err.Error())
			response.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// Create session info from token
		session := &model.SessionInfo{
			UserID:      tokenInfo.UserId,
			CompanyID:   tokenInfo.CompanyId,
			Role:        tokenInfo.Roles,
			SessionID:   tokenInfo.TokenId,
			ClientIP:    c.ClientIP(),
			ClientAgent: c.GetHeader("User-Agent"),
		}

		// Store session in context
		c.Set("session", session)

		// Continue to next handler
		c.Next()
	}
}

// OptionalAuthMiddleware is like AuthMiddleware but doesn't abort if token is missing
// Useful for endpoints that work both with and without authentication
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if auth client is available
		if global.AuthClient == nil {
			c.Next()
			return
		}

		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		// Token should be in format: "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.Next()
			return
		}

		token := parts[1]
		if token == "" {
			c.Next()
			return
		}

		// Parse token via auth service with timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		tokenInfo, err := global.AuthClient.ParseUserToken(ctx, token)
		if err != nil {
			// Log but don't block
			global.Logger.Warn("Optional auth failed: " + err.Error())
			c.Next()
			return
		}

		// Create session info from token
		session := &model.SessionInfo{
			UserID:      tokenInfo.UserId,
			CompanyID:   tokenInfo.CompanyId,
			Role:        tokenInfo.Roles,
			SessionID:   tokenInfo.TokenId,
			ClientIP:    c.ClientIP(),
			ClientAgent: c.GetHeader("User-Agent"),
		}

		// Store session in context
		c.Set("session", session)

		// Continue to next handler
		c.Next()
	}
}
