package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/denkhaus/open-notebook-cli/pkg/config"
	"github.com/denkhaus/open-notebook-cli/pkg/models"
	"github.com/denkhaus/open-notebook-cli/pkg/shared"
	"github.com/samber/do/v2"
)

// Private auth implementation
type auth struct {
	config   config.Service
	logger   shared.Logger
	http     shared.HTTPClient
	mu       sync.RWMutex
	token    string
	tokenEx  time.Time
	password string
}

// NewAuth creates a new auth service
func NewAuth(injector do.Injector) (shared.Auth, error) {
	cfg := do.MustInvoke[config.Service](injector)
	logger := do.MustInvoke[shared.Logger](injector)
	http := do.MustInvoke[shared.HTTPClient](injector)

	a := &auth{
		config: cfg,
		logger: logger,
		http:   http,
	}

	// Set password from config
	a.password = cfg.GetPassword()

	return a, nil
}

// Interface implementation

func (a *auth) Authenticate(ctx context.Context) error {
	if a.password == "" {
		return fmt.Errorf("no password provided")
	}

	// If we have a valid cached token, return success
	a.mu.RLock()
	if a.token != "" && time.Now().Before(a.tokenEx) {
		a.mu.RUnlock()
		return nil
	}
	a.mu.RUnlock()

	// Generate token hash for authentication
	tokenHash := a.generateTokenHash(a.password)

	// Create auth endpoint request
	endpoint := "/auth/status"

	resp, err := a.http.Get(ctx, endpoint)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("authentication failed: status %d", resp.StatusCode)
	}

	// Cache the token
	a.mu.Lock()
	a.token = tokenHash
	a.tokenEx = time.Now().Add(1 * time.Hour) // Cache for 1 hour
	a.mu.Unlock()

	a.logger.Debug("Authentication successful", "expires_at", a.tokenEx)
	return nil
}

func (a *auth) GetToken(ctx context.Context) (string, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.token == "" {
		return "", fmt.Errorf("no token available, authenticate first")
	}

	if time.Now().After(a.tokenEx) {
		return "", fmt.Errorf("token expired, re-authenticate")
	}

	return a.token, nil
}

func (a *auth) InvalidateToken(ctx context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.token = ""
	a.tokenEx = time.Time{}

	a.logger.Debug("Token invalidated")
	return nil
}

func (a *auth) IsAuthenticated(ctx context.Context) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.token != "" && time.Now().Before(a.tokenEx)
}

func (a *auth) RefreshToken(ctx context.Context) error {
	return a.Authenticate(ctx)
}

func (a *auth) SetPassword(password string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.password = password
	// Invalidate current token when password changes
	a.token = ""
	a.tokenEx = time.Time{}
}

// Helper methods

func (a *auth) generateTokenHash(password string) string {
	hash := sha256.Sum256([]byte(password + "open-notebook-auth"))
	return hex.EncodeToString(hash[:])
}

func (a *auth) setToken(token string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.token = token
	a.tokenEx = time.Now().Add(1 * time.Hour)
}

// HTTPClient decorator that adds authentication
type authenticatedHTTPClient struct {
	http shared.HTTPClient
	auth shared.Auth
}

func NewAuthenticatedHTTPClient(base shared.HTTPClient, auth shared.Auth) shared.HTTPClient {
	return &authenticatedHTTPClient{
		http: base,
		auth: auth,
	}
}

func (a *authenticatedHTTPClient) Get(ctx context.Context, endpoint string) (*models.Response, error) {
	if err := a.ensureAuthenticated(ctx); err != nil {
		return nil, err
	}
	return a.http.Get(ctx, endpoint)
}

func (a *authenticatedHTTPClient) Post(ctx context.Context, endpoint string, body interface{}) (*models.Response, error) {
	if err := a.ensureAuthenticated(ctx); err != nil {
		return nil, err
	}
	return a.http.Post(ctx, endpoint, body)
}

func (a *authenticatedHTTPClient) Put(ctx context.Context, endpoint string, body interface{}) (*models.Response, error) {
	if err := a.ensureAuthenticated(ctx); err != nil {
		return nil, err
	}
	return a.http.Put(ctx, endpoint, body)
}

func (a *authenticatedHTTPClient) Delete(ctx context.Context, endpoint string) (*models.Response, error) {
	if err := a.ensureAuthenticated(ctx); err != nil {
		return nil, err
	}
	return a.http.Delete(ctx, endpoint)
}

func (a *authenticatedHTTPClient) PostMultipart(ctx context.Context, endpoint string, fields map[string]string, files map[string]io.Reader) (*models.Response, error) {
	if err := a.ensureAuthenticated(ctx); err != nil {
		return nil, err
	}
	return a.http.PostMultipart(ctx, endpoint, fields, files)
}

func (a *authenticatedHTTPClient) Stream(ctx context.Context, endpoint string, body interface{}) (<-chan []byte, error) {
	if err := a.ensureAuthenticated(ctx); err != nil {
		return nil, err
	}
	return a.http.Stream(ctx, endpoint, body)
}

func (a *authenticatedHTTPClient) SetAuth(token string) {
	a.http.SetAuth(token)
}

func (a *authenticatedHTTPClient) WithTimeout(timeout time.Duration) shared.HTTPClient {
	return &authenticatedHTTPClient{
		http: a.http.WithTimeout(timeout),
		auth: a.auth,
	}
}

func (a *authenticatedHTTPClient) ensureAuthenticated(ctx context.Context) error {
	if !a.auth.IsAuthenticated(ctx) {
		return a.auth.Authenticate(ctx)
	}

	token, err := a.auth.GetToken(ctx)
	if err != nil {
		return err
	}

	a.http.SetAuth(token)
	return nil
}

// Mock implementation for testing
type mockAuth struct {
	mu       sync.RWMutex
	token    string
	password string
}

func NewMockAuth() shared.Auth {
	return &mockAuth{}
}

func (m *mockAuth) Authenticate(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.password == "" {
		return fmt.Errorf("no password set")
	}

	m.token = "mock-token-" + m.password
	return nil
}

func (m *mockAuth) GetToken(ctx context.Context) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.token == "" {
		return "", fmt.Errorf("no token available")
	}

	return m.token, nil
}

func (m *mockAuth) InvalidateToken(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.token = ""
	return nil
}

func (m *mockAuth) IsAuthenticated(ctx context.Context) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.token != ""
}

func (m *mockAuth) RefreshToken(ctx context.Context) error {
	return m.Authenticate(ctx)
}

func (m *mockAuth) SetPassword(password string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.password = password
	m.token = ""
}

func (m *mockAuth) SetToken(token string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.token = token
}
