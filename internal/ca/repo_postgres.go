package ca

import (
	"context"

	"github.com/lamassuiot/lamassuiot/v4/pkg/models"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/logger"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/resources"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/storage"
	"gorm.io/gorm"
)

type PostgresCAStore struct {
	db      *gorm.DB
	querier *storage.PostgresDBQuerier[models.CACertificate]
}

func NewCAPostgresRepository(log *logger.Logger, db *gorm.DB) (CARepository, error) {
	querier, err := storage.TableQuery(log, db, "cas", "id", models.CACertificate{})
	if err != nil {
		return nil, err
	}

	return &PostgresCAStore{
		db:      db,
		querier: querier,
	}, nil
}

func (db *PostgresCAStore) Insert(ctx context.Context, u *models.CACertificate) (*models.CACertificate, error) {
	return db.querier.Insert(ctx, u)
}

func (db *PostgresCAStore) SelectAll(ctx context.Context, req resources.StorageListRequest[models.CACertificate]) (string, error) {
	return db.querier.SelectAll(ctx, req.QueryParams, []storage.GormExtraOps{}, req.ExhaustiveRun, req.ApplyFunc)
}
