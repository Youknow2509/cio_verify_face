package start

import (
	"fmt"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/domain/config"
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/global"
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/interfaces/http"
)

func initGinRouter(cfg *config.ServerSetting) error {
	// Set Gin mode
	router := gin.New()
	if cfg.Mode != "dev" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Setup routes
	http.SetupRoutes(router)

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Port)
	global.Logger.Info(fmt.Sprintf("Starting HTTP server on %s", addr))
	global.Logger.Info(fmt.Sprintf("Swagger documentation available at http://localhost%s/swagger/index.html", addr))

	return router.Run(addr)
}
