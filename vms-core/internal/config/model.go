package config

type ServerConfig struct {
	CorsAllowedOrigins []string `yaml:"corsAllowedOrigins" env:"CORS_ALLOWED_ORIGINS" env-delim:","`
}

type DownsamplingConfig struct {
	Enabled bool `yaml:"enabled" env:"DOWNSAMPLING_ENABLED"`
}
