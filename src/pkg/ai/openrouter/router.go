package ai

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type AIOpenRouter struct {
	client *OpenRouterClient
	model  string
}

func NewAIOpenRouter(apiKey, model, baseURL string) *AIOpenRouter {
	return &AIOpenRouter{
		client: NewOpenRouterClient(apiKey, baseURL),
		model:  model,
	}
}

func (r *AIOpenRouter) UseApiKey(apiKey string) error {
	if strings.EqualFold(r.client.apiKey, apiKey) {
		return errors.New("этот API ключ уже стоит, нечего менять")
	}

	r.client.apiKey = apiKey
	return nil
}

func (r *AIOpenRouter) SendRequest(requstData string) (string, error) {
	message := []Message{
		{
			Role:    "user",
			Content: requstData,
		},
	}

	response, err := r.client.ChatCompletion(r.model, message)
	if err != nil {
		if os.IsTimeout(err) {
			return "", fmt.Errorf("timeout, LLM не ответила за предоставленное ей время: %w", err)
		}
		return "", fmt.Errorf("ошибка обращения к LLM: %w", err)
	}

	return response.Choices[0].Message.Content, nil
}
