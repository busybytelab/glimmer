package llm

import (
	"errors"
	"fmt"

	"github.com/busybytelab.com/glimmer/internal/domain"
	"github.com/pocketbase/pocketbase/core"
	"github.com/rs/zerolog/log"
)

// ChatService provides chat-related operations
type ChatService interface {
	// CreateChat creates a new chat with the given user ID and system prompt
	CreateChat(userID, systemPrompt, model string) (*domain.Chat, error)

	// GetChat retrieves a chat by ID
	GetChat(chatID string) (*domain.Chat, error)

	// GetChats retrieves chats for a user, sorted by most recent
	GetChats(userID string, limit, offset int) ([]*domain.Chat, error)

	// UpdateChatLabel updates a chat's label
	UpdateChatLabel(chatID, label string) error

	// ChatCompletion sends a user message to the LLM and stores the result in the chat history
	ChatCompletion(chatID, userMessage string, opts ...ChatOption) (string, *domain.Usage, error)

	// AddChatMessage adds a message to a chat
	AddChatMessage(chatID, role, content string, usage *domain.Usage) (*domain.ChatItem, error)

	// GetChatMessages retrieves messages for a chat
	GetChatMessages(chatID string, limit, offset int) ([]*domain.ChatItem, error)
}

// chatService implements ChatService interface
type chatService struct {
	app        core.App
	llmService Service
}

// NewChatService creates a new ChatService instance
func NewChatService(app core.App, llmService Service) ChatService {
	return &chatService{
		app:        app,
		llmService: llmService,
	}
}

// CreateChat creates a new chat
func (s *chatService) CreateChat(userID, systemPrompt, model string) (*domain.Chat, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	// Default label is "New chat"
	label := "New chat"

	// Create the chat record
	collection, err := s.app.FindCollectionByNameOrId(domain.CollectionChats)
	if err != nil {
		return nil, fmt.Errorf("failed to find chats collection: %w", err)
	}

	record := core.NewRecord(collection)
	record.Set("user", userID)
	record.Set("label", label)
	record.Set("system_prompt", systemPrompt)
	record.Set("model", model)
	record.Set("total_tokens", 0)
	record.Set("total_cost", 0.0)

	if err := s.app.Save(record); err != nil {
		log.Error().Err(err).Msg("Failed to create chat")
		return nil, fmt.Errorf("failed to create chat: %w", err)
	}

	// If system prompt is provided, add it as the first message
	if systemPrompt != "" {
		_, err = s.AddChatMessage(record.Id, "system", systemPrompt, nil)
		if err != nil {
			log.Error().Err(err).Msg("Failed to add system message to chat")
			// Don't return an error here, the chat was created successfully
		}
	}

	// Convert to domain model
	chat := s.recordToChat(record)
	return chat, nil
}

// GetChat retrieves a chat by ID
func (s *chatService) GetChat(chatID string) (*domain.Chat, error) {
	record, err := s.app.FindRecordById(domain.CollectionChats, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to find chat: %w", err)
	}

	chat := s.recordToChat(record)

	// Get messages
	messages, err := s.GetChatMessages(chatID, 100, 0) // Get up to 100 most recent messages
	if err != nil {
		log.Error().Err(err).Str("chatID", chatID).Msg("Failed to get chat messages")
		// Don't return an error, just return the chat without messages
	}

	chat.Items = messages
	return chat, nil
}

// GetChats retrieves chats for a user, sorted by most recent
func (s *chatService) GetChats(userID string, limit, offset int) ([]*domain.Chat, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	if limit <= 0 {
		limit = 10 // Default limit
	}

	// Ensure offset is not negative
	if offset < 0 {
		offset = 0
	}

	// Find chats
	filter := fmt.Sprintf("user = '%s'", userID)
	records, err := s.app.FindRecordsByFilter(
		domain.CollectionChats,
		filter,
		"-updated", // Sort by most recently updated
		limit,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find chats: %w", err)
	}

	// Convert to domain models
	chats := make([]*domain.Chat, len(records))
	for i, record := range records {
		chats[i] = s.recordToChat(record)
	}

	return chats, nil
}

// UpdateChatLabel updates a chat's label
func (s *chatService) UpdateChatLabel(chatID, label string) error {
	if chatID == "" {
		return errors.New("chat ID is required")
	}

	if label == "" {
		return errors.New("label is required")
	}

	record, err := s.app.FindRecordById(domain.CollectionChats, chatID)
	if err != nil {
		return fmt.Errorf("failed to find chat: %w", err)
	}

	record.Set("label", label)
	if err := s.app.Save(record); err != nil {
		return fmt.Errorf("failed to update chat label: %w", err)
	}

	return nil
}

