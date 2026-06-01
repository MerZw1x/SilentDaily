package digest

import (
	"context"
	"fmt"
	"strings"
	"time"

	"silent/src/internal/async/status"
	workerconfig "silent/src/internal/async/worker/config"
	"silent/src/internal/domain"
	ai "silent/src/pkg/ai/openrouter"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func DigestWorker(ctx context.Context, cfg workerconfig.WorkerConfig, teamID int, memberChatIDs map[int]int64) {
	yesterday := time.Now().AddDate(0, 0, -1)

	updates, err := cfg.DailyUpdateRepository.GetByTeamIDAndDate(cfg.Conn, teamID, yesterday)
	if err != nil || len(updates) == 0 {
		return
	}

	updateIDs := make([]int, len(updates))
	for i, u := range updates {
		updateIDs[i] = u.ID
	}

	structured, err := cfg.StructuredUpdateRepository.GetByTeamIDAndDate(cfg.Conn, teamID, updateIDs)
	if err != nil || len(structured) == 0 {
		return
	}

	members, err := cfg.MemberRepository.GetByTeamID(cfg.Conn, teamID)
	if err != nil {
		return
	}

	memberMap := make(map[int]*domain.Member)
	for _, m := range members {
		memberMap[m.ID] = m
	}

	updateMemberMap := make(map[int]int)
	for _, u := range updates {
		updateMemberMap[u.ID] = u.MemberID
	}

	apiKey, err := getLeastLoadedKey(cfg)
	if err != nil {
		return
	}
	client := ai.NewAIOpenRouter(apiKey, cfg.AiModel, cfg.OpenRouterBaseURL)

	leadDigest := buildLeadDigest(structured, memberMap, updateMemberMap, client)
	_ = cfg.AiApiRepository.IncreaseRequests(cfg.Conn, apiKey)

	digest := &domain.Digest{
		TeamID:     teamID,
		Date:       yesterday,
		LeadDigest: leadDigest,
	}
	_ = cfg.DigestRepository.Upsert(cfg.Conn, digest)

	sendLeadDigest(cfg.Bot, members, leadDigest)
	sendPersonalDigests(cfg.Bot, structured, memberMap, updateMemberMap, memberChatIDs)
}

func buildLeadDigest(
	structured []*domain.StructuredUpdate,
	memberMap map[int]*domain.Member,
	updateMemberMap map[int]int,
	client *ai.AIOpenRouter,
) string {
	var sb strings.Builder
	for _, su := range structured {
		memberID := updateMemberMap[su.DailyUpdateID]
		name := "Unknown"
		if m, ok := memberMap[memberID]; ok {
			name = m.Name
		}
		sb.WriteString(fmt.Sprintf("Разработчик: %s\n", name))
		sb.WriteString(fmt.Sprintf("Прогресс: %s\n", strings.Join(su.Progress, "; ")))
		sb.WriteString(fmt.Sprintf("Планы: %s\n", strings.Join(su.Plans, "; ")))
		if len(su.Blockers) > 0 {
			sb.WriteString(fmt.Sprintf("Блокеры: %s\n", strings.Join(su.Blockers, "; ")))
		}
		sb.WriteString("\n")
	}

	prompt := fmt.Sprintf(`Ты — ассистент тимлида. На основе отчётов команды сформируй краткий дайджест в формате 3-2-1.

Отчёты команды:
%s

Ответь СТРОГО в формате:
ТОП-3 ПРОГРЕССА:
1. ...
2. ...
3. ...

2 ГЛАВНЫХ ПЛАНА НА СЕГОДНЯ:
1. ...
2. ...

 1 КЛЮЧЕВОЙ БЛОКЕР:
...`, sb.String())

	result, err := client.SendRequest(prompt)
	if err != nil {
		return sb.String()
	}
	return result
}

func sendLeadDigest(bot *tgbotapi.BotAPI, members []*domain.Member, digest string) {
	if bot == nil {
		return
	}
	for _, m := range members {
		if m.IsLead {
			msg := tgbotapi.NewMessage(int64(m.TelegramUserID), " *Утренний дайджест команды*\n\n"+digest)
			msg.ParseMode = "Markdown"
			_, _ = bot.Send(msg)
		}
	}
}

func sendPersonalDigests(
	bot *tgbotapi.BotAPI,
	structured []*domain.StructuredUpdate,
	memberMap map[int]*domain.Member,
	updateMemberMap map[int]int,
	memberChatIDs map[int]int64,
) {
	if bot == nil {
		return
	}
	for _, su := range structured {
		memberID := updateMemberMap[su.DailyUpdateID]
		m, ok := memberMap[memberID]
		if !ok {
			continue
		}
		chatID := int64(m.TelegramUserID)

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Привет, %s! Вот твой план на сегодня:\n\n", m.Name))
		if len(su.Plans) > 0 {
			sb.WriteString("*Твои планы:*\n")
			for i, p := range su.Plans {
				sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, p))
			}
		}
		if len(su.Blockers) > 0 {
			sb.WriteString("\n*Твои блокеры:*\n")
			for _, b := range su.Blockers {
				sb.WriteString(fmt.Sprintf("• %s\n", b))
			}
		}

		msg := tgbotapi.NewMessage(chatID, sb.String())
		msg.ParseMode = "Markdown"
		_, _ = bot.Send(msg)
	}
}

func getLeastLoadedKey(cfg workerconfig.WorkerConfig) (string, error) {
	apis, err := cfg.AiApiRepository.GetAllRequestsCount(cfg.Conn, cfg.ApiKeys)
	if err != nil || len(apis) == 0 {
		if len(cfg.ApiKeys) > 0 {
			return cfg.ApiKeys[0], nil
		}
		return "", fmt.Errorf("no api keys configured")
	}
	best := apis[0]
	for _, a := range apis[1:] {
		if a.Requests < best.Requests {
			best = a
		}
	}
	return best.Token, nil
}

func buildStatus(structured []*domain.StructuredUpdate) status.CompletionStatus {
	return status.CompletionStatus{}
}
