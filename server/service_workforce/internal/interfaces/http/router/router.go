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
	shiftRouterV1 := group.Group("/v1/shift")
	shiftRouterV1.Use(infraMiddleware.GetAuthAdminAccessTokenJwtMiddleware().Apply())
	{
		shiftRouterV1.GET("", handler.NewHandler().GetListShift) // Lay danh sach ca lam viec
		shiftRouterV1.POST("", handler.NewHandler().CreateShift)              // Tao ca lam viec
		shiftRouterV1.GET("/:id", handler.NewHandler().GetDetailShift)        // Xem chi tiet thong tin ca lam viec
		shiftRouterV1.POST("/edit", handler.NewHandler().EditShift)           // Chinh sua ca lam viec
		shiftRouterV1.DELETE("/:id", handler.NewHandler().DeleteShift)        // Xoa ca lam viec
		shiftRouterV1.POST("/status", handler.NewHandler().ChangeStatusShift) // Thay doi trang thai ca lam viec
	}
	shiftRouterV1Employee := group.Group("/v1/shift")
	shiftRouterV1Employee.Use(infraMiddleware.GetAuthAccessTokenJwtMiddleware().Apply())
	{
		shiftRouterV1Employee.GET("/employee", handler.NewHandler().GetListShiftEmployee) // Lay danh sach ca lam viec nhan vien
	}
	shiftEmployeeRouterV1 := group.Group("/v1/employee/shift")
	shiftEmployeeRouterV1.Use(infraMiddleware.GetAuthAdminAccessTokenJwtMiddleware().Apply())
	{
		shiftEmployeeRouterV1.POST("", handler.NewHandler().GetShiftUserWithEffectiveDate)                 // Get shift for user with effective date
		shiftEmployeeRouterV1.POST("/edit/effective", handler.NewHandler().EditShiftUserWithEffectiveDate) // Edit shift for user with effective date
		shiftEmployeeRouterV1.POST("/enable", handler.NewHandler().EnableShiftUser)                        // Enable shift for user
		shiftEmployeeRouterV1.POST("/disable", handler.NewHandler().DisableShiftUser)                      // Disable shift for user
		shiftEmployeeRouterV1.POST("/delete", handler.NewHandler().DeleteShiftUser)                        // Delete shift for user
		shiftEmployeeRouterV1.POST("/add", handler.NewHandler().AddShiftEmployee)                          // Add shift employee
		shiftEmployeeRouterV1.POST("/add/list", handler.NewHandler().AddShiftEmployeeList)                 // Add shift employee list
		shiftEmployeeRouterV1.POST("/not_in", handler.NewHandler().GetInfoEmployeeDonotInShift)            // Get info employee donot in shift
		shiftEmployeeRouterV1.POST("/in", handler.NewHandler().GetInfoEmployeeInShift)                     // Get info employee in shift
	}
}
