package llm

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/busybytelab.com/glimmer/internal/domain"
	"github.com/rs/zerolog/log"
)

// cachedPlatform wraps a Platform and adds caching functionality
type cachedPlatform struct {
	delegate Platform
	storage  CacheStorage
	models   []*ModelInfo
}

// cacheEntry represents a cached response
type cacheEntry struct {
	Response string
	Usage    *domain.Usage
}

// newCachedPlatform creates a new cached platform wrapper
func newCachedPlatform(delegate Platform, storage CacheStorage) Platform {
	return &cachedPlatform{
		delegate: delegate,
		storage:  storage,
	}
}

// Type returns the type of the delegate platform
func (c *cachedPlatform) Type() PlatformType {
	return c.delegate.Type()
}

func (c *cachedPlatform) Models() ([]*ModelInfo, error) {
	if c.models == nil {
		models, err := c.delegate.Models()
		if err != nil {
			return nil, err
		}
		c.models = models
	}
	return c.models, nil
}

func (c *cachedPlatform) ClearModelsCache() {
	c.models = nil
}

// Chat implements the Platform interface with caching
func (c *cachedPlatform) Chat(params *ChatParameters) (*ChatResponse, error) {
	// Generate cache key from parameters
	cacheKey := c.storage.GetChatCacheKey(params)

	// Check if we should use cache
	shouldUseCache := true
	if params.Cache != nil {
		if params.Cache.DisableCache {
			log.Debug().Str("cacheKey", cacheKey).Msg("Cache disabled for this request")
			shouldUseCache = false
		}
	}

	// Check if we have a cached response and should use it
	if shouldUseCache {
		if response, err := c.storage.GetChatResponse(cacheKey); err == nil {
			log.Debug().Str("cacheKey", cacheKey).Msg("Cache hit")
			// If IgnoreCache is set, ignore the cached response
			if params.Cache != nil && params.Cache.IgnoreCache {
				log.Debug().Str("cacheKey", cacheKey).Msg("Ignoring cache hit due to IgnoreCache parameter")
			} else {
				// Make sure CacheHit is set to true for cached responses
				if response.Usage != nil {
					response.Usage.CacheHit = true
				}
				return response, nil
			}
		}
	}

	log.Debug().Str("cacheKey", cacheKey).Msg("Cache miss")

	// If not cached or ignoring cache, call the delegate platform
	response, err := c.delegate.Chat(params)
	if err != nil {
		return nil, err
	}

	// Mark as not a cache hit
	if response.Usage != nil {
		response.Usage.CacheHit = false
	}

	// Store the response if caching is not disabled
	if shouldUseCache {
		if err := c.storage.SetChatResponse(cacheKey, params, response); err != nil {
			// Just log the error but don't fail the request
			log.Error().Err(err).Msg("Failed to cache LLM response")
		}
	}

	return response, nil
}

// DescribeImage implements the Platform interface with caching
func (c *cachedPlatform) DescribeImage(params *DescribeImageParameters) (*DescribeImageResponse, error) {
	// Generate cache key
	cacheKey := c.storage.GetDescribeImageCacheKey(params)

	// Check if we should use cache
	shouldUseCache := true
	if params.Cache != nil {
		if params.Cache.DisableCache {
			log.Debug().Str("cacheKey", cacheKey).Msg("Cache disabled for this request")
			shouldUseCache = false
		}
	}

	// Check if we have a cached response and should use it
	if shouldUseCache {
		if response, err := c.storage.GetDescribeImageResponse(cacheKey); err == nil {
			log.Debug().Str("cacheKey", cacheKey).Msg("Cache hit")
			// If IgnoreCache is set, ignore the cached response
			if params.Cache != nil && params.Cache.IgnoreCache {
				log.Debug().Str("cacheKey", cacheKey).Msg("Ignoring cache hit due to IgnoreCache parameter")
			} else {
				// Make sure CacheHit is set to true for cached responses
				if response.Usage != nil {
					response.Usage.CacheHit = true
				}
				return response, nil
			}
		}
	}

	// If not cached or ignoring cache, call the delegate platform
	result, err := c.delegate.DescribeImage(params)
	if err != nil {
		return nil, err
	}

	// Mark as not a cache hit
	if result.Usage != nil {
		result.Usage.CacheHit = false
	}

	// Store the response if caching is not disabled
	if shouldUseCache {
		if err := c.storage.SetDescribeImageResponse(cacheKey, params, result); err != nil {
			// Just log the error but don't fail the request
			log.Error().Err(err).Msg("Failed to cache LLM image description")
		}
	}

	return result, nil
}

// ChatWithHistory implements conversation history chat with caching
func (c *cachedPlatform) ChatWithHistory(messages []*domain.ChatItem, params *ChatParameters) (*ChatResponse, error) {
	if len(messages) == 0 {
		return nil, errors.New("no messages provided for chat with history")
	}

	// Generate cache key for this conversation history
	cacheKey := c.storage.GetChatWithHistoryCacheKey(messages, params.SystemPrompt, params.Model)

	// Check if we should use cache
	shouldUseCache := true
	if params.Cache != nil {
		if params.Cache.DisableCache {
			log.Debug().Str("cacheKey", cacheKey).Msg("Cache disabled for this chat with history request")
			shouldUseCache = false
		}
	}

	// Check if we have a cached response and should use it
	if shouldUseCache {
		if response, err := c.storage.GetChatWithHistoryResponse(cacheKey); err == nil {
			log.Debug().
				Str("cacheKey", cacheKey).
				Int("messagesCount", len(messages)).
				Bool("hasSystemPrompt", params.SystemPrompt != "").
				Msg("Cache hit for chat with history")

			// If IgnoreCache is set, ignore the cached response
			if params.Cache != nil && params.Cache.IgnoreCache {
				log.Debug().Str("cacheKey", cacheKey).Msg("Ignoring cache hit due to IgnoreCache parameter")
			} else {
				// Make sure CacheHit is set to true for cached responses
				if response.Usage != nil {
					response.Usage.CacheHit = true
				}
				return response, nil
			}
		}
	}

	log.Debug().
		Str("cacheKey", cacheKey).
		Int("messagesCount", len(messages)).
		Bool("hasSystemPrompt", params.SystemPrompt != "").
		Msg("Cache miss for chat with history")

	// If not cached or ignoring cache, call the delegate platform
	response, err := c.delegate.ChatWithHistory(messages, params)
	if err != nil {
		return nil, err
	}

	// Mark as not a cache hit
	if response.Usage != nil {
		response.Usage.CacheHit = false
	}

	// Store the response if caching is not disabled
	if shouldUseCache {
		if err := c.storage.SetChatWithHistoryResponse(cacheKey, messages, params.SystemPrompt, params.Model, response); err != nil {
			// Just log the error but don't fail the request
			log.Error().Err(err).Msg("Failed to cache LLM response with history")
		} else {
			log.Debug().
				Str("cacheKey", cacheKey).
				Int("messagesCount", len(messages)).
				Bool("hasSystemPrompt", params.SystemPrompt != "").
				Msg("Stored chat with history response in cache")
		}
	}

	return response, nil
}

// generateCacheKey creates a hash from the prompt, system prompt, model name and backend
func generateCacheKey(prompt, systemPrompt string, modelName string, backend PlatformType) string {
	hasher := sha256.New()
	hasher.Write([]byte(prompt))
	hasher.Write([]byte(systemPrompt))
	hasher.Write([]byte(modelName))
	hasher.Write([]byte(string(backend)))
	return hex.EncodeToString(hasher.Sum(nil))
}
