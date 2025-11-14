package utils

import (
	"strings"
	"time"
)

// TruncateString truncates a string to the specified maximum length,
// adding "..." if the string is longer than maxLen.
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// ParseCommaSeparated parses a comma-separated string into a slice of strings.
// Empty input returns an empty slice.
func ParseCommaSeparated(input string) []string {
	if input == "" {
		return []string{}
	}

	parts := strings.Split(input, ",")
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

// FormatTimestamp formats a timestamp string for display.
// If the timestamp is empty, returns "N/A".
func FormatTimestamp(timestamp string) string {
	if timestamp == "" {
		return "N/A"
	}

	// Try to parse as ISO 8601 format
	if t, err := time.Parse(time.RFC3339, timestamp); err == nil {
		return t.Format("2006-01-02 15:04:05")
	}

	// Try other common formats
	formats := []string{
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timestamp); err == nil {
			return t.Format("2006-01-02 15:04:05")
		}
	}

	// If parsing fails, return the original string
	return timestamp
}

// SafeDereferenceString safely dereferences a string pointer.
// If the pointer is nil, returns an empty string.
func SafeDereferenceString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// SafeDereferenceBool safely dereferences a bool pointer.
// If the pointer is nil, returns false.
func SafeDereferenceBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

// JoinNonEmpty joins non-empty strings with the specified separator.
func JoinNonEmpty(separator string, parts ...string) string {
	var nonEmpty []string
	for _, part := range parts {
		if part != "" {
			nonEmpty = append(nonEmpty, part)
		}
	}
	return strings.Join(nonEmpty, separator)
}

// Helper functions for Source model migration between V1 and V2

