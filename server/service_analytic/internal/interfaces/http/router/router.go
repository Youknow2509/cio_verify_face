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
)

// SetupRouter sets up the HTTP router with authentication middleware
func SetupRouter() *gin.Engine {
	// Create gin router
	router := gin.Default()

	// Setup CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	router.Use(cors.New(config))

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
		}
	}

	return router
}
