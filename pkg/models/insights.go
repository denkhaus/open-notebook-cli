package models

// Insights API models

// InsightType represents insight type for sources
type InsightType string

const (
	InsightTypeSummary    InsightType = "summary"
	InsightTypeAnalysis   InsightType = "analysis"
	InsightTypeExtraction InsightType = "extraction"
	InsightTypeQuestion   InsightType = "question"
	InsightTypeReflection InsightType = "reflection"
)

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
