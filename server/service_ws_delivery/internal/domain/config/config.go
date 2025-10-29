package config

import (
	"context"
	"errors"
)

/**
 * Interface for configuration settings
 */
type IConfig interface {
	LoadConfig(ctx context.Context, pathConfig string) error
	GetConfig(ctx context.Context) (Setting, error)
}

/**
 * Manager for configuration settings
 */
var _vIConfig IConfig

func SetConfig(vIConfig IConfig) error {
	if vIConfig == nil {
		return errors.New("instance of IConfig is nil")
	}
	if _vIConfig != nil {
		return errors.New("instance of IConfig already set")
	}
	_vIConfig = vIConfig
	return nil
}

func GetConfig() IConfig {
	return _vIConfig
}