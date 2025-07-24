package kms

import (
	"context"

	"github.com/lamassuiot/lamassuiot/v4/pkg/models"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/resources"
)

type KMSService interface {
	CreateKMSKey(ctx context.Context, input CreateKMSInput) error
	GetKMSKeys(ctx context.Context, input GetKMSKeysInput) (string, error)
}

type CreateKMSInput struct {
	Alias     string
	Algorithm string
	Size      int
	PublicKey string
}

type GetKMSKeysInput struct {
	QueryParameters *resources.QueryParameters

	ExhaustiveRun bool //wether to iter all elems
	ApplyFunc     func(kmsKey models.KMSKey)
}
