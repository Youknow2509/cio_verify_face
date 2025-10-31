package global

import (
	"context"
	"sync"

	"github.com/go-playground/validator/v10"
	libsConfig "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/config"
	libsPkgLogger "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/logger"
	libsDomainRateLimit "github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/ratelimit"
)

// =========================================
//
//	Global instances
//
// =========================================
var (
	ServerSetting          libsConfig.ServerSetting
	ServerWsSetting        libsConfig.WsSetting
	ServerGrpcSetting      libsConfig.GrpcSetting
	RateLimitPolicyManager libsDomainRateLimit.IPolicytRegistry
	RateLimitWsRead        libsDomainRateLimit.ILimiter
	Logger                 libsPkgLogger.ILogger
	WsContext              context.Context
	WaitGroup              *sync.WaitGroup
	Validate               *validator.Validate
)
