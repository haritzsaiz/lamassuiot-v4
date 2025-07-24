package ca

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/lamassuiot/lamassuiot/v4/pkg/ca"
	"github.com/lamassuiot/lamassuiot/v4/pkg/models"
	fiber_context_mw "github.com/lamassuiot/lamassuiot/v4/pkg/shared/http/server/middleware/context"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/resources"
)

var CAFiltrableFields = map[string]resources.FilterFieldType{
	"id":                   resources.StringFilterFieldType,
	"level":                resources.NumberFilterFieldType,
	"type":                 resources.EnumFilterFieldType,
	"serial_number":        resources.StringFilterFieldType,
	"status":               resources.EnumFilterFieldType,
	"engine_id":            resources.StringFilterFieldType,
	"valid_to":             resources.DateFilterFieldType,
	"valid_from":           resources.DateFilterFieldType,
	"revocation_timestamp": resources.DateFilterFieldType,
	"revocation_reason":    resources.EnumFilterFieldType,
	"subject.common_name":  resources.StringFilterFieldType,
	"subject_key_id":       resources.StringFilterFieldType,
}

var CARequestFiltrableFields = map[string]resources.FilterFieldType{
	"id":                  resources.StringFilterFieldType,
	"level":               resources.NumberFilterFieldType,
	"status":              resources.EnumFilterFieldType,
	"engine_id":           resources.StringFilterFieldType,
	"subject_common_name": resources.StringFilterFieldType,
	"issuer_metadata_id":  resources.StringFilterFieldType,
}

var validate = validator.New()

type caHttpRoutes struct {
	svc ca.CAService
}

func NewCAHttpRoutes(svc ca.CAService) *caHttpRoutes {
	return &caHttpRoutes{
		svc: svc,
	}
}

func (r *caHttpRoutes) CreateCA(ctx *fiber.Ctx) error {
	var requestBody ca.CreateCAInput

	if err := ctx.BodyParser(&requestBody); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"err": err.Error()})
	}

	if err := validate.Struct(&requestBody); err != nil {
		errs := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			errs[e.Field()] = e.Tag()
		}
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"errors": errs})
	}

	err := r.svc.CreateCA(ctx.UserContext(), requestBody)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"err": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(nil)
}

func (r *caHttpRoutes) GetAllCAs(ctx *fiber.Ctx) error {
	queryParams := resources.FilterQuery(ctx, CAFiltrableFields)

	cas := []models.CACertificate{}

	nextBookmark, err := r.svc.GetCAs(fiber_context_mw.GetRequestContext(ctx), ca.GetCAsInput{
		QueryParameters: queryParams,
		ExhaustiveRun:   false,
		ApplyFunc: func(ca models.CACertificate) {
			cas = append(cas, ca)
		},
	})

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"err": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(GetCAsResponse{
		IterableList: resources.IterableList[models.CACertificate]{
			NextBookmark: nextBookmark,
			List:         cas,
		},
	})
}
