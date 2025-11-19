package services

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/denkhaus/open-notebook-cli/pkg/config"
	"github.com/denkhaus/open-notebook-cli/pkg/models"
	"github.com/denkhaus/open-notebook-cli/pkg/shared"
	"github.com/samber/do/v2"
)

// Private HTTP client implementation
type httpService struct {
	config     config.Service
	logger     shared.Logger
	httpClient *http.Client
	authToken  string
}

// NewHTTPClient creates a new HTTP client service
func NewHTTPClient(injector do.Injector) (shared.HTTPClient, error) {
	cfg := do.MustInvoke[config.Service](injector)
	logger := do.MustInvoke[shared.Logger](injector)

	// Create HTTP client with configuration
	httpClient := &http.Client{
		Timeout: time.Duration(cfg.GetTimeout()) * time.Second,
	}

	return &httpService{
		config:     cfg,
		logger:     logger,
		httpClient: httpClient,
	}, nil
}

// Interface implementation

func (h *httpService) Get(ctx context.Context, endpoint string) (*models.Response, error) {
	return h.request(ctx, "GET", endpoint, nil)
}

func (h *httpService) Post(ctx context.Context, endpoint string, body interface{}) (*models.Response, error) {
	return h.request(ctx, "POST", endpoint, body)
}

func (h *httpService) Put(ctx context.Context, endpoint string, body interface{}) (*models.Response, error) {
	return h.request(ctx, "PUT", endpoint, body)
}

func (h *httpService) Delete(ctx context.Context, endpoint string) (*models.Response, error) {
	return h.request(ctx, "DELETE", endpoint, nil)
}

func (h *httpService) PostMultipart(ctx context.Context, endpoint string, fields map[string]string, files map[string]io.Reader) (*models.Response, error) {
	return h.requestMultipart(ctx, "POST", endpoint, fields, files)
}

func (h *httpService) Stream(ctx context.Context, endpoint string, body interface{}) (<-chan []byte, error) {
	// For streaming, we'll need to implement Server-Sent Events handling
	ch := make(chan []byte, 100) // Buffered channel

	go func() {
		defer close(ch)

		reqBody, err := h.marshalBody(body)
		if err != nil {
			h.logger.Error("Failed to marshal streaming request body", "error", err)
			ch <- []byte(fmt.Sprintf(`{"error": "Failed to marshal body: %s"}`, err.Error()))
			return
		}

		req, err := http.NewRequestWithContext(ctx, "POST", h.buildURL(endpoint), reqBody)
		if err != nil {
			h.logger.Error("Failed to create streaming request", "error", err)
			ch <- []byte(fmt.Sprintf(`{"error": "Failed to create request: %s"}`, err.Error()))
			return
		}

		h.setHeaders(req, false)

		resp, err := h.httpClient.Do(req)
		if err != nil {
			h.logger.Error("Streaming request failed", "error", err)
			ch <- []byte(fmt.Sprintf(`{"error": "Request failed: %s"}`, err.Error()))
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			h.logger.Error("Streaming request returned error status", "status", resp.StatusCode)
			ch <- []byte(fmt.Sprintf(`{"error": "HTTP %d"}`, resp.StatusCode))
			return
		}

		// Handle Server-Sent Events
		scanner := newSSEScanner(resp.Body)
		for scanner.Scan() {
			data := scanner.Bytes()
			if len(data) > 0 {
				ch <- data
			}
		}

		if err := scanner.Err(); err != nil {
			h.logger.Error("Error reading streaming response", "error", err)
			ch <- []byte(fmt.Sprintf(`{"error": "Stream error: %s"}`, err.Error()))
		}
	}()

	return ch, nil
}

func (h *httpService) SetAuth(token string) {
	h.authToken = token
}

func (h *httpService) WithTimeout(timeout time.Duration) shared.HTTPClient {
	newClient := *h.httpClient // Copy the HTTP client
	newClient.Timeout = timeout

	newHTTPService := *h // Copy the service
	newHTTPService.httpClient = &newClient

	return &newHTTPService
}

// Private helper methods

func (h *httpService) request(ctx context.Context, method, endpoint string, body interface{}) (*models.Response, error) {
	var reqBody io.Reader
	var err error

	if body != nil {
		reqBody, err = h.marshalBody(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, h.buildURL(endpoint), reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	h.setHeaders(req, false)

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Create response
	response := &models.Response{
		StatusCode: resp.StatusCode,
		Body:       respBody,
		Header:     resp.Header,
	}

	// Log request/response for debugging
	h.logger.Debug("HTTP request completed",
		"method", method,
		"endpoint", endpoint,
		"status", resp.StatusCode,
		"body_size", len(respBody),
	)

	return response, nil
}

func (h *httpService) requestMultipart(ctx context.Context, method, endpoint string, fields map[string]string, files map[string]io.Reader) (*models.Response, error) {
	reqBody, contentType, err := h.createMultipartBody(fields, files)
	if err != nil {
		return nil, fmt.Errorf("failed to create multipart body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, h.buildURL(endpoint), reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	h.setHeaders(req, true)
	req.Header.Set("Content-Type", contentType)

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Create response
	response := &models.Response{
		StatusCode: resp.StatusCode,
		Body:       respBody,
		Header:     resp.Header,
	}

	return response, nil
}

func (h *httpService) buildURL(endpoint string) string {
	baseURL := h.config.GetAPIURL()
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}
	if !strings.HasPrefix(endpoint, "/") {
		endpoint = "/" + endpoint
	}
	return baseURL + "api" + endpoint
}

func (h *httpService) setHeaders(req *http.Request, isMultipart bool) {
	// Set User-Agent
	req.Header.Set("User-Agent", "open-notebook-cli/1.0.0")

	// Set Authorization header if token is available
	if h.authToken != "" {
		req.Header.Set("Authorization", "Bearer "+h.authToken)
	}

	// Set common headers
	req.Header.Set("Accept", "application/json")
	
	// Set Content-Type for non-multipart requests
	if !isMultipart && req.Method != "GET" && req.Method != "DELETE" {
		req.Header.Set("Content-Type", "application/json")
	}
}

func (h *httpService) marshalBody(body interface{}) (io.Reader, error) {
	if body == nil {
		return nil, nil
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(jsonBody), nil
}

func (h *httpService) createMultipartBody(fields map[string]string, files map[string]io.Reader) (io.Reader, string, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add form fields
	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			return nil, "", err
		}
	}

	// Add files
	for fieldName, fileReader := range files {
		part, err := writer.CreateFormFile(fieldName, "upload")
		if err != nil {
			return nil, "", err
		}
		if _, err := io.Copy(part, fileReader); err != nil {
			return nil, "", err
		}
	}

	if err := writer.Close(); err != nil {
		return nil, "", err
	}

	contentType := writer.FormDataContentType()
	return &buf, contentType, nil
}

// SSE Scanner for Server-Sent Events
type sseScanner struct {
	scanner *bufio.Scanner
}

func newSSEScanner(r io.Reader) *sseScanner {
	return &sseScanner{
		scanner: bufio.NewScanner(r),
	}
}

func (s *sseScanner) Scan() bool {
	return s.scanner.Scan()
}

func (s *sseScanner) Bytes() []byte {
	line := s.scanner.Bytes()

	// Skip empty lines and "data:" prefix
	if len(line) == 0 {
		return nil
	}

	lineStr := string(line)
	if strings.HasPrefix(lineStr, "data: ") {
		return []byte(strings.TrimPrefix(lineStr, "data: "))
	}

	return nil
}

func (s *sseScanner) Err() error {
	return s.scanner.Err()
}
