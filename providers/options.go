package providers

import (
	"net/http"
	"time"

	"github.com/yevklym/etfscraper"
)

type Option func(*providerOptions)

type providerOptions struct {
	httpConfig etfscraper.HTTPConfig
}

func WithHTTPConfig(cfg etfscraper.HTTPConfig) Option {
	return func(o *providerOptions) {
		if cfg.Client == nil {
			cfg.Client = o.httpConfig.Client
		}
		o.httpConfig = cfg
	}
}

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

func WithHTTPClient(client etfscraper.HTTPClient) Option {
	return func(o *providerOptions) {
		if client != nil {
			o.httpConfig.Client = client
		}
	}
}

func WithDebug(enabled bool) Option {
	return func(o *providerOptions) {
		o.httpConfig.Debug = enabled
	}
}
