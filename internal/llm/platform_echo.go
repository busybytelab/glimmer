package llm

import (
	"fmt"
	"strings"
	"time"

	"github.com/busybytelab.com/glimmer/internal/domain"
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

func (e *echoPlatform) Models() ([]*ModelInfo, error) {
	return nil, nil
}

func (e *echoPlatform) Chat(params *ChatParameters) (*ChatResponse, error) {
	// Just echo back the prompt and system prompt
	response := fmt.Sprintf("Echo response to: %s\nSystem context: %s", params.Prompt, params.SystemPrompt)

	return &ChatResponse{
		Response: response,
		Usage: &domain.Usage{
			LlmModelName:     "echo",
			CacheHit:         false,
			Cost:             0,
			PromptTokens:     len(params.Prompt) / 4,
			CompletionTokens: len(response) / 4,
			TotalTokens:      (len(params.Prompt) + len(response)) / 4,
		},
	}, nil
}

// ChatWithHistory echoes back information about the chat history
func (e *echoPlatform) ChatWithHistory(messages []*domain.ChatItem, params *ChatParameters) (*ChatResponse, error) {
	// Count messages by role
	userCount, assistantCount, systemCount := 0, 0, 0

	lastUserMessage := ""

	for _, msg := range messages {
		switch msg.Role {
		case domain.ChatItemRoleUser:
			userCount++
			lastUserMessage = msg.Content
		case domain.ChatItemRoleAssistant:
			assistantCount++
		case domain.ChatItemRoleSystem:
			systemCount++
		}
	}

	// Build response
	responseParts := []string{
		fmt.Sprintf("Echo response to message history with %d total messages:", len(messages)),
		fmt.Sprintf("- %d user messages", userCount),
		fmt.Sprintf("- %d assistant messages", assistantCount),
		fmt.Sprintf("- %d system messages", systemCount),
	}

	// Add system prompt info
	if params.SystemPrompt != "" {
		responseParts = append(responseParts, fmt.Sprintf("\nSystem prompt: %s", params.SystemPrompt))
	}

	// Add info about the last user message
	if lastUserMessage != "" {
		if len(lastUserMessage) > 50 {
			lastUserMessage = lastUserMessage[:50] + "..."
		}
		responseParts = append(responseParts, fmt.Sprintf("\nLast user message: %s", lastUserMessage))
	}

	response := strings.Join(responseParts, "\n")

	return &ChatResponse{
		Response: response,
		Usage: &domain.Usage{
			LlmModelName:     "echo",
			CacheHit:         false,
			Cost:             0,
			PromptTokens:     len(messages) * 10, // Rough estimate
			CompletionTokens: len(response) / 4,
			TotalTokens:      (len(messages) * 10) + (len(response) / 4),
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
		Usage: &domain.Usage{
			LlmModelName:     "echo-model",
			CacheHit:         false,
			Cost:             0,
			PromptTokens:     promptTokens,
			CompletionTokens: completionTokens,
			TotalTokens:      totalTokens,
		},
	}, nil
}
