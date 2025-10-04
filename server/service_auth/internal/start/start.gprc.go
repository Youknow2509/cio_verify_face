package start

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/youknow2509/cio_verify_face/server/service_auth/internal/application/service"
	"github.com/youknow2509/cio_verify_face/server/service_auth/internal/domain/logger"
	global "github.com/youknow2509/cio_verify_face/server/service_auth/internal/global"
	grpcHandler "github.com/youknow2509/cio_verify_face/server/service_auth/internal/interfaces/grpc/handler"
	pb "github.com/youknow2509/cio_verify_face/server/service_auth/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

// init server grpc
func initServerGrpc() error {
	config := global.SettingServer.GrpcServer

	// Initialize the gRPC server
	lis, err := net.Listen(
		"tcp", fmt.Sprintf(":%d", config.Port),
	)
	if err != nil {
		global.Logger.Error("failed to listen", "error", err)
		return err
	}

	// ServerOption
	var opts []grpc.ServerOption

	// TLS
	if config.Tls.Enabled {
		// check file existence
		if _, err := os.Stat(config.Tls.CertFile); os.IsNotExist(err) {
			global.Logger.Error("TLS cert file does not exist", "error", err)
			return err
		}
		if _, err := os.Stat(config.Tls.KeyFile); os.IsNotExist(err) {
			global.Logger.Error("TLS key file does not exist", "error", err)
			return err
		}

		// create TLS credentials
		creds, err := credentials.NewServerTLSFromFile(
			config.Tls.CertFile,
			config.Tls.KeyFile,
		)
		if err != nil {
			global.Logger.Error("failed to generate credentials", "error", err)
			return err
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	// Add middleware interceptors
	loggerService := logger.GetLogger()

	// Recovery interceptor
	recoveryInterceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		defer func() {
			if r := recover(); r != nil {
				loggerService.Error("gRPC handler panicked", "method", info.FullMethod, "panic", r)
			}
		}()
		return handler(ctx, req)
	}

	// Logging interceptor
	loggingInterceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		loggerService.Info("gRPC request started", "method", info.FullMethod)
		resp, err := handler(ctx, req)
		duration := time.Since(start)
		if err != nil {
			loggerService.Error("gRPC request failed", "method", info.FullMethod, "duration", duration.String(), "error", err.Error())
		} else {
			loggerService.Info("gRPC request completed", "method", info.FullMethod, "duration", duration.String())
		}
		return resp, err
	}

	// Chain interceptors using ChainUnaryInterceptor
	opts = append(opts, grpc.ChainUnaryInterceptor(
		recoveryInterceptor,
		loggingInterceptor,
	)) // init server
	grpcServer := grpc.NewServer(opts...)

	// Get services
	authCacheService := service.GetAuthCacheService()
	coreAuthService := service.GetCoreAuthService()
	loggerService2 := logger.GetLogger()

	// Create and register gRPC handler
	authHandler := grpcHandler.NewAuthGRPCHandler(
		authCacheService,
		coreAuthService,
		loggerService2,
	)

	// Register service
	pb.RegisterAuthServiceServer(grpcServer, authHandler)

	// Enable reflection for development/debugging
	reflection.Register(grpcServer)

	// start server
	global.Logger.Info("gRPC server starting", "address", lis.Addr().String())

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			global.Logger.Error("failed to start gRPC server", "error", err)
		}
	}()

	return nil
}
