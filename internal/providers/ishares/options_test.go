package ishares

import (
	"testing"
	"time"

	"github.com/yevklym/etfscraper"
	"github.com/yevklym/etfscraper/internal/testutil"
)

func TestClientOptions(t *testing.T) {
	t.Run("default configuration", func(t *testing.T) {
		c, err := New("us")
		if err != nil {
			t.Fatal(err)
		}

		if c.httpConfig.Timeout != 15*time.Second {
			t.Errorf("expected default timeout 15s, got %v", c.httpConfig.Timeout)
		}

		if c.httpConfig.Debug != false {
			t.Error("expected debug to be false by default")
		}
	})

	t.Run("WithTimeout option", func(t *testing.T) {
		c, err := New("us", WithTimeout(30*time.Second))
		if err != nil {
			t.Fatal(err)
		}

		if c.httpConfig.Timeout != 30*time.Second {
			t.Errorf("expected timeout 30s, got %v", c.httpConfig.Timeout)
		}
	})

	t.Run("WithDebug option", func(t *testing.T) {
		c, err := New("us", WithDebug(true))
		if err != nil {
			t.Fatal(err)
		}

		if !c.httpConfig.Debug {
			t.Error("expected debug to be true")
		}
	})

	t.Run("multiple options", func(t *testing.T) {
		mockClient := &testutil.MockHTTPClient{StatusCode: 200}

		c, err := New("us",
			WithHTTPClient(mockClient),
			WithTimeout(45*time.Second),
			WithDebug(true),
		)
		if err != nil {
			t.Fatal(err)
		}

		if c.httpConfig.Client != mockClient {
			t.Error("expected custom HTTP client")
		}
		if c.httpConfig.Timeout != 45*time.Second {
			t.Errorf("expected timeout 45s, got %v", c.httpConfig.Timeout)
		}
		if !c.httpConfig.Debug {
			t.Error("expected debug to be true")
		}
	})

	t.Run("WithLogger option", func(t *testing.T) {
		nop := etfscraper.NopLogger()
		c, err := New("us", WithLogger(nop))
		if err != nil {
			t.Fatal(err)
		}

		if c.httpConfig.Logger != nop {
			t.Error("expected NopLogger to be set")
		}
	})

	t.Run("WithLogger nil ignored", func(t *testing.T) {
		c, err := New("us")
		if err != nil {
			t.Fatal(err)
		}
		original := c.httpConfig.Logger

		WithLogger(nil)(c)
		if c.httpConfig.Logger != original {
			t.Error("expected nil logger to be ignored")
		}
	})
}
