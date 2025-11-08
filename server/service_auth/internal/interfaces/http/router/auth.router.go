package router

import (

	"github.com/gin-gonic/gin"
	infraMiddleware "github.com/youknow2509/cio_verify_face/server/service_auth/internal/infrastructure/middleware"
	"github.com/youknow2509/cio_verify_face/server/service_auth/internal/interfaces/http/handler"
)

/**
 * Auth router
 */
type AuthRouter struct {
}

/**
 * Core Auth
 */
func (r *AuthRouter) InitializeCoreAuth(g *gin.RouterGroup) {
	routerV1Public := g.Group("/v1/auth")
	{
		// Login
		routerV1Public.POST("/login", handler.GetAuthBaseHandler().Login)
		// Login admin
		routerV1Public.POST("/login/admin", handler.GetAuthBaseHandler().LoginAdmin)
		// Refresh token
		routerV1Public.POST("/refresh", handler.GetAuthBaseHandler().RefreshToken)
	}
	routerV1Private := g.Group("/v1/auth")
	routerV1Private.Use(infraMiddleware.GetAuthAccessTokenJwtMiddleware().Apply())
	{
		// Logout
		routerV1Private.POST("/logout", handler.GetAuthBaseHandler().Logout)
		// Get my info
		routerV1Private.GET("/me", handler.GetAuthBaseHandler().GetMyInfo)
		// Update device token
		routerV1Private.POST("/device", handler.GetAuthBaseHandler().UpdateDeviceSession)
		// Delete device token
		routerV1Private.DELETE("/device", handler.GetAuthBaseHandler().DeleteDeviceSession)
	}
}

/**
 * Initialize Auth Routers
 */
func (r *AuthRouter) InitializeAuthRoutes(g *gin.RouterGroup) {
	// Core Auth
	r.InitializeCoreAuth(g)
}
