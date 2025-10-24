package start

import (
	"context"
	"errors"

	"github.com/youknow2509/cio_verify_face/server/service_workforce/internal/constants"
	domainConfig "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/domain/config"
	infraConfig "github.com/youknow2509/cio_verify_face/server/service_workforce/internal/infrastructure/config"
)

func loadConfig() (*domainConfig.Setting, error) {
	implConfig := infraConfig.NewViperConfig()
	if error := domainConfig.SetConfig(implConfig); error != nil {
		return nil, error
	}
	setting := domainConfig.GetConfig()
	if setting == nil {
		return nil, errors.New("failed to get config")
	}
	ctx := context.Background()
	err := setting.LoadConfig(ctx, constants.FilePathConfigProd)
	if err != nil {
		return nil, err
	}
	config, err := setting.GetConfig(ctx)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
