package etfscraper

import (
	"net/http"
	"time"
)

// HTTPClient is the interface used to execute HTTP requests.
// It is satisfied by *http.Client and can be replaced for testing.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// HTTPConfig holds the HTTP transport settings shared by all providers.
type HTTPConfig struct {
	Client  HTTPClient
	Timeout time.Duration
	Debug   bool
}

// DefaultHTTPConfig returns an HTTPConfig with sensible defaults:
// a 15-second timeout and debug logging disabled.
func DefaultHTTPConfig() HTTPConfig {
	return HTTPConfig{
		Client:  &http.Client{Timeout: 15 * time.Second},
		Timeout: 15 * time.Second,
		Debug:   false,
	}
}
