package ca

import (
	"github.com/lamassuiot/lamassuiot/v4/pkg/models"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/resources"
)

type GetCAsResponse struct {
	resources.IterableList[models.CACertificate]
}

type GetItemsResponse[T models.CACertificate] struct {
	resources.IterableList[T]
}

type GetCertsResponse struct {
	resources.IterableList[models.CACertificate]
}

type SignResponse struct {
	SignedData string `json:"signed_data"`
}

type VerifyResponse struct {
	Valid bool `json:"valid"`
}
