package llm

import (
	"fmt"

	"github.com/busybytelab.com/glimmer/internal/domain"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/rs/zerolog/log"
)

// PocketBaseCacheStorage implements CacheStorage using PocketBase
type PocketBaseCacheStorage struct {
	app core.App
	dao *LLMResponseRecordDao
}

// NewPocketBaseCacheStorage creates a new PocketBase-backed cache storage
func NewPocketBaseCacheStorage(app core.App) CacheStorage {
	return &PocketBaseCacheStorage{
		app: app,
		dao: NewLLMResponseRecordDao(app),
	}
}

// GetChatResponse retrieves a cached response by its key
func (p *PocketBaseCacheStorage) GetChatResponse(key string) (*ChatResponse, error) {
	cached, err := p.dao.FindLLMResponseRecordByKey(key)
	if err != nil {
		return nil, err
	}

	if len(cached.Response) < 1 {
		log.Error().
			Str("cacheKey", key).
			Msg("Cached response is empty")
		return nil, fmt.Errorf("cached response is empty")
	}

	log.Debug().
		Str("cacheKey", key).
		Str("response_preview", cached.Response[:min(20, len(cached.Response))]).
		Msg("Cache hit")

	// Create usage information
	usage := &domain.Usage{
		LlmModelName:     cached.ModelName,
		CacheHit:         true,
		Cost:             0, // No cost for cache hit
		TotalTokens:      cached.TotalTokens,
		PromptTokens:     cached.PromptTokens,
		CompletionTokens: cached.CompletionTokens,
	}

	return &ChatResponse{
		Response: cached.Response,
		Usage:    usage,
	}, nil
}

// GetDescribeImageResponse retrieves a cached image description by its key
func (p *PocketBaseCacheStorage) GetDescribeImageResponse(key string) (*DescribeImageResponse, error) {
	// We don't cache image descriptions
	return nil, fmt.Errorf("cache miss")
}

// SetChatResponse stores a chat response with the given key
func (p *PocketBaseCacheStorage) SetChatResponse(key string, params *ChatParameters, response *ChatResponse) error {
	// Prepare the record
	resp := &LLMResponseRecord{
		Key:              key,
		Prompt:           params.Prompt,
		SystemPrompt:     params.SystemPrompt,
		Response:         response.Response,
		ModelName:        response.Usage.LlmModelName,
		Backend:          string(response.Usage.LlmModelName),
		PromptTokens:     response.Usage.PromptTokens,
		CompletionTokens: response.Usage.CompletionTokens,
		TotalTokens:      response.Usage.TotalTokens,
		Cost:             response.Usage.Cost,
	}

	// Save the record
	if err := p.dao.SaveLLMResponseRecord(resp); err != nil {
		log.Error().Err(err).Msg("Failed to save LLM response in pocketbase")
		return err
	}

	return nil
}

// SetDescribeImageResponse stores an image description response
func (p *PocketBaseCacheStorage) SetDescribeImageResponse(key string, params *DescribeImageParameters, response *DescribeImageResponse) error {
	// We don't cache image descriptions
	return nil
}

// GetChatCacheKey generates a cache key from the chat parameters
func (p *PocketBaseCacheStorage) GetChatCacheKey(params *ChatParameters) string {
	modelName := params.Model
	return generateCacheKey(params.Prompt, params.SystemPrompt, modelName, "")
}

// GetChatWithHistoryCacheKey generates a cache key for a chat request with message history
func (p *PocketBaseCacheStorage) GetChatWithHistoryCacheKey(messages []*domain.ChatItem, systemPrompt, model string) string {
	// Create a hash of:
	// 1. System prompt
	// 2. Model name
	// 3. All messages sequentially

	// Create a unique key based on message content
	messagesKey := ""
	for _, msg := range messages {
		// Include role and content to make key unique
		messagesKey += fmt.Sprintf("%s:%s|", msg.Role, msg.Content)
	}

	// Generate a hash that includes all the conversation history
	return generateCacheKey(messagesKey, systemPrompt, model, "history")
}

// GetChatWithHistoryResponse retrieves a cached chat with history response
func (p *PocketBaseCacheStorage) GetChatWithHistoryResponse(key string) (*ChatResponse, error) {
	// Reuse the existing GetChatResponse method since the storage mechanism is the same
	return p.GetChatResponse(key)
}

// SetChatWithHistoryResponse stores a chat with history response in the cache
func (p *PocketBaseCacheStorage) SetChatWithHistoryResponse(key string, messages []*domain.ChatItem, systemPrompt, model string, response *ChatResponse) error {
	// Find the last user message for the prompt field
	lastUserMsg := ""
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Role == domain.ChatItemRoleUser {
			lastUserMsg = messages[i].Content
			break
		}
	}

	// Create a summary prefix for the system prompt field
	// We'll include a marker to indicate this is a conversation history cache
	enhancedSystemPrompt := fmt.Sprintf("[CHAT_HISTORY] Messages: %d | %s",
		len(messages),
		systemPrompt)

	// Trim it if needed to fit in the database field
	if len(enhancedSystemPrompt) > 1000 {
		enhancedSystemPrompt = enhancedSystemPrompt[:1000]
	}

	// Prepare the record
	resp := &LLMResponseRecord{
		Key:              key,
		Prompt:           lastUserMsg,          // Store the last user message
		SystemPrompt:     enhancedSystemPrompt, // Store system prompt with some history info
		Response:         response.Response,
		ModelName:        response.Usage.LlmModelName,
		Backend:          string(response.Usage.LlmModelName),
		PromptTokens:     response.Usage.PromptTokens,
		CompletionTokens: response.Usage.CompletionTokens,
		TotalTokens:      response.Usage.TotalTokens,
		Cost:             response.Usage.Cost,
	}

	// Save the record
	if err := p.dao.SaveLLMResponseRecord(resp); err != nil {
		log.Error().Err(err).Msg("Failed to save LLM response with history in pocketbase")
		return err
	}

	return nil
}

// GetDescribeImageCacheKey generates a cache key for image description parameters
func (p *PocketBaseCacheStorage) GetDescribeImageCacheKey(params *DescribeImageParameters) string {
	modelName := params.Model
	return generateCacheKey(params.Prompt, params.SystemPrompt, modelName, "image")
}

// CleanableStorage is an extension of CacheStorage that can clean up expired entries
type CleanableStorage interface {
	CacheStorage
	// CleanExpired removes expired entries from the storage
	CleanExpired() error
}

// Ensure PocketBaseCacheStorage implements CleanableStorage
var _ CleanableStorage = (*PocketBaseCacheStorage)(nil)

// CleanExpired removes expired entries from the storage
func (p *PocketBaseCacheStorage) CleanExpired() error {
	var expired []*LLMResponseRecord

	err := p.app.RecordQuery(domain.CollectionLLMResponses).
		AndWhere(dbx.NewExp("ttl > 0 AND created + ttl < NOW()")).
		All(&expired)

	if err != nil {
		return err
	}

	for _, cache := range expired {
		if err := p.dao.DeleteLLMResponseRecord(cache); err != nil {
			return err
		}
	}

	return nil
}

// Helper function to get the min of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
