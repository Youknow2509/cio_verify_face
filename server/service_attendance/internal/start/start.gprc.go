package start

import (
	"fmt"
	"net"
	"os"
	"time"

	global "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/global"
	interfaceGrpc "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/interfaces/grpc"
	pb "github.com/youknow2509/cio_verify_face/server/service_attendance/proto"
	"github.com/youknow2509/cio_verify_face/server/pkg/observability"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

var (
	grpcClient pb.AuthServiceClient
)

// init server grpc
func initServerGrpc() error {
	config := global.SettingServer.Grpc
	serviceName := global.SettingServer.Server.Name

	// Initialize the gRPC server
	lis, err := net.Listen(
		config.Network,
		fmt.Sprintf(":%d", config.Port),
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
	// Server Parameters
	kaParams := grpc.KeepaliveParams(keepalive.ServerParameters{
		Time:    time.Duration(config.KeepaliveTimeMs) * time.Millisecond,
		Timeout: time.Duration(config.KeepaliveTimeoutMs) * time.Millisecond,
	})
	opts = append(opts, kaParams)
	// Enforcement Policy
	kaEnforcement := grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
		MinTime:             time.Duration(config.Http2MinTimeBetweenPingsMs) * time.Millisecond,
		PermitWithoutStream: config.KeepalivePermitWithoutCalls,
	})
	opts = append(opts, kaEnforcement)

	// Add observability interceptors
	if grpcMetrics := GetGRPCMetrics(); grpcMetrics != nil {
		opts = append(opts,
			grpc.ChainUnaryInterceptor(grpcMetrics.UnaryServerInterceptor()),
			grpc.ChainStreamInterceptor(grpcMetrics.StreamServerInterceptor()),
		)
	}

	// Add tracing interceptors if tracing is enabled
	if GetTracerProvider() != nil {
		opts = append(opts,
			grpc.ChainUnaryInterceptor(observability.GRPCTracingUnaryServerInterceptor(serviceName)),
			grpc.ChainStreamInterceptor(observability.GRPCTracingStreamServerInterceptor(serviceName)),
		)
	}

	grpcServer := grpc.NewServer(opts...)

	// Register service
	pb.RegisterAttendanceServiceServer(grpcServer, interfaceGrpc.NewAttendanceGRPCServer())

	// start server
	global.Logger.Info("gRPC server starting", "address", lis.Addr().String())

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			global.Logger.Error("failed to start gRPC server", "error", err)
		}
	}()

	return nil
}

// init client grpc
func initClientGrpc() error {
	config := global.SettingServer.ServiceAuth
	// load configuration
	var opts []grpc.DialOption
	if config.Tls.Enabled {
		creds, err := credentials.NewClientTLSFromFile(config.Tls.CertFile, "")
		if err != nil {
			return fmt.Errorf("failed to load TLS credentials: %w", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	// Keepalive parameters
	kaParams := grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:                time.Duration(config.KeepaliveTimeMs) * time.Millisecond,
		Timeout:             time.Duration(config.KeepaliveTimeoutMs) * time.Millisecond,
		PermitWithoutStream: config.KeepalivePermitWithoutCalls,
	})
	opts = append(opts, kaParams)
	// HTTP/2 Ping Policy
	// http2PingPolicy := grpc.WithDefaultCallOptions(
	// 	grpc.MaxCallRecvMsgSize(config.Http2MaxPingsWithoutData),
	// )
	// opts = append(opts, http2PingPolicy)
	// create connection
	conn, err := grpc.Dial(config.GrpcAddr, opts...)
	if err != nil {
		return fmt.Errorf("failed to connect to gRPC server: %w", err)
	}
	grpcClient = pb.NewAuthServiceClient(conn)
	return nil
}
