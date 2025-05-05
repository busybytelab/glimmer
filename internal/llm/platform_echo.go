package llm

import (
	"fmt"
	"time"
)

// echoPlatform is a simple platform that echoes the prompt back
// It's useful for testing without calling external APIs
type echoPlatform struct{}

// newEchoPlatform creates a new echo platform
func newEchoPlatform() Platform {
	return &echoPlatform{}
}

// Type returns the platform type
func (e *echoPlatform) Type() PlatformType {
	return EchoPlatform
}

// Chat simply echoes the prompt back with some metadata
func (e *echoPlatform) Chat(params *ChatParameters) (*ChatResponse, error) {
	if params.Prompt == "" {
		return nil, ErrPromptEmpty
	}

	// Create a simple echo response
	response := fmt.Sprintf("[ECHO] Prompt: %s\n\nSystem: %s\n\nTimestamp: %s",
		params.Prompt,
		params.SystemPrompt,
		time.Now().Format(time.RFC3339))

	// Estimate token count
	promptTokens := estimateTokenCount(params.Prompt)
	completionTokens := estimateTokenCount(response)
	totalTokens := promptTokens + completionTokens

	return &ChatResponse{
		Response: response,
		Usage: &Usage{
			LlmModelName:     "echo-model",
			CacheHit:         false,
			Cost:             0,
			PromptTokens:     promptTokens,
			CompletionTokens: completionTokens,
			TotalTokens:      totalTokens,
		},
	}, nil
}

// DescribeImage returns a simple echo response for image description
func (e *echoPlatform) DescribeImage(params *DescribeImageParameters) (*DescribeImageResponse, error) {
	if params.Reader == nil {
		return nil, ErrContextMissing
	}

	// Create a simple echo response
	description := fmt.Sprintf("[ECHO] Image filename: %s\n\nSystem: %s\n\nPrompt: %s\n\nTimestamp: %s",
		params.FileName,
		params.SystemPrompt,
		params.Prompt,
		time.Now().Format(time.RFC3339))

	// Estimate token count
	promptTokens := estimateTokenCount(params.Prompt) + 50 // Add 50 tokens to represent the image
	completionTokens := estimateTokenCount(description)
	totalTokens := promptTokens + completionTokens

	return &DescribeImageResponse{
		Description: description,
		Usage: &Usage{
			LlmModelName:     "echo-model",
			CacheHit:         false,
			Cost:             0,
			PromptTokens:     promptTokens,
			CompletionTokens: completionTokens,
			TotalTokens:      totalTokens,
		},
	}, nil
}
