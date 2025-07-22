package config

import "github.com/lamassuiot/lamassuiot/v4/pkg/shared/logger"

type PluggableStorageEngine struct {
	LogLevel logger.Level `mapstructure:"log_level"`

	Provider StorageProvider        `mapstructure:"provider"`
	Config   map[string]interface{} `mapstructure:"config,remain"`
}

type StorageProvider string

const (
	Postgres StorageProvider = "postgres"
)
