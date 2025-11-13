package integration

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSSEScanner tests the Server-Sent Events scanner implementation
func TestSSEScanner(t *testing.T) {
	t.Run("Valid SSE data parsing", func(t *testing.T) {
		sseData := `data: {"message": "hello"}
data: {"progress": 50}
data: {"status": "complete"}

invalid line
data: {"final": true}`

		scanner := newSSEScanner(strings.NewReader(sseData))
		var parsedData []string

		for scanner.Scan() {
			if bytes := scanner.Bytes(); bytes != nil {
				parsedData = append(parsedData, string(bytes))
			}
		}

		assert.Len(t, parsedData, 4)
		assert.Equal(t, `{"message": "hello"}`, parsedData[0])
		assert.Equal(t, `{"progress": 50}`, parsedData[1])
		assert.Equal(t, `{"status": "complete"}`, parsedData[2])
		assert.Equal(t, `{"final": true}`, parsedData[3])
	})

	t.Run("Empty lines handling", func(t *testing.T) {
		sseData := `data: first message

data: second message

data: third message`

		scanner := newSSEScanner(strings.NewReader(sseData))
		var parsedData []string

		for scanner.Scan() {
			if bytes := scanner.Bytes(); bytes != nil {
				parsedData = append(parsedData, string(bytes))
			}
		}

		assert.Len(t, parsedData, 3)
		assert.Equal(t, "first message", parsedData[0])
		assert.Equal(t, "second message", parsedData[1])
		assert.Equal(t, "third message", parsedData[2])
	})

	t.Run("No data prefix handling", func(t *testing.T) {
		sseData := `regular line
data: {"valid": true}
another regular line
data: {"also_valid": true}`

		scanner := newSSEScanner(strings.NewReader(sseData))
		var parsedData []string

		for scanner.Scan() {
			if bytes := scanner.Bytes(); bytes != nil {
				parsedData = append(parsedData, string(bytes))
			}
		}

		// Should only extract lines with "data: " prefix
		assert.Len(t, parsedData, 2)
		assert.Equal(t, `{"valid": true}`, parsedData[0])
		assert.Equal(t, `{"also_valid": true}`, parsedData[1])
	})

	t.Run("Empty data lines", func(t *testing.T) {
		sseData := `data:

data: regular message
data:
data: another message`

		scanner := newSSEScanner(strings.NewReader(sseData))
		var parsedData []string

		for scanner.Scan() {
			if bytes := scanner.Bytes(); bytes != nil {
				parsedData = append(parsedData, string(bytes))
			}
		}

		// Should skip empty data lines
		assert.Len(t, parsedData, 2)
		assert.Equal(t, "regular message", parsedData[0])
		assert.Equal(t, "another message", parsedData[1])
	})
}

// TestStreamingServer tests streaming HTTP server scenarios
func TestStreamingServer(t *testing.T) {
	t.Run("Basic SSE response", func(t *testing.T) {
		// Create a test server that streams SSE data
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Connection", "keep-alive")
			w.WriteHeader(http.StatusOK)

			// Send SSE events
			fmt.Fprint(w, "data: {\"type\": \"start\", \"message\": \"Processing started\"}\n\n")
			if flusher, ok := w.(http.Flusher); ok {
				flusher.Flush()
			}

			time.Sleep(50 * time.Millisecond)
			fmt.Fprint(w, "data: {\"type\": \"progress\", \"value\": 50}\n\n")
			if flusher, ok := w.(http.Flusher); ok {
				flusher.Flush()
			}

			time.Sleep(50 * time.Millisecond)
			fmt.Fprint(w, "data: {\"type\": \"complete\", \"result\": \"success\"}\n\n")
			if flusher, ok := w.(http.Flusher); ok {
				flusher.Flush()
			}
		}))
		defer server.Close()

		// Test HTTP request to streaming endpoint
		resp, err := http.Get(server.URL)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "text/event-stream", resp.Header.Get("Content-Type"))

		// Read streaming response
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		responseStr := string(body)
		assert.Contains(t, responseStr, "data: {\"type\": \"start\"}")
		assert.Contains(t, responseStr, "data: {\"type\": \"progress\"}")
		assert.Contains(t, responseStr, "data: {\"type\": \"complete\"}")
	})

	t.Run("Streaming with context cancellation", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/event-stream")
			w.WriteHeader(http.StatusOK)

			// Send data slowly to allow cancellation
			for i := 0; i < 100; i++ {
				select {
				case <-r.Context().Done():
					// Client disconnected
					return
				default:
					fmt.Fprintf(w, "data: {\"chunk\": %d}\n\n", i)
					if flusher, ok := w.(http.Flusher); ok {
						flusher.Flush()
					}
					time.Sleep(10 * time.Millisecond)
				}
			}
		}))
		defer server.Close()

		// Create request with context that will be cancelled
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, "GET", server.URL, nil)
		require.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		// Read partial response
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		responseStr := string(body)
		// Should have received some chunks but not all 100
		assert.Contains(t, responseStr, "data: {\"chunk\":")
		assert.NotContains(t, responseStr, "data: {\"chunk\": 99\"}") // Should not reach the end
	})

	t.Run("Error handling in streaming", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Internal Server Error")
		}))
		defer server.Close()

		resp, err := http.Get(server.URL)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.NotEqual(t, "text/event-stream", resp.Header.Get("Content-Type"))
	})
}

