package router

import (
	"github.com/gin-gonic/gin"
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
	}
	routerV1Private := g.Group("/v1/auth")
	{
		// Logout
		routerV1Private.POST("/logout", handler.GetAuthBaseHandler().Logout)
		// Refresh token
		routerV1Private.POST("/refresh", handler.GetAuthBaseHandler().RefreshToken)
		// Get my info
		routerV1Private.GET("/me", handler.GetAuthBaseHandler().GetMyInfo)
		// Create device token
		routerV1Private.POST("/device", handler.GetAuthBaseHandler().CreateDevice)
		// Delete device token
		routerV1Private.DELETE("/device/:id", handler.GetAuthBaseHandler().DeleteDevice)
		// Refresh device token
		routerV1Private.POST("/device/:id/refresh", handler.GetAuthBaseHandler().RefreshTokenDevice)
	}
}


/**
 * Initialize Auth Routers
 */
func (r *AuthRouter) InitializeAuthRoutes(g *gin.RouterGroup) {
	// Core Auth
	r.InitializeCoreAuth(g)
}
