package observability

import (
	"context"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// GRPCMetrics holds Prometheus metrics for gRPC
type GRPCMetrics struct {
	RequestsTotal   *prometheus.CounterVec
	RequestDuration *prometheus.HistogramVec
	ErrorsTotal     *prometheus.CounterVec
}

var (
	grpcMetricsOnce     sync.Once
	grpcMetricsInstance *GRPCMetrics
)

// NewGRPCMetrics creates new gRPC metrics
// This function is safe to call multiple times - it uses a singleton pattern
func NewGRPCMetrics(serviceName string) *GRPCMetrics {
	grpcMetricsOnce.Do(func() {
		namespace := serviceName

		grpcMetricsInstance = &GRPCMetrics{
			RequestsTotal: prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Namespace: namespace,
					Name:      "grpc_requests_total",
					Help:      "Total number of gRPC requests",
				},
				[]string{"method", "status"},
			),
			RequestDuration: prometheus.NewHistogramVec(
				prometheus.HistogramOpts{
					Namespace: namespace,
					Name:      "grpc_request_duration_seconds",
					Help:      "gRPC request duration in seconds",
					Buckets:   []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
				},
				[]string{"method", "status"},
			),
			ErrorsTotal: prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Namespace: namespace,
					Name:      "grpc_errors_total",
					Help:      "Total number of gRPC errors",
				},
				[]string{"method", "code"},
			),
		}

		prometheus.MustRegister(
			grpcMetricsInstance.RequestsTotal,
			grpcMetricsInstance.RequestDuration,
			grpcMetricsInstance.ErrorsTotal,
		)
	})

	return grpcMetricsInstance
}

// UnaryServerInterceptor returns a gRPC unary server interceptor for metrics
func (m *GRPCMetrics) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(start).Seconds()
		statusCode := "OK"
		if err != nil {
			if st, ok := status.FromError(err); ok {
				statusCode = st.Code().String()
			} else {
				statusCode = "UNKNOWN"
			}
			m.ErrorsTotal.WithLabelValues(info.FullMethod, statusCode).Inc()
		}

		m.RequestsTotal.WithLabelValues(info.FullMethod, statusCode).Inc()
		m.RequestDuration.WithLabelValues(info.FullMethod, statusCode).Observe(duration)

		return resp, err
	}
}

// StreamServerInterceptor returns a gRPC stream server interceptor for metrics
func (m *GRPCMetrics) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()

		err := handler(srv, ss)

		duration := time.Since(start).Seconds()
		statusCode := "OK"
		if err != nil {
			if st, ok := status.FromError(err); ok {
				statusCode = st.Code().String()
			} else {
				statusCode = "UNKNOWN"
			}
			m.ErrorsTotal.WithLabelValues(info.FullMethod, statusCode).Inc()
		}

		m.RequestsTotal.WithLabelValues(info.FullMethod, statusCode).Inc()
		m.RequestDuration.WithLabelValues(info.FullMethod, statusCode).Observe(duration)

		return err
	}
}

// GRPCTracingUnaryServerInterceptor returns a gRPC unary server interceptor for tracing
func GRPCTracingUnaryServerInterceptor(serviceName string) grpc.UnaryServerInterceptor {
	tracer := otel.Tracer(serviceName)
	
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Extract trace context from metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.MD{}
		}
		carrier := &metadataCarrier{md: md}
		ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

		// Start a new span
		ctx, span := tracer.Start(ctx, info.FullMethod,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				semconv.RPCSystemGRPC,
				semconv.RPCMethod(info.FullMethod),
			),
		)
		defer span.End()

		// Process request
		resp, err := handler(ctx, req)

		// Record error if any
		if err != nil {
			span.SetAttributes(attribute.Bool("error", true))
			if st, ok := status.FromError(err); ok {
				span.SetAttributes(attribute.String("grpc.status_code", st.Code().String()))
			}
		}

		return resp, err
	}
}

// GRPCTracingStreamServerInterceptor returns a gRPC stream server interceptor for tracing
func GRPCTracingStreamServerInterceptor(serviceName string) grpc.StreamServerInterceptor {
	tracer := otel.Tracer(serviceName)

	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		
		// Extract trace context from metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.MD{}
		}
		carrier := &metadataCarrier{md: md}
		ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

		// Start a new span
		ctx, span := tracer.Start(ctx, info.FullMethod,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				semconv.RPCSystemGRPC,
				semconv.RPCMethod(info.FullMethod),
			),
		)
		defer span.End()

		// Wrap the server stream with the new context
		wrapped := &wrappedServerStream{
			ServerStream: ss,
			ctx:          ctx,
		}

		// Process stream
		err := handler(srv, wrapped)

		// Record error if any
		if err != nil {
			span.SetAttributes(attribute.Bool("error", true))
			if st, ok := status.FromError(err); ok {
				span.SetAttributes(attribute.String("grpc.status_code", st.Code().String()))
			}
		}

		return err
	}
}

// metadataCarrier implements TextMapCarrier for gRPC metadata
type metadataCarrier struct {
	md metadata.MD
}

func (c *metadataCarrier) Get(key string) string {
	values := c.md.Get(key)
	if len(values) > 0 {
		return values[0]
	}
	return ""
}

func (c *metadataCarrier) Set(key, value string) {
	c.md.Set(key, value)
}

func (c *metadataCarrier) Keys() []string {
	keys := make([]string, 0, len(c.md))
	for k := range c.md {
		keys = append(keys, k)
	}
	return keys
}

// wrappedServerStream wraps a grpc.ServerStream with a custom context
type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}
