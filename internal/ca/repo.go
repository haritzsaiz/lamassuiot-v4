package ca

import (
	"context"

	"github.com/lamassuiot/lamassuiot/v4/pkg/models"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/resources"
)

type CARepository interface {
	Insert(ctx context.Context, cert *models.CACertificate) (*models.CACertificate, error)
	SelectAll(ctx context.Context, req resources.StorageListRequest[models.CACertificate]) (string, error)
}
