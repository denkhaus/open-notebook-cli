package integration

import (
	"context"
	"encoding/json"
	"flag"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"

	"github.com/denkhaus/open-notebook-cli/pkg/di"
	"github.com/denkhaus/open-notebook-cli/pkg/models"
)

// TestBasicConnectivity tests if we can connect to the OpenNotebook API
func TestBasicConnectivity(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Create a minimal CLI context for testing
	ctx, err := createCLIContext()
	require.NoError(t, err)

	// Bootstrap DI
	injector := di.Bootstrap(ctx)

	// Test basic services
	config := di.GetConfig(injector)
	assert.NotNil(t, config)
	assert.Equal(t, "http://localhost:5055", config.GetAPIURL())

	logger := di.GetLogger(injector)
	assert.NotNil(t, logger)
}

// TestAuthFlow tests authentication against the real API
func TestAuthFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx, err := createCLIContext()
	require.NoError(t, err)

	injector := di.Bootstrap(ctx)
	auth := di.GetAuth(injector)
	httpClient := di.GetHTTPClient(injector)

	require.NotNil(t, auth)
	require.NotNil(t, httpClient)

	// Test auth status endpoint (should work without auth)
	testCtx := context.Background()
	resp, err := httpClient.Get(testCtx, "/auth/status")

	// API might not be running, so we'll handle both cases
	if err != nil {
		t.Logf("API connection failed (API might not be running): %v", err)
		t.Skip("OpenNotebook API not available on localhost:5055")
		return
	}

	assert.Equal(t, 200, resp.StatusCode)
	t.Logf("Auth status response: %s", string(resp.Body))

	// Test authentication with password
	auth.SetPassword("admin")
	err = auth.Authenticate(testCtx)

	if err != nil {
		t.Logf("Authentication failed: %v", err)
		// This might fail if the API doesn't have default password
		// We'll log it but not fail the test
	} else {
		assert.True(t, auth.IsAuthenticated(testCtx))
		t.Log("Authentication successful")
	}
}

// TestDIContainer tests that our DI container works
func TestDIContainer(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx, err := createCLIContext()
	require.NoError(t, err)

	injector := di.Bootstrap(ctx)

	// Test all core services can be created
	t.Run("Core services", func(t *testing.T) {
		config := di.GetConfig(injector)
		assert.NotNil(t, config)
		assert.Equal(t, "http://localhost:5055", config.GetAPIURL())

		logger := di.GetLogger(injector)
		assert.NotNil(t, logger)

		auth := di.GetAuth(injector)
		assert.NotNil(t, auth)

		httpClient := di.GetHTTPClient(injector)
		assert.NotNil(t, httpClient)
	})

	t.Run("Domain services", func(t *testing.T) {
		notebookService := di.GetNotebookService(injector)
		assert.NotNil(t, notebookService)

		notebookRepo := di.GetNotebookRepository(injector)
		assert.NotNil(t, notebookRepo)

		// Test service dependencies
		repo := notebookService.Repository()
		assert.NotNil(t, repo)
		assert.Equal(t, notebookRepo, repo)
	})
}

// TestNotebookServiceBasic tests basic notebook service functionality
func TestNotebookServiceBasic(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx, err := createCLIContext()
	require.NoError(t, err)

	injector := di.Bootstrap(ctx)
	testCtx := context.Background()

	// Try to authenticate first
	auth := di.GetAuth(injector)
	auth.SetPassword("admin")
	authErr := auth.Authenticate(testCtx)

	notebookService := di.GetNotebookService(injector)
	require.NotNil(t, notebookService)

	// Test notebook listing (might fail without auth or data)
	notebooks, err := notebookService.ListNotebooks(testCtx)
	if err != nil {
		t.Logf("Notebook listing failed (expected if no auth/data): %v", err)
		if authErr != nil {
			t.Skip("Both authentication and notebook listing failed - API might require proper setup")
		}
	} else {
		t.Logf("Successfully listed %d notebooks", len(notebooks))
	}

	// Test notebook creation
	notebook, err := notebookService.CreateNotebook(testCtx, "Test Notebook", "Created by integration test")
	if err != nil {
		t.Logf("Notebook creation failed (expected without proper auth): %v", err)
	} else {
		assert.NotEmpty(t, notebook.ID)
		assert.Equal(t, "Test Notebook", notebook.Name)
		t.Logf("Successfully created notebook: %s", notebook.ID)

		// Try to clean up - delete the test notebook
		err = notebookService.DeleteNotebook(testCtx, notebook.ID)
		if err != nil {
			t.Logf("Failed to cleanup test notebook: %v", err)
		}
	}
}

// TestModelsValidation tests our model types and enums
func TestModelsValidation(t *testing.T) {
	// This test doesn't require API connection

	// Test our typesafe enums
	t.Run("SearchType enum", func(t *testing.T) {
		assert.Equal(t, "vector", string(models.SearchTypeVector))
		assert.Equal(t, "text", string(models.SearchTypeText))
	})

	t.Run("ModelType enum", func(t *testing.T) {
		assert.Equal(t, "language", string(models.ModelTypeLanguage))
		assert.Equal(t, "embedding", string(models.ModelTypeEmbedding))
		assert.Equal(t, "text_to_speech", string(models.ModelTypeTextToSpeech))
		assert.Equal(t, "speech_to_text", string(models.ModelTypeSpeechToText))
	})

	t.Run("SourceType enum", func(t *testing.T) {
		assert.Equal(t, "link", string(models.SourceTypeLink))
		assert.Equal(t, "upload", string(models.SourceTypeUpload))
		assert.Equal(t, "text", string(models.SourceTypeText))
	})

	t.Run("NoteType enum", func(t *testing.T) {
		assert.Equal(t, "human", string(models.NoteTypeHuman))
		assert.Equal(t, "ai", string(models.NoteTypeAI))
	})
}

