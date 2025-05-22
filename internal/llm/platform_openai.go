package llm

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/busybytelab.com/glimmer/internal/domain"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/rs/zerolog/log"
)

type (
	openAIPlatform struct {
		cfg    *OpenAIConfig
		client *openai.Client
	}

	openaiRequestMessage struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

	openaiRequest struct {
		Model    string                 `json:"model"`
		Messages []openaiRequestMessage `json:"messages"`
	}

	openaiResponseChoice struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	}

	openaiResponseUsage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	}

	openaiResponse struct {
		ID      string                 `json:"id"`
		Created int64                  `json:"created"`
		Choices []openaiResponseChoice `json:"choices"`
		Usage   openaiResponseUsage    `json:"usage"`
	}
)

const (
	openAIBaseURL      = "https://api.openai.com/v1"
	openAITimeout      = 60 * time.Second
	defaultOpenAIModel = "gpt-4.1-nano"
)

// newOpenAIPlatform creates a new OpenAI platform
func newOpenAIPlatform(cfg OpenAIConfig) Platform {
	if cfg.Model == "" {
		cfg.Model = defaultOpenAIModel
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = openAIBaseURL
	}

	client := openai.NewClient(
		option.WithAPIKey(cfg.APIKey),
		option.WithBaseURL(cfg.BaseURL),
	)

	return &openAIPlatform{
		cfg:    &cfg,
		client: &client,
	}
}

// Type returns the platform type
func (o *openAIPlatform) Type() PlatformType {
	return OpenAIPlatform
}

func (o *openAIPlatform) Models() ([]*ModelInfo, error) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), openAITimeout)
	defer cancel()

	// Fetch models from OpenAI
	resp, err := o.client.Models.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching models: %w", err)
	}

	// Convert OpenAI models to our ModelInfo format
	models := make([]*ModelInfo, 0, len(resp.Data))
	for _, model := range resp.Data {
		// Skip models that are not available for chat completions
		if !strings.Contains(model.ID, "gpt") {
			continue
		}

		models = append(models, &ModelInfo{
			Name:      model.ID,
			SizeHuman: "N/A",
			IsDefault: model.ID == o.cfg.Model,
		})
	}
	// Filter models to only include allowed models
	models = filterModels(models, o.cfg.AllowedModels)

	sortModels(models)

	return models, nil
}

// filterModels filters the models to only include allowed models
func filterModels(models []*ModelInfo, allowedModels []string) []*ModelInfo {
	filteredModels := make([]*ModelInfo, 0, len(models))
	for _, model := range models {
		if slices.Contains(allowedModels, model.Name) {
			filteredModels = append(filteredModels, model)
		}
	}
	return filteredModels
}

// Chat sends a chat request to OpenAI
func (o *openAIPlatform) Chat(params *ChatParameters) (*ChatResponse, error) {
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

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), openAITimeout)
	defer cancel()

	// Create chat completion request
	req := openai.ChatCompletionNewParams{
		Model: model,
		Messages: []openai.ChatCompletionMessageParamUnion{
			{
				OfSystem: &openai.ChatCompletionSystemMessageParam{
					Content: openai.ChatCompletionSystemMessageParamContentUnion{
						OfString: openai.String(params.SystemPrompt),
					},
				},
			},
			{
				OfUser: &openai.ChatCompletionUserMessageParam{
					Content: openai.ChatCompletionUserMessageParamContentUnion{
						OfString: openai.String(params.Prompt),
					},
				},
			},
		},
	}

	// Send the request
	resp, err := o.client.Chat.Completions.New(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error sending chat request: %w", err)
	}

	// Ensure we have at least one choice
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from API")
	}

	// Calculate cost
	cost := float64(resp.Usage.TotalTokens) * o.cfg.CostPerMillionToken / 1_000_000

	log.Debug().
		Str("model", model).
		Int64("promptTokens", resp.Usage.PromptTokens).
		Int64("completionTokens", resp.Usage.CompletionTokens).
		Int64("totalTokens", resp.Usage.TotalTokens).
		Float64("cost", cost).
		Msg("OpenAI chat response received")

	// Create the response
	return &ChatResponse{
		Response: resp.Choices[0].Message.Content,
		Usage: &domain.Usage{
			LlmModelName:     model,
			CacheHit:         false,
			Cost:             cost,
			PromptTokens:     int(resp.Usage.PromptTokens),
			CompletionTokens: int(resp.Usage.CompletionTokens),
			TotalTokens:      int(resp.Usage.TotalTokens),
		},
	}, nil
}

