package models

// Additional Source API models for new OpenNotebook API endpoints

// === New API Source Models ===
// These models are based on the current OpenNotebook API structure
// and extend/replace the legacy models in types.go

// SourceCreate represents source creation request (new API format)
type SourceCreateV2 struct {
	NotebookIDs []string  `json:"notebook_ids,omitempty"`
	Title       *string   `json:"title,omitempty"`
	Topics      []string  `json:"topics,omitempty"`
	Type        SourceType `json:"type"`
	Content     *string   `json:"content,omitempty"`
	URL         *string   `json:"url,omitempty"`
	Async       *bool     `json:"async,omitempty"`
}

// SourceUpdate represents source update request (new API format)
type SourceUpdateV2 struct {
	Title  *string  `json:"title,omitempty"`
	Topics []string `json:"topics,omitempty"`
}

// SourceResponse represents source response (new API format)
type SourceResponseV2 struct {
	ID           string            `json:"id"`
	Notebooks    []NotebookRef     `json:"notebooks"`
	Title        string            `json:"title"`
	Topics       []string          `json:"topics"`
	Type         SourceType        `json:"type"`
	Content      *string           `json:"content,omitempty"`
	URL          *string           `json:"url,omitempty"`
	Status       ProcessingStatus  `json:"status"`
	Processing   ProcessingInfo    `json:"processing"`
	CommandID    *string           `json:"command_id,omitempty"`
	ErrorMessage *string           `json:"error_message,omitempty"`
	InsightsCount int              `json:"insights_count"`
	WordCount    int               `json:"word_count"`
	ProcessingDuration *float64    `json:"processing_duration,omitempty"`
	Created      string            `json:"created"`
	Updated      *string           `json:"updated,omitempty"`
}

// SourceListResponseV2 represents source in list view (new API format)
type SourceListResponseV2 struct {
	ID            string           `json:"id"`
	Notebooks     []NotebookRef    `json:"notebooks"`
	Title         string           `json:"title"`
	Type          SourceType       `json:"type"`
	Status        ProcessingStatus `json:"status"`
	InsightsCount int              `json:"insights_count"`
	WordCount     int              `json:"word_count"`
	Created       string           `json:"created"`
}

// SourceStatusResponseV2 represents source status response (new API format)
type SourceStatusResponseV2 struct {
	Status        ProcessingStatus  `json:"status"`
	Processing    ProcessingInfo    `json:"processing"`
	CommandID     *string           `json:"command_id,omitempty"`
	ErrorMessage  *string           `json:"error_message,omitempty"`
	Progress      *float64          `json:"progress,omitempty"`
	ProcessedAt   *string           `json:"processed_at,omitempty"`
}

// ProcessingInfo represents processing information
type ProcessingInfo struct {
	Retries       int     `json:"retries"`
	LastRetry     *string `json:"last_retry,omitempty"`
	NextRetry     *string `json:"next_retry,omitempty"`
	Estimated     *string `json:"estimated,omitempty"`
}

// SourceInsightResponseV2 represents source insight (new API format)
type SourceInsightResponseV2 struct {
	ID          string  `json:"id"`
	Content     string  `json:"content"`
	SourceID    string  `json:"source_id"`
	Created     string  `json:"created"`
	Score       *float64 `json:"score,omitempty"`
}

// CreateSourceInsightRequestV2 represents insight creation request (new API format)
type CreateSourceInsightRequestV2 struct {
	Content string `json:"content"`
}

// NotebookRef represents notebook reference for sources
type NotebookRef struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// === File Upload Models ===

// FileUploadRequest represents file upload request
type FileUploadRequest struct {
	NotebookIDs []string  `json:"notebook_ids,omitempty"`
	Title       *string   `json:"title,omitempty"`
	Topics      []string  `json:"topics,omitempty"`
	Async       *bool     `json:"async,omitempty"`
}

// FileUploadResponse represents file upload response
type FileUploadResponse struct {
	SourceID string `json:"source_id"`
	Message  string `json:"message"`
	Success  bool   `json:"success"`
}

// === Source Query Parameters ===

// SourcesListParams represents parameters for listing sources
type SourcesListParams struct {
	NotebookID *string `json:"notebook_id,omitempty"`
	Limit      *int    `json:"limit,omitempty"`
	Offset     *int    `json:"offset,omitempty"`
	SortBy     *string `json:"sort_by,omitempty"`     // "created" | "updated"
	SortOrder  *string `json:"sort_order,omitempty"`  // "asc" | "desc"
	Type       *SourceType `json:"type,omitempty"`
	Status     *ProcessingStatus `json:"status,omitempty"`
}

// === Source Command Models ===

// SourceRetryRequest represents source retry request
type SourceRetryRequest struct {
	ForceReprocess bool `json:"force_reprocess,omitempty"`
}

// SourceRetryResponse represents source retry response
type SourceRetryResponse struct {
	Message   string `json:"message"`
	CommandID string `json:"command_id"`
}