package grpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/domain/config"
	pb "github.com/youknow2509/cio_verify_face/server/service_profile_update/proto"
)

// AuthClient wraps the gRPC client for auth service
type AuthClient struct {
	conn   *grpc.ClientConn
	client pb.AuthServiceClient
}

// NewAuthClient creates a new auth service gRPC client
func NewAuthClient(cfg *config.ServiceAuthSetting) (*AuthClient, error) {
	if !cfg.Enabled {
		return nil, fmt.Errorf("auth service is disabled")
	}

	var opts []grpc.DialOption

	// TLS configuration
	if cfg.TLS.Enabled {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: false,
		}
		creds := credentials.NewTLS(tlsConfig)
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// Keepalive configuration
	kacp := keepalive.ClientParameters{
		Time:                time.Duration(cfg.KeepaliveTimeMs) * time.Millisecond,
		Timeout:             time.Duration(cfg.KeepaliveTimeoutMs) * time.Millisecond,
		PermitWithoutStream: cfg.KeepalivePermitWithoutCalls,
	}
	opts = append(opts, grpc.WithKeepaliveParams(kacp))

	// Connect to auth service
	conn, err := grpc.Dial(cfg.GrpcAddr, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service: %w", err)
	}

	client := pb.NewAuthServiceClient(conn)

	return &AuthClient{
		conn:   conn,
		client: client,
	}, nil
}

// ParseUserToken parses and validates a user token
func (c *AuthClient) ParseUserToken(ctx context.Context, token string) (*pb.ParseUserTokenResponse, error) {
	req := &pb.ParseUserTokenRequest{
		Token: token,
	}
	return c.client.ParseUserToken(ctx, req)
}

// ParseServiceToken parses and validates a service token
func (c *AuthClient) ParseServiceToken(ctx context.Context, serviceID string) (*pb.ParseServiceTokenResponse, error) {
	req := &pb.ParseServiceTokenRequest{
		ServiceId: serviceID,
	}
	return c.client.ParseServiceToken(ctx, req)
}

// ParseDeviceToken parses and validates a device token
func (c *AuthClient) ParseDeviceToken(ctx context.Context, token, deviceID string) (*pb.ParseDeviceTokenResponse, error) {
	req := &pb.ParseDeviceTokenRequest{
		Token:    token,
		DeviceId: deviceID,
	}
	return c.client.ParseDeviceToken(ctx, req)
}

// HealthCheck checks if auth service is healthy
func (c *AuthClient) HealthCheck(ctx context.Context) error {
	_, err := c.client.HealthCheck(ctx, &emptypb.Empty{})
	return err
}

// Close closes the gRPC connection
func (c *AuthClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
