package kms

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/lamassuiot/lamassuiot/v4/pkg/kms"
	"github.com/lamassuiot/lamassuiot/v4/pkg/models"
	fiber_context_mw "github.com/lamassuiot/lamassuiot/v4/pkg/shared/http/server/middleware/context"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/resources"
)

var KMSFiltrableFields = map[string]resources.FilterFieldType{}

var validate = validator.New()

type kmsHttpRoutes struct {
	svc kms.KMSService
}

func NewKMSHttpRoutes(svc kms.KMSService) *kmsHttpRoutes {
	return &kmsHttpRoutes{
		svc: svc,
	}
}

func (r *kmsHttpRoutes) CreateKMSKey(ctx *fiber.Ctx) error {
	var requestBody kms.CreateKMSInput

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

	err := r.svc.CreateKMSKey(ctx.UserContext(), requestBody)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"err": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(nil)
}

func (r *kmsHttpRoutes) GetAllKMSKeys(ctx *fiber.Ctx) error {
	queryParams := resources.FilterQuery(ctx, KMSFiltrableFields)

	kmsKeys := []models.KMSKey{}

	nextBookmark, err := r.svc.GetKMSKeys(fiber_context_mw.GetRequestContext(ctx), kms.GetKMSKeysInput{
		QueryParameters: queryParams,
		ExhaustiveRun:   false,
		ApplyFunc: func(kmsKey models.KMSKey) {
			kmsKeys = append(kmsKeys, kmsKey)
		},
	})

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"err": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(kms.GetKMSKeysResponse{
		IterableList: resources.IterableList[models.KMSKey]{
			NextBookmark: nextBookmark,
			List:         kmsKeys,
		},
	})
}
