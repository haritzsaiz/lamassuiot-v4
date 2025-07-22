package config

import "github.com/lamassuiot/lamassuiot/v4/pkg/shared/logger"

type Logging struct {
	Level logger.Level `mapstructure:"level"`
}
