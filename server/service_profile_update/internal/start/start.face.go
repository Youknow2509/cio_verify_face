package start

import (
	"fmt"

	domainGrpc "github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/domain/grpc"
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/global"
	grpcClient "github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/infrastructure/grpc"
)

// initFaceServiceClient initializes the face verification service gRPC client.
func initFaceServiceClient() error {
	cfg := &global.SettingServer.ServiceFace

	if !cfg.Enabled {
		global.Logger.Info("Face service is disabled, skipping initialization")
		return nil
	}

	global.Logger.Info(fmt.Sprintf("Initializing face service client to %s", cfg.GrpcAddr))

	client, err := grpcClient.NewFaceServiceClient(cfg)
	if err != nil {
		return fmt.Errorf("failed to create face service client: %w", err)
	}

	domainGrpc.SetFaceServiceClient(client)

	global.Logger.Info("Face service client initialized successfully")
	return nil
}
