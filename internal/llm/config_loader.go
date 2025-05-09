package llm

import (
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

// CacheBackendType defines the supported cache backend types
type CacheBackendType string

const (
	// MemoryCache uses in-memory storage for caching
	MemoryCache CacheBackendType = "memory"

	// PocketBaseCache uses PocketBase database for caching
	PocketBaseCache CacheBackendType = "pocketbase"
)

// LoadConfig loads LLM configuration from environment variables
func LoadConfig() *Config {
	// Set default configuration
	config := &Config{
		Platform: OllamaPlatform,
		OpenAI: OpenAIConfig{
			Model:               "gpt-4o-mini",
			CostPerMillionToken: 0.15, // Approximate cost for GPT-4o mini
			BaseURL:             "https://api.openai.com/v1",
		},
		Ollama: OllamaConfig{
			Model: "gemma3:1b",
			URL:   "http://localhost:11434",
		},
		Cache: CacheConfig{
			Enabled: true,
			Backend: string(MemoryCache),
		},
	}

	// Override with environment variables if provided
	if platform := os.Getenv("LLM_PLATFORM"); platform != "" {
		switch strings.ToLower(platform) {
		case "openai":
			config.Platform = OpenAIPlatform
		case "ollama":
			config.Platform = OllamaPlatform
		case "echo":
			config.Platform = EchoPlatform
		default:
			log.Warn().Str("platform", platform).Msg("Unknown LLM platform, defaulting to Ollama")
		}
	}

	// OpenAI configuration
	if apiKey := os.Getenv("OPENAI_API_KEY"); apiKey != "" {
		config.OpenAI.APIKey = apiKey
	}

	if model := os.Getenv("OPENAI_MODEL"); model != "" {
		config.OpenAI.Model = model
	}

	if baseURL := os.Getenv("OPENAI_BASE_URL"); baseURL != "" {
		config.OpenAI.BaseURL = baseURL
	}

	// Ollama configuration
	if url := os.Getenv("OLLAMA_URL"); url != "" {
		config.Ollama.URL = url
	}

	if fallbackURL := os.Getenv("OLLAMA_FALLBACK_URL"); fallbackURL != "" {
		config.Ollama.FallbackURL = fallbackURL
	}

	if model := os.Getenv("OLLAMA_MODEL"); model != "" {
		config.Ollama.Model = model
	}

	// Cache configuration
	if cacheEnabled := os.Getenv("LLM_CACHE_ENABLED"); cacheEnabled != "" {
		config.Cache.Enabled = cacheEnabled != "false" && cacheEnabled != "0"
	}

	if cacheBackend := os.Getenv("LLM_CACHE_BACKEND"); cacheBackend != "" {
		switch strings.ToLower(cacheBackend) {
		case "memory":
			config.Cache.Backend = string(MemoryCache)
		case "pocketbase":
			config.Cache.Backend = string(PocketBaseCache)
		default:
			log.Warn().Str("cacheBackend", cacheBackend).Msg("Unknown cache backend, defaulting to memory")
			config.Cache.Backend = string(MemoryCache)
		}
	}

	log.Info().
		Str("platform", string(config.Platform)).
		Str("ollamaURL", config.Ollama.URL).
		Str("ollamaModel", config.Ollama.Model).
		Str("openaiModel", config.OpenAI.Model).
		Bool("cacheEnabled", config.Cache.Enabled).
		Str("cacheBackend", config.Cache.Backend).
		Msg("LLM configuration loaded")

	return config
}
