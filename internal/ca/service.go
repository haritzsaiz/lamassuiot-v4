package ca

import (
	"context"
	"fmt"

	"github.com/lamassuiot/lamassuiot/v4/pkg/ca"
	"github.com/lamassuiot/lamassuiot/v4/pkg/kms"
	"github.com/lamassuiot/lamassuiot/v4/pkg/models"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/logger"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/resources"
)

type CAServiceBackend struct {
	logger     *logger.Logger
	caStorage  CARepository
	kmsService kms.KMSService
}

type CAServiceBuilder struct {
	Logger     *logger.Logger
	CAStorage  CARepository
	KMSService kms.KMSService
}

func NewCAService(builder CAServiceBuilder) ca.CAService {
	svc := CAServiceBackend{
		logger:     builder.Logger,
		caStorage:  builder.CAStorage,
		kmsService: builder.KMSService,
	}

	return &svc
}

func (svc *CAServiceBackend) CreateCA(ctx context.Context, input ca.CreateCAInput) error {
	svc.kmsService.CreateKMSKey(ctx, kms.CreateKMSInput{
		Name: input.Name,
	})

	ca, err := svc.caStorage.Insert(ctx, &models.CACertificate{
		Name: input.Name,
	})

	fmt.Println(ca)
	svc.logger.Info("CA created", "name", input.Name)

	if err != nil {
		return err
	}

	return nil
}

func (svc *CAServiceBackend) GetCAs(ctx context.Context, input ca.GetCAsInput) (string, error) {
	cas := []models.CACertificate{}
	bookmark, err := svc.caStorage.SelectAll(ctx, resources.StorageListRequest[models.CACertificate]{
		ApplyFunc: func(item models.CACertificate) {
			cas = append(cas, item)
		},
	})
	if err != nil {
		return "", err
	}

	return bookmark, nil

}
