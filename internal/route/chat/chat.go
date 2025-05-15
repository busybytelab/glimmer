package chat

// ChatRoutes handles chat operations. This route specifically handles chat creation on demand.
// For other operations such as update, search, list, etc., use the standard PocketBase collection API.

import (
	"net/http"

	"github.com/busybytelab.com/glimmer/internal/domain"
	"github.com/busybytelab.com/glimmer/internal/llm"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/rs/zerolog/log"
)

type (
	// ChatRequest defines the request body for the chat endpoint
	ChatRequest struct {
		ChatID       string `json:"chatId" form:"chatId"`
		UserMessage  string `json:"userMessage" form:"userMessage"`
		SystemPrompt string `json:"systemPrompt" form:"systemPrompt"`
		Model        string `json:"model" form:"model"`
	}

	// ChatResponse defines the response body for the chat endpoint
	ChatResponse struct {
		Response string        `json:"response"`
		Usage    *domain.Usage `json:"usage,omitempty"`
		Chat     *domain.Chat  `json:"chat,omitempty"`
	}

	ChatRoutes interface {
		HandleChatRequest(e *core.RequestEvent) error
	}

	chatRoutes struct {
		chatService llm.ChatService
	}
)

func New(chatService llm.ChatService) ChatRoutes {
	return &chatRoutes{
		chatService: chatService,
	}
}

// HandleChatRequest handles chat requests, creating a chat if needed
func (r *chatRoutes) HandleChatRequest(e *core.RequestEvent) error {
	// Parse request body
	var req ChatRequest
	if err := e.BindBody(&req); err != nil {
		return e.BadRequestError("Invalid request body", err)
	}

	// Validate request
	if req.UserMessage == "" {
		return e.BadRequestError("User message is required", nil)
	}

	// Get user ID from auth record
	if e.Auth == nil {
		return apis.NewUnauthorizedError("You must be logged in", nil)
	}
	userID := e.Auth.Id

	var chatID string
	var err error
	var chat *domain.Chat

	// Create a new chat if no chat ID is provided
	if req.ChatID == "" {
		// Set default system prompt if not provided
		systemPrompt := req.SystemPrompt
		if systemPrompt == "" {
			systemPrompt = "You are a helpful assistant."
		}

		// Create a new chat
		chat, err = r.chatService.CreateChat(userID, systemPrompt, req.Model)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create chat")
			return e.InternalServerError("Failed to create chat", err)
		}
		log.Debug().Str("chatID", chat.ID).Msg("Created chat")
		chatID = chat.ID
	} else {
		// Use existing chat
		chatID = req.ChatID

		// Verify the chat belongs to the user
		chat, err = r.chatService.GetChat(chatID)
		if err != nil {
			log.Error().Err(err).Str("chatID", chatID).Msg("Failed to get chat")
			return e.NotFoundError("Chat not found", err)
		}

		if chat.UserID != userID {
			log.Warn().Str("chatID", chatID).Str("userID", userID).Msg("User tried to access chat they don't own")
			return e.UnauthorizedError("Not authorized to access this chat", nil)
		}
	}

	// Add options from the request if provided
	var opts []llm.ChatOption
	if req.Model != "" {
		opts = append(opts, llm.WithModel(req.Model))
	}

	// Process chat request
	response, usage, err := r.chatService.ChatCompletion(chatID, req.UserMessage, opts...)
	if err != nil {
		log.Error().Err(err).Msg("Failed to process chat request")
		return e.InternalServerError("Failed to process chat request", err)
	}

	// Get the updated chat (with latest messages)
	chat, err = r.chatService.GetChat(chatID)
	if err != nil {
		log.Warn().Err(err).Str("chatID", chatID).Msg("Failed to get updated chat")
	}

	// Return response
	return e.JSON(http.StatusOK, ChatResponse{
		Response: response,
		Usage:    usage,
		Chat:     chat,
	})
}
