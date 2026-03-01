package etfscraper

import (
	"fmt"
	"testing"
)

func TestDefaultLogger(t *testing.T) {
	logger := DefaultLogger()
	if logger == nil {
		t.Fatal("DefaultLogger() returned nil")
	}

	// Should not panic when called.
	logger.Printf("test message: %d", 42)
}

func TestNopLogger(t *testing.T) {
	logger := NopLogger()
	if logger == nil {
		t.Fatal("NopLogger() returned nil")
	}

	// Should not panic when called.
	logger.Printf("this should be discarded: %d", 42)
}

func TestDefaultHTTPConfig_Logger(t *testing.T) {
	cfg := DefaultHTTPConfig()
	if cfg.Logger == nil {
		t.Fatal("DefaultHTTPConfig().Logger is nil, expected DefaultLogger")
	}
}

// capturingLogger records Printf calls for testing.
type capturingLogger struct {
	messages []string
}

func (l *capturingLogger) Printf(format string, v ...any) {
	l.messages = append(l.messages, fmt.Sprintf(format, v...))
}

func TestLoggerInterface(t *testing.T) {
	logger := &capturingLogger{}

	// Verify it satisfies the Logger interface.
	var _ Logger = logger

	logger.Printf("hello %s", "world")
	if len(logger.messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(logger.messages))
	}
	if logger.messages[0] != "hello world" {
		t.Errorf("message = %q, want %q", logger.messages[0], "hello world")
	}
}
