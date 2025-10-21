package start

import (
	"fmt"
	global "github.com/youknow2509/cio_verify_face/server/service_device/internal/global"
	pb "github.com/youknow2509/cio_verify_face/server/service_device/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	grpcClient pb.AuthServiceClient
)

// init client grpc
func initClientGrpc() error {
	config := global.SettingServer.GrpcClient
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
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", config.Network, config.Port), opts...)
	if err != nil {
		return fmt.Errorf("failed to connect to gRPC server: %w", err)
	}
	grpcClient = pb.NewAuthServiceClient(conn)
	return nil
}
