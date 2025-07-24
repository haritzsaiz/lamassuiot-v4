package controllers

import (
	"github.com/gofiber/fiber/v2"
)

type hcheckRoute struct {
	info APIServiceInfo
}

type APIServiceInfo struct {
	Version   string `json:"version"`
	BuildSHA  string `json:"build_sha"`
	BuildTime string `json:"build_time"`
}

func NewHealthCheckRoute(info APIServiceInfo) *hcheckRoute {
	return &hcheckRoute{
		info: info,
	}
}

func (r *hcheckRoute) HealthCheck(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"health":     true,
		"version":    r.info.Version,
		"build":      r.info.BuildSHA,
		"build_time": r.info.BuildTime,
	})
}
