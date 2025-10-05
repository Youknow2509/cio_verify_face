package start

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	libsMiddleware "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/infrastructure/middleware"
	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/global"
	httpRouter "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/interfaces/http/routes"
	wsRouter "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/interfaces/ws/router"
)

func initGinRouter() error {
	ginEngine := gin.Default()
	serverHttpSetting := global.ServerSetting
	// Set Gin mode
	if serverHttpSetting.Mode != "dev" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	// Set port for Gin
	portGin := fmt.Sprintf(":%d", serverHttpSetting.Port)

	// Initialize routes
	if err := initRouter(ginEngine); err != nil {
		return err
	}
	// Start Gin server
	if err := ginEngine.Run(portGin); err != nil {
		return err
	}
	return nil
}

func initRouter(ginEngine *gin.Engine) error {
	// Global middleware
	ginEngine.Use(cors.Default()) // TODO: customize CORS settings - Now allow all origins
	ginEngine.Use(libsMiddleware.GetValidateMiddleware().Apply())
	// Initialize ws routes
	wsRouter.GetBaseRouter().Initialize(ginEngine)
	// Initialize http routes
	apiGroup := ginEngine.Group("/api")
	httpRouter.GetRouteManager().HealthRoute.InitHealthRoute(apiGroup)
	return nil
}
