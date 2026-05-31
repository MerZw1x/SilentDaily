package config

import (
	"os"
)

type Config struct {
	DB       DBConfig
	Server   ServerConfig
	Telegram TelegramConfig
	AI       AIConfig
}

type DBConfig struct {
	Host, User, Port, Password, Name string
}

type ServerConfig struct {
	Port string
}

type TelegramConfig struct {
	Token string
}

type AIConfig struct {
	OpenRouterKey string
	OpenRouterURL string
	Model         string
}

func Load() *Config {
	dbCfg := DBConfig{
		User:     "postgres",
		Password: "",
		Name:     "postgres",
		Host:     "localhost",
		Port:     "5432",
	}

	serverCfg := ServerConfig{
		Port: "80",
	}

	tgCfg := TelegramConfig{}

	aiCfg := AIConfig{
		OpenRouterURL: "https://openrouter.ai/api/v1",
		Model:         "openrouter/free",
	}

	cfg := &Config{
		DB:       dbCfg,
		Server:   serverCfg,
		Telegram: tgCfg,
		AI:       aiCfg,
	}

	if value := os.Getenv("DATABASE_USER"); value != "" {
		cfg.DB.User = value
	}

	if value := os.Getenv("DATABASE_PASSWORD"); value != "" {
		cfg.DB.Password = value
	}

	if value := os.Getenv("DATABASE_DBNAME"); value != "" {
		cfg.DB.Name = value
	}

	if value := os.Getenv("DATABASE_HOST"); value != "" {
		cfg.DB.Host = value
	}

	if value := os.Getenv("DATABASE_PORT"); value != "" {
		cfg.DB.Port = value
	}

	if value := os.Getenv("SERVER_PORT"); value != "" {
		cfg.Server.Port = value
	}

	if value := os.Getenv("OPENROUTER_BASE_URL"); value != "" {
		cfg.AI.OpenRouterURL = value
	}

	if value := os.Getenv("OPENROUTER_API_KEYS"); value != "" {
		cfg.AI.OpenRouterKey = value
	}

	if value := os.Getenv("OPENROUTER_MODEL"); value != "" {
		cfg.AI.Model = value
	}

	if value := os.Getenv("TELEGRAM_TOKEN"); value != "" {
		cfg.Telegram.Token = value
	}

	return cfg
}

func (cfg *Config) GetDBDSN() string {
	return "host=" + cfg.DB.Host +
		" port=" + cfg.DB.Port +
		" user=" + cfg.DB.User +
		" password=" + cfg.DB.Password +
		" dbname=" + cfg.DB.Name +
		" sslmode=disable"
}
