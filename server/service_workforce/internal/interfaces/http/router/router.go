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
	shiftRouterV1.Use(infraMiddleware.GetAuthAdminAccessTokenJwtMiddleware().Apply())
	{
		shiftRouterV1.POST("", handler.NewHandler().CreateShift)       // Tao ca lam viec
		shiftRouterV1.GET("/:id", handler.NewHandler().GetDetailShift) // Xem chi tiet thong tin ca lam viec
		shiftRouterV1.POST("/edit", handler.NewHandler().EditShift)    // Chinh sua ca lam viec
		shiftRouterV1.DELETE("/:id", handler.NewHandler().DeleteShift) // Xoa ca lam viec
	}

	shiftEmployeeRouterV1 := group.Group("/employee/shift")
	shiftEmployeeRouterV1.Use(infraMiddleware.GetAuthAdminAccessTokenJwtMiddleware().Apply())
	{
		shiftEmployeeRouterV1.POST("", handler.NewHandler().GetShiftForUserWithEffectiveDate)                 // Get shift for user with effective date
		shiftEmployeeRouterV1.POST("/edit/effective", handler.NewHandler().EditShiftForUserWithEffectiveDate) // Edit shift for user with effective date
		shiftEmployeeRouterV1.POST("/enable", handler.NewHandler().EnableShiftForUser)                        // Enable shift for user
		shiftEmployeeRouterV1.POST("/disable", handler.NewHandler().DisableShiftForUser)                      // Disable shift for user
		shiftEmployeeRouterV1.DELETE("/:id", handler.NewHandler().DeleteShiftForUser)                             // Delete shift for user
		shiftEmployeeRouterV1.POST("/add", handler.NewHandler().AddShiftEmployee)                              // Add shift employee
	}
}
