package llm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMemoryCacheStorage tests the memory cache storage implementation
func TestMemoryCacheStorage(t *testing.T) {
	// Create a new memory cache storage
	storage := NewMemoryCacheStorage()

	// Test parameters and responses for cache testing
	testPrompt := "What is the capital of France?"
	testSystemPrompt := "You are a helpful assistant"
	testModel := "llama3.2:1b"
	testResponse := "Paris is the capital of France"

	// Create chat parameters for testing
	params := &ChatParameters{
		Prompt:       testPrompt,
		SystemPrompt: testSystemPrompt,
		Model:        testModel,
	}

	// Create a chat response for testing
	usage := &Usage{
		LlmModelName:     testModel,
		CacheHit:         false,
		Cost:             0.0,
		PromptTokens:     10,
		CompletionTokens: 5,
		TotalTokens:      15,
	}

	response := &ChatResponse{
		Response: testResponse,
		Usage:    usage,
	}

	// Test getting a cache key
	t.Run("GetChatCacheKey", func(t *testing.T) {
		key := storage.GetChatCacheKey(params)
		assert.NotEmpty(t, key)

		// Getting the key again with the same parameters should give the same key
		key2 := storage.GetChatCacheKey(params)
		assert.Equal(t, key, key2)

		// Modifying parameters should give a different key
		paramsModified := &ChatParameters{
			Prompt:       testPrompt + " Modified",
			SystemPrompt: testSystemPrompt,
			Model:        testModel,
		}
		keyModified := storage.GetChatCacheKey(paramsModified)
		assert.NotEqual(t, key, keyModified)
	})

	// Test setting and getting a chat response
	t.Run("SetGetChatResponse", func(t *testing.T) {
		// Generate a key for the parameters
		key := storage.GetChatCacheKey(params)

		// Initially, there should be no response in the cache
		cachedResponse, err := storage.GetChatResponse(key)
		assert.Error(t, err)
		assert.Nil(t, cachedResponse)

		// Set the response in the cache
		err = storage.SetChatResponse(key, params, response)
		assert.NoError(t, err)

		// Now we should be able to get the response from the cache
		cachedResponse, err = storage.GetChatResponse(key)
		assert.NoError(t, err)
		assert.NotNil(t, cachedResponse)

		// The cached response should match the original
		assert.Equal(t, testResponse, cachedResponse.Response)
		assert.True(t, cachedResponse.Usage.CacheHit) // Should be marked as a cache hit
		assert.Equal(t, testModel, cachedResponse.Usage.LlmModelName)
		assert.Equal(t, usage.PromptTokens, cachedResponse.Usage.PromptTokens)
		assert.Equal(t, usage.CompletionTokens, cachedResponse.Usage.CompletionTokens)
		assert.Equal(t, usage.TotalTokens, cachedResponse.Usage.TotalTokens)
	})

	// Test that a non-existent key returns an error
	t.Run("GetNonExistentKey", func(t *testing.T) {
		cachedResponse, err := storage.GetChatResponse("non-existent-key")
		assert.Error(t, err)
		assert.Nil(t, cachedResponse)
	})

	t.Run("ImageDescriptionCaching", func(t *testing.T) {
		imgParams := &DescribeImageParameters{
			ChatParameters: ChatParameters{
				Prompt:       "Describe this image",
				SystemPrompt: "You are a helpful assistant",
				Model:        testModel,
			},
			FileName: "test.jpg",
		}

		imgKey := storage.GetDescribeImageCacheKey(imgParams)
		assert.NotEmpty(t, imgKey)

		// Cache miss
		missResponse, err := storage.GetDescribeImageResponse(imgKey)
		assert.Error(t, err)
		assert.Nil(t, missResponse)

		err = storage.SetDescribeImageResponse(imgKey, imgParams, &DescribeImageResponse{
			Description: "This is a test image description",
			Usage:       usage,
		})

		assert.NoError(t, err)

		// Cache hit
		hitUsage := &Usage{
			LlmModelName:     usage.LlmModelName,
			CacheHit:         true,
			Cost:             0.0,
			PromptTokens:     usage.PromptTokens,
			CompletionTokens: usage.CompletionTokens,
			TotalTokens:      usage.TotalTokens,
		}
		expectedResponse := &DescribeImageResponse{
			Description: "This is a test image description",
			Usage:       hitUsage,
		}

		hitResponse, err := storage.GetDescribeImageResponse(imgKey)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, hitResponse)
	})
}
