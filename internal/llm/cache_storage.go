package llm

import "github.com/busybytelab.com/glimmer/internal/domain"

// CacheStorage interface for storing and retrieving LLM responses
type CacheStorage interface {
	// GetChatCacheKey generates a cache key for a chat request
	GetChatCacheKey(params *ChatParameters) string

	// GetChatResponse retrieves a cached chat response
	GetChatResponse(cacheKey string) (*ChatResponse, error)

	// SetChatResponse stores a chat response in the cache
	SetChatResponse(cacheKey string, params *ChatParameters, response *ChatResponse) error

	// GetChatWithHistoryCacheKey generates a cache key for a chat request with message history
	GetChatWithHistoryCacheKey(messages []*domain.ChatItem, systemPrompt, model string) string

	// GetChatWithHistoryResponse retrieves a cached chat with history response
	GetChatWithHistoryResponse(cacheKey string) (*ChatResponse, error)

	// SetChatWithHistoryResponse stores a chat with history response in the cache
	SetChatWithHistoryResponse(cacheKey string, messages []*domain.ChatItem, systemPrompt, model string, response *ChatResponse) error

	// GetDescribeImageCacheKey generates a cache key for an image description request
	GetDescribeImageCacheKey(params *DescribeImageParameters) string

	// GetDescribeImageResponse retrieves a cached image description
	GetDescribeImageResponse(cacheKey string) (*DescribeImageResponse, error)

	// SetDescribeImageResponse stores an image description in the cache
	SetDescribeImageResponse(cacheKey string, params *DescribeImageParameters, response *DescribeImageResponse) error
}
