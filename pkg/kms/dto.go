package kms

import (
	"github.com/lamassuiot/lamassuiot/v4/pkg/models"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/resources"
)

type CreateKMSRequestBody struct {
	Alias     string `json:"alias"`
	Algorithm string `json:"algorithm"`
	Size      int    `json:"size"`
	PublicKey string `json:"public_key"`
}

type GetKMSKeysResponse struct {
	resources.IterableList[models.KMSKey]
}

type GetItemsResponse[T models.KMSKey] struct {
	resources.IterableList[T]
}
