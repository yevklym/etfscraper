package ishares

import (
	"net/http"
	"time"

	"github.com/yevklym/etfscraper"
)

type ClientOption func(*Client)

func WithHTTPConfig(cfg etfscraper.HTTPConfig) ClientOption {
	return func(c *Client) {
		c.httpConfig = cfg
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		if timeout <= 0 {
			timeout = 15 * time.Second
		}

		c.httpConfig.Timeout = timeout
		if hc, ok := c.httpConfig.Client.(*http.Client); ok {
			hc.Timeout = timeout
		}
	}
}

func WithHTTPClient(client etfscraper.HTTPClient) ClientOption {
	return func(c *Client) {
		if client != nil {
			c.httpConfig.Client = client
		}
	}
}

func WithDebug(enabled bool) ClientOption {
	return func(c *Client) {
		c.httpConfig.Debug = enabled
	}
}
