package llm

import (
	"net/http"

	"github.com/busybytelab.com/glimmer/internal/llm"
	"github.com/pocketbase/pocketbase/core"
	"github.com/rs/zerolog/log"
)

type (
	// ChatRequest defines the request body for the LLM chat endpoint
	ChatRequest struct {
		Prompt       string `json:"prompt" form:"prompt"`
		SystemPrompt string `json:"systemPrompt" form:"systemPrompt"`
		Model        string `json:"model" form:"model"`
	}

	// ChatResponse defines the response body for the LLM chat endpoint
	ChatResponse struct {
		Response string     `json:"response"`
		Usage    *llm.Usage `json:"usage,omitempty"`
	}

	LLMRoutes interface {
		HandleChatRequest(e *core.RequestEvent) error
		HandleInfoRequest(e *core.RequestEvent) error
	}

	llmRoutes struct {
		llmService llm.Service
	}
)

func New(llmService llm.Service) LLMRoutes {
	return &llmRoutes{
		llmService: llmService,
	}
}

// HandleChatRequest handles LLM chat requests
func (r *llmRoutes) HandleChatRequest(e *core.RequestEvent) error {
	// Parse request body
	var req ChatRequest
	if err := e.BindBody(&req); err != nil {
		return e.BadRequestError("Invalid request body", err)
	}

	// Validate request
	if req.Prompt == "" {
		return e.BadRequestError("Prompt is required", nil)
	}

	// Set default system prompt if not provided
	if req.SystemPrompt == "" {
		req.SystemPrompt = "You are a helpful assistant."
	}

	// Process chat request
	var opts []llm.ChatOption
	if req.Model != "" {
		opts = append(opts, llm.WithModel(req.Model))
	}

	// Send chat request to LLM service
	response, usage, err := r.llmService.Chat(req.Prompt, req.SystemPrompt, opts...)
	if err != nil {
		log.Error().Err(err).Msg("Failed to process LLM chat request")
		return e.InternalServerError("Failed to process chat request", err)
	}

	// Return response
	return e.JSON(http.StatusOK, ChatResponse{
		Response: response,
		Usage:    usage,
	})
}

// HandleInfoRequest handles requests for LLM platform and model information
func (r *llmRoutes) HandleInfoRequest(e *core.RequestEvent) error {
	log.Debug().Msg("Processing LLM models info request")

	info := r.llmService.Info()

	return e.JSON(http.StatusOK, info)
}
