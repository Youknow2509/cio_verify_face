package start

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	domainConfig "github.com/youknow2509/cio_verify_face/server/service_notify/internal/domain/config"
	"github.com/youknow2509/cio_verify_face/server/service_notify/internal/global"
	infraMiddleware "github.com/youknow2509/cio_verify_face/server/service_notify/internal/infrastructure/middleware"
)

func initGinRouter(setting *domainConfig.ServerSetting) error {
	var ginEngine *gin.Engine
	// Set Gin mode
	if setting.Mode != "dev" {
		gin.SetMode(gin.ReleaseMode)
		ginEngine = gin.New()
	} else {
		gin.SetMode(gin.DebugMode)
		gin.ForceConsoleColor()
		ginEngine = gin.Default()
	}

	// Set port for Gin
	portGin := fmt.Sprintf(":%d", setting.Port)

	// Health check endpoint
	ginEngine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
	// Initialize routes
	if err := initRouter(ginEngine); err != nil {
		return err
	}
	// Start Gin server
	global.WaitGroup.Add(1)
	go func() {
		defer global.WaitGroup.Done()
		err := ginEngine.Run(portGin)
		if err != nil {
			global.Logger.Error(err.Error())
		}
	}()

	return nil
}

func initRouter(ginEngine *gin.Engine) error {
	// global middleware
	ginEngine.Use(getConfigCors())
	ginEngine.Use(infraMiddleware.GetValidateMiddleware().Apply())
	// Initialize routes
	
	return nil
}

func getConfigCors() gin.HandlerFunc {
	corsConfig := cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"*"}, // allow all request headers
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	return cors.New(corsConfig)
}