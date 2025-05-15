package llm

import (
	"errors"
	"io"
	"sort"

	"github.com/busybytelab.com/glimmer/internal/domain"
	"github.com/rs/zerolog/log"
)

type (
	PlatformType string

	// ModelInfo represents information about a specific LLM model
	ModelInfo struct {
		Name      string `json:"name"`
		SizeHuman string `json:"sizeHuman"`
		IsDefault bool   `json:"isDefault"`
	}

	DescribeImageResult struct {
		Description string        `json:"description"`
		Usage       *domain.Usage `json:"usage"`
	}

	Platform interface {
		Type() PlatformType
		Chat(params *ChatParameters) (*ChatResponse, error)
		ChatWithHistory(messages []*domain.ChatItem, params *ChatParameters) (*ChatResponse, error)
		DescribeImage(params *DescribeImageParameters) (*DescribeImageResponse, error)
		Models() ([]*ModelInfo, error)
	}

	CacheParameters struct {
		IgnoreCache  bool `json:"ignoreCache"`
		DisableCache bool
	}

	ChatParameters struct {
		Prompt       string           `json:"prompt"`
		SystemPrompt string           `json:"systemPrompt"`
		Model        string           `json:"model"`
		Cache        *CacheParameters `json:"cache"`
	}

	ChatResponse struct {
		Response string
		Usage    *domain.Usage
	}

	DescribeImageParameters struct {
		ChatParameters
		Reader   io.Reader
		FileName string
	}

	DescribeImageResponse struct {
		Description string
		Usage       *domain.Usage
	}
)

const (
	OpenAIPlatform PlatformType = "openai"
	EchoPlatform   PlatformType = "echo"
	OllamaPlatform PlatformType = "ollama"
)

var (
	ErrPlatformNotImplemented = errors.New("platform not implemented")
	ErrModelNotSpecified      = errors.New("model not specified")
	ErrPromptEmpty            = errors.New("prompt cannot be empty")
	ErrContextMissing         = errors.New("context missing required values")
)

// NewPlatform creates a new LLM platform based on the configuration
func NewPlatform(cfg *Config, cacheStorage CacheStorage) Platform {
	var platform Platform

	switch cfg.Platform {
	case OpenAIPlatform:
		platform = newOpenAIPlatform(cfg.OpenAI)
	case OllamaPlatform:
		platform = newOllamaPlatform(cfg.Ollama)
	case EchoPlatform:
		platform = newEchoPlatform()
	default:
		log.Fatal().Msgf("Unknown platform: %s", cfg.Platform)
		return nil
	}

	// Wrap with cache if enabled
	if cfg.Cache.Enabled && cacheStorage != nil {
		log.Info().Msg("Using provided cache storage for LLM")
		return newCachedPlatform(platform, cacheStorage)
	}

	return platform
}

func sortModels(models []*ModelInfo) {
	sort.Slice(models, func(i, j int) bool {
		return models[i].Name < models[j].Name
	})
}
