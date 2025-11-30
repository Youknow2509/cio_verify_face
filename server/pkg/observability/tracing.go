package observability

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

// TracingConfig holds configuration for tracing
type TracingConfig struct {
	ServiceName    string
	JaegerEndpoint string // e.g., "http://jaeger:4318/v1/traces"
	Environment    string
	Version        string
}

// TracerProvider wraps the OpenTelemetry TracerProvider
type TracerProvider struct {
	provider *sdktrace.TracerProvider
	tracer   trace.Tracer
}

// InitTracer initializes OpenTelemetry tracing with Jaeger exporter
func InitTracer(cfg TracingConfig) (*TracerProvider, error) {
	ctx := context.Background()

	// Create OTLP HTTP exporter (Jaeger supports OTLP)
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpointURL(cfg.JaegerEndpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	// Create resource with service information
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion(cfg.Version),
			attribute.String("environment", cfg.Environment),
		),
	)
	if err != nil {
		return nil, err
	}

	// Create TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	// Set global TracerProvider
	otel.SetTracerProvider(tp)

	// Set global propagator
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return &TracerProvider{
		provider: tp,
		tracer:   tp.Tracer(cfg.ServiceName),
	}, nil
}

// Shutdown cleanly shuts down the tracer provider
func (tp *TracerProvider) Shutdown(ctx context.Context) error {
	return tp.provider.Shutdown(ctx)
}

// Tracer returns the tracer instance
func (tp *TracerProvider) Tracer() trace.Tracer {
	return tp.tracer
}

// GinTracingMiddleware returns a Gin middleware that traces HTTP requests
func (tp *TracerProvider) GinTracingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract trace context from incoming request
		ctx := otel.GetTextMapPropagator().Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))

		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		// Start a new span
		ctx, span := tp.tracer.Start(ctx, c.Request.Method+" "+path,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				semconv.HTTPMethod(c.Request.Method),
				semconv.HTTPURL(c.Request.URL.String()),
				semconv.HTTPRoute(path),
				semconv.HTTPUserAgent(c.Request.UserAgent()),
				attribute.String("http.client_ip", c.ClientIP()),
			),
		)
		defer span.End()

		// Store span in context
		c.Request = c.Request.WithContext(ctx)

		// Process request
		startTime := time.Now()
		c.Next()
		duration := time.Since(startTime)

		// Add response attributes
		span.SetAttributes(
			semconv.HTTPStatusCode(c.Writer.Status()),
			attribute.Int64("http.response_content_length", int64(c.Writer.Size())),
			attribute.Float64("http.duration_ms", float64(duration.Milliseconds())),
		)

		// Record error if status code >= 400
		if c.Writer.Status() >= 400 {
			span.SetAttributes(attribute.Bool("error", true))
		}
	}
}

// GetSpanFromContext extracts the current span from context
func GetSpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

// StartSpan starts a new span as a child of the current span in context
func StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return otel.Tracer("").Start(ctx, name, opts...)
}
