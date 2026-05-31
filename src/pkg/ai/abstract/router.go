package ai

type IAIOpenRouter interface {
	SendRequest(requstData string) (string, error)
	UseApiKey(apiKey string) error
}
