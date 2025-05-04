package app

import (
	"os"

	"github.com/rs/zerolog/log"
)

type (
	DBConfig struct {
		AutoMigrate bool `json:"autoMigrate"`
	}

	OllamaConfig struct {
		URL string `json:"url"`
	}
	LLMConfig struct {
		OllamaConfig OllamaConfig `json:"ollama"`
	}

	Config struct {
		DB  DBConfig  `json:"db"`
		LLM LLMConfig `json:"llm"`
	}
)

// EnvConfig returns the config from the environment variables
func EnvConfig() *Config {
	autoMigrate := true

	if os.Getenv("DB_DISABLE_AUTO_MIGRATE") == "true" {
		autoMigrate = false
	}

	log.Info().Bool("autoMigrate", autoMigrate).Str("DB_DISABLE_AUTO_MIGRATE", os.Getenv("DB_DISABLE_AUTO_MIGRATE")).Msg("Auto-migration")

	return &Config{
		DB: DBConfig{AutoMigrate: autoMigrate},
		LLM: LLMConfig{
			OllamaConfig: OllamaConfig{
				URL: os.Getenv("OLLAMA_URL"),
			},
		},
	}
}
