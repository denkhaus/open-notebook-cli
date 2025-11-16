package models

// Search models from OpenNotebook API

// SearchType represents search type with type safety
type SearchType string

const (
	SearchTypeVector SearchType = "vector"
	SearchTypeText   SearchType = "text"
)

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
