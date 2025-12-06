package start

import (
	"fmt"
	"time"

	global "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/global"
	pb "github.com/youknow2509/cio_verify_face/server/service_workforce/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

var (
	grpcClient pb.AuthServiceClient
)

// init client grpc
func initClientGrpc() error {
	config := global.SettingServer.AuthService
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
	// create connection
	conn, err := grpc.Dial(config.GrpcAddr, opts...)
	if err != nil {
		return fmt.Errorf("failed to connect to gRPC server: %w", err)
	}
	grpcClient = pb.NewAuthServiceClient(conn)
	return nil
}
