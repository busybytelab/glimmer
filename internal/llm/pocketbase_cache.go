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
	usage := &Usage{
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
