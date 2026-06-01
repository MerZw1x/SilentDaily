package extract

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"silent/src/internal/async/status"
	workerconfig "silent/src/internal/async/worker/config"
	"silent/src/internal/domain"
	ai "silent/src/pkg/ai/openrouter"
)

func ExtractWorker(
	ctx context.Context,
	ch chan status.CompletionStatus,
	cfg workerconfig.WorkerConfig,
	updateID int,
) {
	defer func() {
		if r := recover(); r != nil {
			ch <- status.CompletionStatus{UpdateID: updateID, Err: fmt.Errorf("panic: %v", r)}
		}
	}()

	update, err := cfg.DailyUpdateRepository.GetOneQueued(cfg.Conn)
	if err != nil || update == nil {
		ch <- status.CompletionStatus{UpdateID: updateID, Err: errors.New("update not found")}
		return
	}

	apiKey, err := getLeastLoadedKey(cfg)
	if err != nil {
		ch <- status.CompletionStatus{UpdateID: updateID, Err: err}
		return
	}

	client := ai.NewAIOpenRouter(apiKey, cfg.AiModel, cfg.OpenRouterBaseURL)
	prompt := buildExtractionPrompt(update.RawText)

	response, err := client.SendRequest(prompt)
	if err != nil {
		ch <- status.CompletionStatus{UpdateID: updateID, Err: err}
		return
	}

	structured := parseResponse(response, update.ID)
	if err = cfg.StructuredUpdateRepository.Create(cfg.Conn, structured); err != nil {
		ch <- status.CompletionStatus{UpdateID: updateID, Err: err}
		return
	}

	_ = cfg.AiApiRepository.IncreaseRequests(cfg.Conn, apiKey)

	ch <- status.CompletionStatus{UpdateID: updateID, Err: nil}
}

func getLeastLoadedKey(cfg workerconfig.WorkerConfig) (string, error) {
	apis, err := cfg.AiApiRepository.GetAllRequestsCount(cfg.Conn, cfg.ApiKeys)
	if err != nil || len(apis) == 0 {
		if len(cfg.ApiKeys) > 0 {
			return cfg.ApiKeys[0], nil
		}
		return "", errors.New("no api keys configured")
	}
	best := apis[0]
	for _, a := range apis[1:] {
		if a.Requests < best.Requests {
			best = a
		}
	}
	return best.Token, nil
}

func buildExtractionPrompt(rawText string) string {
	return fmt.Sprintf(`Ты — ассистент для анализа ежедневных отчётов разработчиков.

Проанализируй следующий текст и извлеки структурированную информацию.

Текст разработчика:
"%s"

Ответь СТРОГО в следующем формате (каждый пункт с новой строки):
PROGRESS: <что сделано сегодня, через ; если несколько>
PLANS: <что планируется завтра, через ; если несколько>
BLOCKERS: <что мешает работе, через ; если несколько, или "нет" если блокеров нет>

Пример:
PROGRESS: реализовал авторизацию; написал тесты для API
PLANS: задеплоить на staging; провести code review
BLOCKERS: нет`, rawText)
}

func parseResponse(response string, dailyUpdateID int) *domain.StructuredUpdate {
	result := &domain.StructuredUpdate{DailyUpdateID: dailyUpdateID}
	lines := strings.Split(response, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "PROGRESS:") {
			result.Progress = splitItems(strings.TrimPrefix(line, "PROGRESS:"))
		} else if strings.HasPrefix(line, "PLANS:") {
			result.Plans = splitItems(strings.TrimPrefix(line, "PLANS:"))
		} else if strings.HasPrefix(line, "BLOCKERS:") {
			val := strings.TrimSpace(strings.TrimPrefix(line, "BLOCKERS:"))
			if val != "нет" && val != "" {
				result.Blockers = splitItems(val)
			}
		}
	}
	return result
}

func splitItems(s string) []string {
	parts := strings.Split(s, ";")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
