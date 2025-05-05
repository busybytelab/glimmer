package llm

// CacheStorage interface for storing and retrieving LLM responses
type CacheStorage interface {
	// GetChatCacheKey generates a cache key for a chat request
	GetChatCacheKey(params *ChatParameters) string

	// GetChatResponse retrieves a cached chat response
	GetChatResponse(cacheKey string) (*ChatResponse, error)

	// SetChatResponse stores a chat response in the cache
	SetChatResponse(cacheKey string, params *ChatParameters, response *ChatResponse) error

	// GetDescribeImageCacheKey generates a cache key for an image description request
	GetDescribeImageCacheKey(params *DescribeImageParameters) string

	// GetDescribeImageResponse retrieves a cached image description
	GetDescribeImageResponse(cacheKey string) (*DescribeImageResponse, error)

	// SetDescribeImageResponse stores an image description in the cache
	SetDescribeImageResponse(cacheKey string, params *DescribeImageParameters, response *DescribeImageResponse) error
}
