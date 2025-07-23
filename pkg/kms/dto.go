package kms

import (
	"github.com/lamassuiot/lamassuiot/v4/pkg/models"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/resources"
)

type CreateKMSInput struct {
	Name string `json:"name"`
}

type GetKMSKeysInput struct {
	QueryParameters *resources.QueryParameters

	ExhaustiveRun bool //wether to iter all elems
	ApplyFunc     func(kmsKey models.KMSKey)
}
