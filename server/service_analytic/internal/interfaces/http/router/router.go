package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/youknow2509/cio_verify_face/server/service_analytic/internal/docs"
	"github.com/youknow2509/cio_verify_face/server/service_analytic/internal/global"
	"github.com/youknow2509/cio_verify_face/server/service_analytic/internal/infrastructure/middleware"
	"github.com/youknow2509/cio_verify_face/server/service_analytic/internal/interfaces/http/handler"
	"github.com/youknow2509/cio_verify_face/server/pkg/observability"
)

// SetupRouter sets up the HTTP router with authentication middleware
func SetupRouter() *gin.Engine {
	// Create gin router
	router := gin.Default()
	if global.SettingServer.Server.Mode != "dev" {
		gin.SetMode(gin.ReleaseMode)
	}
	// Setup CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	router.Use(cors.New(config))

	// Add observability middleware if enabled
	if global.SettingServer.Observability.Enabled {
		if metrics := observability.GetHTTPMetrics(); metrics != nil {
			router.Use(metrics.GinMiddleware())
		}
	}

	// Swagger only in dev mode
	if global.SettingServer.Server.Mode == "dev" {
		// Ensure Swagger metadata aligns with our routing base
		docs.SwaggerInfo.BasePath = "/api/v1"
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Health check (no auth required)
	healthHandler := handler.NewHealthHandler()
	router.GET("/health", healthHandler.HealthCheck)

	// API v1 routes with authentication
	v1 := router.Group("/api/v1")

	// Apply authentication middleware if auth service is enabled
	if global.SettingServer.AuthService.Enabled {
		authMiddleware := middleware.NewHTTPAuthMiddleware()
		if authMiddleware != nil {
			v1.Use(authMiddleware.Apply())
		} else {
			global.Logger.Warn("HTTP auth middleware not initialized - routes will be unprotected")
		}
	} else {
		global.Logger.Warn("Auth service disabled - API routes will be unprotected")
	}

	{
		// Analytics/Reports routes (protected by auth middleware)
		analyticHandler := handler.NewAnalyticHandler()
		reports := v1.Group("/reports")
		{
			reports.GET("/daily", analyticHandler.GetDailyReport)
			reports.GET("/summary", analyticHandler.GetSummaryReport)
			reports.POST("/export", analyticHandler.ExportReport)
			reports.GET("/download/:filename", analyticHandler.DownloadExport)
		}

		// ScyllaDB data access routes (protected by auth middleware)
		scyllaHandler := handler.NewScyllaHandler()

		// Attendance Records routes
		attendanceRecords := v1.Group("/attendance-records")
		{
			attendanceRecords.GET("", scyllaHandler.GetAttendanceRecords)
			attendanceRecords.GET("/range", scyllaHandler.GetAttendanceRecordsByTimeRange)
			attendanceRecords.GET("/employee/:employee_id", scyllaHandler.GetAttendanceRecordsByEmployee)
			attendanceRecords.GET("/user/:employee_id", scyllaHandler.GetAttendanceRecordsByUser)
		}

		// Daily Summaries routes
		dailySummaries := v1.Group("/daily-summaries")
		{
			dailySummaries.GET("", scyllaHandler.GetDailySummaries)
			dailySummaries.GET("/user/:employee_id", scyllaHandler.GetDailySummariesByUser)
		}

		// Audit Logs routes
		auditLogs := v1.Group("/audit-logs")
		{
			auditLogs.GET("", scyllaHandler.GetAuditLogs)
			auditLogs.GET("/range", scyllaHandler.GetAuditLogsByTimeRange)
			auditLogs.POST("", scyllaHandler.CreateAuditLog)
		}

		// Face Enrollment Logs routes
		faceEnrollmentLogs := v1.Group("/face-enrollment-logs")
		{
			faceEnrollmentLogs.GET("", scyllaHandler.GetFaceEnrollmentLogs)
			faceEnrollmentLogs.GET("/employee/:employee_id", scyllaHandler.GetFaceEnrollmentLogsByEmployee)
		}

		// Attendance Records No Shift routes
		attendanceRecordsNoShift := v1.Group("/attendance-records-no-shift")
		{
			attendanceRecordsNoShift.GET("", scyllaHandler.GetAttendanceRecordsNoShift)
		}

		// ============================================
		// Company Admin - Advanced Analytics routes
		// ============================================
		company := v1.Group("/company")
		{
			// Daily attendance status
			company.GET("/daily-attendance-status", scyllaHandler.GetDailyAttendanceStatus)
			company.GET("/attendance-status/range", scyllaHandler.GetAttendanceStatusByTimeRange)

			// Monthly detailed summary
			company.GET("/monthly-summary", scyllaHandler.GetMonthlyDetailedSummary)

			// Export endpoints
			company.POST("/export-daily-status", scyllaHandler.ExportDailyStatus)
			company.POST("/export-monthly-summary", scyllaHandler.ExportMonthlySummary)
		}

		// ============================================
		// Employee Self-Service routes (for employees to view their own data)
		// ============================================
		employeeHandler := handler.NewEmployeeHandler()
		employee := v1.Group("/employee")
		{
			// Attendance Records - Employee's own records
			employee.GET("/my-attendance-records", employeeHandler.GetMyAttendanceRecords)
			employee.GET("/my-attendance-records/range", employeeHandler.GetMyAttendanceRecordsInRange)

			// Daily Summaries - Employee's own summaries
			employee.GET("/my-daily-summaries", employeeHandler.GetMyDailySummaries)
			employee.GET("/my-daily-summary/:date", employeeHandler.GetMyDailySummaryByDate)

			// Statistics - Employee's own statistics
			employee.GET("/my-stats", employeeHandler.GetMyAttendanceStats)

			// Detailed status endpoints
			employee.GET("/my-daily-status", employeeHandler.GetMyDailyStatus)
			employee.GET("/my-status/range", employeeHandler.GetMyStatusByTimeRange)
			employee.GET("/my-monthly-summary", employeeHandler.GetMyDetailedMonthlySummary)

			// Export endpoints
			employee.POST("/export-daily-status", employeeHandler.ExportMyDailyStatus)
			employee.POST("/export-monthly-summary", employeeHandler.ExportMyMonthlySummary)
		}
	}

	return router
}
