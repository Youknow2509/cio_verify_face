package router

import (
	"github.com/gin-gonic/gin"
	"github.com/youknow2509/cio_verify_face/server/service_attendance/internal/infrastructure/middleware"
	httpHandler "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/interfaces/http/handler/attendance"
)

// ============================================
// Attendance router
// ============================================
type AttendanceRouter struct{}

// Deploy attendance routes
func (r *AttendanceRouter) Deploy(g *gin.RouterGroup) {
	// group api router v1
	v1Admin := g.Group("/v1/attendance")
	v1Admin.Use(middleware.GetAuthAdminMiddleware().Apply())
	{
		v1Admin.POST("/check_in", httpHandler.NewAttendanceHandler().CheckIn)
		v1Admin.POST("/check_out", httpHandler.NewAttendanceHandler().CheckOut)
		v1Admin.POST("/records", httpHandler.NewAttendanceHandler().GetRecords)
	}
	v1User := g.Group("/v1/attendance")
	v1User.Use(middleware.GetAuthAccessTokenJwtMiddleware().Apply())
	{
		v1User.GET("/records/:record_id", httpHandler.NewAttendanceHandler().GetRecordByID)
		v1User.POST("/history/my", httpHandler.NewAttendanceHandler().GetMyHistory)
	}
}
