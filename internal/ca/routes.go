package ca

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lamassuiot/lamassuiot/v4/pkg/ca"
)

func NewCAHTTPLayer(parentRouterGroup *fiber.Router, svc ca.CAService) {
	routes := NewCAHttpRoutes(svc)

	router := parentRouterGroup
	rv1 := (*router).Group("/v1")

	rv1.Get("/ca", routes.GetAllCAs)
	rv1.Post("/ca", routes.CreateCA)
}
