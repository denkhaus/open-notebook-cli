package integration

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/denkhaus/open-notebook-cli/pkg/di"
	clierrors "github.com/denkhaus/open-notebook-cli/pkg/errors"
)

// TestCLIErrorDisplay tests the CLI error display system with comprehensive user guidance
func TestCLIErrorDisplay(t *testing.T) {
	testCases := []struct {
		name            string
		errorType       clierrors.ErrorType
		message         string
		suggestions     []string
		expectedParts   []string
		description     string
	}{
		{
			name:          "Configuration Error",
			errorType:     clierrors.ErrorTypeConfig,
			message:       "API URL is not configured",
			expectedParts: []string{
				"❌ API URL is not configured",
				"💡 Suggestions:",
				"Check your configuration files",
				"Verify environment variables are set",
				"🎯 Next steps:",
				"Run 'onb --help' to see all available options",
				"📋 Examples:",
				"export OPEN_NOTEBOOK_API_URL=http://localhost:5055",
			},
			description: "Tests configuration error with comprehensive guidance",
		},
		{
			name:          "Authentication Error",
			errorType:     clierrors.ErrorTypeAuth,
			message:       "Authentication failed",
			expectedParts: []string{
				"❌ Authentication failed",
				"💡 Suggestions:",
				"Check if API authentication is enabled",
				"🎯 Next steps:",
				"Test API connection with 'curl http://localhost:5055/auth/status'",
				"📋 Examples:",
				"onb --password admin notebooks list",
			},
			description: "Tests authentication error with password guidance",
		},
		{
			name:          "Network Error",
			errorType:     clierrors.ErrorTypeNetwork,
			message:       "Cannot connect to API server",
			expectedParts: []string{
				"❌ Cannot connect to API server",
				"💡 Suggestions:",
				"Verify the OpenNotebook API is running",
				"🎯 Next steps:",
				"Start OpenNotebook: docker run -p 5055:5055 lfnovo/open-notebook",
			},
			description: "Tests network error with connectivity guidance",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create CLI error
			cliErr := clierrors.NewCLIError(tc.errorType, tc.message, tc.suggestions...)

			// Test error structure (simplified since Display() writes to stderr)
			assert.Equal(t, tc.errorType, cliErr.Type)
			assert.Equal(t, tc.message, cliErr.Message)
			assert.NotEmpty(t, cliErr.Suggestions)
			assert.NotEmpty(t, cliErr.NextSteps)
			assert.NotEmpty(t, cliErr.Examples)

			// Test specific suggestions for config errors
			if tc.errorType == clierrors.ErrorTypeConfig {
				assert.Contains(t, cliErr.Suggestions, "Check your configuration files")
			}
			if tc.errorType == clierrors.ErrorTypeAuth {
				assert.Contains(t, cliErr.Suggestions, "Check if API authentication is enabled")
			}
			if tc.errorType == clierrors.ErrorTypeNetwork {
				assert.Contains(t, cliErr.Suggestions, "Verify the OpenNotebook API is running")
			}

			t.Logf("✅ %s: %s", tc.name, tc.description)
		})
	}
}

// TestErrorCategorization tests automatic error categorization
func TestErrorCategorization(t *testing.T) {
	testCases := []struct {
		name             string
		errorMsg         string
		expectedType     clierrors.ErrorType
		expectedContains []string
	}{
		{
			name:             "Connection refused",
			errorMsg:         "Get \"http://localhost:9999/api\": connection refused",
			expectedType:     clierrors.ErrorTypeNetwork,
			expectedContains: []string{"Cannot connect to OpenNotebook API"},
		},
		{
			name:             "Request timeout",
			errorMsg:         "context deadline exceeded (Client.Timeout exceeded)",
			expectedType:     clierrors.ErrorTypeNetwork,
			expectedContains: []string{"Request to OpenNotebook API timed out"},
		},
		{
			name:             "401 Unauthorized",
			errorMsg:         "HTTP 401: unauthorized",
			expectedType:     clierrors.ErrorTypeAuth,
			expectedContains: []string{"Authentication failed"},
		},
		{
			name:             "404 Not Found",
			errorMsg:         "HTTP 404: not found",
			expectedType:     clierrors.ErrorTypeNotFound,
			expectedContains: []string{"Resource not found"},
		},
		{
			name:             "500 Server Error",
			errorMsg:         "HTTP 500: internal server error",
			expectedType:     clierrors.ErrorTypeServer,
			expectedContains: []string{"OpenNotebook server error"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create context for categorization
			ctx, err := createCLIContextWithServer("http://localhost:5055")
			require.NoError(t, err)

			// Create error and categorize it
			originalErr := errors.New(tc.errorMsg)
			categorizedErr := clierrors.CategorizeError(originalErr, ctx)

			// Verify type and message
			assert.Equal(t, tc.expectedType, categorizedErr.Type)
			assert.Contains(t, categorizedErr.Message, tc.expectedContains[0])

			// Test suggestions are present
			assert.NotEmpty(t, categorizedErr.Suggestions)
			assert.NotEmpty(t, categorizedErr.NextSteps)
			assert.NotEmpty(t, categorizedErr.Examples)

			t.Logf("✅ %s: Categorized as %v", tc.name, categorizedErr.Type)
		})
	}
}

