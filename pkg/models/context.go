package models

// Context API models

// ContextLevel represents relevance level for context items
type ContextLevel string

const (
	ContextLevelLow      ContextLevel = "low"
	ContextLevelMedium   ContextLevel = "medium"
	ContextLevelHigh     ContextLevel = "high"
	ContextLevelCritical ContextLevel = "critical"
)

// ContextConfig represents context configuration
type ContextConfig struct {
	Sources map[string]ContextLevel `json:"sources,omitempty"` // {source_id: level}
	Notes   map[string]ContextLevel `json:"notes,omitempty"`   // {note_id: level}
}

// ContextRequest represents context request
type ContextRequest struct {
	NotebookID    *string        `json:"notebook_id,omitempty"`
	ContextConfig *ContextConfig `json:"context_config,omitempty"`
}

// ContextResponse represents context response
type ContextResponse struct {
	NotebookID  string           `json:"notebook_id"`
	Sources     []map[string]any `json:"sources"`
	Notes       []map[string]any `json:"notes"`
	TotalTokens *int             `json:"total_tokens,omitempty"`
}
