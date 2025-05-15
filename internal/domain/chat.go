package domain

import (
	"time"
)

// ChatItem represents a single message in a chat conversation
type ChatItem struct {
	ID      string    `json:"id"`
	ChatID  string    `json:"chat"`
	Role    string    `json:"role"`
	Content string    `json:"content"`
	Usage   *Usage    `json:"usage,omitempty"`
	Order   int       `json:"order"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

// Chat represents a chat conversation
type Chat struct {
	ID           string      `json:"id"`
	UserID       string      `json:"user"`
	Label        string      `json:"label"`
	SystemPrompt string      `json:"system_prompt"`
	Model        string      `json:"model"`
	TotalTokens  int         `json:"total_tokens"`
	TotalCost    float64     `json:"total_cost"`
	Created      time.Time   `json:"created"`
	Updated      time.Time   `json:"updated"`
	Items        []*ChatItem `json:"items,omitempty"`
}

// ChatItemRole defines the possible roles for chat messages
const (
	ChatItemRoleUser      = "user"
	ChatItemRoleAssistant = "assistant"
	ChatItemRoleSystem    = "system"
)
