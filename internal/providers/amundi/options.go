package amundi

import (
	"net/http"
	"time"

	"github.com/yevklym/etfscraper"
)

type ClientOption func(*Client)

const (
	// DefaultTimeout is the default HTTP timeout (15 seconds)
	DefaultTimeout = 15 * time.Second
)

// WithHTTPConfig sets the complete HTTP configuration.
//
// Example:
//
//	cfg := etfscraper.DefaultHTTPConfig()
//	cfg.Timeout = 30 * time.Second
//	client, _ := amundi.New("de", amundi.WithHTTPConfig(cfg))
func WithHTTPConfig(cfg etfscraper.HTTPConfig) ClientOption {
	return func(c *Client) {
		if cfg.Client == nil {
			cfg.Client = c.httpConfig.Client
		}
		c.httpConfig = cfg
	}
}

// WithTimeout sets the HTTP request timeout.
// If timeout <= 0, defaults to 15 seconds.
//
// Example:
//
//	client, _ := amundi.New("de", amundi.WithTimeout(30*time.Second))
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		if timeout <= 0 {
			timeout = DefaultTimeout
		}

		c.httpConfig.Timeout = timeout
		if hc, ok := c.httpConfig.Client.(*http.Client); ok {
			hc.Timeout = timeout
		}
	}
}

// WithHTTPClient sets a custom HTTP client implementation.
//
// Example:
//
//	customClient := &http.Client{Timeout: 60*time.Second}
//	client, _ := amundi.New("de", amundi.WithHTTPClient(customClient))
func WithHTTPClient(client etfscraper.HTTPClient) ClientOption {
	return func(c *Client) {
		if client != nil {
			c.httpConfig.Client = client
		}
	}
}

// WithDebug enables debug logging of HTTP requests and responses.
// Should only be used during development.
func WithDebug(enabled bool) ClientOption {
	return func(c *Client) {
		c.httpConfig.Debug = enabled
	}
}

// WithCacheTTL sets the time-to-live for the discovery cache.
// Cached fund data is reused across FundInfo and Holdings calls
// within this duration, avoiding repeated API requests.
// Default is 5 minutes. Set to 0 to disable caching.
func WithCacheTTL(ttl time.Duration) ClientOption {
	return func(c *Client) {
		c.cacheTTL = ttl
	}
}
