package ca

import (
	"fmt"

	"github.com/lamassuiot/lamassuiot/v4/pkg/ca"
	"github.com/lamassuiot/lamassuiot/v4/pkg/kms"
	"github.com/lamassuiot/lamassuiot/v4/pkg/models"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/config"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/logger"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/otel"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/storage"
)

const (
	DB_NAME = "ca"
)

func AssembleCAService(conf *CAConfig) (*ca.CAService, error) {
	otel.InitTracer("CA Service")

	lSvc := logger.SetupLogger(conf.AppConfig.Logs.Level, "CA", "Service")
	lStorage := logger.SetupLogger(conf.Storage.LogLevel, "CA", "Storage")

	caStorage, err := createCAStorageInstance(lStorage, conf.Storage)
	if err != nil {
		return nil, fmt.Errorf("could not create CA storage instance: %s", err)
	}

	svc := NewCAService(CAServiceBuilder{
		Logger:     lSvc,
		CAStorage:  caStorage,
		KMSService: kms.NewKMSSdkService(),
	})

	return &svc, nil
}

func createCAStorageInstance(logger *logger.Logger, conf config.PluggableStorageEngine) (CARepository, error) {
	pconf, err := config.DecodeStruct[config.PostgresConfig](conf.Config)
	if err != nil {
		return nil, fmt.Errorf("could not decode storage config: %s", err)
	}

	psqlCli, err := storage.CreatePostgresDBConnection(logger, pconf, DB_NAME)
	if err != nil {
		return nil, fmt.Errorf("could not create storage engine: %s", err)
	}

	err = psqlCli.AutoMigrate(&models.CACertificate{})
	if err != nil {
		return nil, fmt.Errorf("could not migrate CA certificate model: %s", err)
	}

	userStorage, err := NewCAPostgresRepository(logger, psqlCli)
	if err != nil {
		return nil, err
	}

	return userStorage, nil
}
