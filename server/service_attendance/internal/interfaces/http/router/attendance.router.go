package router

import (
	"github.com/gin-gonic/gin"
	"github.com/youknow2509/cio_verify_face/server/service_attendance/internal/infrastructure/middleware"
	httpHandler "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/interfaces/http/handler"
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
		// Add attendance record
		v1Admin.POST("/", httpHandler.NewAttendanceHandler().AddAttendance)
		// Get attendance records for company
		v1Admin.POST("/records", httpHandler.NewAttendanceHandler().GetAttendanceRecords)
		// Get daily attendance summary for company
		v1Admin.POST("/records/summary/daily", httpHandler.NewAttendanceHandler().GetDailyAttendanceSummary)
	}
	//
	v1User := g.Group("/v1/attendance")
	v1User.Use(middleware.GetAuthAccessTokenJwtMiddleware().Apply())
	{
		// Get attendance records for employee
		v1User.POST("/records/employee", httpHandler.NewAttendanceHandler().GetAttendanceRecordsEmployee)
		// Get daily attendance summary for employee
		v1User.POST("/records/employee/summary/daily", httpHandler.NewAttendanceHandler().GetDailyAttendanceSummaryEmployee)
	}
}
