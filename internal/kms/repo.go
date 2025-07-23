package kms

import (
	"context"

	"github.com/lamassuiot/lamassuiot/v4/pkg/models"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/resources"
)

type KMSRepository interface {
	Insert(ctx context.Context, key *models.KMSKey) (*models.KMSKey, error)
	SelectAll(ctx context.Context, req resources.StorageListRequest[models.KMSKey]) (string, error)
}
