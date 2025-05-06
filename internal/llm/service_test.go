package llm

import (
	"os"
	"testing"

	"github.com/ollama/ollama/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestLLMServiceWithRealOllama tests the LLM service using Ollama platform with memory caching
func TestLLMServiceWithRealOllama(t *testing.T) {
	// Skip tests if OLLAMA_URL is not set
	ollamaURL := os.Getenv("OLLAMA_URL")
	if ollamaURL == "" {
		t.Skip("OLLAMA_URL is not set, skipping tests")
	}

	// Create a configuration
	config := &Config{
		Platform: OllamaPlatform,
		Ollama: OllamaConfig{
			URL:   ollamaURL,
			Model: "llama3.2:1b", // Use the default model or override
		},
		Cache: CacheConfig{
			Enabled: true,
			Backend: string(MemoryCache),
		},
	}

	// Create a new service
	service := MemoryCacheService(config)

	// Test cases
	tests := []struct {
		name         string
		prompt       string
		systemPrompt string
		options      []ChatOption
	}{
		{
			name:         "Simple prompt",
			prompt:       "Hello, how are you?",
			systemPrompt: "",
			options:      nil,
		},
		{
			name:         "With system prompt",
			prompt:       "Tell me about yourself",
			systemPrompt: "You are a helpful assistant. Keep your responses short and concise.",
			options:      nil,
		},
		{
			name:         "With model option",
			prompt:       "What time is it?",
			systemPrompt: "",
			options:      []ChatOption{WithModel("llama3.2:1b")},
		},
		{
			name:         "With cache options",
			prompt:       "What's the weather today?",
			systemPrompt: "",
			options:      []ChatOption{WithCache(false, false)},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Call the Chat method
			response, usage, err := service.Chat(tc.prompt, tc.systemPrompt, tc.options...)

			// Assert no errors
			assert.NoError(t, err)

			// Assert response is not empty
			assert.NotEmpty(t, response)

			// Assert usage is present and valid
			assert.NotNil(t, usage)
			assert.Equal(t, "llama3.2:1b", usage.LlmModelName) // Ensure model name is correct
			assert.GreaterOrEqual(t, usage.PromptTokens, 1)
			assert.GreaterOrEqual(t, usage.CompletionTokens, 1)
			assert.GreaterOrEqual(t, usage.TotalTokens, usage.PromptTokens+usage.CompletionTokens)
			assert.Equal(t, float64(0), usage.Cost) // Ollama is free, so cost should be 0
		})
	}

	// Test cache hit by repeating the first request
	t.Run("Cache hit test", func(t *testing.T) {
		// Get the first test case
		tc := tests[0]

		// First call should be a cache miss
		_, firstUsage, err := service.Chat(tc.prompt, tc.systemPrompt)
		assert.NoError(t, err)
		assert.False(t, firstUsage.CacheHit)

		// Second call with same parameters should be a cache hit
		_, secondUsage, err := service.Chat(tc.prompt, tc.systemPrompt)
		assert.NoError(t, err)
		assert.True(t, secondUsage.CacheHit)
	})

	// Test models
	info := service.Info()
	assert.NotNil(t, info)
	assert.Equal(t, 1, len(info.Platforms))
	assert.Equal(t, OllamaPlatform, info.Platforms[0].Name)
	assert.Greater(t, len(info.Platforms[0].Models), 0)
}

func TestLLMServiceWithMockOllama(t *testing.T) {
	mockClient := new(MockOllamaClient)

	// Set up the expected behavior
	mockClient.On("ChatWithModel",
		mock.Anything, // context
		"llama3.2:1b", // model name
		mock.Anything, // messages
		false,         // stream
		mock.Anything, // options
	).Return(&api.ChatResponse{
		Message: api.Message{
			Role:    "assistant",
			Content: "This is a mock response",
		},
		Model: "llama3.2:1b",
		Done:  true,
	}, nil)

	config := &Config{
		Platform: OllamaPlatform,
		Ollama: OllamaConfig{
			URL:   "http://localhost:11434", // Mock URL
			Model: "llama3.2:1b",
		},
		Cache: CacheConfig{
			Enabled: true,
			Backend: string(MemoryCache),
		},
	}

	// Create service
	s := MemoryCacheService(config)

	// Cast the platform to ollamaPlatform to set the mock client
	cachedPlatform := s.(*service).platform.(*cachedPlatform)
	ollamaPlatform := cachedPlatform.delegate.(*ollamaPlatform)
	ollamaPlatform.client = mockClient

	// Test Chat
	response, usage, err := s.Chat("Hello", "You are a helpful assistant")

	assert.NoError(t, err)
	assert.Equal(t, "This is a mock response", response)
	assert.NotNil(t, usage)
	assert.Equal(t, "llama3.2:1b", usage.LlmModelName)

	mockClient.AssertExpectations(t)
}
