package kms

import (
	"fmt"

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

func AssembleKMSService(conf *KMSConfig) (*kms.KMSService, error) {
	otel.InitTracer("KMS Service")

	lSvc := logger.SetupLogger(conf.AppConfig.Logs.Level, "KMS", "Service")
	lStorage := logger.SetupLogger(conf.Storage.LogLevel, "KMS", "Storage")

	kmsStorage, err := createKMSStorageInstance(lStorage, conf.Storage)
	if err != nil {
		return nil, fmt.Errorf("could not create KMS storage instance: %s", err)
	}

	svc := NewKMSService(KMSServiceBuilder{
		Logger:     lSvc,
		KMSStorage: kmsStorage,
	})

	return &svc, nil
}

func createKMSStorageInstance(logger *logger.Logger, conf config.PluggableStorageEngine) (KMSRepository, error) {
	pconf, err := config.DecodeStruct[config.PostgresConfig](conf.Config)
	if err != nil {
		return nil, fmt.Errorf("could not decode storage config: %s", err)
	}

	psqlCli, err := storage.CreatePostgresDBConnection(logger, pconf, DB_NAME)
	if err != nil {
		return nil, fmt.Errorf("could not create storage engine: %s", err)
	}

	err = psqlCli.AutoMigrate(&models.KMSKey{})
	if err != nil {
		return nil, fmt.Errorf("could not migrate KMS key model: %s", err)
	}

	store, err := NewKMSPostgresRepository(logger, psqlCli)
	if err != nil {
		return nil, err
	}

	return store, nil
}
