package ca

import (
	"fmt"

	"github.com/lamassuiot/lamassuiot/v4/pkg/ca"
	"github.com/lamassuiot/lamassuiot/v4/pkg/models"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/config"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/logger"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/storage"
)

const (
	CA_DB_NAME = "ca"
)

func AssembleCAService(conf *CAConfig) (*ca.CAService, error) {
	lSvc := logger.SetupLogger(conf.AppConfig.Logs.Level, "CA", "Service")
	lStorage := logger.SetupLogger(conf.Storage.LogLevel, "CA", "Storage")

	caStorage, err := createCAStorageInstance(lStorage, conf.Storage)
	if err != nil {
		return nil, fmt.Errorf("could not create CA storage instance: %s", err)
	}

	svc := NewCAService(CAServiceBuilder{
		Logger:    lSvc,
		CAStorage: caStorage,
	})

	return &svc, nil
}

func createCAStorageInstance(logger *logger.Logger, conf config.PluggableStorageEngine) (CARepository, error) {
	pconf, err := config.DecodeStruct[config.PostgresConfig](conf.Config)
	if err != nil {
		return nil, fmt.Errorf("could not decode storage config: %s", err)
	}

	psqlCli, err := storage.CreatePostgresDBConnection(logger, pconf, CA_DB_NAME)
	if err != nil {
		return nil, fmt.Errorf("could not create storage engine: %s", err)
	}

	psqlCli.AutoMigrate(&models.CACertificate{})

	userStorage, err := NewCAPostgresRepository(logger, psqlCli)
	if err != nil {
		return nil, err
	}

	return userStorage, nil
}
