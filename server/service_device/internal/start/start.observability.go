package start

import (
	"context"
	"fmt"
	"net/http"

	"github.com/youknow2509/cio_verify_face/server/pkg/observability"
	domainConfig "github.com/youknow2509/cio_verify_face/server/service_device/internal/domain/config"
	"github.com/youknow2509/cio_verify_face/server/service_device/internal/global"
)

var (
	httpMetrics    *observability.Metrics
	grpcMetrics    *observability.GRPCMetrics
	tracerProvider *observability.TracerProvider
)

// initObservability initializes Prometheus metrics and Jaeger tracing
func initObservability(setting *domainConfig.ObservabilitySetting, serverName string) error {
	if !setting.Enabled {
		global.Logger.Info("Observability is disabled")
		return nil
	}

	// Initialize HTTP metrics
	httpMetrics = observability.NewMetrics(serverName)
	global.Logger.Info("Prometheus HTTP metrics initialized")

	// Initialize gRPC metrics
	grpcMetrics = observability.NewGRPCMetrics(serverName)
	global.Logger.Info("Prometheus gRPC metrics initialized")

	// Initialize tracing if enabled
	if setting.TracingEnabled && setting.OTLPEndpoint != "" {
		cfg := observability.TracingConfig{
			ServiceName:    serverName,
			JaegerEndpoint: setting.OTLPEndpoint,
			Environment:    global.SettingServer.Server.Mode,
			Version:        "1.0.0",
		}

		tp, err := observability.InitTracer(cfg)
		if err != nil {
			global.Logger.Warn(fmt.Sprintf("Failed to initialize tracing: %v", err))
		} else {
			tracerProvider = tp
			global.Logger.Info("Jaeger tracing initialized")
		}
	}

	// Start metrics server on a separate port
	if setting.MetricsPort > 0 {
		go startMetricsServer(setting.MetricsPort, setting.MetricsPath)
	}

	return nil
}

// startMetricsServer starts a dedicated HTTP server for Prometheus metrics
func startMetricsServer(port int, path string) {
	if path == "" {
		path = "/metrics"
	}

	mux := http.NewServeMux()
	mux.Handle(path, observability.HTTPMetricsHandler())
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	addr := fmt.Sprintf(":%d", port)
	global.Logger.Info(fmt.Sprintf("Starting metrics server on %s%s", addr, path))

	if err := http.ListenAndServe(addr, mux); err != nil {
		global.Logger.Error(fmt.Sprintf("Metrics server error: %v", err))
	}
}

// GetHTTPMetrics returns the HTTP metrics instance
func GetHTTPMetrics() *observability.Metrics {
	return httpMetrics
}

// GetGRPCMetrics returns the gRPC metrics instance
func GetGRPCMetrics() *observability.GRPCMetrics {
	return grpcMetrics
}

// GetTracerProvider returns the tracer provider instance
func GetTracerProvider() *observability.TracerProvider {
	return tracerProvider
}

// ShutdownObservability gracefully shuts down the observability components
func ShutdownObservability() {
	if tracerProvider != nil {
		if err := tracerProvider.Shutdown(context.Background()); err != nil {
			global.Logger.Error(fmt.Sprintf("Error shutting down tracer: %v", err))
		}
	}
}
