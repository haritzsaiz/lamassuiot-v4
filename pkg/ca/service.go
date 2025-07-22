package ca

import (
	"context"
)

type CAService interface {
	CreateCA(ctx context.Context, input CreateCAInput) error
	GetCAs(ctx context.Context, input GetCAsInput) (string, error)
}
