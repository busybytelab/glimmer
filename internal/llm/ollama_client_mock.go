package llm

import (
	"context"

	"github.com/ollama/ollama/api"
	"github.com/stretchr/testify/mock"
)

// MockOllamaClient is a mock implementation of OllamaClient for testing
type MockOllamaClient struct {
	mock.Mock
}

// NewMockOllamaClient creates a new mock Ollama client
func NewMockOllamaClient() *MockOllamaClient {
	return &MockOllamaClient{}
}

// ChatWithModel implements the OllamaClient interface for testing
func (m *MockOllamaClient) ChatWithModel(ctx context.Context, modelName string, messages []api.Message, stream bool, options map[string]interface{}) (*api.ChatResponse, error) {
	args := m.Called(ctx, modelName, messages, stream, options)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*api.ChatResponse), args.Error(1)
}
