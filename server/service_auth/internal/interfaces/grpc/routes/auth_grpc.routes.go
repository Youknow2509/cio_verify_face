package routes

import (
	"github.com/youknow2509/cio_verify_face/server/service_auth/internal/application/service"
	"github.com/youknow2509/cio_verify_face/server/service_auth/internal/domain/logger"
	grpcHandler "github.com/youknow2509/cio_verify_face/server/service_auth/internal/interfaces/grpc/handler"
	pb "github.com/youknow2509/cio_verify_face/server/service_auth/proto"
	"google.golang.org/grpc"
)

// RegisterGRPCRoutes registers gRPC service routes
func RegisterGRPCRoutes(
	s *grpc.Server,
	authCacheService service.IAuthCacheService,
	coreAuthService service.ICoreAuthService,
	logger logger.ILogger,
) {
	// Create gRPC handler
	authHandler := grpcHandler.NewAuthGRPCHandler(
		authCacheService,
		coreAuthService,
		logger,
	)

	// Register service
	pb.RegisterAuthServiceServer(s, authHandler)
}