package services

import "context"

// Mock implementation for testing
type mockLogger struct {
	messages []string
}

func NewMockLogger() Logger {
	return &mockLogger{
		messages: make([]string, 0),
	}
}

func (m *mockLogger) Debug(msg string, fields ...interface{}) {
	m.messages = append(m.messages, "DEBUG: "+msg)
}

func (m *mockLogger) Info(msg string, fields ...interface{}) {
	m.messages = append(m.messages, "INFO: "+msg)
}

func (m *mockLogger) Warn(msg string, fields ...interface{}) {
	m.messages = append(m.messages, "WARN: "+msg)
}

func (m *mockLogger) Error(msg string, fields ...interface{}) {
	m.messages = append(m.messages, "ERROR: "+msg)
}

func (m *mockLogger) Fatal(msg string, fields ...interface{}) {
	m.messages = append(m.messages, "FATAL: "+msg)
}

func (m *mockLogger) Sync() error {
	return nil
}

func (m *mockLogger) With(fields ...interface{}) Logger {
	return m
}

func (m *mockLogger) WithContext(ctx context.Context) Logger {
	return m
}

func (m *mockLogger) GetMessages() []string {
	return m.messages
}

func (m *mockLogger) ClearMessages() {
	m.messages = make([]string, 0)
}
