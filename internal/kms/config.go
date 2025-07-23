package kms

import "github.com/lamassuiot/lamassuiot/v4/pkg/shared/config"

type KMSConfig struct {
	AppConfig config.AppConfig              `mapstructure:"app"`
	Storage   config.PluggableStorageEngine `mapstructure:"storage"`
}
