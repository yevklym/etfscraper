package etfscraper

import (
	"log"
	"net/http"
	"time"
)

// HTTPClient is the interface used to execute HTTP requests.
// It is satisfied by *http.Client and can be replaced for testing.
type HTTPClient interface {
	// Do sends an HTTP request and returns an HTTP response.
	// Implementations should follow the same contract as http.Client.Do.
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
	// Client is the HTTP client used for requests. Defaults to an *http.Client
	// with a 15-second timeout.
	Client HTTPClient
	// Timeout is the HTTP request timeout. This is kept in sync with Client.Timeout
	// when using the default *http.Client or WithTimeout.
	Timeout time.Duration
	// Debug enables verbose diagnostic logging (request URLs and response metadata).
	Debug bool
	// Logger receives diagnostic and warning messages. Defaults to DefaultLogger.
	// Set to NopLogger() to silence all output.
	Logger Logger
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
