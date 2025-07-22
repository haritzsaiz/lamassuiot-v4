package config

type AppConfig struct {
	Logs   Logging    `mapstructure:"logs"`
	Server HttpServer `mapstructure:"server"`
}
