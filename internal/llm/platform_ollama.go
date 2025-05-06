package llm

import (
	"context"
	"fmt"
	"time"

	"github.com/ollama/ollama/api"
	"github.com/rs/zerolog/log"
)

type ollamaPlatform struct {
	cfg    *OllamaConfig
	client OllamaClient
}

const (
	defaultOllamaTimeout = 30 * time.Minute
	defaultOllamaModel   = "llama3.2:1b"
)

func newOllamaPlatform(cfg OllamaConfig) Platform {
	if cfg.Model == "" {
		cfg.Model = defaultOllamaModel
	}

	log.Debug().
		Str("model", cfg.Model).
		Str("url", cfg.URL).
		Msg("Creating new Ollama platform")

	// Create the Ollama client
	client, err := NewOllamaClient(cfg.URL, defaultOllamaTimeout)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create Ollama client, will attempt to create on first use")
	}

	return &ollamaPlatform{
		cfg:    &cfg,
		client: client,
	}
}

// Type returns the platform type
func (o *ollamaPlatform) Type() PlatformType {
	return OllamaPlatform
}

func (o *ollamaPlatform) Models() ([]*ModelInfo, error) {
	models, err := o.client.ListModels()
	if err != nil {
		return nil, fmt.Errorf("failed to list Ollama models: %w", err)
	}
	// find the model with name matching in the config and set its IsDefault to true
	for _, model := range models {
		if model.Name == o.cfg.Model {
			model.IsDefault = true
		}
	}

	return models, nil
}

// getClient gets the client, creating it if necessary
func (o *ollamaPlatform) getClient() (OllamaClient, error) {
	if o.client != nil {
		return o.client, nil
	}

	// Try to create the client
	client, err := NewOllamaClient(o.cfg.URL, defaultOllamaTimeout)
	if err != nil {
		return nil, fmt.Errorf("failed to create Ollama client: %w", err)
	}

	o.client = client
	return client, nil
}

// Chat sends a chat request to Ollama
func (o *ollamaPlatform) Chat(params *ChatParameters) (*ChatResponse, error) {
	if params.Prompt == "" {
		return nil, ErrPromptEmpty
	}

	model := params.Model
	if model == "" {
		model = o.cfg.Model
	}

	if model == "" {
		return nil, ErrModelNotSpecified
	}

	// Get or create the client
	client, err := o.getClient()
	if err != nil {
		// If primary URL fails and fallback is configured, try the fallback
		if o.cfg.FallbackURL != "" {
			log.Warn().
				Str("primary", o.cfg.URL).
				Str("fallback", o.cfg.FallbackURL).
				Err(err).
				Msg("Primary Ollama URL failed, attempting fallback")

			fallbackClient, fallbackErr := NewOllamaClient(o.cfg.FallbackURL, defaultOllamaTimeout)
			if fallbackErr != nil {
				return nil, fmt.Errorf("failed to create fallback client: %w", fallbackErr)
			}

			client = fallbackClient
		} else {
			return nil, err
		}
	}

	// Create messages array with user's prompt
	messages := []api.Message{
		{
			Role:    "user",
			Content: params.Prompt,
		},
	}

	// Add system prompt if provided
	if params.SystemPrompt != "" {
		// Prepend system message
		messages = append([]api.Message{
			{
				Role:    "system",
				Content: params.SystemPrompt,
			},
		}, messages...)
	}

	// Set options (if any in the future)
	options := map[string]interface{}{}

	// Non-streaming mode for now
	stream := false

	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), defaultOllamaTimeout)
	defer cancel()

	log.Debug().
		Str("model", model).
		Int("messagesCount", len(messages)).
		Bool("hasSystemPrompt", params.SystemPrompt != "").
		Msg("Sending request to Ollama")

	// Send the chat request
	resp, err := client.ChatWithModel(ctx, model, messages, stream, options)
	if err != nil {
		// If primary URL fails and fallback is configured, try the fallback
		if o.cfg.FallbackURL != "" && client != o.client {
			log.Warn().
				Str("primary", o.cfg.URL).
				Str("fallback", o.cfg.FallbackURL).
				Err(err).
				Msg("Primary Ollama URL failed, attempting fallback")

			fallbackClient, fallbackErr := NewOllamaClient(o.cfg.FallbackURL, defaultOllamaTimeout)
			if fallbackErr != nil {
				return nil, fmt.Errorf("failed to create fallback client: %w", fallbackErr)
			}

			resp, err = fallbackClient.ChatWithModel(ctx, model, messages, stream, options)
			if err != nil {
				return nil, fmt.Errorf("failed to use fallback: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to chat with Ollama: %w", err)
		}
	}

	// Get the response content
	responseText := ""
	if resp != nil {
		responseText = resp.Message.Content
	}

	// Estimate token count based on word count (rough approximation)
	// For Ollama, we don't have direct token count, so we estimate
	promptTokens := estimateTokenCount(params.Prompt)
	completionTokens := estimateTokenCount(responseText)
	totalTokens := promptTokens + completionTokens

	log.Debug().
		Str("model", model).
		Int("estimatedPromptTokens", promptTokens).
		Int("estimatedCompletionTokens", completionTokens).
		Int("estimatedTotalTokens", totalTokens).
		Str("response", responseText).
		Msg("Ollama chat response received")

	// Create the response
	return &ChatResponse{
		Response: responseText,
		Usage: &Usage{
			LlmModelName:     model,
			CacheHit:         false,
			Cost:             0, // Ollama is free, so cost is 0
			PromptTokens:     promptTokens,
			CompletionTokens: completionTokens,
			TotalTokens:      totalTokens,
		},
	}, nil
}

// DescribeImage sends an image to Ollama for description
// Note: This is a simplified implementation as Ollama may have limited image capabilities
func (o *ollamaPlatform) DescribeImage(params *DescribeImageParameters) (*DescribeImageResponse, error) {
	// For now, we'll just return a not implemented error for Ollama
	// This can be expanded in the future if Ollama adds better image capabilities
	return nil, ErrPlatformNotImplemented
}

// estimateTokenCount provides a rough estimate of tokens from text
// This is just an approximation - tokens aren't exactly words
func estimateTokenCount(text string) int {
	// A very rough approximation is 4 characters per token on average
	if text == "" {
		return 0
	}
	return len(text) / 4
}
