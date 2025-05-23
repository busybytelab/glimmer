package llm

type (
	// Config for LLM service
	Config struct {
		Platform PlatformType `json:"platform"`
		OpenAI   OpenAIConfig `json:"openai"`
		Ollama   OllamaConfig `json:"ollama"`
		Cache    CacheConfig  `json:"cache"`
	}

	// OpenAIConfig holds configuration for OpenAI services
	OpenAIConfig struct {
		APIKey              string   `json:"apiKey"`
		CostPerMillionToken float64  `json:"costPerMillionToken"`
		Model               string   `json:"model"`
		BaseURL             string   `json:"baseUrl"`
		AllowedModels       []string `json:"allowedModels"` // List of models that are allowed to be used
	}

	// OllamaConfig holds configuration for Ollama services
	OllamaConfig struct {
		URL         string `json:"url"`
		FallbackURL string `json:"fallbackUrl"` // Fallback URL to use if the primary URL is unavailable
		Model       string `json:"model"`
	}

	// CacheConfig holds caching configuration for LLM responses
	CacheConfig struct {
		Enabled bool   `json:"enabled"`
		Backend string `json:"backend"` //
	}
)
