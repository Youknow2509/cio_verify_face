package global

import (
	"sync"

	domainConfig "github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/domain/config"
	domainLogger "github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/domain/logger"
	grpcClient "github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/infrastructure/grpc"
)

var (
	WaitGroup     *sync.WaitGroup
	Logger        domainLogger.ILogger
	SettingServer domainConfig.Setting
	AuthClient    *grpcClient.AuthClient
)
