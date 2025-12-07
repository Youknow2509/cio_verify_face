package router

import (
	"github.com/gin-gonic/gin"
	infraMiddleware "github.com/youknow2509/cio_verify_face/server/service_device/internal/infrastructure/middleware"
	"github.com/youknow2509/cio_verify_face/server/service_device/internal/interfaces/http/handler"
)

/**
 * Http router manager
 */
type HttpRouterManager struct {
}

/**
 * Initialize routes
 */
func (r *HttpRouterManager) InitRoutes(group *gin.RouterGroup) {
	deviceV1 := group.Group("/v1/device")
	deviceV1.Use(infraMiddleware.GetAuthAdminAccessTokenJwtMiddleware().Apply())
	{
		deviceV1.GET("", handler.NewHandler().GetListDevices)
		deviceV1.POST("", handler.NewHandler().CreateNewDevice)
		deviceV1.GET("/:device_id", handler.NewHandler().GetDeviceById)
		deviceV1.GET("/token/:device_id", handler.NewHandler().GetDeviceToken)
		deviceV1.POST("/token/refresh/:device_id", handler.NewHandler().RefreshDeviceToken)
		deviceV1.PUT("/:device_id", handler.NewHandler().UpdateDeviceById)
		deviceV1.DELETE("/:device_id", handler.NewHandler().DeleteDeviceById)
		deviceV1.POST("/location", handler.NewHandler().UpdateLocationDevice)
		deviceV1.POST("/name", handler.NewHandler().UpdateNameDevice)
		deviceV1.POST("/info", handler.NewHandler().UpdateInfoDevice)
		deviceV1.POST("/status", handler.NewHandler().UpdateStatusDevice)
	}
	deviceSelf := group.Group("/v1/device")
	deviceSelf.Use(infraMiddleware.GetAuthDeviceAccessTokenJwtMiddleware().Apply())
	{
		deviceSelf.GET("/me", handler.NewHandler().GetInfoDevice)
	}
}