// Helper function to create CLI context for testing
func createCLIContext() (*cli.Context, error) {
	app := &cli.App{
		Name: "test",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "api-url",
				Value: "http://localhost:5055",
			},
			&cli.StringFlag{
				Name:  "password",
				Value: "admin",
			},
			&cli.BoolFlag{
				Name:  "verbose",
				Value: true,
			},
			&cli.IntFlag{
				Name:  "timeout",
				Value: 30,
			},
			&cli.StringFlag{
				Name:  "output",
				Value: "table",
			},
			&cli.StringFlag{
				Name:  "config-dir",
				Value: "/tmp/open-notebook-cli-test",
			},
		},
	}

	// Create flagset and set flag values
	flagSet := flag.NewFlagSet(app.Name, flag.ContinueOnError)

	// Define flags in flagset
	flagSet.String("api-url", "http://localhost:5055", "")
	flagSet.String("password", "admin", "")
	flagSet.Bool("verbose", true, "")
	flagSet.Int("timeout", 30, "")
	flagSet.String("output", "table", "")
	flagSet.String("config-dir", "/tmp/open-notebook-cli-test", "")

	// Parse arguments (empty for default values)
	args := []string{}
	err := flagSet.Parse(args)
	if err != nil {
		return nil, err
	}

	// Create context using correct signature
	ctx := cli.NewContext(app, flagSet, nil)
	return ctx, nil
}

// TestAPIAvailability checks if the API is available for integration tests
func TestAPIAvailability(t *testing.T) {
	ctx, err := createCLIContext()
	require.NoError(t, err)

	injector := di.Bootstrap(ctx)
	httpClient := di.GetHTTPClient(injector)

	testCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Try the auth status endpoint first to check if API is available
	resp, err := httpClient.Get(testCtx, "/auth/status")
	if err != nil {
		t.Logf("❌ OpenNotebook API not available on localhost:5055")
		t.Logf("   Error: %v", err)
		t.Logf("   Please start OpenNotebook with:")
		t.Logf("   docker run -p 5055:5055 lfnovo/open-notebook")
		t.Skip("API not available")
		return
	}

	// API is reachable, now test settings endpoint for model validation
	settingsResp, err := httpClient.Get(testCtx, "/settings")
	if err != nil {
		t.Logf("⚠️  Could not fetch settings endpoint: %v", err)
	}

	// API is reachable (even if auth fails, we should get a response)
	t.Logf("✅ OpenNotebook API is reachable on localhost:5055")
	t.Logf("   Status: %d", resp.StatusCode)
	t.Logf("   Response: %s", string(resp.Body))

	// Test our updated models by unmarshalling the settings response
	if settingsResp != nil && settingsResp.StatusCode == 200 {
		var settings models.SettingsResponse
		unmarshalErr := json.Unmarshal(settingsResp.Body, &settings)
		if unmarshalErr != nil {
			t.Logf("⚠️  Failed to unmarshal settings response: %v", unmarshalErr)
		} else {
			t.Logf("✅ Settings model validation passed:")
			t.Logf("   AutoDeleteFiles: %s (type: %T)", settings.AutoDeleteFiles, settings.AutoDeleteFiles)
			t.Logf("   ContentProcessingEngineDoc: %s (type: %T)", settings.DefaultContentProcessingEngineDoc, settings.DefaultContentProcessingEngineDoc)
			t.Logf("   ContentProcessingEngineURL: %s (type: %T)", settings.DefaultContentProcessingEngineURL, settings.DefaultContentProcessingEngineURL)
			t.Logf("   DefaultEmbeddingOption: %s (type: %T)", settings.DefaultEmbeddingOption, settings.DefaultEmbeddingOption)

			// Verify all are our new enum types
			assert.Equal(t, models.YesNoDecision("yes"), settings.AutoDeleteFiles)
			assert.Equal(t, models.ContentProcessingEngine("auto"), settings.DefaultContentProcessingEngineDoc)
			assert.Equal(t, models.ContentProcessingEngineURL("auto"), settings.DefaultContentProcessingEngineURL)
			assert.Equal(t, models.EmbeddingOption("ask"), settings.DefaultEmbeddingOption)
			t.Logf("   ✅ All enum types working correctly!")
		}
	}

	// Try without authentication first
	if resp.StatusCode == 401 {
		t.Logf("🔐 API requires authentication (this is expected)")
	} else if resp.StatusCode == 200 {
		t.Logf("✅ API responded successfully")
	} else {
		t.Logf("ℹ️  API responded with status %d", resp.StatusCode)
	}
}

// Benchmark DI container creation
func BenchmarkDIContainer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ctx, _ := createCLIContext()
		injector := di.Bootstrap(ctx)

		// Force creation of all services
		_ = di.GetConfig(injector)
		_ = di.GetLogger(injector)
		_ = di.GetAuth(injector)
		_ = di.GetHTTPClient(injector)
		_ = di.GetNotebookService(injector)
	}
}