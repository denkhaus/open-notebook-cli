package errors

import (
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

// Error types for different categories of CLI errors
type ErrorType int

const (
	ErrorTypeConfig ErrorType = iota
	ErrorTypeAuth
	ErrorTypeNetwork
	ErrorTypeAPI
	ErrorTypeValidation
	ErrorTypePermission
	ErrorTypeNotFound
	ErrorTypeServer
	ErrorTypeUsage
)

// CLIError represents a structured CLI error with user guidance
type CLIError struct {
	Type        ErrorType
	Message     string
	Suggestions []string
	NextSteps   []string
	Examples    []string
	ExitCode    int
}

// Error implements the error interface
func (e *CLIError) Error() string {
	return e.Message
}

// NewCLIError creates a new CLI error with user-friendly guidance
func NewCLIError(errorType ErrorType, message string, suggestions ...string) *CLIError {
	error := &CLIError{
		Type:        errorType,
		Message:     message,
		Suggestions: suggestions,
		ExitCode:    1,
	}

	// Add type-specific suggestions and next steps
	error.addTypeSpecificGuidance()

	return error
}

// addTypeSpecificGuidance adds contextual help based on error type
func (e *CLIError) addTypeSpecificGuidance() {
	switch e.Type {
	case ErrorTypeConfig:
		if len(e.Suggestions) == 0 {
			e.Suggestions = []string{
				"Check your configuration files",
				"Verify environment variables are set",
				"Use --help to see configuration options",
			}
		}
		e.NextSteps = []string{
			"Run 'onb --help' to see all available options",
			"Check environment variables with 'env | grep OPEN_NOTEBOOK'",
			"Try using default configuration values",
		}
		e.Examples = []string{
			"export OPEN_NOTEBOOK_API_URL=http://localhost:5055",
			"onb --api-url http://localhost:5055 notebooks list",
		}

	case ErrorTypeAuth:
		if len(e.Suggestions) == 0 {
			e.Suggestions = []string{
				"Check if API authentication is enabled",
				"Verify your password is correct",
				"Ensure you have access to the OpenNotebook instance",
			}
		}
		e.NextSteps = []string{
			"Test API connection with 'curl http://localhost:5055/auth/status'",
			"Set password via environment: export OPEN_NOTEBOOK_PASSWORD=your-password",
			"Or use --password flag: onb --password your-password <command>",
		}
		e.Examples = []string{
			"onb --password admin notebooks list",
			"export OPEN_NOTEBOOK_PASSWORD=admin && onb notebooks list",
		}

	case ErrorTypeNetwork:
		if len(e.Suggestions) == 0 {
			e.Suggestions = []string{
				"Verify the OpenNotebook API is running",
				"Check network connectivity to the API server",
				"Confirm the API URL is correct",
			}
		}
		e.NextSteps = []string{
			"Test API availability: curl http://localhost:5055/api/notebooks",
			"Start OpenNotebook: docker run -p 5055:5055 lfnovo/open-notebook",
			"Check if API is running on different port: --api-url http://localhost:8080",
		}
		e.Examples = []string{
			"curl -f http://localhost:5055/api/notebooks || echo 'API not reachable'",
			"onb --api-url http://localhost:8080 notebooks list",
		}

	case ErrorTypeAPI:
		if len(e.Suggestions) == 0 {
			e.Suggestions = []string{
				"Check API server logs for details",
				"Verify API version compatibility",
				"Try the operation again later",
			}
		}
		e.NextSteps = []string{
			"Check OpenNotebook logs: docker logs <container-name>",
			"Verify API health: curl http://localhost:5055/health",
			"Try with increased timeout: --timeout 60",
		}
		e.Examples = []string{
			"onb --timeout 60 --verbose search \"test query\"",
			"curl -X GET http://localhost:5055/api/notebooks -H 'Accept: application/json'",
		}

	case ErrorTypeValidation:
		if len(e.Suggestions) == 0 {
			e.Suggestions = []string{
				"Check command syntax and required parameters",
				"Verify input data format and values",
				"Use --help for command-specific requirements",
			}
		}
		e.NextSteps = []string{
			"Run 'onb <command> --help' to see required parameters",
			"Check the OpenNotebook API documentation",
			"Validate input data before submitting",
		}
		e.Examples = []string{
			"onb notebooks create --help",
			"onb search --help",
			"onb sources create --help",
		}

	case ErrorTypePermission:
		if len(e.Suggestions) == 0 {
			e.Suggestions = []string{
				"Check user permissions for the requested operation",
				"Verify notebook access rights",
				"Contact OpenNotebook administrator",
			}
		}
		e.NextSteps = []string{
			"Check available notebooks: onb notebooks list",
			"Verify you have access to the specific notebook",
			"Contact administrator if you need additional permissions",
		}
		e.Examples = []string{
			"onb notebooks list --output json | jq '.[] | {id, name}'",
			"onb --verbose notebooks list",
		}

	case ErrorTypeNotFound:
		if len(e.Suggestions) == 0 {
			e.Suggestions = []string{
				"Verify the resource exists",
				"Check resource ID and spelling",
				"List available resources first",
			}
		}
		e.NextSteps = []string{
			"List all notebooks: onb notebooks list",
			"Search for notes: onb search 'your query'",
			"Check sources: onb sources list",
		}
		e.Examples = []string{
			"onb notebooks list",
			"onb search 'example' --limit 10",
			"onb notes list --notebook-id nb-123",
		}

	case ErrorTypeServer:
		if len(e.Suggestions) == 0 {
			e.Suggestions = []string{
				"Check OpenNotebook server status",
				"Verify server logs for detailed errors",
				"Try the operation again later",
			}
		}
		e.NextSteps = []string{
			"Restart OpenNotebook server if needed",
			"Check server logs for error details",
			"Contact support if issue persists",
		}
		e.Examples = []string{
			"docker logs <open-notebook-container>",
			"curl -X GET http://localhost:5055/health",
			"onb --verbose --timeout 60 <command>",
		}

	case ErrorTypeUsage:
		if len(e.Suggestions) == 0 {
			e.Suggestions = []string{
				"Check command usage with --help",
				"Verify required parameters are provided",
				"Ensure parameter values are in correct format",
			}
		}
		e.NextSteps = []string{
			"Run 'onb --help' to see all commands",
			"Run 'onb <command> --help' for command-specific help",
			"Check the examples below for correct usage",
		}
		e.Examples = []string{
			"onb --help",
			"onb notebooks --help",
			"onb search --help",
		}
	}
}

// Display formats and prints the error with user guidance
func (e *CLIError) Display() {
	fmt.Fprintf(os.Stderr, "\n❌ %s\n\n", e.Message)

	// Show suggestions
	if len(e.Suggestions) > 0 {
		fmt.Fprintf(os.Stderr, "💡 Suggestions:\n")
		for _, suggestion := range e.Suggestions {
			fmt.Fprintf(os.Stderr, "   • %s\n", suggestion)
		}
		fmt.Fprintln(os.Stderr)
	}

	// Show next steps
	if len(e.NextSteps) > 0 {
		fmt.Fprintf(os.Stderr, "🎯 Next steps:\n")
		for _, step := range e.NextSteps {
			fmt.Fprintf(os.Stderr, "   → %s\n", step)
		}
		fmt.Fprintln(os.Stderr)
	}

	// Show examples
	if len(e.Examples) > 0 {
		fmt.Fprintf(os.Stderr, "📋 Examples:\n")
		for _, example := range e.Examples {
			fmt.Fprintf(os.Stderr, "   %s\n", example)
		}
		fmt.Fprintln(os.Stderr)
	}

	// Add help hint
	fmt.Fprintf(os.Stderr, "💬 For help getting started with OpenNotebook CLI, run: onb --help\n")
	fmt.Fprintf(os.Stderr, "🔍 For command-specific help, run: onb <command> --help\n")
}

// HandleCLIError displays a CLI error and exits with appropriate code
func HandleCLIError(err error, ctx *cli.Context) {
	if err == nil {
		return
	}

	// Check if it's already a CLIError
	if cliErr, ok := err.(*CLIError); ok {
		cliErr.Display()
		os.Exit(cliErr.ExitCode)
	}

	// Categorize and wrap standard errors
	cliErr := CategorizeError(err, ctx)
	cliErr.Display()
	os.Exit(cliErr.ExitCode)
}

// CategorizeError converts standard errors to CLIErrors with appropriate guidance
func CategorizeError(err error, ctx *cli.Context) *CLIError {
	errMsg := err.Error()

	// Network errors
	if strings.Contains(errMsg, "connection refused") ||
	   strings.Contains(errMsg, "no such host") ||
	   strings.Contains(errMsg, "network is unreachable") ||
	   strings.Contains(errMsg, "invalid port") ||
	   strings.Contains(errMsg, "dial tcp") {
		return NewCLIError(ErrorTypeNetwork,
			"Cannot connect to OpenNotebook API",
			"Verify the OpenNotebook server is running and accessible")
	}

	// Timeout errors
	if strings.Contains(errMsg, "timeout") || strings.Contains(errMsg, "deadline exceeded") {
		// Check if it's a gateway timeout (504) first
		if strings.Contains(errMsg, "504") || strings.Contains(errMsg, "gateway timeout") {
			return NewCLIError(ErrorTypeServer,
				"OpenNotebook server error",
				"Gateway timeout occurred")
		}
		return NewCLIError(ErrorTypeNetwork,
			"Request to OpenNotebook API timed out",
			"Increase timeout with --timeout flag or check network connectivity")
	}

	// Authentication errors
	if strings.Contains(errMsg, "401") || strings.Contains(errMsg, "unauthorized") {
		return NewCLIError(ErrorTypeAuth,
			"Authentication failed",
			"Check your password and API authentication settings")
	}

	// Permission errors
	if strings.Contains(errMsg, "403") || strings.Contains(errMsg, "forbidden") {
		return NewCLIError(ErrorTypePermission,
			"Permission denied",
			"Verify you have access to the requested resource")
	}

	// Not found errors
	if strings.Contains(errMsg, "404") || strings.Contains(errMsg, "not found") {
		return NewCLIError(ErrorTypeNotFound,
			"Resource not found",
			"Check the resource ID and verify it exists")
	}

	// Server errors
	if strings.Contains(errMsg, "500") || strings.Contains(errMsg, "502") ||
	   strings.Contains(errMsg, "503") || strings.Contains(errMsg, "504") {
		return NewCLIError(ErrorTypeServer,
			"OpenNotebook server error",
			"Check server status and try again later")
	}

	// Rate limiting
	if strings.Contains(errMsg, "429") || strings.Contains(errMsg, "too many requests") {
		return NewCLIError(ErrorTypeAPI,
			"Rate limit exceeded",
			"Wait a moment and try again, or use retry functionality")
	}

	// Default: generic error
	return NewCLIError(ErrorTypeAPI,
		fmt.Sprintf("An error occurred: %s", errMsg),
		"Check the error message and try the suggested solutions")
}

// Convenience functions for common error types

// ConfigError creates a configuration-related error
func ConfigError(message string, suggestions ...string) *CLIError {
	return NewCLIError(ErrorTypeConfig, message, suggestions...)
}

// AuthError creates an authentication-related error
func AuthError(message string, suggestions ...string) *CLIError {
	return NewCLIError(ErrorTypeAuth, message, suggestions...)
}

// NetworkError creates a network-related error
func NetworkError(message string, suggestions ...string) *CLIError {
	return NewCLIError(ErrorTypeNetwork, message, suggestions...)
}

// ValidationError creates a validation-related error
func ValidationError(message string, suggestions ...string) *CLIError {
	return NewCLIError(ErrorTypeValidation, message, suggestions...)
}

// UsageError creates a usage-related error
func UsageError(message string, suggestions ...string) *CLIError {
	return NewCLIError(ErrorTypeUsage, message, suggestions...)
}

// APIError creates an API-related error
func APIError(message string, suggestions ...string) *CLIError {
	return NewCLIError(ErrorTypeAPI, message, suggestions...)
}

// NotFoundError creates a not found error
func NotFoundError(message string, suggestions ...string) *CLIError {
	return NewCLIError(ErrorTypeNotFound, message, suggestions...)
}