package start

import (
	"context"
	"fmt"
	"net/http"
	"time"

	domainConfig "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/domain/config"
	"github.com/youknow2509/cio_verify_face/server/service_analytic/internal/global"
	httpRouter "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/interfaces/http/router"
)

// initHTTPServer initializes and starts the HTTP server
func initHTTPServer(serverConfig *domainConfig.ServerConfig) error {
	// Setup router
	router := httpRouter.SetupRouter()

	// Create HTTP server
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", serverConfig.Port),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	// Start server in goroutine
	global.WaitGroup.Add(1)
	go func() {
		defer global.WaitGroup.Done()
		
		global.Logger.Info("Starting HTTP server", "port", serverConfig.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			global.Logger.Error("HTTP server error", "error", err.Error())
		}
	}()

	// Graceful shutdown handler
	go func() {
		<-make(chan struct{}) // This would be replaced with actual shutdown signal
		
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		if err := server.Shutdown(ctx); err != nil {
			global.Logger.Error("HTTP server shutdown error", "error", err.Error())
		}
	}()

	return nil
}
