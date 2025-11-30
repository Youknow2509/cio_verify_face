// Package observability provides Prometheus metrics and Jaeger tracing for services
package observability

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics holds all Prometheus metrics
type Metrics struct {
	RequestsTotal    *prometheus.CounterVec
	RequestDuration  *prometheus.HistogramVec
	RequestsInFlight prometheus.Gauge
	ResponseSize     *prometheus.HistogramVec
	ErrorsTotal      *prometheus.CounterVec
}

var (
	metricsOnce     sync.Once
	metricsInstance *Metrics
)

// NewMetrics creates a new Metrics instance with all metrics registered
// This function is safe to call multiple times - it uses a singleton pattern
func NewMetrics(serviceName string) *Metrics {
	metricsOnce.Do(func() {
		namespace := serviceName

		metricsInstance = &Metrics{
			RequestsTotal: prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Namespace: namespace,
					Name:      "http_requests_total",
					Help:      "Total number of HTTP requests",
				},
				[]string{"method", "path", "status"},
			),
			RequestDuration: prometheus.NewHistogramVec(
				prometheus.HistogramOpts{
					Namespace: namespace,
					Name:      "http_request_duration_seconds",
					Help:      "HTTP request duration in seconds",
					Buckets:   []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
				},
				[]string{"method", "path", "status"},
			),
			RequestsInFlight: prometheus.NewGauge(
				prometheus.GaugeOpts{
					Namespace: namespace,
					Name:      "http_requests_in_flight",
					Help:      "Current number of HTTP requests being processed",
				},
			),
			ResponseSize: prometheus.NewHistogramVec(
				prometheus.HistogramOpts{
					Namespace: namespace,
					Name:      "http_response_size_bytes",
					Help:      "HTTP response size in bytes",
					Buckets:   prometheus.ExponentialBuckets(100, 10, 8),
				},
				[]string{"method", "path", "status"},
			),
			ErrorsTotal: prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Namespace: namespace,
					Name:      "http_errors_total",
					Help:      "Total number of HTTP errors",
				},
				[]string{"method", "path", "status"},
			),
		}

		// Register all metrics
		prometheus.MustRegister(
			metricsInstance.RequestsTotal,
			metricsInstance.RequestDuration,
			metricsInstance.RequestsInFlight,
			metricsInstance.ResponseSize,
			metricsInstance.ErrorsTotal,
		)
	})

	return metricsInstance
}

// GinMiddleware returns a Gin middleware that records metrics for each request
func (m *Metrics) GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		method := c.Request.Method

		// Track in-flight requests
		m.RequestsInFlight.Inc()
		defer m.RequestsInFlight.Dec()

		// Process request
		c.Next()

		// Record metrics after request is processed
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		size := float64(c.Writer.Size())

		m.RequestsTotal.WithLabelValues(method, path, status).Inc()
		m.RequestDuration.WithLabelValues(method, path, status).Observe(duration)
		m.ResponseSize.WithLabelValues(method, path, status).Observe(size)

		// Record errors (4xx and 5xx status codes)
		if c.Writer.Status() >= 400 {
			m.ErrorsTotal.WithLabelValues(method, path, status).Inc()
		}
	}
}

// MetricsHandler returns the Prometheus metrics HTTP handler
func MetricsHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// HTTPMetricsHandler returns a standard http.Handler for Prometheus metrics
func HTTPMetricsHandler() http.Handler {
	return promhttp.Handler()
}

// GetHTTPMetrics returns the singleton metrics instance
func GetHTTPMetrics() *Metrics {
	return metricsInstance
}

// GetGRPCMetrics returns the singleton gRPC metrics instance
func GetGRPCMetrics() *GRPCMetrics {
	return grpcMetricsInstance
}
