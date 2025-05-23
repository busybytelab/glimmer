package llm

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sync"
	"time"

	"github.com/busybytelab.com/glimmer/internal/domain"
)

// MemoryCacheStorage implements the CacheStorage interface using an in-memory map
type MemoryCacheStorage struct {
	chatCache            map[string]*ChatCacheEntry
	chatWithHistoryCache map[string]*ChatCacheEntry
	imageDscCache        map[string]*ImageCacheEntry
	mutex                sync.RWMutex
}

// ChatCacheEntry represents a cached chat response
type ChatCacheEntry struct {
	Response  string
	Usage     *domain.Usage
	CreatedAt time.Time
	TTL       time.Duration // 0 means no expiration
}

// ImageCacheEntry represents a cached image description
type ImageCacheEntry struct {
	Description string
	Usage       *domain.Usage
	CreatedAt   time.Time
	TTL         time.Duration // 0 means no expiration
}

// NewMemoryCacheStorage creates a new memory-backed cache storage
func NewMemoryCacheStorage() CacheStorage {
	return &MemoryCacheStorage{
		chatCache:            make(map[string]*ChatCacheEntry),
		chatWithHistoryCache: make(map[string]*ChatCacheEntry),
		imageDscCache:        make(map[string]*ImageCacheEntry),
	}
}

// GetChatCacheKey generates a cache key for a chat request
func (m *MemoryCacheStorage) GetChatCacheKey(params *ChatParameters) string {
	hasher := sha256.New()
	hasher.Write([]byte(params.Prompt))
	hasher.Write([]byte(params.SystemPrompt))
	model := params.Model
	if model == "" {
		model = "default"
	}
	hasher.Write([]byte(model))
	return hex.EncodeToString(hasher.Sum(nil))
}

// GetChatResponse retrieves a cached chat response
func (m *MemoryCacheStorage) GetChatResponse(cacheKey string) (*ChatResponse, error) {
	m.mutex.RLock()
	entry, exists := m.chatCache[cacheKey]
	m.mutex.RUnlock()

	if !exists {
		return nil, errors.New("record not found")
	}

	// Check expiration
	if entry.TTL > 0 && time.Since(entry.CreatedAt) > entry.TTL {
		// Remove expired entry
		m.mutex.Lock()
		delete(m.chatCache, cacheKey)
		m.mutex.Unlock()
		return nil, errors.New("record expired")
	}

	// Create a copy of usage data to avoid concurrent access issues
	usageCopy := &domain.Usage{
		LlmModelName:     entry.Usage.LlmModelName,
		CacheHit:         true, // Override to indicate a cache hit
		Cost:             entry.Usage.Cost,
		PromptTokens:     entry.Usage.PromptTokens,
		CompletionTokens: entry.Usage.CompletionTokens,
		TotalTokens:      entry.Usage.TotalTokens,
	}

	return &ChatResponse{
		Response: entry.Response,
		Usage:    usageCopy,
	}, nil
}

// SetChatResponse stores a chat response in the cache
func (m *MemoryCacheStorage) SetChatResponse(cacheKey string, params *ChatParameters, response *ChatResponse) error {
	// Default TTL (24 hours)
	ttl := 24 * time.Hour

	// Create a copy of usage data to avoid concurrent access issues
	var usageCopy *domain.Usage
	if response.Usage != nil {
		usageCopy = &domain.Usage{
			LlmModelName:     response.Usage.LlmModelName,
			CacheHit:         response.Usage.CacheHit,
			Cost:             response.Usage.Cost,
			PromptTokens:     response.Usage.PromptTokens,
			CompletionTokens: response.Usage.CompletionTokens,
			TotalTokens:      response.Usage.TotalTokens,
		}
	}

	entry := &ChatCacheEntry{
		Response:  response.Response,
		Usage:     usageCopy,
		CreatedAt: time.Now(),
		TTL:       ttl,
	}

	m.mutex.Lock()
	m.chatCache[cacheKey] = entry
	m.mutex.Unlock()

	return nil
}

// GetChatWithHistoryCacheKey generates a cache key for a chat request with message history
func (m *MemoryCacheStorage) GetChatWithHistoryCacheKey(messages []*domain.ChatItem, systemPrompt, model string) string {
	hasher := sha256.New()

	// Add system prompt
	hasher.Write([]byte(systemPrompt))

	// Add model
	if model == "" {
		model = "default"
	}
	hasher.Write([]byte(model))

	// Add all messages in sequence
	for _, msg := range messages {
		hasher.Write([]byte(msg.Role))
		hasher.Write([]byte(msg.Content))
	}

	return "history-" + hex.EncodeToString(hasher.Sum(nil))
}

