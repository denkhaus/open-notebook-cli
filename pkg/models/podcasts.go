package models

// Podcast API models

// PodcastGenerationRequest represents podcast generation request
type PodcastGenerationRequest struct {
	SourceIDs   []string `json:"source_ids,omitempty"`
	NotebookIDs []string `json:"notebook_ids,omitempty"`
	Query       string   `json:"query,omitempty"`
	ModelID     *string  `json:"model_id,omitempty"`
	Voice       string   `json:"voice,omitempty"`
	Language    string   `json:"language,omitempty"`
	Style       string   `json:"style,omitempty"`
}

// PodcastGenerationResponse represents podcast generation response
type PodcastGenerationResponse struct {
	JobID   string `json:"job_id"`
	Message string `json:"message"`
}

// PodcastEpisode represents a podcast episode
type PodcastEpisode struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Duration    float64 `json:"duration"`
	AudioURL    string  `json:"audio_url"`
	JobID       string  `json:"job_id"`
	ModelID     string  `json:"model_id"`
	Voice       string  `json:"voice"`
	Language    string  `json:"language"`
	Style       string  `json:"style"`
	Created     string  `json:"created"`
	Updated     string  `json:"updated"`
}

// PodcastEpisodeResponse represents podcast episode response
type PodcastEpisodeResponse struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Duration    float64 `json:"duration"`
	AudioURL    string  `json:"audio_url"`
	JobID       string  `json:"job_id"`
	ModelID     string  `json:"model_id"`
	Voice       string  `json:"voice"`
	Language    string  `json:"language"`
	Style       string  `json:"style"`
	Created     string  `json:"created"`
	Updated     string  `json:"updated"`
}

// PodcastEpisodesListResponse represents episodes list response
type PodcastEpisodesListResponse struct {
	Episodes []PodcastEpisodeResponse `json:"episodes"`
	Total    int                      `json:"total"`
}

// PodcastJobStatus represents podcast generation job status (extends JobStatus)
type PodcastJobStatus struct {
	ID        string   `json:"id"`
	Status    string   `json:"status"`
	Progress  *float64 `json:"progress,omitempty"`
	Message   *string  `json:"message,omitempty"`
	EpisodeID *string  `json:"episode_id,omitempty"`
	Created   string   `json:"created"`
	Updated   *string  `json:"updated,omitempty"`
}
