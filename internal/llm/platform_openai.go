package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type (
	openAIPlatform struct {
		cfg *OpenAIConfig
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

	return &openAIPlatform{
		cfg: &cfg,
	}
}

// Type returns the platform type
func (o *openAIPlatform) Type() PlatformType {
	return OpenAIPlatform
}

func (o *openAIPlatform) Models() ([]*ModelInfo, error) {
	return nil, nil
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

	// Create request payload
	messages := []openaiRequestMessage{
		{
			Role:    "system",
			Content: params.SystemPrompt,
		},
		{
			Role:    "user",
			Content: params.Prompt,
		},
	}

	payload := openaiRequest{
		Model:    model,
		Messages: messages,
	}

	// Serialize the request
	requestBody, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	// Prepare the request
	url := fmt.Sprintf("%s/chat/completions", o.cfg.BaseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", o.cfg.APIKey))

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), openAITimeout)
	defer cancel()
	req = req.WithContext(ctx)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s, %s", resp.Status, string(respBody))
	}

	// Parse response
	var response openaiResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	// Ensure we have at least one choice
	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no response from API")
	}

	// Calculate cost
	cost := float64(response.Usage.TotalTokens) * o.cfg.CostPerMillionToken / 1_000_000

	log.Debug().
		Str("model", model).
		Int("promptTokens", response.Usage.PromptTokens).
		Int("completionTokens", response.Usage.CompletionTokens).
		Int("totalTokens", response.Usage.TotalTokens).
		Float64("cost", cost).
		Msg("OpenAI chat response received")

	// Create the response
	return &ChatResponse{
		Response: response.Choices[0].Message.Content,
		Usage: &Usage{
			LlmModelName:     model,
			CacheHit:         false,
			Cost:             cost,
			PromptTokens:     response.Usage.PromptTokens,
			CompletionTokens: response.Usage.CompletionTokens,
			TotalTokens:      response.Usage.TotalTokens,
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

	// Create a buffer to store the multipart form data
	var requestBody bytes.Buffer
	multipartWriter := multipart.NewWriter(&requestBody)

	// Add the model field
	if err := multipartWriter.WriteField("model", model); err != nil {
		return nil, fmt.Errorf("error writing model field: %w", err)
	}

	// Add system content if provided
	if params.SystemPrompt != "" {
		systemContent := map[string]interface{}{
			"type": "text",
			"text": params.SystemPrompt,
		}
		systemContentJSON, err := json.Marshal(systemContent)
		if err != nil {
			return nil, fmt.Errorf("error marshaling system content: %w", err)
		}
		if err := multipartWriter.WriteField("system", string(systemContentJSON)); err != nil {
			return nil, fmt.Errorf("error writing system field: %w", err)
		}
	}

	// Add the prompt field (user content)
	promptContent := "What's in this image?"
	if params.Prompt != "" {
		promptContent = params.Prompt
	}

	userContent := []map[string]interface{}{
		{
			"type": "text",
			"text": promptContent,
		},
		{
			"type": "image",
			"image": map[string]interface{}{
				"data": params.Reader,
			},
		},
	}
	userContentJSON, err := json.Marshal(userContent)
	if err != nil {
		return nil, fmt.Errorf("error marshaling user content: %w", err)
	}
	if err := multipartWriter.WriteField("user", string(userContentJSON)); err != nil {
		return nil, fmt.Errorf("error writing user field: %w", err)
	}

	// Close the multipart writer
	if err := multipartWriter.Close(); err != nil {
		return nil, fmt.Errorf("error closing multipart writer: %w", err)
	}

	// Create the HTTP request
	url := fmt.Sprintf("%s/chat/completions", o.cfg.BaseURL)
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", o.cfg.APIKey))

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), openAITimeout)
	defer cancel()
	req = req.WithContext(ctx)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s, %s", resp.Status, string(respBody))
	}

	// Parse response
	var response openaiResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	// Ensure we have at least one choice
	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no response from API")
	}

	// Calculate cost
	cost := float64(response.Usage.TotalTokens) * o.cfg.CostPerMillionToken / 1_000_000

	log.Debug().
		Str("model", model).
		Int("promptTokens", response.Usage.PromptTokens).
		Int("completionTokens", response.Usage.CompletionTokens).
		Int("totalTokens", response.Usage.TotalTokens).
		Float64("cost", cost).
		Msg("OpenAI image description received")

	// Create the response
	return &DescribeImageResponse{
		Description: response.Choices[0].Message.Content,
		Usage: &Usage{
			LlmModelName:     model,
			CacheHit:         false,
			Cost:             cost,
			PromptTokens:     response.Usage.PromptTokens,
			CompletionTokens: response.Usage.CompletionTokens,
			TotalTokens:      response.Usage.TotalTokens,
		},
	}, nil
}
