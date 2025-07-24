package kms

import (
	"context"
	"fmt"
	"time"

	"github.com/lamassuiot/lamassuiot/v4/pkg/kms"
	"github.com/lamassuiot/lamassuiot/v4/pkg/models"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/logger"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/resources"
)

type KMSServiceBackend struct {
	logger     *logger.Logger
	kmsStorage KMSRepository
}

type KMSServiceBuilder struct {
	Logger     *logger.Logger
	KMSStorage KMSRepository
}

func NewKMSService(builder KMSServiceBuilder) kms.KMSService {
	svc := KMSServiceBackend{
		logger:     builder.Logger,
		kmsStorage: builder.KMSStorage,
	}

	return &svc
}

func (svc *KMSServiceBackend) CreateKMSKey(ctx context.Context, input kms.CreateKMSInput) error {
	kmsKey, err := svc.kmsStorage.Insert(ctx, &models.KMSKey{
		Alias:      input.Alias,
		Algorithm:  input.Algorithm,
		Size:       input.Size,
		PublicKey:  input.PublicKey,
		CreationTS: time.Now(),
		Metadata:   map[string]any{},
	})

	fmt.Println(kmsKey)
	svc.logger.Info("KMS key created", "name", input.Alias)

	if err != nil {
		return err
	}

	return nil
}

func (svc *KMSServiceBackend) GetKMSKeys(ctx context.Context, input kms.GetKMSKeysInput) (string, error) {
	kmsKeys := []models.KMSKey{}
	bookmark, err := svc.kmsStorage.SelectAll(ctx, resources.StorageListRequest[models.KMSKey]{
		ApplyFunc: func(item models.KMSKey) {
			kmsKeys = append(kmsKeys, item)
		},
	})
	if err != nil {
		return "", err
	}

	return bookmark, nil

}
