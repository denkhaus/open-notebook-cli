package models

// Embedding API models

// EmbedRequest represents embedding request
type EmbedRequest struct {
	ItemID          string   `json:"item_id"`
	ItemType        ItemType `json:"item_type"` // typesafe enum from common.go
	AsyncProcessing bool     `json:"async_processing"`
}

// EmbedResponse represents embedding response
type EmbedResponse struct {
	Success   bool     `json:"success"`
	Message   string   `json:"message"`
	ItemID    string   `json:"item_id"`
	ItemType  ItemType `json:"item_type"`
	CommandID *string  `json:"command_id,omitempty"`
}

// Rebuild API models

// RebuildMode represents rebuild mode with type safety
type RebuildMode string

const (
	RebuildModeExisting RebuildMode = "existing"
	RebuildModeAll      RebuildMode = "all"
)

// RebuildRequest represents rebuild request
type RebuildRequest struct {
	Mode            RebuildMode `json:"mode"` // typesafe enum
	IncludeSources  bool        `json:"include_sources"`
	IncludeNotes    bool        `json:"include_notes"`
	IncludeInsights bool        `json:"include_insights"`
}

// RebuildResponse represents rebuild response
type RebuildResponse struct {
	CommandID  string `json:"command_id"`
	TotalItems int    `json:"total_items"`
	Message    string `json:"message"`
}

// RebuildProgress represents rebuild progress
type RebuildProgress struct {
	Processed  int     `json:"processed"`
	Total      int     `json:"total"`
	Percentage float64 `json:"percentage"`
}

// RebuildStats represents rebuild statistics
type RebuildStats struct {
	Sources  int `json:"sources"`
	Notes    int `json:"notes"`
	Insights int `json:"insights"`
	Failed   int `json:"failed"`
}

// RebuildStatusResponse represents rebuild status response
type RebuildStatusResponse struct {
	CommandID    string           `json:"command_id"`
	Status       RebuildStatus    `json:"status"` // typesafe enum from common.go
	Progress     *RebuildProgress `json:"progress,omitempty"`
	Stats        *RebuildStats    `json:"stats,omitempty"`
	StartedAt    *string          `json:"started_at,omitempty"`
	CompletedAt  *string          `json:"completed_at,omitempty"`
	ErrorMessage *string          `json:"error_message,omitempty"`
}
