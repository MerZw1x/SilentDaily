package config

import (
	"silent/src/internal/db/abstract"
	repoabstract "silent/src/internal/repository/abstract"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type WorkerConfig struct {
	Conn                      abstract.IDBConnection
	DailyUpdateRepository     repoabstract.IDailyUpdateRepository
	StructuredUpdateRepository repoabstract.IStructuredUpdateRepository
	MemberRepository          repoabstract.IMemberRepository
	DigestRepository          repoabstract.IDigestRepository
	AiApiRepository           repoabstract.IAiApiRepository
	Bot                       *tgbotapi.BotAPI
	ApiKeys                   []string
	AiModel                   string
	OpenRouterBaseURL         string
	MaxAttempts               int
}
