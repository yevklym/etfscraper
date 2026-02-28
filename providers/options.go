package providers

import (
	"net/http"
	"time"

	"github.com/yevklym/etfscraper"
)

// Option configures a provider created via Open, OpenSpec, or OpenNameRegion.
type Option func(*providerOptions)

type providerOptions struct {
	httpConfig etfscraper.HTTPConfig
	cacheTTL   *time.Duration // nil means use provider default
}

// WithHTTPConfig applies a complete HTTP configuration to provider clients.
func WithHTTPConfig(cfg etfscraper.HTTPConfig) Option {
	return func(o *providerOptions) {
		if cfg.Client == nil {
			cfg.Client = o.httpConfig.Client
		}
		o.httpConfig = cfg
	}
}

// WithTimeout sets the HTTP request timeout for provider clients.
func WithTimeout(timeout time.Duration) Option {
	return func(o *providerOptions) {
		if timeout <= 0 {
			timeout = o.httpConfig.Timeout
		}

		o.httpConfig.Timeout = timeout
		if hc, ok := o.httpConfig.Client.(*http.Client); ok {
			hc.Timeout = timeout
		}
	}
}

// WithHTTPClient sets a custom HTTP client for provider requests.
func WithHTTPClient(client etfscraper.HTTPClient) Option {
	return func(o *providerOptions) {
		if client != nil {
			o.httpConfig.Client = client
		}
	}
}

// WithDebug enables provider debug logging.
func WithDebug(enabled bool) Option {
	return func(o *providerOptions) {
		o.httpConfig.Debug = enabled
	}
}

// WithCacheTTL sets the time-to-live for the provider's discovery cache.
// Cached fund data is reused across FundInfo and Holdings calls within
// this duration, avoiding repeated API requests.
// Default is 5 minutes. Set to 0 to disable caching.
func WithCacheTTL(ttl time.Duration) Option {
	return func(o *providerOptions) {
		o.cacheTTL = &ttl
	}
}
