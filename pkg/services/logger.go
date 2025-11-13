package services

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/samber/do/v2"
	"github.com/denkhaus/open-notebook-cli/pkg/config"
)

// Private logger implementation
type logger struct {
	zap *zap.Logger
}

// NewLogger creates a new logger service with zap integration
func NewLogger(injector do.Injector) (Logger, error) {
	cfg := do.MustInvoke[config.Service](injector)

	// Configure zap based on verbosity
	var zapConfig zap.Config
	if cfg.IsVerbose() {
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		zapConfig = zap.NewProductionConfig()
		zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		zapConfig.EncoderConfig.TimeKey = "timestamp"
		zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	// Create the logger
	zapLogger, err := zapConfig.Build(
		zap.AddCallerSkip(1), // Skip the wrapper to show correct caller
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		return nil, err
	}

	return &logger{
		zap: zapLogger,
	}, nil
}

// Interface implementation

func (l *logger) Debug(msg string, fields ...interface{}) {
	l.zap.Sugar().Debugw(msg, fields...)
}

func (l *logger) Info(msg string, fields ...interface{}) {
	l.zap.Sugar().Infow(msg, fields...)
}

func (l *logger) Warn(msg string, fields ...interface{}) {
	l.zap.Sugar().Warnw(msg, fields...)
}

func (l *logger) Error(msg string, fields ...interface{}) {
	l.zap.Sugar().Errorw(msg, fields...)
}

func (l *logger) Fatal(msg string, fields ...interface{}) {
	l.zap.Sugar().Fatalw(msg, fields...)
}

func (l *logger) Sync() error {
	return l.zap.Sync()
}

func (l *logger) With(fields ...interface{}) Logger {
	sugar := l.zap.Sugar().With(fields...)
	return &logger{
		zap: sugar.Desugar(),
	}
}

func (l *logger) WithContext(ctx context.Context) Logger {
	// Extract common fields from context like request ID, user ID, etc.
	var fields []interface{}

	// Add trace ID if available
	if traceID := ctx.Value("trace_id"); traceID != nil {
		fields = append(fields, "trace_id", traceID)
	}

	// Add request ID if available
	if requestID := ctx.Value("request_id"); requestID != nil {
		fields = append(fields, "request_id", requestID)
	}

	// Add user ID if available
	if userID := ctx.Value("user_id"); userID != nil {
		fields = append(fields, "user_id", userID)
	}

	if len(fields) > 0 {
		return l.With(fields...)
	}

	return l
}

// Internal access to zap logger for advanced usage
func (l *logger) Zap() *zap.Logger {
	return l.zap
}

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