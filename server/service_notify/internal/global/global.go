package global

import (
	"context"
	"sync"

	"github.com/go-playground/validator/v10"
	domainConfig "github.com/youknow2509/cio_verify_face/server/service_notify/internal/domain/config"
	domainLogger "github.com/youknow2509/cio_verify_face/server/service_notify/internal/domain/logger"
)

var (
	WaitGroup     *sync.WaitGroup
	Logger        domainLogger.ILogger
	SettingServer domainConfig.Setting
	Validator     *validator.Validate
	ContextSystem context.Context
	CancelFunc    context.CancelFunc
)
