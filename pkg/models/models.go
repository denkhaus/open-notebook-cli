package models

// Models API models

// ModelType represents AI model type with type safety
type ModelType string

const (
	ModelTypeLanguage     ModelType = "language"
	ModelTypeEmbedding    ModelType = "embedding"
	ModelTypeTextToSpeech ModelType = "text_to_speech"
	ModelTypeSpeechToText ModelType = "speech_to_text"
)

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

// ModelsListResponse represents models list response with pagination
type ModelsListResponse struct {
	Models []Model `json:"models"`
	Total  int     `json:"total"`
}
