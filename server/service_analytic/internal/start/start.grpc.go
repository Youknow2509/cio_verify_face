package start

import (
	"fmt"
	"net"

	"github.com/youknow2509/cio_verify_face/server/service_analytic/internal/global"
	"github.com/youknow2509/cio_verify_face/server/service_analytic/internal/infrastructure/middleware"
	// grpcRouter "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/interfaces/grpc/router"
	// pb "github.com/youknow2509/cio_verify_face/server/service_analytic/proto/pb/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// initGrpcServer initializes and starts the gRPC server
func initGrpcServer() error {
	config := global.SettingServer.Grpc

	// Create listener
	lis, err := net.Listen(config.Network, fmt.Sprintf("%s:%d", config.Host, config.Port))
	if err != nil {
		global.Logger.Error("failed to listen for gRPC", "error", err)
		return err
	}

	// Create gRPC server with session interceptor
	// For inter-service calls, session info is passed directly in metadata (already authenticated)
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(middleware.SessionInterceptor()),
	}

	if config.TLS.Enabled {
		creds, err := credentials.NewServerTLSFromFile(config.TLS.CertFile, config.TLS.KeyFile)
		if err != nil {
			return err
		}
		opts = append(opts, grpc.Creds(creds))
	}

	grpcServer := grpc.NewServer(opts...)

	// Register services
	// NOTE: Actual registration will be done after proto generation:
	// pb.RegisterAnalyticServiceServer(grpcServer, grpcRouter.NewAnalyticRouter())

	// Start server in goroutine
	global.WaitGroup.Add(1)
	go func() {
		defer global.WaitGroup.Done()
		global.Logger.Info("Starting gRPC server", "address", lis.Addr().String())
		if err := grpcServer.Serve(lis); err != nil {
			global.Logger.Error("gRPC server error", "error", err)
		}
	}()

	return nil
}

// initAuthGrpcClient initializes the auth service gRPC client
func initAuthGrpcClient() error {
	config := global.SettingServer.AuthService
	
	if !config.Enabled {
		global.Logger.Info("Auth gRPC client disabled in config")
		return nil
	}

	// This will be implemented after proto generation
	// return infraGrpc.InitAuthClient(config.GrpcAddr)
	
	global.Logger.Info("Auth gRPC client placeholder", "addr", config.GrpcAddr)
	return nil
}
