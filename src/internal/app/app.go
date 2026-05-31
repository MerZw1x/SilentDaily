package app

import (
	"log"

	"backend/src/internal/config"
	"backend/src/internal/db/postgres"
	"backend/src/internal/handler/api"
	"backend/src/internal/handler/health"
	"backend/src/internal/handler/public"
	"backend/src/internal/middleware"
	"backend/src/internal/provider"
	rimpl "backend/src/internal/repository/impl"
	simpl "backend/src/internal/service/impl"
	"backend/src/internal/validator"

	sabst "backend/src/internal/service/abstract"

	"github.com/gofiber/fiber/v3"
)

func Run() {
	config := config.Load()
	conn := postgres.NewPostgresConnection(config.GetDBDSN())
	serviceProvider := provider.NewServiceProvider()

	rosstatRepo := rimpl.NewRosstatRepository()
	rosstatAgeRepo := rimpl.NewRosstatAgeRepository()
	serviceProvider.Register((*sabst.IRosstatService)(nil), simpl.NewRosstatService(conn, rosstatRepo, rosstatAgeRepo))

	geoRepo := rimpl.NewGeoRepository()
	serviceProvider.Register((*sabst.IGeoService)(nil), simpl.NewGeoService(conn, geoRepo))

	app := fiber.New(fiber.Config{
		EnableSplittingOnParsers: true,
		StructValidator:          validator.NewFiberStructValidator(),
	})

	// CORS middleware
	app.Use(func(c fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "https://midray.ru")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Authorization")

		if c.Method() == fiber.MethodOptions {
			return c.SendStatus(fiber.StatusNoContent)
		}

		return c.Next()
	})

	app.Get("/ping", health.PingHandler)

	app.Get("/api/v1/rosstat", middleware.Adapt(public.GetRosstatHandler, serviceProvider))
	app.Get("/api/v1/geo", middleware.Adapt(public.GetGeoHandler, serviceProvider))

	app.Get("/openapi.yaml", api.OpenapiYamlHandler)
	app.Get("/api/*", api.ApiHandler())

	app.Get("/api/v1/ai/report/:code", middleware.Adapt(public.GetReportAsyncHandler, serviceProvider))
	app.Get("/api/v1/ai/report/:code/request/:hash", middleware.Adapt(public.GetRequestStatusHandler, serviceProvider))
	app.Post("/api/v1/ai/report", middleware.Adapt(public.PostReportAsyncHandler, serviceProvider))

	log.Fatal(app.Listen(":80"))
}
