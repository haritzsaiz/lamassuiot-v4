package kms

import (
	"context"

	"github.com/lamassuiot/lamassuiot/v4/pkg/models"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/logger"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/resources"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/storage"
	"gorm.io/gorm"
)

type PostgresKMSStore struct {
	db      *gorm.DB
	querier *storage.PostgresDBQuerier[models.KMSKey]
}

func NewKMSPostgresRepository(log *logger.Logger, db *gorm.DB) (KMSRepository, error) {
	querier, err := storage.TableQuery(log, db, "kms", "id", models.KMSKey{})
	if err != nil {
		return nil, err
	}

	return &PostgresKMSStore{
		db:      db,
		querier: querier,
	}, nil
}

func (db *PostgresKMSStore) Insert(ctx context.Context, u *models.KMSKey) (*models.KMSKey, error) {
	return db.querier.Insert(ctx, u)
}

func (db *PostgresKMSStore) SelectAll(ctx context.Context, req resources.StorageListRequest[models.KMSKey]) (string, error) {
	return db.querier.SelectAll(ctx, req.QueryParams, []storage.GormExtraOps{}, req.ExhaustiveRun, req.ApplyFunc)
}
