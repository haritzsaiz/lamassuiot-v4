package ca

import (
	"context"
	"fmt"

	"github.com/lamassuiot/lamassuiot/v4/pkg/ca"
	"github.com/lamassuiot/lamassuiot/v4/pkg/models"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/logger"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/resources"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
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
	ctx, span := otel.GetTracerProvider().Tracer("svc").Start(ctx, "CreateCA", oteltrace.WithAttributes(attribute.String("ca.name", input.Name)))
	defer span.End()

	ca, err := svc.caStorage.Insert(ctx, &models.CACertificate{
		ID: input.Name,
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
