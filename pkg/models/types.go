package models

// Typesafe enums for all model fields

// SearchType represents search type with type safety
type SearchType string

const (
	SearchTypeVector SearchType = "vector"
	SearchTypeText   SearchType = "text"
)

// ModelType represents AI model type with type safety
type ModelType string

const (
	ModelTypeLanguage     ModelType = "language"
	ModelTypeEmbedding    ModelType = "embedding"
	ModelTypeTextToSpeech ModelType = "text_to_speech"
	ModelTypeSpeechToText ModelType = "speech_to_text"
)

// SourceType represents source type with type safety
type SourceType string

const (
	SourceTypeLink   SourceType = "link"
	SourceTypeUpload SourceType = "upload"
	SourceTypeText   SourceType = "text"
)

// NoteType represents note type with type safety
type NoteType string

const (
	NoteTypeHuman NoteType = "human"
	NoteTypeAI    NoteType = "ai"
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

// RebuildMode represents rebuild mode with type safety
type RebuildMode string

const (
	RebuildModeExisting RebuildMode = "existing"
	RebuildModeAll      RebuildMode = "all"
)

// RebuildStatus represents rebuild status with type safety
type RebuildStatus string

const (
	RebuildStatusQueued    RebuildStatus = "queued"
	RebuildStatusRunning   RebuildStatus = "running"
	RebuildStatusCompleted RebuildStatus = "completed"
	RebuildStatusFailed    RebuildStatus = "failed"
)

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

// InsightType represents insight type for sources
type InsightType string

const (
	InsightTypeSummary    InsightType = "summary"
	InsightTypeAnalysis   InsightType = "analysis"
	InsightTypeExtraction InsightType = "extraction"
	InsightTypeQuestion   InsightType = "question"
	InsightTypeReflection InsightType = "reflection"
)

// ContextLevel represents relevance level for context items
type ContextLevel string

const (
	ContextLevelLow      ContextLevel = "low"
	ContextLevelMedium   ContextLevel = "medium"
	ContextLevelHigh     ContextLevel = "high"
	ContextLevelCritical ContextLevel = "critical"
)

// Core domain types and data models based on actual OpenNotebook API

// Search models from OpenNotebook API

// SearchRequest represents search request
type SearchRequest struct {
	Query         string     `json:"query"`
	Type          SearchType `json:"type"` // typesafe enum
	Limit         int        `json:"limit"`
	SearchSources bool       `json:"search_sources"`
	SearchNotes   bool       `json:"search_notes"`
	MinimumScore  float64    `json:"minimum_score"`
}

// SearchResponse represents search response
type SearchResponse struct {
	Results    []SearchResult `json:"results"`
	TotalCount int            `json:"total_count"`
	SearchType string         `json:"search_type"`
}

// SearchResult represents a single search result item
type SearchResult struct {
	ID        string  `json:"id"`
	ParentID  string  `json:"parent_id"`
	Relevance float64 `json:"relevance"`
	Title     string  `json:"title"`
}

// AskRequest represents ask request
type AskRequest struct {
	Question         string `json:"question"`
	StrategyModel    string `json:"strategy_model"`
	AnswerModel      string `json:"answer_model"`
	FinalAnswerModel string `json:"final_answer_model"`
}

// AskResponse represents ask response
type AskResponse struct {
	Answer   string `json:"answer"`
	Question string `json:"question"`
}

// Models API models

// ModelCreate represents model creation request
type ModelCreate struct {
	Name     string    `json:"name"`
	Provider string    `json:"provider"`
	Type     ModelType `json:"type"` // typesafe enum
}

// Model represents a model from the API
type Model struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Provider string    `json:"provider"`
	Type     ModelType `json:"type"` // typesafe enum
	Created  string    `json:"created"`
	Updated  string    `json:"updated"`
}

// DefaultModelsResponse represents default models response
type DefaultModelsResponse struct {
	DefaultChatModel           *string `json:"default_chat_model"`
	DefaultTransformationModel *string `json:"default_transformation_model"`
	LargeContextModel          *string `json:"large_context_model"`
	DefaultTextToSpeechModel   *string `json:"default_text_to_speech_model"`
	DefaultSpeechToTextModel   *string `json:"default_speech_to_text_model"`
	DefaultEmbeddingModel      *string `json:"default_embedding_model"`
	DefaultToolsModel          *string `json:"default_tools_model"`
}

// ProviderAvailabilityResponse represents provider availability response
type ProviderAvailabilityResponse struct {
	Available      []string            `json:"available"`
	Unavailable    []string            `json:"unavailable"`
	SupportedTypes map[string][]string `json:"supported_types"`
}

