package etfscraper

import (
	"net/http"
	"time"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type HTTPConfig struct {
	Client  HTTPClient
	Timeout time.Duration
	Debug   bool
}

func DefaultHTTPConfig() HTTPConfig {
	return HTTPConfig{
		Client:  &http.Client{Timeout: 15 * time.Second},
		Timeout: 15 * time.Second,
		Debug:   false,
	}
}
