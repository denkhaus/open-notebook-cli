package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/denkhaus/open-notebook-cli/pkg/models"
	"github.com/samber/do/v2"
)

type chatRepository struct {
	httpClient HTTPClient
	logger     Logger
}

// NewChatRepository creates a new chat repository
func NewChatRepository(injector do.Injector) (ChatRepository, error) {
	httpClient := do.MustInvoke[HTTPClient](injector)
	logger := do.MustInvoke[Logger](injector)

	return &chatRepository{
		httpClient: httpClient,
		logger:     logger,
	}, nil
}

// ListSessions implements ChatRepository interface
func (r *chatRepository) ListSessions(ctx context.Context) (*models.ChatSessionsResponse, error) {
	return r.ListSessionsForNotebook(ctx, "")
}

// ListSessionsForNotebook lists chat sessions for a specific notebook
func (r *chatRepository) ListSessionsForNotebook(ctx context.Context, notebookID string) (*models.ChatSessionsResponse, error) {
	r.logger.Info("Listing chat sessions", "notebook_id", notebookID)

	endpoint := "/chat/sessions"
	if notebookID != "" {
		endpoint += "?notebook_id=" + notebookID
	}

	resp, err := r.httpClient.Get(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to list chat sessions: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(resp.Body))
	}

	var response models.ChatSessionsResponse
	if err := json.Unmarshal(resp.Body, &response); err != nil {
		return nil, fmt.Errorf("failed to decode sessions response: %w", err)
	}

	r.logger.Info("Retrieved chat sessions", "count", len(response))
	return &response, nil
}

// CreateSession implements ChatRepository interface
func (r *chatRepository) CreateSession(ctx context.Context, req *models.ChatCreateRequest) (*models.ChatSession, error) {
	r.logger.Info("Creating chat session", "title", req.Title, "model", req.ModelID)

	endpoint := "/chat/sessions"
	resp, err := r.httpClient.Post(ctx, endpoint, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create chat session: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(resp.Body))
	}

	var session models.ChatSession
	if err := json.Unmarshal(resp.Body, &session); err != nil {
		return nil, fmt.Errorf("failed to decode session response: %w", err)
	}

	r.logger.Info("Created chat session", "session_id", session.ID)
	return &session, nil
}

// GetSession implements ChatRepository interface
func (r *chatRepository) GetSession(ctx context.Context, sessionID string) (*models.ChatSession, error) {
	r.logger.Info("Getting chat session", "session_id", sessionID)

	endpoint := fmt.Sprintf("/chat/sessions/%s", sessionID)
	resp, err := r.httpClient.Get(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat session: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(resp.Body))
	}

	var session models.ChatSession
	if err := json.Unmarshal(resp.Body, &session); err != nil {
		return nil, fmt.Errorf("failed to decode session response: %w", err)
	}

	r.logger.Info("Retrieved chat session", "session_id", session.ID)
	return &session, nil
}

// DeleteSession implements ChatRepository interface
func (r *chatRepository) DeleteSession(ctx context.Context, sessionID string) error {
	r.logger.Info("Deleting chat session", "session_id", sessionID)

	endpoint := fmt.Sprintf("/chat/sessions/%s", sessionID)
	resp, err := r.httpClient.Delete(ctx, endpoint)
	if err != nil {
		return fmt.Errorf("failed to delete chat session: %w", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("API error: %d - %s", resp.StatusCode, string(resp.Body))
	}

	r.logger.Info("Deleted chat session", "session_id", sessionID)
	return nil
}

// ExecuteChat implements ChatRepository interface
func (r *chatRepository) ExecuteChat(ctx context.Context, req *models.ChatExecuteRequest) (*models.ChatExecuteResponse, error) {
	r.logger.Info("Executing chat", "session_id", req.SessionID, "stream", req.Stream)

	endpoint := "/chat/execute"
	resp, err := r.httpClient.Post(ctx, endpoint, req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute chat: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(resp.Body))
	}

	var response models.ChatExecuteResponse
	if err := json.Unmarshal(resp.Body, &response); err != nil {
		return nil, fmt.Errorf("failed to decode chat response: %w", err)
	}

	r.logger.Info("Chat executed successfully", "session_id", response.SessionID, "message_id", response.MessageID)
	return &response, nil
}

// StreamChat implements ChatRepository interface
func (r *chatRepository) StreamChat(ctx context.Context, req *models.ChatExecuteRequest) (<-chan *models.StreamChunk, error) {
	r.logger.Info("Starting chat stream", "session_id", req.SessionID)

	endpoint := "/chat/execute"
	chunkChan, err := r.httpClient.Stream(ctx, endpoint, req)
	if err != nil {
		return nil, fmt.Errorf("failed to start chat stream: %w", err)
	}

	// Convert response chunks to model chunks
	resultChan := make(chan *models.StreamChunk)
	go func() {
		defer close(resultChan)

		for chunk := range chunkChan {
			var streamChunk models.StreamChunk
			if err := json.Unmarshal(chunk, &streamChunk); err != nil {
				// If JSON parsing fails, treat as raw content
				streamChunk = models.StreamChunk{
					Content: string(chunk),
					Done:    false,
				}
			}
			resultChan <- &streamChunk
		}
	}()

	r.logger.Info("Chat stream started")
	return resultChan, nil
}

// GetMessages implements ChatRepository interface
func (r *chatRepository) GetMessages(ctx context.Context, sessionID string) ([]*models.ChatMessage, error) {
	r.logger.Info("Getting chat messages", "session_id", sessionID)

	endpoint := fmt.Sprintf("/chat/sessions/%s/messages", sessionID)
	resp, err := r.httpClient.Get(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat messages: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(resp.Body))
	}

	var messages []*models.ChatMessage
	if err := json.Unmarshal(resp.Body, &messages); err != nil {
		return nil, fmt.Errorf("failed to decode messages response: %w", err)
	}

	r.logger.Info("Retrieved chat messages", "session_id", sessionID, "count", len(messages))
	return messages, nil
}
