package context

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

const CtxKey = "context"

func WithContext(timeout time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)

		// Guardar en c.Locals para acceso posterior
		c.Locals(CtxKey, ctx)

		// Cancelar el contexto cuando termine la request
		defer cancel()

		return c.Next()
	}
}

func GetRequestContext(c *fiber.Ctx) context.Context {
	if ctx, ok := c.Locals(CtxKey).(context.Context); ok {
		return ctx
	}
	return context.Background() // fallback
}
