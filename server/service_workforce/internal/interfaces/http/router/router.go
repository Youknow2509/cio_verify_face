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

	shiftEmployeeRouterV1 := group.Group("/employee/shift")
	shiftEmployeeRouterV1.Use(infraMiddleware.GetAuthAccessTokenJwtMiddleware().Apply())
	{
		shiftEmployeeRouterV1.POST("/user", nil)                // Get shift for user with effective date
		shiftEmployeeRouterV1.POST("/user/edit/effective", nil) // Edit shift for user with effective date
		shiftEmployeeRouterV1.POST("/user/enable", nil)         // Enable shift for user
		shiftEmployeeRouterV1.POST("/user/disable", nil)        // Disable shift for user
		shiftEmployeeRouterV1.DELETE("", nil) 					// Delete shift for user
		shiftEmployeeRouterV1.GET("/user/add", nil)             // Add shift employee
	}
}
