package models

// Job management models

// JobStatus represents background job status
type JobStatus struct {
	ID       string   `json:"id"`
	Status   string   `json:"status"` // queued, running, completed, failed
	Progress *float64 `json:"progress,omitempty"`
	Message  *string  `json:"message,omitempty"`
	Created  string   `json:"created"`
	Updated  *string  `json:"updated,omitempty"`
}

// JobsListResponse represents jobs list response
type JobsListResponse struct {
	Jobs []JobStatus `json:"jobs"`
}
