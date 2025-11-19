package errors

import (
	"fmt"
)

// CommandError represents structured command-level errors with proper context
type CommandError struct {
	Type      string
	Message   string
	Command   string
	Argument  string
	Cause     error
	Context   map[string]interface{}
}

// Error implements the error interface
func (e *CommandError) Error() string {
	if e.Argument != "" {
		return fmt.Sprintf("%s: %s", e.Message, e.Argument)
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *CommandError) Unwrap() error {
	return e.Cause
}

// NewCommandError creates a new command error
func NewCommandError(errorType, message, command string) *CommandError {
	return &CommandError{
		Type:    errorType,
		Message: message,
		Command: command,
		Context: make(map[string]interface{}),
	}
}

// WithArgument adds argument context to the error
func (e *CommandError) WithArgument(arg string) *CommandError {
	e.Argument = arg
	return e
}

// WithContext adds additional context to the error
func (e *CommandError) WithContext(key string, value interface{}) *CommandError {
	e.Context[key] = value
	return e
}

// Command error constructors for common patterns

// MissingArgument creates a standardized missing argument error
func MissingArgument(argName, command string) *CommandError {
	return NewCommandError("missing_argument", "Missing "+argName, command).WithArgument(argName)
}

// TooManyArguments creates a standardized too many arguments error
func TooManyArguments(expected, command string) *CommandError {
	return NewCommandError("too_many_arguments", "Too many arguments. Expected only "+expected, command)
}

// InvalidArgument creates a standardized invalid argument error
func InvalidArgument(argName, command string) *CommandError {
	return NewCommandError("invalid_argument", "Invalid "+argName, command).WithArgument(argName)
}

// RequiredField creates a standardized required field error
func RequiredField(fieldName, command string) *CommandError {
	return NewCommandError("required_field", fieldName+" is required", command).WithArgument(fieldName)
}

// AtLeastOneField creates a standardized "at least one field" error
func AtLeastOneField(fieldNames []string, command string) *CommandError {
	return NewCommandError("at_least_one_field",
		fmt.Sprintf("You must specify at least one field to update (%s)", formatFieldList(fieldNames)),
		command)
}

// CommandExecution creates a standardized command execution error
func CommandExecution(operation, command string, cause error) *CommandError {
	return NewCommandError("execution_error",
		fmt.Sprintf("Failed to %s in command %s", operation, command),
		command)
}

// formatFieldList formats a slice of field names into a readable string
func formatFieldList(fields []string) string {
	if len(fields) == 0 {
		return ""
	}
	if len(fields) == 1 {
		return "--" + fields[0]
	}

	result := "--" + fields[0]
	for i := 1; i < len(fields)-1; i++ {
		result += ", --" + fields[i]
	}
	if len(fields) > 1 {
		result += " or --" + fields[len(fields)-1]
	}
	return result
}

// WrapCommandError wraps any error into a CommandError with context
func WrapCommandError(err error, operation, command string) *CommandError {
	if err == nil {
		return nil
	}

	// If it's already a CommandError, return it
	if cmdErr, ok := err.(*CommandError); ok {
		return cmdErr
	}

	return CommandExecution(operation, command, err)
}