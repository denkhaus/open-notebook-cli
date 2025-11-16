package models

// Sources API models

// SourceType represents source type with type safety
type SourceType string

const (
	SourceTypeLink   SourceType = "link"
	SourceTypeUpload SourceType = "upload"
	SourceTypeText   SourceType = "text"
)

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
	ID             *string        `json:"id"`
	Title          *string        `json:"title"`
	Topics         []string       `json:"topics"`
	Asset          *AssetModel    `json:"asset"`
	FullText       *string        `json:"full_text"`
	Embedded       bool           `json:"embedded"`
	EmbeddedChunks int            `json:"embedded_chunks"`
	FileAvailable  *bool          `json:"file_available,omitempty"`
	Created        string         `json:"created"`
	Updated        string         `json:"updated"`
	CommandID      *string        `json:"command_id,omitempty"`
	Status         *SourceStatus  `json:"status,omitempty"`
	ProcessingInfo map[string]any `json:"processing_info,omitempty"`
	Notebooks      []string       `json:"notebooks,omitempty"`
}

// SourceListResponse represents source in list response
type SourceListResponse struct {
	ID             *string        `json:"id"`
	Title          *string        `json:"title"`
	Topics         []string       `json:"topics"`
	Asset          *AssetModel    `json:"asset"`
	Embedded       bool           `json:"embedded"`
	EmbeddedChunks int            `json:"embedded_chunks"`
	InsightsCount  int            `json:"insights_count"`
	Created        string         `json:"created"`
	Updated        string         `json:"updated"`
	FileAvailable  *bool          `json:"file_available,omitempty"`
	CommandID      *string        `json:"command_id,omitempty"`
	Status         *SourceStatus  `json:"status,omitempty"`
	ProcessingInfo map[string]any `json:"processing_info,omitempty"`
}

// SourcesListResponse represents sources list response
type SourcesListResponse []SourceListResponse

// SourceStatusResponse represents source status response
type SourceStatusResponse struct {
	Status         *SourceStatus  `json:"status,omitempty"`
	Message        string         `json:"message"`
	ProcessingInfo map[string]any `json:"processing_info,omitempty"`
	CommandID      *string        `json:"command_id,omitempty"`
}
