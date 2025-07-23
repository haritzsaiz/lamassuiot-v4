package kms

import (
	"context"
)

type KMSService interface {
	CreateKMSKey(ctx context.Context, input CreateKMSInput) error
	GetKMSKeys(ctx context.Context, input GetKMSKeysInput) (string, error)
}
