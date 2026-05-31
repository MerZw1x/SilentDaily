package app

import (
	"log"
	"silent/src/internal/config"
	"silent/src/internal/db/postgres"
	"silent/src/internal/handler/health"
	"silent/src/internal/handler/public"
	"silent/src/internal/provider"
	"silent/src/internal/validator"

	rimpl "silent/src/internal/repository/impl"
	"silent/src/internal/service/abstract"
	simpl "silent/src/internal/service/impl"

	"github.com/gofiber/fiber/v3"
)

func Run() {
	config := config.Load()
	conn := postgres.NewPostgresConnection(config.GetDBDSN())

	serviceProvider := provider.NewServiceProvider()

	updateRepository := rimpl.NewDailyUpdateRepository()
	memberRepository := rimpl.NewMemberRepository()
	serviceProvider.Register((*abstract.IUpdateService)(nil), simpl.NewUpdateService(memberRepository, updateRepository))

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

	app.Post("/api/v1/updates", public.SubmitProvider)

	log.Fatal(app.Listen(":" + config.Server.Port))
}
