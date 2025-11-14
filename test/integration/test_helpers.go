package integration

import (
	"encoding/json"
	"io"

	"github.com/denkhaus/open-notebook-cli/pkg/models"
)

// Helper functions to bridge between our response format and standard interfaces

// ParseResponseJSON parses JSON from our []byte response body
func ParseResponseJSON(resp *models.Response, target interface{}) error {
	return json.Unmarshal(resp.Body, target)
}

// ResponseBodyReader converts []byte to io.Reader for compatibility
func ResponseBodyReader(resp *models.Response) io.Reader {
	if resp == nil {
		return nil
	}
	return &byteReader{data: resp.Body, pos: 0}
}

// byteReader implements io.Reader for []byte data
type byteReader struct {
	data []byte
	pos  int
}

func (r *byteReader) Read(p []byte) (n int, err error) {
	if r.pos >= len(r.data) {
		return 0, io.EOF
	}
	n = copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}

// TestResponse wraps our response to provide Close() method for compatibility
type TestResponse struct {
	*models.Response
}

// Close is a no-op for our response format (compatibility method)
func (r *TestResponse) Close() error {
	return nil
}

// Helper function to create string pointer
func StringPtr(s string) *string {
	return &s
}

// Helper function to check API availability
func IsAPIAvailable() bool {
	// This would be implemented to check if localhost:5055 is reachable
	// For now, return true to allow tests to run
	return true
}