// TestStreamingResponseFormat tests different streaming response formats
func TestStreamingResponseFormat(t *testing.T) {
	t.Run("AI response streaming format", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/event-stream")
			w.WriteHeader(http.StatusOK)

			// Simulate AI token streaming
			tokens := []string{"This", "is", "a", "streamed", "AI", "response"}

			for i, token := range tokens {
				eventType := "token"
				if i == len(tokens)-1 {
					eventType = "final"
				}
				fmt.Fprintf(w, "data: {\"type\": \"%s\", \"content\": \"%s\"}\n\n", eventType, token)
				if flusher, ok := w.(http.Flusher); ok {
					flusher.Flush()
				}
				time.Sleep(20 * time.Millisecond)
			}
		}))
		defer server.Close()

		resp, err := http.Get(server.URL)
		require.NoError(t, err)
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		responseStr := string(body)
		assert.Contains(t, responseStr, "\"type\": \"token\"")
		assert.Contains(t, responseStr, "\"type\": \"final\"")
		assert.Contains(t, responseStr, "\"content\": \"AI\"")
	})

	t.Run("Job progress streaming format", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/event-stream")
			w.WriteHeader(http.StatusOK)

			// Simulate job progress updates
			progressUpdates := []struct {
				step    string
				value   int
				status  string
			}{
				{"starting", 0, "running"},
				{"processing", 25, "running"},
				{"analyzing", 50, "running"},
				{"finalizing", 75, "running"},
				{"complete", 100, "completed"},
			}

			for _, update := range progressUpdates {
				fmt.Fprintf(w, "data: {\"step\": \"%s\", \"progress\": %d, \"status\": \"%s\", \"timestamp\": %d}\n\n",
					update.step, update.value, update.status, time.Now().Unix())
				if flusher, ok := w.(http.Flusher); ok {
					flusher.Flush()
				}
				time.Sleep(30 * time.Millisecond)
			}
		}))
		defer server.Close()

		resp, err := http.Get(server.URL)
		require.NoError(t, err)
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		responseStr := string(body)
		assert.Contains(t, responseStr, "\"step\": \"starting\"")
		assert.Contains(t, responseStr, "\"progress\": 100")
		assert.Contains(t, responseStr, "\"status\": \"completed\"")
		assert.Contains(t, responseStr, "\"timestamp\"")
	})
}

// TestStreamingPerformance tests streaming performance characteristics
func TestStreamingPerformance(t *testing.T) {
	t.Run("High frequency small messages", func(t *testing.T) {
		messageCount := 50

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/event-stream")
			w.WriteHeader(http.StatusOK)

			for i := 0; i < messageCount; i++ {
				fmt.Fprintf(w, "data: {\"id\": %d, \"data\": \"message %d\"}\n\n", i, i)
				if flusher, ok := w.(http.Flusher); ok {
					flusher.Flush()
				}
			}
		}))
		defer server.Close()

		start := time.Now()
		resp, err := http.Get(server.URL)
		require.NoError(t, err)
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		elapsed := time.Since(start)
		responseStr := string(body)

		// Verify all messages were received
		for i := 0; i < messageCount; i++ {
			assert.Contains(t, responseStr, fmt.Sprintf("\"id\": %d", i))
		}

		// Should complete quickly
		assert.Less(t, elapsed, 5*time.Second)
	})
}

// Copy SSE Scanner implementation for testing
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