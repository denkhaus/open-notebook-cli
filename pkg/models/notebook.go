package models

// Notebook models from OpenNotebook API

// Notebook represents a notebook from the API
type Notebook struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Archived    bool   `json:"archived"`
	Created     string `json:"created"`
	Updated     string `json:"updated"`
	SourceCount int    `json:"source_count"`
	NoteCount   int    `json:"note_count"`
}

// NotebookCreate represents notebook creation request
type NotebookCreate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// NotebookUpdate represents notebook update request
type NotebookUpdate struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Archived    *bool   `json:"archived,omitempty"`
}