// ChatWithHistory sends a chat request with message history to OpenAI
func (o *openAIPlatform) ChatWithHistory(messages []*domain.ChatItem, params *ChatParameters) (*ChatResponse, error) {
	if len(messages) == 0 {
		return nil, fmt.Errorf("no messages provided")
	}

	model := params.Model
	if model == "" {
		log.Debug().Str("model", o.cfg.Model).Msg("params model is empty, using default model")
		model = o.cfg.Model
	}

	if model == "" {
		return nil, ErrModelNotSpecified
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), openAITimeout)
	defer cancel()

	// Convert domain.ChatItem to openai.ChatCompletionMessageParamUnion
	openaiMessages := make([]openai.ChatCompletionMessageParamUnion, 0, len(messages)+1) // +1 for system prompt

	// Add system prompt as first message if provided
	if params.SystemPrompt != "" {
		openaiMessages = append(openaiMessages, openai.ChatCompletionMessageParamUnion{
			OfSystem: &openai.ChatCompletionSystemMessageParam{
				Content: openai.ChatCompletionSystemMessageParamContentUnion{
					OfString: openai.String(params.SystemPrompt),
				},
			},
		})
	}

	// Add the rest of the messages
	for _, msg := range messages {
		var openAIMsg openai.ChatCompletionMessageParamUnion
		switch msg.Role {
		case domain.ChatItemRoleSystem:
			openAIMsg = openai.ChatCompletionMessageParamUnion{
				OfSystem: &openai.ChatCompletionSystemMessageParam{
					Content: openai.ChatCompletionSystemMessageParamContentUnion{
						OfString: openai.String(msg.Content),
					},
				},
			}
		case domain.ChatItemRoleUser:
			openAIMsg = openai.ChatCompletionMessageParamUnion{
				OfUser: &openai.ChatCompletionUserMessageParam{
					Content: openai.ChatCompletionUserMessageParamContentUnion{
						OfString: openai.String(msg.Content),
					},
				},
			}
		case domain.ChatItemRoleAssistant:
			openAIMsg = openai.ChatCompletionMessageParamUnion{
				OfAssistant: &openai.ChatCompletionAssistantMessageParam{
					Content: openai.ChatCompletionAssistantMessageParamContentUnion{
						OfString: openai.String(msg.Content),
					},
				},
			}
		default:
			log.Warn().Str("role", msg.Role).Msg("Unknown message role encountered while converting to OpenAI format")
			return nil, fmt.Errorf("unknown message role: %s", msg.Role)
		}
		openaiMessages = append(openaiMessages, openAIMsg)
	}

	// Create chat completion request
	req := openai.ChatCompletionNewParams{
		Model:    model,
		Messages: openaiMessages,
	}

	// Send the request
	resp, err := o.client.Chat.Completions.New(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error sending chat request: %w", err)
	}

	// Ensure we have at least one choice
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from API")
	}

	// Calculate cost
	cost := float64(resp.Usage.TotalTokens) * o.cfg.CostPerMillionToken / 1_000_000

	log.Debug().
		Str("model", model).
		Int64("promptTokens", resp.Usage.PromptTokens).
		Int64("completionTokens", resp.Usage.CompletionTokens).
		Int64("totalTokens", resp.Usage.TotalTokens).
		Float64("cost", cost).
		Int("messageCount", len(messages)).
		Bool("hasSystemPrompt", params.SystemPrompt != "").
		Msg("OpenAI chat with history response received")

	// Create the response
	return &ChatResponse{
		Response: resp.Choices[0].Message.Content,
		Usage: &domain.Usage{
			LlmModelName:     model,
			CacheHit:         false,
			Cost:             cost,
			PromptTokens:     int(resp.Usage.PromptTokens),
			CompletionTokens: int(resp.Usage.CompletionTokens),
			TotalTokens:      int(resp.Usage.TotalTokens),
		},
	}, nil
}

// DescribeImage sends an image to OpenAI for description
func (o *openAIPlatform) DescribeImage(params *DescribeImageParameters) (*DescribeImageResponse, error) {
	if params.Reader == nil {
		return nil, ErrContextMissing
	}

	model := params.Model
	if model == "" {
		model = o.cfg.Model
	}

	if model == "" {
		return nil, ErrModelNotSpecified
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), openAITimeout)
	defer cancel()

	// Create chat completion request with image
	req := openai.ChatCompletionNewParams{
		Model: model,
		Messages: []openai.ChatCompletionMessageParamUnion{
			{
				OfSystem: &openai.ChatCompletionSystemMessageParam{
					Content: openai.ChatCompletionSystemMessageParamContentUnion{
						OfString: openai.String(params.SystemPrompt),
					},
				},
			},
			{
				OfUser: &openai.ChatCompletionUserMessageParam{
					Content: openai.ChatCompletionUserMessageParamContentUnion{
						OfString: openai.String(params.Prompt),
					},
				},
			},
		},
	}

	// Send the request
	resp, err := o.client.Chat.Completions.New(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error sending image description request: %w", err)
	}

	// Ensure we have at least one choice
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from API")
	}

	// Calculate cost
	cost := float64(resp.Usage.TotalTokens) * o.cfg.CostPerMillionToken / 1_000_000

	log.Debug().
		Str("model", model).
		Int64("promptTokens", resp.Usage.PromptTokens).
		Int64("completionTokens", resp.Usage.CompletionTokens).
		Int64("totalTokens", resp.Usage.TotalTokens).
		Float64("cost", cost).
		Msg("OpenAI image description received")

	// Create the response
	return &DescribeImageResponse{
		Description: resp.Choices[0].Message.Content,
		Usage: &domain.Usage{
			LlmModelName:     model,
			CacheHit:         false,
			Cost:             cost,
			PromptTokens:     int(resp.Usage.PromptTokens),
			CompletionTokens: int(resp.Usage.CompletionTokens),
			TotalTokens:      int(resp.Usage.TotalTokens),
		},
	}, nil
}
