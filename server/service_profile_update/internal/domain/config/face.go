package config

// ServiceFaceSetting defines gRPC client configuration for the face service
// Fields mirror auth client knobs for consistency.
type ServiceFaceSetting struct {
	Enabled                           bool                  `mapstructure:"enabled"`
	GrpcAddr                          string                `mapstructure:"grpc_addr"`
	KeepaliveTimeMs                   int                   `mapstructure:"keepalive_time_ms"`
	KeepaliveTimeoutMs                int                   `mapstructure:"keepalive_timeout_ms"`
	KeepalivePermitWithoutCalls       bool                  `mapstructure:"keepalive_permit_without_calls"`
	Http2MaxPingsWithoutData          int                   `mapstructure:"http2_max_pings_without_data"`
	Http2MinTimeBetweenPingsMs        int                   `mapstructure:"http2_min_time_between_pings_ms"`
	Http2MinPingIntervalWithoutDataMs int                   `mapstructure:"http2_min_ping_interval_without_data_ms"`
	TLS                               ServiceFaceTLSSetting `mapstructure:"tls"`
}

type ServiceFaceTLSSetting struct {
	Enabled  bool   `mapstructure:"enabled"`
	CertFile string `mapstructure:"cert_file"`
	KeyFile  string `mapstructure:"key_file"`
}
