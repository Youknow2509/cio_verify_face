package start

import (
	"fmt"
	"net"

	"github.com/youknow2509/cio_verify_face/server/pkg/observability"
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

	// Create gRPC server with interceptors
	unaryInterceptors := []grpc.UnaryServerInterceptor{middleware.SessionInterceptor()}
	streamInterceptors := []grpc.StreamServerInterceptor{}
	var opts []grpc.ServerOption

	if config.TLS.Enabled {
		creds, err := credentials.NewServerTLSFromFile(config.TLS.CertFile, config.TLS.KeyFile)
		if err != nil {
			return err
		}
		opts = append(opts, grpc.Creds(creds))
	}

	if grpcMetrics := GetGRPCMetrics(); grpcMetrics != nil {
		unaryInterceptors = append(unaryInterceptors, grpcMetrics.UnaryServerInterceptor())
		streamInterceptors = append(streamInterceptors, grpcMetrics.StreamServerInterceptor())
	}
	if GetTracerProvider() != nil {
		unaryInterceptors = append(unaryInterceptors, observability.GRPCTracingUnaryServerInterceptor(global.SettingServer.Server.Name))
		streamInterceptors = append(streamInterceptors, observability.GRPCTracingStreamServerInterceptor(global.SettingServer.Server.Name))
	}

	opts = append(opts, grpc.ChainUnaryInterceptor(unaryInterceptors...))
	if len(streamInterceptors) > 0 {
		opts = append(opts, grpc.ChainStreamInterceptor(streamInterceptors...))
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
