package mocks

import (
	"context"
	"crypto/rand"
	"fmt"
	"sync"
	"time"

	"github.com/denkhaus/open-notebook-cli/pkg/shared"
)

// MockBase provides common functionality for all mocks
type MockBase struct {
	mu             sync.RWMutex
	calls          map[string][]CallInfo
	delay          time.Duration
	errorOnCall    map[string]error
	shouldFail     bool
	failureError   error
}

// CallInfo records information about each method call
type CallInfo struct {
	Args    []interface{}
	Result  interface{}
	Error   error
	CalledAt time.Time
}

// NewMockBase creates a new mock base with optional delay for testing
func NewMockBase(delay time.Duration) *MockBase {
	return &MockBase{
		calls:       make(map[string][]CallInfo),
		delay:       delay,
		errorOnCall: make(map[string]error),
	}
}

// RecordCall records a method call with its arguments and results
func (m *MockBase) RecordCall(method string, args []interface{}, result interface{}, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	call := CallInfo{
		Args:     make([]interface{}, len(args)),
		Result:   result,
		Error:    err,
		CalledAt: time.Now(),
	}

	// Deep copy arguments to avoid mutation issues
	for i, arg := range args {
		call.Args[i] = arg
	}

	m.calls[method] = append(m.calls[method], call)
}

// GetCalls returns all calls made to a specific method
func (m *MockBase) GetCalls(method string) []CallInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	calls := make([]CallInfo, len(m.calls[method]))
	copy(calls, m.calls[method])
	return calls
}

// CallCount returns the number of times a method was called
func (m *MockBase) CallCount(method string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.calls[method])
}

// WasCalled returns true if a method was called
func (m *MockBase) WasCalled(method string) bool {
	return m.CallCount(method) > 0
}

// ClearCalls clears all recorded calls for a method or all methods
func (m *MockBase) ClearCalls(method ...string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(method) == 0 {
		m.calls = make(map[string][]CallInfo)
	} else {
		for _, methodName := range method {
			delete(m.calls, methodName)
		}
	}
}

// SetError sets an error to be returned on the next call to a method
func (m *MockBase) SetError(method string, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.errorOnCall[method] = err
}

// GetError gets and clears the error for a method
func (m *MockBase) GetError(method string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	err := m.errorOnCall[method]
	delete(m.errorOnCall, method)
	return err
}

// SetFailure causes all subsequent calls to fail with the given error
func (m *MockBase) SetFailure(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.shouldFail = true
	m.failureError = err
}

// ClearFailure clears failure mode
func (m *MockBase) ClearFailure() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.shouldFail = false
	m.failureError = nil
}

// simulateDelay simulates network or processing delay
func (m *MockBase) simulateDelay() {
	if m.delay > 0 {
		time.Sleep(m.delay)
	}
}

// checkFailure checks if the mock should fail
func (m *MockBase) checkFailure() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.shouldFail {
		return m.failureError
	}
	return nil
}

// MockLogger provides a mock implementation of Logger interface
type MockLogger struct {
	*MockBase
	logs    map[string][]LogEntry
	debug   bool
	logChan chan LogEntry
}

// LogEntry represents a single log entry
type LogEntry struct {
	Level   string
	Message string
	Fields  []interface{}
	Time    time.Time
}

// NewMockLogger creates a new mock logger
func NewMockLogger(debug bool) *MockLogger {
	return &MockLogger{
		MockBase: NewMockBase(0),
		logs:     make(map[string][]LogEntry),
		debug:    debug,
		logChan:  make(chan LogEntry, 1000),
	}
}

// Debug implements Logger interface
func (l *MockLogger) Debug(msg string, fields ...interface{}) {
	l.logEntry("DEBUG", msg, fields...)
}

// Info implements Logger interface
func (l *MockLogger) Info(msg string, fields ...interface{}) {
	l.logEntry("INFO", msg, fields...)
}

// Warn implements Logger interface
func (l *MockLogger) Warn(msg string, fields ...interface{}) {
	l.logEntry("WARN", msg, fields...)
}

// Error implements Logger interface
func (l *MockLogger) Error(msg string, fields ...interface{}) {
	l.logEntry("ERROR", msg, fields...)
}

// Fatal implements Logger interface
func (l *MockLogger) Fatal(msg string, fields ...interface{}) {
	l.logEntry("FATAL", msg, fields...)
}

// Sync implements Logger interface
func (l *MockLogger) Sync() error {
	return nil
}

// With implements Logger interface
func (l *MockLogger) With(fields ...interface{}) shared.Logger {
	// Return a copy with additional fields
	return l
}

// WithContext implements Logger interface
func (l *MockLogger) WithContext(ctx context.Context) shared.Logger {
	return l
}

// logEntry records a log entry
func (l *MockLogger) logEntry(level, message string, fields ...interface{}) {
	entry := LogEntry{
		Level:   level,
		Message: message,
		Fields:  fields,
		Time:    time.Now(),
	}

	l.mu.Lock()
	l.logs[level] = append(l.logs[level], entry)
	l.mu.Unlock()

	select {
	case l.logChan <- entry:
	default:
		// Channel is full, skip
	}
}

// GetLogs returns all logs for a specific level
func (l *MockLogger) GetLogs(level string) []LogEntry {
	l.mu.RLock()
	defer l.mu.RUnlock()

	logs := make([]LogEntry, len(l.logs[level]))
	copy(logs, l.logs[level])
	return logs
}

// GetAllLogs returns all logs
func (l *MockLogger) GetAllLogs() []LogEntry {
	l.mu.RLock()
	defer l.mu.RUnlock()

	var allLogs []LogEntry
	for _, levelLogs := range l.logs {
		allLogs = append(allLogs, levelLogs...)
	}
	return allLogs
}

// LogStream returns the log channel for real-time monitoring
func (l *MockLogger) LogStream() <-chan LogEntry {
	return l.logChan
}

// ClearLogs clears all recorded logs
func (l *MockLogger) ClearLogs() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logs = make(map[string][]LogEntry)
}

// Helper functions for mock implementations

// generateID generates a random ID for mock entities
func generateID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

// currentTime returns the current time for mock timestamps
func currentTime() time.Time {
	return time.Now().UTC()
}

// generateShortID generates a short random ID
func generateShortID() string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}