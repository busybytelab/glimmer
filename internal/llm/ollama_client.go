package llm

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/ollama/ollama/api"
	"github.com/rs/zerolog/log"
)

const defaultOllamaTimeout = 2 * time.Minute

// OllamaClient defines the interface for interacting with the Ollama API
type OllamaClient interface {
	// ChatWithModel sends a chat request to the Ollama API
	ChatWithModel(ctx context.Context, modelName string, messages []api.Message, stream bool, options map[string]interface{}) (*api.ChatResponse, error)
}

// DefaultOllamaClient is the default implementation of OllamaClient
type DefaultOllamaClient struct {
	baseURL *url.URL
	timeout time.Duration
}

// NewOllamaClient creates a new Ollama client with the given URL
func NewOllamaClient(baseURLStr string, timeout time.Duration) (OllamaClient, error) {
	if timeout == 0 {
		timeout = defaultOllamaTimeout
	}

	baseURL, err := url.Parse(baseURLStr)
	if err != nil {
		return nil, fmt.Errorf("invalid Ollama URL: %w", err)
	}

	return &DefaultOllamaClient{
		baseURL: baseURL,
		timeout: timeout,
	}, nil
}

// createAPIClient creates a new Ollama API client
func (c *DefaultOllamaClient) createAPIClient() *api.Client {
	transport := &http.Transport{
		DisableKeepAlives: false,
		MaxIdleConns:      100,
		IdleConnTimeout:   90 * time.Second,
	}

	return api.NewClient(c.baseURL, &http.Client{
		Timeout:   c.timeout,
		Transport: transport,
	})
}

// ChatWithModel sends a chat request to the Ollama API
func (c *DefaultOllamaClient) ChatWithModel(ctx context.Context, modelName string, messages []api.Message, stream bool, options map[string]interface{}) (*api.ChatResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	apiClient := c.createAPIClient()
	req := &api.ChatRequest{
		Model:    modelName,
		Messages: messages,
		Stream:   &stream,
		Options:  options,
	}

	log.Debug().
		Str("model", modelName).
		Int("messages", len(messages)).
		Bool("stream", stream).
		Msg("Sending chat request to Ollama")

	var finalResponse *api.ChatResponse
	err := apiClient.Chat(ctx, req, func(response api.ChatResponse) error {
		if response.Done {
			finalResponse = &response
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to chat with Ollama: %w", err)
	}

	if finalResponse == nil {
		return nil, fmt.Errorf("no response received from Ollama")
	}

	log.Debug().
		Str("model", modelName).
		Bool("done", finalResponse.Done).
		Msg("Received chat response from Ollama")

	return finalResponse, nil
}
