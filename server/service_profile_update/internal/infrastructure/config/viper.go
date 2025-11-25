package config

import (
	"context"
	"errors"

	"github.com/spf13/viper"
	"github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/constants"
	domainConfig "github.com/youknow2509/cio_verify_face/server/service_profile_update/internal/domain/config"
)

// ViperConfig struct
type ViperConfig struct {
	settings *domainConfig.Setting `mapstructure:"setting"`
	filePath string                `mapstructure:"file_path"`
}

// GetConfig implements config.IConfig.
func (v *ViperConfig) GetConfig(ctx context.Context) (domainConfig.Setting, error) {
	settings := v.settings
	if settings == nil {
		return domainConfig.Setting{}, errors.New("config not loaded")
	}
	return *settings, nil
}

// LoadConfig implements config.IConfig.
func (v *ViperConfig) LoadConfig(ctx context.Context, filePath string) error {
	client := viper.New()
	// Set the path to the configuration file
	if err := setPathConfigFile(client, filePath); err != nil {
		return err
	}
	// Read the configuration file
	if err := client.ReadInConfig(); err != nil {
		return err
	}
	// Unmarshal the configuration into the Setting struct
	var setting domainConfig.Setting
	if err := client.Unmarshal(&setting); err != nil {
		return err
	}
	// save setting to struct
	v.settings = &setting
	v.filePath = filePath
	return nil
}

// NewViperConfig instance
func NewViperConfig() domainConfig.IConfig {
	return &ViperConfig{}
}

// ========================================
// Helper functions
// ========================================
// setPathConfigFile sets the path to the configuration file
func setPathConfigFile(client *viper.Viper, filePath string) error {
	if filePath == "" {
		filePath = constants.DEFAULT_CONFIG_FILE_PATH
	}

	// Check file extension
	valid := false
	for _, ext := range []string{".yaml", ".yml"} {
		if len(filePath) >= len(ext) && filePath[len(filePath)-len(ext):] == ext {
			valid = true
			break
		}
	}
	if !valid {
		return errors.New("config file must have .yaml or .yml extension")
	}

	client.SetConfigFile(filePath)
	client.SetConfigType("yaml")
	return nil
}
