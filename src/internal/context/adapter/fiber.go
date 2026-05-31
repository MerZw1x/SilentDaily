package adapter

import (
	"silent/src/internal/context/abstract"

	"github.com/gofiber/fiber/v3"
)

type FiberCtxAdapter struct {
	Ctx fiber.Ctx
}

func (fiberAdapter *FiberCtxAdapter) Get(key string, defaultValue ...string) string {
	return fiberAdapter.Ctx.Get(key, defaultValue...)
}

func (fiberAdapter *FiberCtxAdapter) Status(status int) abstract.HandlerContext {
	fiberAdapter.Ctx.Status(status)
	return fiberAdapter
}

func (fiberAdapter *FiberCtxAdapter) BindJSON(data any) error {
	return fiberAdapter.Ctx.Bind().JSON(data)
}
