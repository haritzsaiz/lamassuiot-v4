package kms

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lamassuiot/lamassuiot/v4/pkg/kms"
)

func NewKMSHTTPLayer(parentRouterGroup *fiber.Router, svc kms.KMSService) {
	routes := NewKMSHttpRoutes(svc)

	router := parentRouterGroup
	rv1 := (*router).Group("/v1")

	rv1.Get("/kms", routes.GetAllKMSKeys)
	rv1.Post("/kms", routes.CreateKMSKey)
}
