package middleware

import (
	"context"
	"strings"
	"time"

	"github.com/youknow2509/cio_verify_face/server/service_auth/internal/domain/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// GRPCLoggingInterceptor logs gRPC requests and responses
func GRPCLoggingInterceptor(logger logger.ILogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		// Extract metadata
		md, _ := metadata.FromIncomingContext(ctx)
		clientIP := getClientIPFromMetadata(md)

		// Log request
		logger.Info("gRPC request started",
			"method", info.FullMethod,
			"client_ip", clientIP,
			"request", req,
		)

		// Call handler
		resp, err := handler(ctx, req)

		// Log response
		duration := time.Since(start)
		if err != nil {
			logger.Error("gRPC request failed",
				"method", info.FullMethod,
				"client_ip", clientIP,
				"duration", duration.String(),
				"error", err.Error(),
			)
		} else {
			logger.Info("gRPC request completed",
				"method", info.FullMethod,
				"client_ip", clientIP,
				"duration", duration.String(),
			)
		}

		return resp, err
	}
}

// GRPCRecoveryInterceptor recovers from panics in gRPC handlers
func GRPCRecoveryInterceptor(logger logger.ILogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("gRPC handler panicked",
					"method", info.FullMethod,
					"panic", r,
				)
				err = status.Errorf(codes.Internal, "internal server error")
			}
		}()

		return handler(ctx, req)
	}
}

// GRPCRateLimitInterceptor provides rate limiting for gRPC requests
func GRPCRateLimitInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// TODO: Implement rate limiting logic
		// For now, just pass through
		return handler(ctx, req)
	}
}

// GRPCAuthInterceptor provides authentication for sensitive methods
func GRPCAuthInterceptor(sensitiveMethods []string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Skip auth for public methods
		if !isSensitiveMethod(info.FullMethod, sensitiveMethods) {
			return handler(ctx, req)
		}

		// Extract auth token from metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "metadata not found")
		}

		authHeaders := md.Get("authorization")
		if len(authHeaders) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "authorization token not provided")
		}

		token := strings.TrimPrefix(authHeaders[0], "Bearer ")
		if token == "" {
			return nil, status.Errorf(codes.Unauthenticated, "invalid authorization token format")
		}

		// TODO: Validate token here
		// For now, just pass through

		return handler(ctx, req)
	}
}

// Helper functions

func getClientIPFromMetadata(md metadata.MD) string {
	// Try x-forwarded-for first
	if xff := md.Get("x-forwarded-for"); len(xff) > 0 {
		return strings.Split(xff[0], ",")[0]
	}

	// Try x-real-ip
	if xri := md.Get("x-real-ip"); len(xri) > 0 {
		return xri[0]
	}

	// Try remote-addr
	if ra := md.Get("remote-addr"); len(ra) > 0 {
		return ra[0]
	}

	return "unknown"
}

func isSensitiveMethod(method string, sensitiveList []string) bool {
	for _, sensitive := range sensitiveList {
		if method == sensitive {
			return true
		}
	}
	return false
}
