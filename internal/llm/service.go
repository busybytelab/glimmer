package llm

import (
	"io"

	"github.com/pocketbase/pocketbase/core"
	"github.com/rs/zerolog/log"
)

// Service provides a high-level API for interacting with LLM platforms
type (
	Service interface {
		Chat(prompt string, systemPrompt string, options ...ChatOption) (string, *Usage, error)
		DescribeImage(reader io.Reader, fileName string, prompt string, systemPrompt string) (string, *Usage, error)
		Info() Info
	}

	Info struct {
		Platforms []PlatformInfo `json:"platforms"`
	}

	// PlatformInfo represents information about an LLM platform
	PlatformInfo struct {
		Name      string       `json:"name"`
		IsDefault bool         `json:"isDefault"`
		Models    []*ModelInfo `json:"models"`
	}

	// ChatOption defines a function that can modify ChatParameters
	ChatOption func(*ChatParameters)

	service struct {
		platform Platform
		config   *Config
	}
)

// MemoryCacheService creates a new LLM service with in-memory cache
func MemoryCacheService(config *Config) Service {
	var cacheStorage CacheStorage

	// Create appropriate cache storage
	if config.Cache.Enabled {
		log.Info().Msg("Creating in-memory LLM cache storage")
		cacheStorage = NewMemoryCacheStorage()
	}

	// Create the appropriate platform based on configuration
	platform := NewPlatform(config, cacheStorage)

	return &service{
		platform: platform,
		config:   config,
	}
}

// AppService creates a new LLM service with PocketBase cache if enabled
func AppService(config *Config, app core.App) Service {
	var platform Platform

	// If using PocketBase cache and it's enabled, initialize it now
	if config.Cache.Enabled && config.Cache.Backend == string(PocketBaseCache) {
		log.Info().Msg("Creating PocketBase-backed LLM cache storage")
		cacheStorage := NewPocketBaseCacheStorage(app)

		// Create a new platform with the PocketBase cache
		platform = NewPlatform(config, cacheStorage)
	} else {
		platform = NewPlatform(config, NewMemoryCacheStorage())
	}

	return &service{
		platform: platform,
		config:   config,
	}
}

// Chat sends a chat request to the configured LLM platform
func (s *service) Chat(prompt string, systemPrompt string, options ...ChatOption) (string, *Usage, error) {
	params := &ChatParameters{
		Prompt:       prompt,
		SystemPrompt: systemPrompt,
		Model:        "", // Will use platform default
	}

	// Apply any custom options
	for _, option := range options {
		option(params)
	}

	// Send the chat request
	response, err := s.platform.Chat(params)
	if err != nil {
		return "", nil, err
	}

	log.Debug().
		Str("model", response.Usage.LlmModelName).
		Bool("cacheHit", response.Usage.CacheHit).
		Int("promptTokens", response.Usage.PromptTokens).
		Int("completionTokens", response.Usage.CompletionTokens).
		Int("totalTokens", response.Usage.TotalTokens).
		Float64("cost", response.Usage.Cost).
		Msg("Chat completion performed")

	return response.Response, response.Usage, nil
}

// DescribeImage sends an image to the configured LLM platform for description
func (s *service) DescribeImage(reader io.Reader, fileName string, prompt string, systemPrompt string) (string, *Usage, error) {
	params := &DescribeImageParameters{
		ChatParameters: ChatParameters{
			Prompt:       prompt,
			SystemPrompt: systemPrompt,
			Model:        "", // Will use platform default
		},
		Reader:   reader,
		FileName: fileName,
	}

	// Send the image description request
	response, err := s.platform.DescribeImage(params)
	if err != nil {
		return "", nil, err
	}

	log.Debug().
		Str("model", response.Usage.LlmModelName).
		Bool("cacheHit", response.Usage.CacheHit).
		Int("promptTokens", response.Usage.PromptTokens).
		Int("completionTokens", response.Usage.CompletionTokens).
		Int("totalTokens", response.Usage.TotalTokens).
		Float64("cost", response.Usage.Cost).
		Msg("Image description performed")

	return response.Description, response.Usage, nil
}

func (s *service) Info() Info {
	models, err := s.platform.Models()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get platform models")
	}

	return Info{
		Platforms: []PlatformInfo{
			{
				Name:      string(s.platform.Type()),
				IsDefault: true,
				Models:    models,
			},
		},
	}
}

// WithModel sets a specific model for the chat
func WithModel(model string) ChatOption {
	return func(params *ChatParameters) {
		params.Model = model
	}
}

// WithCache sets cache parameters for the chat
func WithCache(ignoreCache bool, disableCache bool) ChatOption {
	return func(params *ChatParameters) {
		params.Cache = &CacheParameters{
			IgnoreCache:  ignoreCache,
			DisableCache: disableCache,
		}
	}
}