// Transformations API models

// TransformationCreate represents transformation creation request
type TransformationCreate struct {
	Name         string `json:"name"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Prompt       string `json:"prompt"`
	ApplyDefault bool   `json:"apply_default"`
}

// TransformationUpdate represents transformation update request
type TransformationUpdate struct {
	Name         *string `json:"name,omitempty"`
	Title        *string `json:"title,omitempty"`
	Description  *string `json:"description,omitempty"`
	Prompt       *string `json:"prompt,omitempty"`
	ApplyDefault *bool   `json:"apply_default,omitempty"`
}

// Transformation represents a transformation from the API
type Transformation struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Prompt       string `json:"prompt"`
	ApplyDefault bool   `json:"apply_default"`
	Created      string `json:"created"`
	Updated      string `json:"updated"`
}

// TransformationExecuteRequest represents transformation execution request
type TransformationExecuteRequest struct {
	TransformationID string `json:"transformation_id"`
	InputText        string `json:"input_text"`
	ModelID          string `json:"model_id"`
}

// TransformationExecuteResponse represents transformation execution response
type TransformationExecuteResponse struct {
	Output           string `json:"output"`
	TransformationID string `json:"transformation_id"`
	ModelID          string `json:"model_id"`
}

// Notes API models

// NoteCreate represents note creation request
type NoteCreate struct {
	Title      *string   `json:"title,omitempty"`
	Content    string    `json:"content"`
	NoteType   *NoteType `json:"note_type,omitempty"` // typesafe enum
	NotebookID *string   `json:"notebook_id,omitempty"`
}

// NoteUpdate represents note update request
type NoteUpdate struct {
	Title    *string   `json:"title,omitempty"`
	Content  *string   `json:"content,omitempty"`
	NoteType *NoteType `json:"note_type,omitempty"`
}

// Note represents a note from the API
type Note struct {
	ID       *string   `json:"id"`
	Title    *string   `json:"title"`
	Content  *string   `json:"content"`
	NoteType *NoteType `json:"note_type"` // typesafe enum
	Created  string    `json:"created"`
	Updated  string    `json:"updated"`
	// notebook_id is not returned by the API, it's only used for creation
}

// Embedding API models

// EmbedRequest represents embedding request
type EmbedRequest struct {
	ItemID          string   `json:"item_id"`
	ItemType        ItemType `json:"item_type"` // typesafe enum
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
	Status       RebuildStatus    `json:"status"` // typesafe enum
	Progress     *RebuildProgress `json:"progress,omitempty"`
	Stats        *RebuildStats    `json:"stats,omitempty"`
	StartedAt    *string          `json:"started_at,omitempty"`
	CompletedAt  *string          `json:"completed_at,omitempty"`
	ErrorMessage *string          `json:"error_message,omitempty"`
}

// Settings API models

// SettingsResponse represents settings response
type SettingsResponse struct {
	DefaultContentProcessingEngineDoc ContentProcessingEngine    `json:"default_content_processing_engine_doc"`
	DefaultContentProcessingEngineURL ContentProcessingEngineURL `json:"default_content_processing_engine_url"`
	DefaultEmbeddingOption            EmbeddingOption            `json:"default_embedding_option"`
	AutoDeleteFiles                   YesNoDecision              `json:"auto_delete_files"`
	YoutubePreferredLanguages         []string                   `json:"youtube_preferred_languages"`
}

// SettingsUpdate represents settings update request
type SettingsUpdate struct {
	DefaultContentProcessingEngineDoc *ContentProcessingEngine    `json:"default_content_processing_engine_doc,omitempty"`
	DefaultContentProcessingEngineURL *ContentProcessingEngineURL `json:"default_content_processing_engine_url,omitempty"`
	DefaultEmbeddingOption            *EmbeddingOption            `json:"default_embedding_option,omitempty"`
	AutoDeleteFiles                   *YesNoDecision              `json:"auto_delete_files,omitempty"`
	YoutubePreferredLanguages         []string                    `json:"youtube_preferred_languages,omitempty"`
}

// Sources API models

// AssetModel represents asset information
type AssetModel struct {
	FilePath *string `json:"file_path,omitempty"`
	URL      *string `json:"url,omitempty"`
}

// SourceCreate represents source creation request
type SourceCreate struct {
	NotebookID      *string    `json:"notebook_id,omitempty"` // use Notebooks instead
	Notebooks       []string   `json:"notebooks,omitempty"`   // preferred way
	Type            SourceType `json:"type"`                  // typesafe enum
	URL             *string    `json:"url,omitempty"`
	FilePath        *string    `json:"file_path,omitempty"`
	Content         *string    `json:"content,omitempty"`
	Title           *string    `json:"title,omitempty"`
	Transformations []string   `json:"transformations,omitempty"`
	Embed           bool       `json:"embed"`
	DeleteSource    bool       `json:"delete_source"`
	AsyncProcessing bool       `json:"async_processing"`
}

// SourceUpdate represents source update request
type SourceUpdate struct {
	Title  *string  `json:"title,omitempty"`
	Topics []string `json:"topics,omitempty"`
}

// Source represents a source from the API
type Source struct {
	ID             *string                `json:"id"`
	Title          *string                `json:"title"`
	Topics         []string               `json:"topics"`
	Asset          *AssetModel            `json:"asset"`
	FullText       *string                `json:"full_text"`
	Embedded       bool                   `json:"embedded"`
	EmbeddedChunks int                    `json:"embedded_chunks"`
	FileAvailable  *bool                  `json:"file_available,omitempty"`
	Created        string                 `json:"created"`
	Updated        string                 `json:"updated"`
	CommandID      *string                `json:"command_id,omitempty"`
	Status         *SourceStatus          `json:"status,omitempty"`
	ProcessingInfo map[string]any        `json:"processing_info,omitempty"`
	Notebooks      []string               `json:"notebooks,omitempty"`
}

// SourceListResponse represents source in list response
type SourceListResponse struct {
	ID             *string                `json:"id"`
	Title          *string                `json:"title"`
	Topics         []string               `json:"topics"`
	Asset          *AssetModel            `json:"asset"`
	Embedded       bool                   `json:"embedded"`
	EmbeddedChunks int                    `json:"embedded_chunks"`
	InsightsCount  int                    `json:"insights_count"`
	Created        string                 `json:"created"`
	Updated        string                 `json:"updated"`
	FileAvailable  *bool                  `json:"file_available,omitempty"`
	CommandID      *string                `json:"command_id,omitempty"`
	Status         *SourceStatus          `json:"status,omitempty"`
	ProcessingInfo map[string]any        `json:"processing_info,omitempty"`
}

// Context API models

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
	NotebookID  string                   `json:"notebook_id"`
	Sources     []map[string]any `json:"sources"`
	Notes       []map[string]any `json:"notes"`
	TotalTokens *int                     `json:"total_tokens,omitempty"`
}

// Insights API models

// SourceInsightResponse represents source insight response
type SourceInsightResponse struct {
	ID          string      `json:"id"`
	SourceID    string      `json:"source_id"`
	InsightType InsightType `json:"insight_type"`
	Content     string      `json:"content"`
	Created     string      `json:"created"`
	Updated     string      `json:"updated"`
}

// SaveAsNoteRequest represents save as note request
type SaveAsNoteRequest struct {
	NotebookID *string `json:"notebook_id,omitempty"`
}

// CreateSourceInsightRequest represents create source insight request
type CreateSourceInsightRequest struct {
	TransformationID string  `json:"transformation_id"`
	ModelID          *string `json:"model_id,omitempty"`
}

// Error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// Job management models

// JobStatus represents background job status
type JobStatus struct {
	ID       string   `json:"id"`
	Status   string   `json:"status"` // queued, running, completed, failed
	Progress *float64 `json:"progress,omitempty"`
	Message  *string  `json:"message,omitempty"`
	Created  string   `json:"created"`
	Updated  *string  `json:"updated,omitempty"`
}

// JobsListResponse represents jobs list response
type JobsListResponse struct {
	Jobs []JobStatus `json:"jobs"`
}

// ModelsListResponse represents models list response with pagination
type ModelsListResponse struct {
	Models []Model `json:"models"`
	Total  int     `json:"total"`
}

// SourcesListResponse represents sources list response
type SourcesListResponse struct {
	Sources []SourceListResponse `json:"sources"`
}

// SourceStatus represents source status response
type SourceStatusResponse struct {
	Status         *SourceStatus          `json:"status,omitempty"`
	Message        string                 `json:"message"`
	ProcessingInfo map[string]any        `json:"processing_info,omitempty"`
	CommandID      *string                `json:"command_id,omitempty"`
}

// Default Prompt API models

// DefaultPromptResponse represents default prompt response
type DefaultPromptResponse struct {
	TransformationInstructions string `json:"transformation_instructions"`
}

// DefaultPromptUpdate represents default prompt update request
type DefaultPromptUpdate struct {
	TransformationInstructions string `json:"transformation_instructions"`
}

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