// GetChatWithHistoryResponse retrieves a cached chat with history response
func (m *MemoryCacheStorage) GetChatWithHistoryResponse(cacheKey string) (*ChatResponse, error) {
	m.mutex.RLock()
	entry, exists := m.chatWithHistoryCache[cacheKey]
	m.mutex.RUnlock()

	if !exists {
		return nil, errors.New("record not found")
	}

	// Check expiration
	if entry.TTL > 0 && time.Since(entry.CreatedAt) > entry.TTL {
		// Remove expired entry
		m.mutex.Lock()
		delete(m.chatWithHistoryCache, cacheKey)
		m.mutex.Unlock()
		return nil, errors.New("record expired")
	}

	// Create a copy of usage data to avoid concurrent access issues
	usageCopy := &domain.Usage{
		LlmModelName:     entry.Usage.LlmModelName,
		CacheHit:         true, // Override to indicate a cache hit
		Cost:             entry.Usage.Cost,
		PromptTokens:     entry.Usage.PromptTokens,
		CompletionTokens: entry.Usage.CompletionTokens,
		TotalTokens:      entry.Usage.TotalTokens,
	}

	return &ChatResponse{
		Response: entry.Response,
		Usage:    usageCopy,
	}, nil
}

// SetChatWithHistoryResponse stores a chat with history response in the cache
func (m *MemoryCacheStorage) SetChatWithHistoryResponse(cacheKey string, messages []*domain.ChatItem, systemPrompt, model string, response *ChatResponse) error {
	// Shorter TTL for conversation history responses (6 hours) since conversations evolve
	ttl := 6 * time.Hour

	// Create a copy of usage data to avoid concurrent access issues
	var usageCopy *domain.Usage
	if response.Usage != nil {
		usageCopy = &domain.Usage{
			LlmModelName:     response.Usage.LlmModelName,
			CacheHit:         response.Usage.CacheHit,
			Cost:             response.Usage.Cost,
			PromptTokens:     response.Usage.PromptTokens,
			CompletionTokens: response.Usage.CompletionTokens,
			TotalTokens:      response.Usage.TotalTokens,
		}
	}

	entry := &ChatCacheEntry{
		Response:  response.Response,
		Usage:     usageCopy,
		CreatedAt: time.Now(),
		TTL:       ttl,
	}

	m.mutex.Lock()
	m.chatWithHistoryCache[cacheKey] = entry
	m.mutex.Unlock()

	return nil
}

// GetDescribeImageCacheKey generates a cache key for an image description request
func (m *MemoryCacheStorage) GetDescribeImageCacheKey(params *DescribeImageParameters) string {
	// In a real implementation, you'd need to hash the image data too
	// For now, just use a unique timestamp as we don't actually cache images
	return "image-" + time.Now().Format(time.RFC3339Nano)
}

// GetDescribeImageResponse retrieves a cached image description
func (m *MemoryCacheStorage) GetDescribeImageResponse(cacheKey string) (*DescribeImageResponse, error) {
	m.mutex.RLock()
	entry, exists := m.imageDscCache[cacheKey]
	m.mutex.RUnlock()

	if !exists {
		return nil, errors.New("record not found")
	}

	// Check expiration
	if entry.TTL > 0 && time.Since(entry.CreatedAt) > entry.TTL {
		// Remove expired entry
		m.mutex.Lock()
		delete(m.imageDscCache, cacheKey)
		m.mutex.Unlock()
		return nil, errors.New("record expired")
	}

	// Create a copy of usage data to avoid concurrent access issues
	usageCopy := &domain.Usage{
		LlmModelName:     entry.Usage.LlmModelName,
		CacheHit:         true, // Override to indicate a cache hit
		Cost:             entry.Usage.Cost,
		PromptTokens:     entry.Usage.PromptTokens,
		CompletionTokens: entry.Usage.CompletionTokens,
		TotalTokens:      entry.Usage.TotalTokens,
	}

	return &DescribeImageResponse{
		Description: entry.Description,
		Usage:       usageCopy,
	}, nil
}

// SetDescribeImageResponse stores an image description in the cache
func (m *MemoryCacheStorage) SetDescribeImageResponse(cacheKey string, params *DescribeImageParameters, response *DescribeImageResponse) error {
	// Default TTL (24 hours)
	ttl := 24 * time.Hour

	// Create a copy of usage data to avoid concurrent access issues
	var usageCopy *domain.Usage
	if response.Usage != nil {
		usageCopy = &domain.Usage{
			LlmModelName:     response.Usage.LlmModelName,
			CacheHit:         response.Usage.CacheHit,
			Cost:             response.Usage.Cost,
			PromptTokens:     response.Usage.PromptTokens,
			CompletionTokens: response.Usage.CompletionTokens,
			TotalTokens:      response.Usage.TotalTokens,
		}
	}

	entry := &ImageCacheEntry{
		Description: response.Description,
		Usage:       usageCopy,
		CreatedAt:   time.Now(),
		TTL:         ttl,
	}

	m.mutex.Lock()
	m.imageDscCache[cacheKey] = entry
	m.mutex.Unlock()

	return nil
}
