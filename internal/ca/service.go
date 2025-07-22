package ca

import (
	"context"
	"fmt"

	"github.com/lamassuiot/lamassuiot/v4/pkg/ca"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/logger"
)

type CAServiceBackend struct {
	logger    *logger.Logger
	caStorage CARepository
}

type CAServiceBuilder struct {
	Logger    *logger.Logger
	CAStorage CARepository
}

func NewCAService(builder CAServiceBuilder) ca.CAService {
	svc := CAServiceBackend{
		logger:    builder.Logger,
		caStorage: builder.CAStorage,
	}

	return &svc
}

func (svc *CAServiceBackend) CreateCA(ctx context.Context, input ca.CreateCAInput) error {
	return fmt.Errorf("CreateCA not implemented")
}

func (svc *CAServiceBackend) GetCAs(ctx context.Context, input ca.GetCAsInput) (string, error) {
	return "", fmt.Errorf("GetCAs not implemented")
}
