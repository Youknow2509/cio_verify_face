package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	serviceImpl "github.com/youknow2509/cio_verify_face/server/service_auth/internal/application/service/impl"
	"github.com/youknow2509/cio_verify_face/server/service_auth/internal/global"
	utilsContext "github.com/youknow2509/cio_verify_face/server/service_auth/internal/shared/utils/context"
)

// OptimizedAuthMiddleware struct for cached authentication
type OptimizedAuthMiddleware struct {
	authCacheService *serviceImpl.AuthCacheService
}

// NewOptimizedAuthMiddleware creates a new optimized auth middleware
func NewOptimizedAuthMiddleware() (*OptimizedAuthMiddleware, error) {
	authCacheService, err := serviceImpl.NewAuthCacheService()
	if err != nil {
		return nil, err
	}

	// Type assertion to get the concrete implementation
	authCacheImpl, ok := authCacheService.(*serviceImpl.AuthCacheService)
	if !ok {
		return nil, fmt.Errorf("unexpected auth cache service type")
	}

	return &OptimizedAuthMiddleware{
		authCacheService: authCacheImpl,
	}, nil
}

// Apply method for optimized JWT access token processing with multi-level caching
func (m *OptimizedAuthMiddleware) Apply() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the JWT token from the Authorization header
		tokenStr, ok := utilsContext.ExtractBearerToken(c)
		if !ok || tokenStr == "" {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Validate token using cached service
		result, err := m.authCacheService.ValidateAccessToken(c.Request.Context(), tokenStr)
		if err != nil {
			global.Logger.Error("Error validating access token", "error", err)
			c.JSON(500, gin.H{"error": "Internal Server Error"})
			c.Abort()
			return
		}

		if !result.Valid {
			c.JSON(401, gin.H{"error": result.Error})
			c.Abort()
			return
		}

		// Set user information in context for downstream handlers
		c.Set("user_id", result.UserID)
		c.Set("user_role", result.Role)

		// Continue to the next handler
		c.Next()
	}
}

// ApplyWithRoleCheck creates middleware with role checking
func (m *OptimizedAuthMiddleware) ApplyWithRoleCheck(allowedRoles ...int) gin.HandlerFunc {
	return func(c *gin.Context) {
		// First validate the token
		tokenStr, ok := utilsContext.ExtractBearerToken(c)
		if !ok || tokenStr == "" {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		result, err := m.authCacheService.ValidateAccessToken(c.Request.Context(), tokenStr)
		if err != nil {
			global.Logger.Error("Error validating access token", "error", err)
			c.JSON(500, gin.H{"error": "Internal Server Error"})
			c.Abort()
			return
		}

		if !result.Valid {
			c.JSON(401, gin.H{"error": result.Error})
			c.Abort()
			return
		}

		// Check if user role is allowed
		roleAllowed := false
		for _, role := range allowedRoles {
			if role == result.Role {
				roleAllowed = true
				break
			}
		}

		if !roleAllowed {
			c.JSON(403, gin.H{"error": "Forbidden: insufficient permissions"})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", result.UserID)
		c.Set("user_role", result.Role)

		c.Next()
	}
}

// ApplyWithCompanyPermission creates middleware with company permission checking
func (m *OptimizedAuthMiddleware) ApplyWithCompanyPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		// First validate the token
		tokenStr, ok := utilsContext.ExtractBearerToken(c)
		if !ok || tokenStr == "" {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		result, err := m.authCacheService.ValidateAccessToken(c.Request.Context(), tokenStr)
		if err != nil {
			global.Logger.Error("Error validating access token", "error", err)
			c.JSON(500, gin.H{"error": "Internal Server Error"})
			c.Abort()
			return
		}

		if !result.Valid {
			c.JSON(401, gin.H{"error": result.Error})
			c.Abort()
			return
		}

		// Get company ID from request (could be from path param, query, or body)
		companyID := c.Param("company_id")
		if companyID == "" {
			companyID = c.Query("company_id")
		}

		if companyID == "" {
			c.JSON(400, gin.H{"error": "Company ID is required"})
			c.Abort()
			return
		}

		// Check company permission with caching
		hasPermission, err := m.authCacheService.CheckUserPermissionCached(
			c.Request.Context(),
			companyID,
			result.UserID.String(),
		)

		if err != nil {
			global.Logger.Error("Error checking company permission", "error", err)
			c.JSON(500, gin.H{"error": "Internal Server Error"})
			c.Abort()
			return
		}

		if !hasPermission {
			c.JSON(403, gin.H{"error": "Forbidden: no permission for this company"})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", result.UserID)
		c.Set("user_role", result.Role)
		c.Set("company_id", companyID)

		c.Next()
	}
}

// PreloadUserDataMiddleware preloads user data for performance
func (m *OptimizedAuthMiddleware) PreloadUserDataMiddleware(userIDs []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Preload user data asynchronously
		go func() {
			if err := m.authCacheService.PreloadUserData(c.Request.Context(), userIDs); err != nil {
				global.Logger.Warn("Failed to preload user data", "error", err)
			}
		}()

		c.Next()
	}
}
