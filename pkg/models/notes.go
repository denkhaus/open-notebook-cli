package models

// Notes API models

// NoteType represents note type with type safety
type NoteType string

const (
	NoteTypeHuman NoteType = "human"
	NoteTypeAI    NoteType = "ai"
)

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
