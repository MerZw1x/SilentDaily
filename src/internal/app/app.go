package app

import (
	"context"
	"log"
	"reflect"
	"strings"
	"time"

	"silent/src/internal/async/dispatcher"
	workerconfig "silent/src/internal/async/worker/config"
	"silent/src/internal/config"
	"silent/src/internal/db/postgres"
	bothandler "silent/src/internal/handler/bot"
	"silent/src/internal/handler/health"
	"silent/src/internal/handler/public"
	"silent/src/internal/middleware"
	"silent/src/internal/provider"
	rimpl "silent/src/internal/repository/impl"
	"silent/src/internal/service/abstract"
	simpl "silent/src/internal/service/impl"
	"silent/src/internal/validator"
	"silent/src/pkg/scheduler"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofiber/fiber/v3"
)

func Run() {
	cfg := config.Load()
	conn := postgres.NewPostgresConnection(cfg.GetDBDSN())

	// --- repositories ---
	memberRepo := rimpl.NewMemberRepository()
	teamRepo := rimpl.NewTeamRepository()
	updateRepo := rimpl.NewDailyUpdateRepository()
	structuredRepo := rimpl.NewStructuredUpdateRepository()
	digestRepo := rimpl.NewDigestRepository()
	aiApiRepo := rimpl.NewAiApiRepository()

	// --- services ---
	sp := provider.NewServiceProvider()
	sp.Register((*abstract.IUpdateService)(nil), simpl.NewUpdateService(conn, memberRepo, updateRepo))
	sp.Register((*abstract.IMemberService)(nil), simpl.NewMemberRepository(conn, memberRepo))
	sp.Register((*abstract.ITeamService)(nil), simpl.NewTeamService(conn, teamRepo))
	sp.Register((*abstract.IDigestService)(nil), simpl.NewDigestService(digestRepo))

	// --- telegram bot ---
	var bot *tgbotapi.BotAPI
	if cfg.Telegram.Token != "" {
		var err error
		bot, err = tgbotapi.NewBotAPI(cfg.Telegram.Token)
		if err != nil {
			log.Printf("[app] telegram bot init failed: %v", err)
		} else {
			log.Printf("[app] telegram bot authorized as @%s", bot.Self.UserName)
		}
	} else {
		log.Println("[app] TELEGRAM_TOKEN not set — bot disabled")
	}

	// --- async dispatcher ---
	apiKeys := strings.Split(cfg.AI.OpenRouterKey, ",")
	workerCfg := workerconfig.WorkerConfig{
		Conn:                       conn,
		DailyUpdateRepository:      updateRepo,
		StructuredUpdateRepository: structuredRepo,
		MemberRepository:           memberRepo,
		DigestRepository:           digestRepo,
		AiApiRepository:            aiApiRepo,
		Bot:                        bot,
		ApiKeys:                    apiKeys,
		AiModel:                    cfg.AI.Model,
		OpenRouterBaseURL:          cfg.AI.OpenRouterURL,
		MaxAttempts:                3,
	}

	if len(apiKeys) > 0 && apiKeys[0] != "" {
		_ = aiApiRepo.InsertIfNotExist(conn, apiKeys)
	}

	d := dispatcher.NewDispatcher(conn, updateRepo, workerCfg, 2*time.Second, 3, 5)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go d.RunWorkers(ctx)
	go d.ProcessCompletion(ctx)

	// --- digest scheduler ---
	if bot != nil {
		sched := scheduler.NewDigestScheduler(workerCfg, []int{}, map[int]int64{})
		sched.Start()
		defer sched.Stop()
	}

	// --- bot polling ---
	if bot != nil {
		updateSvc, _ := sp.Get(reflect.TypeOf((*abstract.IUpdateService)(nil)).Elem())
		memberSvc, _ := sp.Get(reflect.TypeOf((*abstract.IMemberService)(nil)).Elem())
		botRouter := bothandler.NewRouter(bot, updateSvc.(abstract.IUpdateService), memberSvc.(abstract.IMemberService))
		go botRouter.Start()
	}

	// --- HTTP server ---
	app := fiber.New(fiber.Config{
		EnableSplittingOnParsers: true,
		StructValidator:          validator.NewFiberStructValidator(),
	})

	app.Use(func(c fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Authorization")
		if c.Method() == fiber.MethodOptions {
			return c.SendStatus(fiber.StatusNoContent)
		}
		return c.Next()
	})

	app.Get("/ping", health.PingHandler)
	app.Post("/api/v1/updates", middleware.Adapt(public.SubmitProvider, sp))
	app.Post("/api/v1/members", middleware.Adapt(public.RegisterMemberHandler, sp))
	app.Post("/api/v1/teams", middleware.Adapt(public.CreateTeamHandler, sp))
	app.Get("/api/v1/digest", middleware.Adapt(public.GetDigestHandler, sp))

	log.Fatal(app.Listen(":" + cfg.Server.Port))
}
