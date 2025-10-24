package router

import (
	"github.com/gin-gonic/gin"
	infraMiddleware "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/infrastructure/middleware"
	handler "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/interfaces/http/handler"
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
	shiftRouterV1 := group.Group("/shift")
	shiftRouterV1.Use(infraMiddleware.GetAuthAccessTokenJwtMiddleware().Apply())
	{
		shiftRouterV1.GET("", handler.NewHandler().CreateShiftUser)        // Lay ca lam viec
		shiftRouterV1.POST("", handler.NewHandler().CreateShiftUser)       // Tao ca lam viec
		shiftRouterV1.GET("/:id", handler.NewHandler().GetDetailShiftUser) // Xem chi tiet thong tin ca lam viec
		shiftRouterV1.POST("/edit", handler.NewHandler().EditShiftUser)    // Chinh sua ca lam viec
		shiftRouterV1.DELETE("/:id", handler.NewHandler().DeleteShiftUser) // Xoa ca lam viec
	}

	scheduleRouterV1 := group.Group("/schedule")
	scheduleRouterV1.Use(infraMiddleware.GetAuthAccessTokenJwtMiddleware().Apply())
	{
		scheduleRouterV1.GET("", handler.NewHandler().GetInfoScheduleUser)       // Lay lich lam viec
		scheduleRouterV1.POST("", handler.NewHandler().CreateScheduleUser)       // Phan ca/ lich lam viec
		scheduleRouterV1.GET("/:id", handler.NewHandler().GetDetailScheduleUser) // Xem chi tiet lich lam viec
		scheduleRouterV1.POST("/edit", handler.NewHandler().EditScheduleUser)    // Chinh sua lich lam viec
		scheduleRouterV1.DELETE("/:id", handler.NewHandler().DeleteScheduleUser) // Xoa lich lam viec
	}
}
