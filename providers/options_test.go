package providers

import (
	"net/http"
	"testing"
	"time"

	"github.com/yevklym/etfscraper"
)

func TestWithTimeout(t *testing.T) {
	tests := []struct {
		name    string
		timeout time.Duration
		want    time.Duration
	}{
		{name: "positive timeout", timeout: 30 * time.Second, want: 30 * time.Second},
		{name: "zero keeps default", timeout: 0, want: 15 * time.Second},
		{name: "negative keeps default", timeout: -1, want: 15 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := providerOptions{httpConfig: etfscraper.DefaultHTTPConfig()}
			WithTimeout(tt.timeout)(&opts)
			if opts.httpConfig.Timeout != tt.want {
				t.Errorf("Timeout = %v, want %v", opts.httpConfig.Timeout, tt.want)
			}
		})
	}
}

func TestWithDebug(t *testing.T) {
	opts := providerOptions{httpConfig: etfscraper.DefaultHTTPConfig()}
	if opts.httpConfig.Debug {
		t.Fatal("expected debug to be false by default")
	}

	WithDebug(true)(&opts)
	if !opts.httpConfig.Debug {
		t.Fatal("expected debug to be true after WithDebug(true)")
	}

	WithDebug(false)(&opts)
	if opts.httpConfig.Debug {
		t.Fatal("expected debug to be false after WithDebug(false)")
	}
}

type mockHTTPClient struct{}

func (m *mockHTTPClient) Do(_ *http.Request) (*http.Response, error) {
	return nil, nil
}

func TestWithHTTPClient(t *testing.T) {
	opts := providerOptions{httpConfig: etfscraper.DefaultHTTPConfig()}

	mock := &mockHTTPClient{}
	WithHTTPClient(mock)(&opts)
	if opts.httpConfig.Client != mock {
		t.Fatal("expected custom HTTP client to be set")
	}

	// nil should be ignored
	WithHTTPClient(nil)(&opts)
	if opts.httpConfig.Client != mock {
		t.Fatal("expected nil to be ignored, client should remain unchanged")
	}
}

func TestWithHTTPConfig(t *testing.T) {
	opts := providerOptions{httpConfig: etfscraper.DefaultHTTPConfig()}
	originalClient := opts.httpConfig.Client

	cfg := etfscraper.HTTPConfig{
		Timeout: 60 * time.Second,
		Debug:   true,
	}
	WithHTTPConfig(cfg)(&opts)

	if opts.httpConfig.Timeout != 60*time.Second {
		t.Errorf("Timeout = %v, want 60s", opts.httpConfig.Timeout)
	}
	if !opts.httpConfig.Debug {
		t.Error("expected debug to be true")
	}
	// nil Client should be filled from existing config
	if opts.httpConfig.Client != originalClient {
		t.Error("expected nil Client to be preserved from original config")
	}
}

func TestWithCacheTTL(t *testing.T) {
	opts := providerOptions{httpConfig: etfscraper.DefaultHTTPConfig()}
	if opts.cacheTTL != nil {
		t.Fatal("expected cacheTTL to be nil by default")
	}

	WithCacheTTL(10 * time.Minute)(&opts)
	if opts.cacheTTL == nil || *opts.cacheTTL != 10*time.Minute {
		t.Errorf("cacheTTL = %v, want 10m", opts.cacheTTL)
	}

	WithCacheTTL(0)(&opts)
	if opts.cacheTTL == nil || *opts.cacheTTL != 0 {
		t.Errorf("cacheTTL = %v, want 0 (disable)", opts.cacheTTL)
	}
}
