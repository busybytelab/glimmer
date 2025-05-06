package llm

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/rs/zerolog/log"
)

// cachedPlatform wraps a Platform and adds caching functionality
type cachedPlatform struct {
	delegate Platform
	storage  CacheStorage
}

// cacheEntry represents a cached response
type cacheEntry struct {
	Response string
	Usage    *Usage
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
	return c.delegate.Models()
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

// generateCacheKey creates a hash from the prompt, system prompt, model name and backend
func generateCacheKey(prompt, systemPrompt string, modelName string, backend PlatformType) string {
	hasher := sha256.New()
	hasher.Write([]byte(prompt))
	hasher.Write([]byte(systemPrompt))
	hasher.Write([]byte(modelName))
	hasher.Write([]byte(string(backend)))
	return hex.EncodeToString(hasher.Sum(nil))
}
