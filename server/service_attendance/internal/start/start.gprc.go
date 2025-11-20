package start

import (
	"fmt"
	"net"
	"os"

	global "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/global"
	interfaceGrpc "github.com/youknow2509/cio_verify_face/server/service_attendance/internal/interfaces/grpc"
	pb "github.com/youknow2509/cio_verify_face/server/service_attendance/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	grpcClient pb.AuthServiceClient
)

// init server grpc
func initServerGrpc() error {
	config := global.SettingServer.Grpc

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
	conn, err := grpc.Dial(config.GrpcAddr, opts...)
	if err != nil {
		return fmt.Errorf("failed to connect to gRPC server: %w", err)
	}
	grpcClient = pb.NewAuthServiceClient(conn)
	return nil
}
