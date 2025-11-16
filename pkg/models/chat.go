package models

// Chat API models

// ChatSession represents a chat session
type ChatSession struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	ModelID      string `json:"model_id"`
	Created      string `json:"created"`
	Updated      string `json:"updated"`
	MessageCount int    `json:"message_count"`
	IsActive     bool   `json:"is_active"`
}

// ChatMessage represents a chat message
type ChatMessage struct {
	ID        string `json:"id"`
	SessionID string `json:"session_id"`
	Role      string `json:"role"` // "user", "assistant", "system"
	Content   string `json:"content"`
	Created   string `json:"created"`
}

// ChatExecuteRequest represents a chat execution request
type ChatExecuteRequest struct {
	SessionID string              `json:"session_id"`
	Message   string              `json:"message"`
	ModelID   *string             `json:"model_id,omitempty"`
	Stream    bool                `json:"stream"`
	Context   *ChatContextRequest `json:"context,omitempty"`
}

// ChatExecuteResponse represents a chat execution response
type ChatExecuteResponse struct {
	SessionID string `json:"session_id"`
	MessageID string `json:"message_id"`
	Content   string `json:"content"`
	ModelID   string `json:"model_id"`
	Created   string `json:"created"`
}

// ChatContextRequest represents chat context request
type ChatContextRequest struct {
	NotebookID string   `json:"notebook_id,omitempty"`
	Sources    []string `json:"sources,omitempty"`
	MaxTokens  *int     `json:"max_tokens,omitempty"`
}

// ChatSessionsResponse represents chat sessions list response
type ChatSessionsResponse []ChatSession

// ChatCreateRequest represents chat session creation request
type ChatCreateRequest struct {
	Title   string  `json:"title"`
	ModelID *string `json:"model_id,omitempty"`
}
