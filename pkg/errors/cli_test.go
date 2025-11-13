package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCLIError_Error(t *testing.T) {
	err := NewCLIError(ErrorTypeConfig, "test message")
	assert.Equal(t, "test message", err.Error())
}

func TestCLIError_Display(t *testing.T) {
	err := NewCLIError(ErrorTypeConfig, "test message")
	// Display() writes to stderr, so we just test that it doesn't panic
	assert.NotPanics(t, err.Display)
}

func TestNewCLIError(t *testing.T) {
	suggestions := []string{"suggestion 1", "suggestion 2"}
	err := NewCLIError(ErrorTypeAuth, "auth error", suggestions...)

	assert.Equal(t, ErrorTypeAuth, err.Type)
	assert.Equal(t, "auth error", err.Message)
	assert.Equal(t, suggestions, err.Suggestions)
	assert.Equal(t, 1, err.ExitCode)
	assert.NotEmpty(t, err.NextSteps)
	assert.NotEmpty(t, err.Examples)
}

func TestAddTypeSpecificGuidance_Config(t *testing.T) {
	err := NewCLIError(ErrorTypeConfig, "config error")

	assert.Contains(t, err.Suggestions, "Check your configuration files")
	assert.Contains(t, err.NextSteps, "Run 'onb --help' to see all available options")
	assert.Contains(t, err.Examples, "export OPEN_NOTEBOOK_API_URL=http://localhost:5055")
}

func TestAddTypeSpecificGuidance_Auth(t *testing.T) {
	err := NewCLIError(ErrorTypeAuth, "auth error")

	assert.Contains(t, err.Suggestions, "Check if API authentication is enabled")
	assert.Contains(t, err.NextSteps, "Test API connection with 'curl http://localhost:5055/auth/status'")
	assert.Contains(t, err.Examples, "onb --password admin notebooks list")
}

func TestAddTypeSpecificGuidance_Network(t *testing.T) {
	err := NewCLIError(ErrorTypeNetwork, "network error")

	assert.Contains(t, err.Suggestions, "Verify the OpenNotebook API is running")
	assert.Contains(t, err.NextSteps, "Test API availability: curl http://localhost:5055/api/notebooks")
	assert.NotEmpty(t, err.Examples) // Just check examples exist, content can vary
}

func TestAddTypeSpecificGuidance_API(t *testing.T) {
	err := NewCLIError(ErrorTypeAPI, "api error")

	assert.Contains(t, err.Suggestions, "Check API server logs for details")
	assert.Contains(t, err.NextSteps, "Check OpenNotebook logs: docker logs <container-name>")
	assert.Contains(t, err.Examples, "onb --timeout 60 --verbose search \"test query\"")
}

func TestAddTypeSpecificGuidance_Validation(t *testing.T) {
	err := NewCLIError(ErrorTypeValidation, "validation error")

	assert.Contains(t, err.Suggestions, "Check command syntax and required parameters")
	assert.Contains(t, err.NextSteps, "Run 'onb <command> --help' to see required parameters")
	assert.Contains(t, err.Examples, "onb notebooks create --help")
}

func TestAddTypeSpecificGuidance_Permission(t *testing.T) {
	err := NewCLIError(ErrorTypePermission, "permission error")

	assert.Contains(t, err.Suggestions, "Check user permissions for the requested operation")
	assert.Contains(t, err.NextSteps, "Check available notebooks: onb notebooks list")
	assert.Contains(t, err.Examples, "onb notebooks list --output json | jq '.[] | {id, name}'")
}

func TestAddTypeSpecificGuidance_NotFound(t *testing.T) {
	err := NewCLIError(ErrorTypeNotFound, "not found error")

	assert.Contains(t, err.Suggestions, "Verify the resource exists")
	assert.Contains(t, err.NextSteps, "List all notebooks: onb notebooks list")
	assert.Contains(t, err.Examples, "onb search 'example' --limit 10")
}

func TestAddTypeSpecificGuidance_Server(t *testing.T) {
	err := NewCLIError(ErrorTypeServer, "server error")

	assert.Contains(t, err.Suggestions, "Check OpenNotebook server status")
	assert.Contains(t, err.NextSteps, "Restart OpenNotebook server if needed")
	assert.Contains(t, err.Examples, "docker logs <open-notebook-container>")
}

func TestAddTypeSpecificGuidance_Usage(t *testing.T) {
	err := NewCLIError(ErrorTypeUsage, "usage error")

	assert.Contains(t, err.Suggestions, "Check command usage with --help")
	assert.Contains(t, err.NextSteps, "Run 'onb --help' to see all commands")
	assert.Contains(t, err.Examples, "onb --help")
}

