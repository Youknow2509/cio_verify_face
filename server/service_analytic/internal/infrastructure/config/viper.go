package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	domainConfig "github.com/youknow2509/cio_verify_face/server/service_analytic/internal/domain/config"
)

// LoadConfig loads configuration from file
func LoadConfig() (*domainConfig.Config, error) {
	// Get config path from environment or use default
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config"
	}

	configName := os.Getenv("CONFIG_NAME")
	if configName == "" {
		configName = "config.dev"
	}

	// Setup viper
	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)
	viper.AddConfigPath(".")

	// Read environment variables
	viper.AutomaticEnv()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Unmarshal config
	var config domainConfig.Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}
