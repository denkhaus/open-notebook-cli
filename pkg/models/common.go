package models

// Common enums and types used across multiple model files

// YesNoDecision represents yes/no decision with type safety
type YesNoDecision string

const (
	YesNoDecisionYes YesNoDecision = "yes"
	YesNoDecisionNo  YesNoDecision = "no"
)

// ContentProcessingEngine represents content processing engine for documents
type ContentProcessingEngine string

const (
	ContentProcessingEngineAuto    ContentProcessingEngine = "auto"
	ContentProcessingEngineDocling ContentProcessingEngine = "docling"
	ContentProcessingEngineSimple  ContentProcessingEngine = "simple"
)

// ContentProcessingEngineURL represents content processing engine for URLs
type ContentProcessingEngineURL string

const (
	ContentProcessingEngineURLAuto      ContentProcessingEngineURL = "auto"
	ContentProcessingEngineURLFirecrawl ContentProcessingEngineURL = "firecrawl"
	ContentProcessingEngineURLJina      ContentProcessingEngineURL = "jina"
	ContentProcessingEngineURLSimple    ContentProcessingEngineURL = "simple"
)

// EmbeddingOption represents when to perform embedding
type EmbeddingOption string

const (
	EmbeddingOptionAsk    EmbeddingOption = "ask"
	EmbeddingOptionAlways EmbeddingOption = "always"
	EmbeddingOptionNever  EmbeddingOption = "never"
)

// ItemType represents item type for embedding operations
type ItemType string

const (
	ItemTypeSource ItemType = "source"
	ItemTypeNote   ItemType = "note"
)

// ProcessingStatus represents processing status with type safety
type ProcessingStatus string

const (
	ProcessingStatusPending   ProcessingStatus = "pending"
	ProcessingStatusRunning   ProcessingStatus = "running"
	ProcessingStatusCompleted ProcessingStatus = "completed"
	ProcessingStatusFailed    ProcessingStatus = "failed"
)

// SourceStatus represents source status (legacy, for backward compatibility)
type SourceStatus string

const (
	SourceStatusPending   SourceStatus = "pending"
	SourceStatusRunning   SourceStatus = "running"
	SourceStatusCompleted SourceStatus = "completed"
	SourceStatusFailed    SourceStatus = "failed"
)

// RebuildStatus represents rebuild status with type safety
type RebuildStatus string

const (
	RebuildStatusQueued    RebuildStatus = "queued"
	RebuildStatusRunning   RebuildStatus = "running"
	RebuildStatusCompleted RebuildStatus = "completed"
	RebuildStatusFailed    RebuildStatus = "failed"
)

// Internal service types for DI architecture

// StreamChunk represents a streaming response chunk (for ask/execute streaming)
type StreamChunk struct {
	Content string `json:"content"`
	Done    bool   `json:"done"`
	Error   string `json:"error,omitempty"`
}

// HTTP Response wrapper
type Response struct {
	StatusCode int
	Body       []byte
	Header     map[string][]string
}

// Error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// Service-specific options for enhanced functionality

// SearchOptions represents search options for service layer
type SearchOptions struct {
	Type          string
	Limit         int
	MinimumScore  float64
	SearchSources bool
	SearchNotes   bool
}

// AskOptions represents ask options for service layer
type AskOptions struct {
	StrategyModel    string
	AnswerModel      string
	FinalAnswerModel string
}

// SourceOptions represents source creation options
type SourceOptions struct {
	Transformations []string
	Embed           bool
	DeleteSource    bool
	AsyncProcessing bool
}

// TransformationOptions represents transformation execution options
type TransformationOptions struct {
	ModelID string
}
