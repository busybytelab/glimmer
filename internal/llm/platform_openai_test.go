package llm

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// skipIfNoOpenAI skips the test if OpenAI tests are disabled or API key is not set
func skipIfNoOpenAI(t *testing.T) {
	t.Helper()
	if os.Getenv("DISABLE_REAL_OPENAI_TESTS") != "" {
		t.Skip("Skipping test: OpenAI tests disabled. Unset DISABLE_REAL_OPENAI_TESTS to enable.")
	}
	if os.Getenv("OPENAI_API_KEY") == "" {
		t.Skip("Skipping test: OPENAI_API_KEY not set")
	}
}

func TestNewOpenAIPlatform(t *testing.T) {
	tests := []struct {
		name     string
		cfg      OpenAIConfig
		wantType PlatformType
	}{
		{
			name: "default config",
			cfg: OpenAIConfig{
				APIKey: "test-key",
			},
			wantType: OpenAIPlatform,
		},
		{
			name: "custom model",
			cfg: OpenAIConfig{
				APIKey: "test-key",
				Model:  "gpt-3.5-turbo",
			},
			wantType: OpenAIPlatform,
		},
		{
			name: "custom base URL",
			cfg: OpenAIConfig{
				APIKey:  "test-key",
				BaseURL: "https://custom.openai.com/v1",
			},
			wantType: OpenAIPlatform,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			platform := newOpenAIPlatform(tt.cfg)
			assert.Equal(t, tt.wantType, platform.Type())
		})
	}
}

func TestOpenAIPlatform_Chat(t *testing.T) {
	// Skip if OpenAI tests are disabled or API key is not set
	skipIfNoOpenAI(t)

	platform := newOpenAIPlatform(OpenAIConfig{
		APIKey: os.Getenv("OPENAI_API_KEY"),
		Model:  "gpt-3.5-turbo", // Use the cheapest model for testing
	})

	tests := []struct {
		name    string
		params  *ChatParameters
		wantErr bool
	}{
		{
			name: "simple chat",
			params: &ChatParameters{
				Prompt: "Say hello",
			},
			wantErr: false,
		},
		{
			name: "empty prompt",
			params: &ChatParameters{
				Prompt: "",
			},
			wantErr: true,
		},
		{
			name: "with system prompt",
			params: &ChatParameters{
				Prompt:       "What is 2+2?",
				SystemPrompt: "You are a helpful assistant that only responds with numbers.",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := platform.Chat(tt.params)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, resp.Response)
			assert.NotNil(t, resp.Usage)
			assert.Equal(t, "gpt-3.5-turbo", resp.Usage.LlmModelName)
			assert.False(t, resp.Usage.CacheHit)
			assert.Greater(t, resp.Usage.TotalTokens, 0)
		})
	}
}

func TestOpenAIPlatform_DescribeImage(t *testing.T) {
	t.Skip("Skipping test: TestOpenAIPlatform_DescribeImage disabled for now")

	// Skip if no API key is provided
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping test: OPENAI_API_KEY not set")
	}

	platform := newOpenAIPlatform(OpenAIConfig{
		APIKey: apiKey,
		Model:  "gpt-4-vision-preview", // Required for image analysis
	})

	// Create a simple test image (1x1 pixel)
	imageData := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D,
		0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4, 0x89, 0x00, 0x00, 0x00,
		0x0A, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9C, 0x63, 0x00, 0x01, 0x00, 0x00,
		0x05, 0x00, 0x01, 0x0D, 0x0A, 0x2D, 0xB4, 0x00, 0x00, 0x00, 0x00, 0x49,
		0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82,
	}

	tests := []struct {
		name    string
		params  *DescribeImageParameters
		wantErr bool
	}{
		{
			name: "simple image description",
			params: &DescribeImageParameters{
				ChatParameters: ChatParameters{
					Prompt: "What do you see in this image?",
				},
				Reader:   bytes.NewReader(imageData),
				FileName: "test.png",
			},
			wantErr: false,
		},
		{
			name: "nil reader",
			params: &DescribeImageParameters{
				ChatParameters: ChatParameters{
					Prompt: "What do you see?",
				},
				Reader:   nil,
				FileName: "test.png",
			},
			wantErr: true,
		},
		{
			name: "with system prompt",
			params: &DescribeImageParameters{
				ChatParameters: ChatParameters{
					Prompt:       "What do you see?",
					SystemPrompt: "You are a helpful assistant that describes images.",
				},
				Reader:   bytes.NewReader(imageData),
				FileName: "test.png",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := platform.DescribeImage(tt.params)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, resp.Description)
			assert.NotNil(t, resp.Usage)
			assert.Equal(t, "gpt-4-vision-preview", resp.Usage.LlmModelName)
			assert.False(t, resp.Usage.CacheHit)
			assert.Greater(t, resp.Usage.TotalTokens, 0)
		})
	}
}

func TestOpenAIPlatform_Models(t *testing.T) {
	tests := []struct {
		name         string
		apiKey       string
		defaultModel string
		wantErr      bool
		checkModels  func(t *testing.T, models []*ModelInfo)
		skipIfNoKey  bool
	}{
		{
			name:         "valid API key",
			apiKey:       os.Getenv("OPENAI_API_KEY"),
			defaultModel: "gpt-4",
			wantErr:      false,
			skipIfNoKey:  true,
			checkModels: func(t *testing.T, models []*ModelInfo) {
				require.NotEmpty(t, models, "should return at least one model")

				// Check each model's properties
				for _, model := range models {
					assert.NotEmpty(t, model.Name, "model name should not be empty")
					assert.Contains(t, model.Name, "gpt", "model name should contain 'gpt'")
					assert.Equal(t, "N/A", model.SizeHuman, "model size should be 'N/A'")

					// Check default model marking
					if model.Name == "gpt-4" {
						assert.True(t, model.IsDefault, "gpt-4 should be marked as default")
					}
				}

				// Verify model types
				modelTypes := make(map[string]bool)
				for _, model := range models {
					if strings.Contains(model.Name, "gpt-4") {
						modelTypes["gpt-4"] = true
					} else if strings.Contains(model.Name, "gpt-3.5") {
						modelTypes["gpt-3.5"] = true
					}
				}
				assert.True(t, modelTypes["gpt-4"], "should have at least one GPT-4 model")
				assert.True(t, modelTypes["gpt-3.5"], "should have at least one GPT-3.5 model")
			},
		},
		{
			name:         "invalid API key",
			apiKey:       "invalid-key",
			defaultModel: "gpt-4",
			wantErr:      true,
			skipIfNoKey:  false,
		},
		{
			name:         "empty API key",
			apiKey:       "",
			defaultModel: "gpt-4",
			wantErr:      true,
			skipIfNoKey:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip test if OpenAI tests are disabled or API key is required but not set
			if tt.skipIfNoKey {
				skipIfNoOpenAI(t)
			}

			// Create platform with test configuration
			platform := newOpenAIPlatform(OpenAIConfig{
				APIKey: tt.apiKey,
				Model:  tt.defaultModel,
			})

			// Get models
			models, err := platform.Models()

			// Check error condition
			if tt.wantErr {
				assert.Error(t, err, "expected error")
				return
			}

			// Verify successful response
			require.NoError(t, err, "unexpected error")
			require.NotNil(t, models, "models should not be nil")

			// Run model-specific checks if provided
			if tt.checkModels != nil {
				tt.checkModels(t, models)
			}
		})
	}
}
