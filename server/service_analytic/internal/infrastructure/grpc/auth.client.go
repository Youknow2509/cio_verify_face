package grpc

import (
	"context"
	"errors"
	"fmt"
	"time"

	applicationModel "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/application/model"
	"github.com/youknow2509/cio_verify_face/server/service_analytic/internal/global"
	pb "github.com/youknow2509/cio_verify_face/server/service_analytic/proto/pb/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// AuthServiceClient represents a client for the auth service
type AuthServiceClient struct {
	conn   *grpc.ClientConn
	addr   string
	client pb.AuthServiceClient
}

// NewAuthServiceClient creates a new auth service client
func NewAuthServiceClient() (*AuthServiceClient, error) {
	if global.SettingServer.AuthService.GrpcAddr == "" {
		return nil, errors.New("auth service address not configured")
	}

	addr := global.SettingServer.AuthService.GrpcAddr

	// Connect to auth service with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service at %s: %w", addr, err)
	}

	return &AuthServiceClient{
		conn: conn,
		addr: addr,
		client: pb.NewAuthServiceClient(conn), 
	}, nil
}

// Close closes the auth client connection
func (c *AuthServiceClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// ValidateToken validates a JWT token via auth service and returns session info
// This is a placeholder until proto is generated
func (c *AuthServiceClient) ValidateToken(ctx context.Context, token string) (*applicationModel.SessionInfo, error) {
	req := &pb.ParseUserTokenRequest{Token: token}
	resp, err := c.client.ParseUserToken(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}

	return &applicationModel.SessionInfo{
		UserID:    resp.UserId,
		Role:      resp.Roles,
		SessionID: resp.TokenId,
		CompanyID: resp.CompanyId,
	}, nil
}

// ParseUserTokenResponse represents the response from ParseUserToken
type ParseUserTokenResponse struct {
	UserID    string
	Role      int32
	TokenID   string
	CompanyID string
	ExpiresAt int64
}