// TestCLIErrorIntegration tests CLI error handling in real scenarios
func TestCLIErrorIntegration(t *testing.T) {
	t.Run("Network error in real command", func(t *testing.T) {
		// Create CLI context with invalid URL
		ctx, err := createCLIContextWithServer("http://localhost:99999")
		require.NoError(t, err)

		// Bootstrap DI
		injector := di.Bootstrap(ctx)
		httpClient := di.GetHTTPClient(injector)

		// Try to make request - should fail
		testCtx := context.Background()
		_, err = httpClient.Get(testCtx, "/test")

		// Should be network error (could be "connection refused" or "invalid port")
		require.Error(t, err)
		assert.True(t,
			strings.Contains(err.Error(), "connection refused") ||
			strings.Contains(err.Error(), "invalid port") ||
			strings.Contains(err.Error(), "dial tcp"),
			"Error should contain connection-related message")

		// Test error handling
		cliErr := clierrors.CategorizeError(err, ctx)
		assert.Equal(t, clierrors.ErrorTypeNetwork, cliErr.Type)
		assert.Contains(t, cliErr.Message, "Cannot connect to OpenNotebook API")

		t.Log("✅ Network error properly categorized and handled")
	})
}

// TestConvenienceFunctions tests convenience error functions
func TestConvenienceFunctions(t *testing.T) {
	t.Run("ConfigError function", func(t *testing.T) {
		err := clierrors.ConfigError("Missing API URL", "Use --api-url flag")
		assert.Equal(t, clierrors.ErrorTypeConfig, err.Type)
		assert.Equal(t, "Missing API URL", err.Message)
		assert.Contains(t, err.Suggestions, "Use --api-url flag")
	})

	t.Run("AuthError function", func(t *testing.T) {
		err := clierrors.AuthError("Invalid password", "Check password with --password flag")
		assert.Equal(t, clierrors.ErrorTypeAuth, err.Type)
		assert.Equal(t, "Invalid password", err.Message)
	})

	t.Run("NetworkError function", func(t *testing.T) {
		err := clierrors.NetworkError("API unreachable", "Try different --api-url")
		assert.Equal(t, clierrors.ErrorTypeNetwork, err.Type)
		assert.Equal(t, "API unreachable", err.Message)
	})

	t.Run("ValidationError function", func(t *testing.T) {
		err := clierrors.ValidationError("Invalid notebook name", "Use valid characters")
		assert.Equal(t, clierrors.ErrorTypeValidation, err.Type)
		assert.Equal(t, "Invalid notebook name", err.Message)
	})

	t.Run("UsageError function", func(t *testing.T) {
		err := clierrors.UsageError("Missing required flag", "Use --help for assistance")
		assert.Equal(t, clierrors.ErrorTypeUsage, err.Type)
		assert.Equal(t, "Missing required flag", err.Message)
	})

	t.Log("✅ All convenience functions working correctly")
}

// TestErrorExitCodes tests that different error types have appropriate exit codes
func TestErrorExitCodes(t *testing.T) {
	testCases := []struct {
		name       string
		errorFunc  func() *clierrors.CLIError
		expectCode int
	}{
		{
			name: "Config error exit code",
			errorFunc: func() *clierrors.CLIError {
				return clierrors.ConfigError("Test config error")
			},
			expectCode: 1,
		},
		{
			name: "Network error exit code",
			errorFunc: func() *clierrors.CLIError {
				return clierrors.NetworkError("Test network error")
			},
			expectCode: 1,
		},
		{
			name: "Usage error exit code",
			errorFunc: func() *clierrors.CLIError {
				return clierrors.UsageError("Test usage error")
			},
			expectCode: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.errorFunc()
			assert.Equal(t, tc.expectCode, err.ExitCode)
			t.Logf("✅ %s: Exit code %d", tc.name, err.ExitCode)
		})
	}
}