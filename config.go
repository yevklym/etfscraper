package etfscraper

import (
	"log"
	"net/http"
	"time"
)

// HTTPClient is the interface used to execute HTTP requests.
// It is satisfied by *http.Client and can be replaced for testing.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Logger is the interface used for diagnostic output.
// It matches the signature of log.Printf and slog-style adapters.
type Logger interface {
	Printf(format string, v ...any)
}

// DefaultLogger returns a Logger that writes to the standard log package.
func DefaultLogger() Logger {
	return &stdLogger{}
}

type stdLogger struct{}

func (l *stdLogger) Printf(format string, v ...any) {
	log.Printf(format, v...)
}

// NopLogger returns a Logger that discards all output.
func NopLogger() Logger {
	return &nopLogger{}
}

type nopLogger struct{}

func (l *nopLogger) Printf(string, ...any) {}

// HTTPConfig holds the HTTP transport settings shared by all providers.
type HTTPConfig struct {
	Client  HTTPClient
	Timeout time.Duration
	Debug   bool
	Logger  Logger
}

// DefaultHTTPConfig returns an HTTPConfig with sensible defaults:
// a 15-second timeout, debug logging disabled, and the standard logger.
func DefaultHTTPConfig() HTTPConfig {
	return HTTPConfig{
		Client:  &http.Client{Timeout: 15 * time.Second},
		Timeout: 15 * time.Second,
		Debug:   false,
		Logger:  DefaultLogger(),
	}
}
