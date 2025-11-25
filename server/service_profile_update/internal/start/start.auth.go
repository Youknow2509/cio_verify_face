package start

import (
	"fmt"

	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/global"
	grpcClient "github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/infrastructure/grpc"
)

// initAuthClient initializes the auth service gRPC client
func initAuthClient() error {
	cfg := &global.SettingServer.ServiceAuth

	if !cfg.Enabled {
		global.Logger.Info("Auth service is disabled, skipping initialization")
		return nil
	}

	global.Logger.Info(fmt.Sprintf("Initializing auth service client to %s", cfg.GrpcAddr))

	client, err := grpcClient.NewAuthClient(cfg)
	if err != nil {
		return fmt.Errorf("failed to create auth client: %w", err)
	}

	// Store in global
	global.AuthClient = client

	global.Logger.Info("Auth service client initialized successfully")
	return nil
}
