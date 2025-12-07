package start

import (
	"fmt"
	"time"

	global "github.com/youknow2509/cio_verify_face/server/service_device/internal/global"
	pb "github.com/youknow2509/cio_verify_face/server/service_device/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

var (
	authGrpcClient pb.AuthServiceClient
	faceGrpcClient pb.FaceVerificationServiceClient
)

// init client grpc
func initClientGrpc() error {
	if err := initAuthClientGrpc(); err != nil {
		return err
	}
	if err := initFaceClientGrpc(); err != nil {
		return err
	}
	return nil
}

func initAuthClientGrpc() error {
	config := global.SettingServer.AuthService
	opts, err := buildGrpcDialOptions(config.Tls.Enabled, config.Tls.CertFile, config.KeepaliveTimeMs, config.KeepaliveTimeoutMs, config.KeepalivePermitWithoutCalls)
	if err != nil {
		return err
	}
	conn, err := grpc.Dial(config.GrpcAddr, opts...)
	if err != nil {
		return fmt.Errorf("failed to connect to auth gRPC server: %w", err)
	}
	authGrpcClient = pb.NewAuthServiceClient(conn)
	return nil
}

func initFaceClientGrpc() error {
	config := global.SettingServer.FaceService
	if !config.Enabled {
		return fmt.Errorf("face verification gRPC client is disabled in configuration")
	}
	opts, err := buildGrpcDialOptions(config.Tls.Enabled, config.Tls.CertFile, config.KeepaliveTimeMs, config.KeepaliveTimeoutMs, config.KeepalivePermitWithoutCalls)
	if err != nil {
		return err
	}
	conn, err := grpc.Dial(config.GrpcAddr, opts...)
	if err != nil {
		return fmt.Errorf("failed to connect to face verification gRPC server: %w", err)
	}
	faceGrpcClient = pb.NewFaceVerificationServiceClient(conn)
	return nil
}

func buildGrpcDialOptions(enableTLS bool, certFile string, keepaliveTimeMs, keepaliveTimeoutMs int, permitWithoutCalls bool) ([]grpc.DialOption, error) {
	var opts []grpc.DialOption
	if enableTLS {
		creds, err := credentials.NewClientTLSFromFile(certFile, "")
		if err != nil {
			return nil, fmt.Errorf("failed to load TLS credentials: %w", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	kaParams := grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:                time.Duration(keepaliveTimeMs) * time.Millisecond,
		Timeout:             time.Duration(keepaliveTimeoutMs) * time.Millisecond,
		PermitWithoutStream: permitWithoutCalls,
	})
	opts = append(opts, kaParams)
	return opts, nil
}