// ChatCompletion sends a user message to the LLM and stores the result in the chat history
func (s *chatService) ChatCompletion(chatID, userMessage string, opts ...ChatOption) (string, *domain.Usage, error) {
	if chatID == "" {
		return "", nil, errors.New("chat ID is required")
	}

	if userMessage == "" {
		return "", nil, errors.New("user message is required")
	}

	// Get the chat
	chat, err := s.GetChat(chatID)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get chat: %w", err)
	}

	// Add user message to chat
	_, err = s.AddChatMessage(chatID, "user", userMessage, nil)
	if err != nil {
		return "", nil, fmt.Errorf("failed to add user message to chat: %w", err)
	}

	// Use the model from the chat if not specified in options
	model := chat.Model

	// If model is still empty, use a default model
	info := s.llmService.Info()
	if model == "" && len(info.Platforms) > 0 && len(info.Platforms[0].Models) > 0 {
		model = info.Platforms[0].Models[0].Name
	}

	// Get the LLM response
	var chatOpts []ChatOption
	if model != "" {
		chatOpts = append(chatOpts, WithModel(model))
	}

	// If chat has items, use them for context
	var llmResponse string
	var usage *domain.Usage

	// Get previous messages for context
	previousMessages, err := s.GetChatMessages(chatID, 100, 0) // Limit to recent messages
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get previous messages, proceeding with single message")

		// Fallback to regular Chat without history
		llmResponse, usage, err = s.llmService.Chat(userMessage, chat.SystemPrompt, chatOpts...)
		if err != nil {
			return "", nil, fmt.Errorf("failed to get LLM response: %w", err)
		}
	} else {
		// Use ChatWithHistory for conversation context
		llmResponse, usage, err = s.llmService.ChatWithHistory(previousMessages, chat.SystemPrompt, chatOpts...)
		if err != nil {
			log.Warn().Err(err).Msg("ChatWithHistory failed, falling back to single message Chat")

			// Fallback to regular Chat if conversational context fails
			llmResponse, usage, err = s.llmService.Chat(userMessage, chat.SystemPrompt, chatOpts...)
			if err != nil {
				return "", nil, fmt.Errorf("failed to get LLM response: %w", err)
			}
		}
	}

	// Add assistant response to chat
	_, err = s.AddChatMessage(chatID, "assistant", llmResponse, usage)
	if err != nil {
		log.Error().Err(err).Msg("Failed to add assistant message to chat")
		// Don't return an error, the LLM response was generated successfully
	}

	// Update chat record with the new token usage
	chatRecord, err := s.app.FindRecordById(domain.CollectionChats, chatID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to find chat record for updating token usage")
	} else {
		currentTokens := chatRecord.GetInt("total_tokens")
		currentCost := chatRecord.GetFloat("total_cost")

		if usage != nil {
			chatRecord.Set("total_tokens", currentTokens+usage.TotalTokens)
			chatRecord.Set("total_cost", currentCost+usage.Cost)
			if err := s.app.Save(chatRecord); err != nil {
				log.Error().Err(err).Msg("Failed to update chat token usage")
			}
		}
	}

	return llmResponse, usage, nil
}

// AddChatMessage adds a message to a chat
func (s *chatService) AddChatMessage(chatID, role, content string, usage *domain.Usage) (*domain.ChatItem, error) {
	if chatID == "" {
		return nil, errors.New("chat ID is required")
	}

	if role == "" {
		return nil, errors.New("role is required")
	}

	if content == "" {
		return nil, errors.New("content is required")
	}

	// Get the chat to make sure it exists
	_, err := s.app.FindRecordById(domain.CollectionChats, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to find chat: %w", err)
	}

	// Get the last order number
	order := 0
	lastMessage, err := s.getLastChatMessage(chatID)
	if err == nil && lastMessage != nil {
		order = lastMessage.Order + 1
	}

	// Create the message record
	collection, err := s.app.FindCollectionByNameOrId(domain.CollectionChatItems)
	if err != nil {
		return nil, fmt.Errorf("failed to find chat_items collection: %w", err)
	}

	record := core.NewRecord(collection)
	record.Set("chat", chatID)
	record.Set("role", role)
	record.Set("content", content)
	record.Set("order", order)

	// Add usage data if provided
	if usage != nil {
		record.Set("prompt_tokens", usage.PromptTokens)
		record.Set("completion_tokens", usage.CompletionTokens)
		record.Set("total_tokens", usage.TotalTokens)
		record.Set("cost", usage.Cost)
		record.Set("model", usage.LlmModelName)
	}

	if err := s.app.Save(record); err != nil {
		return nil, fmt.Errorf("failed to create chat message: %w", err)
	}

	// Update the updated_at field of the chat
	chatRecord, err := s.app.FindRecordById(domain.CollectionChats, chatID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to find chat record for updating timestamp")
	} else {
		// PocketBase automatically updates the "updated" field on Save
		if err := s.app.Save(chatRecord); err != nil {
			log.Error().Err(err).Msg("Failed to update chat timestamp")
		}
	}

	// Convert to domain model
	chatItem := s.recordToChatItem(record)
	return chatItem, nil
}

