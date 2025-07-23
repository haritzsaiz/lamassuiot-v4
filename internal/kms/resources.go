package kms

import (
	"github.com/lamassuiot/lamassuiot/v4/pkg/models"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/resources"
)

type GetKMSKeysResponse struct {
	resources.IterableList[models.KMSKey]
}

type GetItemsResponse[T models.KMSKey] struct {
	resources.IterableList[T]
}
