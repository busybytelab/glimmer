package llm

import (
	"context"
	"fmt"
	"strings"
	"time"

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
	defaultOpenAIModel = "gpt-4o-mini"
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
	sortModels(models)

	return models, nil
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
		Usage: &Usage{
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
		Usage: &Usage{
			LlmModelName:     model,
			CacheHit:         false,
			Cost:             cost,
			PromptTokens:     int(resp.Usage.PromptTokens),
			CompletionTokens: int(resp.Usage.CompletionTokens),
			TotalTokens:      int(resp.Usage.TotalTokens),
		},
	}, nil
}
