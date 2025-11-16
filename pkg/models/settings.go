package models

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