// GetChatMessages retrieves messages for a chat
func (s *chatService) GetChatMessages(chatID string, limit, offset int) ([]*domain.ChatItem, error) {
	if chatID == "" {
		return nil, errors.New("chat ID is required")
	}

	if limit <= 0 {
		limit = 50 // Default limit
	}

	// Ensure offset is not negative
	if offset < 0 {
		offset = 0
	}

	// Find messages
	filter := fmt.Sprintf("chat = '%s'", chatID)
	records, err := s.app.FindRecordsByFilter(
		domain.CollectionChatItems,
		filter,
		"order", // Sort by order
		limit,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find chat messages: %w", err)
	}

	// Convert to domain models
	messages := make([]*domain.ChatItem, len(records))
	for i, record := range records {
		messages[i] = s.recordToChatItem(record)
	}

	return messages, nil
}

// getLastChatMessage gets the last message in a chat
func (s *chatService) getLastChatMessage(chatID string) (*domain.ChatItem, error) {
	if chatID == "" {
		return nil, errors.New("chat ID is required")
	}

	filter := fmt.Sprintf("chat = '%s'", chatID)
	records, err := s.app.FindRecordsByFilter(
		domain.CollectionChatItems,
		filter,
		"-order", // Sort by order descending
		1,        // Limit to 1
		0,        // No offset
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find last chat message: %w", err)
	}

	if len(records) == 0 {
		// No previous messages, return a placeholder with order -1
		// This will cause the first message to get order 0 when +1 is applied
		return &domain.ChatItem{Order: -1}, nil
	}

	// Convert to domain model
	chatItem := s.recordToChatItem(records[0])
	return chatItem, nil
}

// recordToChat converts a PocketBase record to a domain.Chat
func (s *chatService) recordToChat(record *core.Record) *domain.Chat {
	created := record.GetDateTime("created")
	updated := record.GetDateTime("updated")

	return &domain.Chat{
		ID:           record.Id,
		UserID:       record.GetString("user"),
		Label:        record.GetString("label"),
		SystemPrompt: record.GetString("system_prompt"),
		Model:        record.GetString("model"),
		TotalTokens:  record.GetInt("total_tokens"),
		TotalCost:    record.GetFloat("total_cost"),
		Created:      created.Time(),
		Updated:      updated.Time(),
		Items:        nil, // Will be populated by GetChat if needed
	}
}

// recordToChatItem converts a PocketBase record to a domain.ChatItem
func (s *chatService) recordToChatItem(record *core.Record) *domain.ChatItem {
	created := record.GetDateTime("created")
	updated := record.GetDateTime("updated")

	// Create usage if the fields exist
	var usage *domain.Usage
	if record.Get("prompt_tokens") != nil && record.Get("completion_tokens") != nil && record.Get("total_tokens") != nil {
		usage = &domain.Usage{
			PromptTokens:     record.GetInt("prompt_tokens"),
			CompletionTokens: record.GetInt("completion_tokens"),
			TotalTokens:      record.GetInt("total_tokens"),
			Cost:             record.GetFloat("cost"),
			LlmModelName:     record.GetString("model"),
		}
	}

	return &domain.ChatItem{
		ID:      record.Id,
		ChatID:  record.GetString("chat"),
		Role:    record.GetString("role"),
		Content: record.GetString("content"),
		Usage:   usage,
		Order:   record.GetInt("order"),
		Created: created.Time(),
		Updated: updated.Time(),
	}
}
