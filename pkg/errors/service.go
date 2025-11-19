package errors

import (
	"fmt"
)

// ServiceErrorType represents different categories of service-level errors
type ServiceErrorType int

const (
	ServiceErrorTypeNetwork ServiceErrorType = iota
	ServiceErrorTypeAPI
	ServiceErrorTypeValidation
	ServiceErrorTypeNotFound
	ServiceErrorTypeServer
	ServiceErrorTypeTimeout
	ServiceErrorTypeAuth
	ServiceErrorTypeParsing
)

// ServiceError represents structured service-level errors with proper context
type ServiceError struct {
	Type      ServiceErrorType
	Operation string
	Resource  string
	Cause     error
	Context   map[string]interface{}
}

// Error implements the error interface
func (e *ServiceError) Error() string {
	baseMsg := fmt.Sprintf("%s %s failed", e.Operation, e.Resource)

	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", baseMsg, e.Cause)
	}
	return baseMsg
}

// Unwrap returns the underlying error
func (e *ServiceError) Unwrap() error {
	return e.Cause
}

// NewServiceError creates a new service error with proper context
func NewServiceError(errorType ServiceErrorType, operation, resource string, cause error) *ServiceError {
	return &ServiceError{
		Type:      errorType,
		Operation: operation,
		Resource:  resource,
		Cause:     cause,
		Context:   make(map[string]interface{}),
	}
}

// WithContext adds context information to the error
func (e *ServiceError) WithContext(key string, value interface{}) *ServiceError {
	e.Context[key] = value
	return e
}

// Service error constructors with consistent messaging patterns

// NetworkServiceError creates a network-related service error
func NetworkServiceError(operation, resource string, cause error) *ServiceError {
	return NewServiceError(ServiceErrorTypeNetwork, operation, resource, cause)
}

// APIServiceError creates an API-related service error
func APIServiceError(operation, resource string, cause error) *ServiceError {
	return NewServiceError(ServiceErrorTypeAPI, operation, resource, cause)
}

// ValidationServiceError creates a validation-related service error
func ValidationServiceError(operation, resource string, cause error) *ServiceError {
	return NewServiceError(ServiceErrorTypeValidation, operation, resource, cause)
}

// NotFoundServiceError creates a not found service error
func NotFoundServiceError(operation, resource string, cause error) *ServiceError {
	return NewServiceError(ServiceErrorTypeNotFound, operation, resource, cause)
}

// ServerServiceError creates a server-related service error
func ServerServiceError(operation, resource string, cause error) *ServiceError {
	return NewServiceError(ServiceErrorTypeServer, operation, resource, cause)
}

// TimeoutServiceError creates a timeout-related service error
func TimeoutServiceError(operation, resource string, cause error) *ServiceError {
	return NewServiceError(ServiceErrorTypeTimeout, operation, resource, cause)
}

// AuthServiceError creates an auth-related service error
func AuthServiceError(operation, resource string, cause error) *ServiceError {
	return NewServiceError(ServiceErrorTypeAuth, operation, resource, cause)
}

// ParsingServiceError creates a parsing-related service error
func ParsingServiceError(operation, resource string, cause error) *ServiceError {
	return NewServiceError(ServiceErrorTypeParsing, operation, resource, cause)
}

// Convenience functions for common service operations

// FailedTo creates standardized error messages for failed operations
func FailedTo(operation, resource string, cause error) *ServiceError {
	// Auto-detect error type based on the cause or use default
	if cause != nil {
		errMsg := cause.Error()

		// Network related
		if containsAny(errMsg, []string{"connection refused", "network", "dial", "timeout"}) {
			return NetworkServiceError(operation, resource, cause)
		}

		// HTTP status codes
		if containsAny(errMsg, []string{"404", "not found"}) {
			return NotFoundServiceError(operation, resource, cause)
		}
		if containsAny(errMsg, []string{"401", "403", "unauthorized", "forbidden"}) {
			return AuthServiceError(operation, resource, cause)
		}
		if containsAny(errMsg, []string{"500", "502", "503", "504"}) {
			return ServerServiceError(operation, resource, cause)
		}
		if containsAny(errMsg, []string{"429", "rate limit"}) {
			return APIServiceError(operation, resource, cause)
		}

		// Parsing related
		if containsAny(errMsg, []string{"json", "unmarshal", "decode", "parse", "invalid"}) {
			return ParsingServiceError(operation, resource, cause)
		}
	}

	// Default to API error
	return APIServiceError(operation, resource, cause)
}

// FailedToList creates standardized "failed to list" error
func FailedToList(resource string, cause error) *ServiceError {
	return FailedTo("list", resource, cause)
}

// FailedToCreate creates standardized "failed to create" error
func FailedToCreate(resource string, cause error) *ServiceError {
	return FailedTo("create", resource, cause)
}

// FailedToGet creates standardized "failed to get" error
func FailedToGet(resource string, cause error) *ServiceError {
	return FailedTo("get", resource, cause)
}

// FailedToUpdate creates standardized "failed to update" error
func FailedToUpdate(resource string, cause error) *ServiceError {
	return FailedTo("update", resource, cause)
}

// FailedToDelete creates standardized "failed to delete" error
func FailedToDelete(resource string, cause error) *ServiceError {
	return FailedTo("delete", resource, cause)
}

// FailedToExecute creates standardized "failed to execute" error
func FailedToExecute(resource string, cause error) *ServiceError {
	return FailedTo("execute", resource, cause)
}

// FailedToDecode creates standardized "failed to decode" error
func FailedToDecode(resource string, cause error) *ServiceError {
	return FailedTo("decode", resource, cause)
}

// FailedToStart creates standardized "failed to start" error
func FailedToStart(resource string, cause error) *ServiceError {
	return FailedTo("start", resource, cause)
}

// FailedToCancel creates standardized "failed to cancel" error
func FailedToCancel(resource string, cause error) *ServiceError {
	return FailedTo("cancel", resource, cause)
}

// Helper function to check if string contains any of the substrings
func containsAny(s string, substrings []string) bool {
	for _, substr := range substrings {
		if len(substr) > 0 && len(s) >= len(substr) {
			// Simple contains check
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
		}
	}
	return false
}

// WrapServiceError wraps any error into a ServiceError with context
func WrapServiceError(err error, operation, resource string) *ServiceError {
	if err == nil {
		return nil
	}

	// If it's already a ServiceError, just add context
	if svcErr, ok := err.(*ServiceError); ok {
		return svcErr
	}

	return FailedTo(operation, resource, err)
}