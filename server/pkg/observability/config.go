package observability

// Config holds configuration for observability
type Config struct {
	// Service information
	ServiceName string `mapstructure:"service_name" yaml:"service_name"`
	Environment string `mapstructure:"environment" yaml:"environment"`
	Version     string `mapstructure:"version" yaml:"version"`

	// Prometheus configuration
	Prometheus PrometheusConfig `mapstructure:"prometheus" yaml:"prometheus"`

	// Jaeger/Tracing configuration
	Tracing TracingConfigYAML `mapstructure:"tracing" yaml:"tracing"`
}

// PrometheusConfig holds Prometheus configuration
type PrometheusConfig struct {
	Enabled     bool   `mapstructure:"enabled" yaml:"enabled"`
	MetricsPath string `mapstructure:"metrics_path" yaml:"metrics_path"`
	Port        int    `mapstructure:"port" yaml:"port"`
}

// TracingConfigYAML holds tracing configuration from YAML
type TracingConfigYAML struct {
	Enabled  bool   `mapstructure:"enabled" yaml:"enabled"`
	Endpoint string `mapstructure:"endpoint" yaml:"endpoint"`
}

// DefaultConfig returns default configuration
func DefaultConfig() Config {
	return Config{
		ServiceName: "unknown-service",
		Environment: "development",
		Version:     "0.0.1",
		Prometheus: PrometheusConfig{
			Enabled:     true,
			MetricsPath: "/metrics",
			Port:        9090,
		},
		Tracing: TracingConfigYAML{
			Enabled:  true,
			Endpoint: "http://jaeger:4318/v1/traces",
		},
	}
}
