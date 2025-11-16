package models

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

// Default Prompt API models

// DefaultPromptResponse represents default prompt response
type DefaultPromptResponse struct {
	TransformationInstructions string `json:"transformation_instructions"`
}

// DefaultPromptUpdate represents default prompt update request
type DefaultPromptUpdate struct {
	TransformationInstructions string `json:"transformation_instructions"`
}
