package xtrackers

import (
	"net/http"
	"time"

	"github.com/yevklym/etfscraper"
)

// ClientOption configures a Client.
type ClientOption func(*Client)

const (
	// DefaultTimeout is the default HTTP timeout (15 seconds).
	DefaultTimeout = 15 * time.Second
)

// WithHTTPConfig sets the complete HTTP configuration.
func WithHTTPConfig(cfg etfscraper.HTTPConfig) ClientOption {
	return func(c *Client) {
		if cfg.Client == nil {
			cfg.Client = c.httpConfig.Client
		}
		if cfg.Logger == nil {
			cfg.Logger = c.httpConfig.Logger
		}
		c.httpConfig = cfg
	}
}

// WithTimeout sets the HTTP request timeout.
// If timeout <= 0, defaults to 15 seconds.
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
func WithHTTPClient(client etfscraper.HTTPClient) ClientOption {
	return func(c *Client) {
		if client != nil {
			c.httpConfig.Client = client
		}
	}
}

// WithDebug enables debug logging of HTTP requests and responses.
func WithDebug(enabled bool) ClientOption {
	return func(c *Client) {
		c.httpConfig.Debug = enabled
	}
}

// WithLogger sets a custom logger for diagnostic output.
func WithLogger(logger etfscraper.Logger) ClientOption {
	return func(c *Client) {
		if logger != nil {
			c.httpConfig.Logger = logger
		}
	}
}

// WithCacheTTL sets the time-to-live for the discovery cache.
// Default is 5 minutes. Set to 0 to disable caching.
func WithCacheTTL(ttl time.Duration) ClientOption {
	return func(c *Client) {
		c.cacheTTL = ttl
	}
}

// withSkipBrowserFetch is an internal option used exclusively by unit tests
// to prevent launching a real headless browser during mocked mock responses.
func withSkipBrowserFetch() ClientOption {
	return func(c *Client) {
		c.skipBrowserFetch = true
	}
}
