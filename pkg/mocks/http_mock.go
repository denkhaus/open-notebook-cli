package mocks

import (
	"context"
	"io"
	"time"

	"github.com/denkhaus/open-notebook-cli/pkg/models"
	"github.com/denkhaus/open-notebook-cli/pkg/shared"
)

// MockHTTPClient is a mock implementation of HTTPClient for testing
type MockHTTPClient struct {
	responses map[string]*models.Response
}

// NewMockHTTPClient creates a new mock HTTP client for testing
func NewMockHTTPClient() shared.HTTPClient {
	return &MockHTTPClient{
		responses: make(map[string]*models.Response),
	}
}

// SetMockResponse sets a mock response for a specific endpoint
func (m *MockHTTPClient) SetMockResponse(endpoint string, response *models.Response) {
	m.responses[endpoint] = response
}

// Get performs a mock HTTP GET request
func (m *MockHTTPClient) Get(ctx context.Context, endpoint string) (*models.Response, error) {
	if resp, ok := m.responses[endpoint]; ok {
		return resp, nil
	}
	return &models.Response{
		StatusCode: 404,
		Body:       []byte(`{"error": "not found"}`),
	}, nil
}

// Post performs a mock HTTP POST request
func (m *MockHTTPClient) Post(ctx context.Context, endpoint string, body interface{}) (*models.Response, error) {
	if resp, ok := m.responses[endpoint]; ok {
		return resp, nil
	}
	return &models.Response{
		StatusCode: 404,
		Body:       []byte(`{"error": "not found"}`),
	}, nil
}

// Put performs a mock HTTP PUT request
func (m *MockHTTPClient) Put(ctx context.Context, endpoint string, body interface{}) (*models.Response, error) {
	if resp, ok := m.responses[endpoint]; ok {
		return resp, nil
	}
	return &models.Response{
		StatusCode: 404,
		Body:       []byte(`{"error": "not found"}`),
	}, nil
}

// Delete performs a mock HTTP DELETE request
func (m *MockHTTPClient) Delete(ctx context.Context, endpoint string) (*models.Response, error) {
	if resp, ok := m.responses[endpoint]; ok {
		return resp, nil
	}
	return &models.Response{
		StatusCode: 404,
		Body:       []byte(`{"error": "not found"}`),
	}, nil
}

// PostMultipart performs a mock HTTP multipart POST request
func (m *MockHTTPClient) PostMultipart(ctx context.Context, endpoint string, fields map[string]string, files map[string]io.Reader) (*models.Response, error) {
	if resp, ok := m.responses[endpoint]; ok {
		return resp, nil
	}
	return &models.Response{
		StatusCode: 404,
		Body:       []byte(`{"error": "not found"}`),
	}, nil
}

// Stream performs a mock HTTP streaming request
func (m *MockHTTPClient) Stream(ctx context.Context, endpoint string, body interface{}) (<-chan []byte, error) {
	ch := make(chan []byte, 1)

	go func() {
		defer close(ch)
		if resp, ok := m.responses[endpoint]; ok {
			ch <- resp.Body
		} else {
			ch <- []byte(`{"error": "not found"}`)
		}
	}()

	return ch, nil
}

// SetAuth mock implementation
func (m *MockHTTPClient) SetAuth(token string) {
	// Mock implementation - no-op
}

// WithTimeout mock implementation
func (m *MockHTTPClient) WithTimeout(timeout time.Duration) shared.HTTPClient {
	return m
}