func TestCategorizeError_Network(t *testing.T) {
	testCases := []struct {
		name     string
		errorMsg string
		expected string
	}{
		{"connection refused", "Get \"http://localhost:9999\": connection refused", "Cannot connect to OpenNotebook API"},
		{"no such host", "dial tcp: no such host", "Cannot connect to OpenNotebook API"},
		{"invalid port", "dial tcp: address 99999: invalid port", "Cannot connect to OpenNotebook API"},
		{"network unreachable", "network is unreachable", "Cannot connect to OpenNotebook API"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := CategorizeError(&testError{msg: tc.errorMsg}, nil)
			assert.Equal(t, ErrorTypeNetwork, err.Type)
			assert.Contains(t, err.Message, tc.expected)
		})
	}
}

func TestCategorizeError_Timeout(t *testing.T) {
	err := CategorizeError(&testError{msg: "context deadline exceeded"}, nil)
	assert.Equal(t, ErrorTypeNetwork, err.Type)
	assert.Contains(t, err.Message, "Request to OpenNotebook API timed out")
}

func TestCategorizeError_Auth(t *testing.T) {
	err := CategorizeError(&testError{msg: "HTTP 401: unauthorized"}, nil)
	assert.Equal(t, ErrorTypeAuth, err.Type)
	assert.Contains(t, err.Message, "Authentication failed")
}

func TestCategorizeError_Permission(t *testing.T) {
	err := CategorizeError(&testError{msg: "HTTP 403: forbidden"}, nil)
	assert.Equal(t, ErrorTypePermission, err.Type)
	assert.Contains(t, err.Message, "Permission denied")
}

func TestCategorizeError_NotFound(t *testing.T) {
	err := CategorizeError(&testError{msg: "HTTP 404: not found"}, nil)
	assert.Equal(t, ErrorTypeNotFound, err.Type)
	assert.Contains(t, err.Message, "Resource not found")
}

func TestCategorizeError_Server(t *testing.T) {
	testCases := []string{
		"HTTP 500: internal server error",
		"HTTP 502: bad gateway",
		"HTTP 503: service unavailable",
		"HTTP 504: gateway timeout",
	}

	for _, errMsg := range testCases {
		t.Run(errMsg, func(t *testing.T) {
			err := CategorizeError(&testError{msg: errMsg}, nil)
			assert.Equal(t, ErrorTypeServer, err.Type)
			assert.Contains(t, err.Message, "OpenNotebook server error")
		})
	}
}

func TestCategorizeError_RateLimit(t *testing.T) {
	err := CategorizeError(&testError{msg: "HTTP 429: too many requests"}, nil)
	assert.Equal(t, ErrorTypeAPI, err.Type)
	assert.Contains(t, err.Message, "Rate limit exceeded")
}

func TestCategorizeError_Generic(t *testing.T) {
	err := CategorizeError(&testError{msg: "some random error"}, nil)
	assert.Equal(t, ErrorTypeAPI, err.Type)
	assert.Contains(t, err.Message, "An error occurred")
	assert.Contains(t, err.Message, "some random error")
}

func TestConfigError(t *testing.T) {
	err := ConfigError("missing config", "use config file")
	assert.Equal(t, ErrorTypeConfig, err.Type)
	assert.Equal(t, "missing config", err.Message)
	assert.Contains(t, err.Suggestions, "use config file")
}

func TestAuthError(t *testing.T) {
	err := AuthError("invalid password", "check credentials")
	assert.Equal(t, ErrorTypeAuth, err.Type)
	assert.Equal(t, "invalid password", err.Message)
	assert.Contains(t, err.Suggestions, "check credentials")
}

func TestNetworkError(t *testing.T) {
	err := NetworkError("api down", "check server")
	assert.Equal(t, ErrorTypeNetwork, err.Type)
	assert.Equal(t, "api down", err.Message)
	assert.Contains(t, err.Suggestions, "check server")
}

func TestValidationError(t *testing.T) {
	err := ValidationError("invalid input", "validate data")
	assert.Equal(t, ErrorTypeValidation, err.Type)
	assert.Equal(t, "invalid input", err.Message)
	assert.Contains(t, err.Suggestions, "validate data")
}

func TestUsageError(t *testing.T) {
	err := UsageError("missing flag", "use --help")
	assert.Equal(t, ErrorTypeUsage, err.Type)
	assert.Equal(t, "missing flag", err.Message)
	assert.Contains(t, err.Suggestions, "use --help")
}

// testError is a simple error implementation for testing
type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}