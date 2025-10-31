package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/youknow2509/cio_verify_face/server/service_auth/internal/application/service"
	serviceImpl "github.com/youknow2509/cio_verify_face/server/service_auth/internal/application/service/impl"
	"github.com/youknow2509/cio_verify_face/server/service_auth/internal/global"
)

// OptimizedUserHandler demonstrates optimized user operations with caching
type OptimizedUserHandler struct {
	authCacheService service.IAuthCacheService
}

// NewOptimizedUserHandler creates a new optimized user handler
func NewOptimizedUserHandler() (*OptimizedUserHandler, error) {
	authCacheService, err := serviceImpl.NewAuthCacheService()
	if err != nil {
		return nil, err
	}

	return &OptimizedUserHandler{
		authCacheService: authCacheService,
	}, nil
}

// GetUserProfile gets user profile with caching optimization
func (h *OptimizedUserHandler) GetUserProfile(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Get user info using cached service
	userInfo, err := h.authCacheService.GetUserInfoCached(c.Request.Context(), userID.String())
	if err != nil {
		global.Logger.Error("Error getting user info", "error", err, "userID", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user information"})
		return
	}

	if userInfo == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    userInfo,
		"message": "User profile retrieved successfully",
	})
}

// CheckCompanyAccess checks if user has access to company with caching
func (h *OptimizedUserHandler) CheckCompanyAccess(c *gin.Context) {
	// Get user ID from context
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// Get company ID from path parameter
	companyID := c.Param("company_id")
	if companyID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Company ID is required"})
		return
	}

	// Check permission using cached service
	hasPermission, err := h.authCacheService.CheckUserPermissionCached(
		c.Request.Context(),
		companyID,
		userID.String(),
	)

	if err != nil {
		global.Logger.Error("Error checking company permission", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check permission"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"has_permission": hasPermission,
		"message":        "Permission check completed",
	})
}

// CheckDeviceAccess checks if device exists in company with caching
func (h *OptimizedUserHandler) CheckDeviceAccess(c *gin.Context) {
	// Get company ID and device ID from path parameters
	companyID := c.Param("company_id")
	deviceID := c.Param("device_id")

	if companyID == "" || deviceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Company ID and Device ID are required"})
		return
	}

	// Check device existence using cached service
	deviceExists, err := h.authCacheService.CheckDeviceInCompanyCached(
		c.Request.Context(),
		companyID,
		deviceID,
	)

	if err != nil {
		global.Logger.Error("Error checking device existence", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check device"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"device_exists": deviceExists,
		"message":       "Device check completed",
	})
}

// GetCacheStatistics returns cache performance statistics
func (h *OptimizedUserHandler) GetCacheStatistics(c *gin.Context) {
	stats, err := h.authCacheService.GetCacheStats(c.Request.Context())
	if err != nil {
		global.Logger.Error("Error getting cache stats", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cache statistics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    stats,
		"message": "Cache statistics retrieved successfully",
	})
}

// InvalidateUserCache manually invalidates user cache
func (h *OptimizedUserHandler) InvalidateUserCache(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// Invalidate user cache
	h.authCacheService.InvalidateUserCache(c.Request.Context(), userID)

	c.JSON(http.StatusOK, gin.H{
		"message": "User cache invalidated successfully",
	})
}

// WarmupCache pre-loads frequently accessed data into cache
func (h *OptimizedUserHandler) WarmupCache(c *gin.Context) {
	err := h.authCacheService.WarmupCache(c.Request.Context())
	if err != nil {
		global.Logger.Error("Error warming up cache", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to warmup cache"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cache warmup completed successfully",
	})
}
