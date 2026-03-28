package xtrackers

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/yevklym/etfscraper"
	"github.com/yevklym/etfscraper/internal/testutil"
)

type trackingMockClient struct {
	callCount int
	mock      *testutil.MockHTTPClient
}

func (m *trackingMockClient) Do(req *http.Request) (*http.Response, error) {
	m.callCount++
	return m.mock.Do(req)
}

func TestFundInfo(t *testing.T) {
	client, err := New("uk",
		WithHTTPClient(&testutil.MockHTTPClient{
			StatusCode:   http.StatusOK,
			ResponseBody: discoveryResponseGB,
		}),
		withSkipBrowserFetch(),
	)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	ctx := context.Background()

	t.Run("valid identifier", func(t *testing.T) {
		fund, err := client.FundInfo(ctx, "IE00BK1PV551")
		if err != nil {
			t.Fatalf("FundInfo() failed: %v", err)
		}
		if fund == nil || fund.ISIN != "IE00BK1PV551" {
			t.Fatalf("expected fund with ISIN IE00BK1PV551")
		}
	})

	t.Run("lowercase identifier", func(t *testing.T) {
		fund, err := client.FundInfo(ctx, "ie00bk1pv551")
		if err != nil {
			t.Fatalf("FundInfo() failed: %v", err)
		}
		if fund == nil || fund.ISIN != "IE00BK1PV551" {
			t.Fatalf("expected fund matched by lowercase")
		}
	})

	t.Run("invalid identifier", func(t *testing.T) {
		_, err := client.FundInfo(ctx, "INVALID123")
		if err == nil {
			t.Fatalf("expected error for invalid fund")
		}
	})

	t.Run("empty identifier", func(t *testing.T) {
		_, err := client.FundInfo(ctx, "   ")
		if err == nil {
			t.Fatalf("expected error for empty identifier")
		}
	})
}

func TestCacheReuse(t *testing.T) {
	tracker := &trackingMockClient{
		mock: &testutil.MockHTTPClient{
			StatusCode:   http.StatusOK,
			ResponseBody: discoveryResponseGB,
		},
	}
	client, err := New("uk",
		WithHTTPClient(tracker),
		withSkipBrowserFetch(),
	)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	ctx := context.Background()

	// First call should trigger HTTP response
	funds1, err := client.DiscoverETFs(ctx)
	if err != nil {
		t.Fatalf("DiscoverETFs(1) failed: %v", err)
	}
	if tracker.callCount != 1 {
		t.Fatalf("expected 1 call, got %d", tracker.callCount)
	}

	// Second call within TTL should hit cache
	funds2, err := client.DiscoverETFs(ctx)
	if err != nil {
		t.Fatalf("DiscoverETFs(2) failed: %v", err)
	}
	if tracker.callCount != 1 {
		t.Fatalf("expected 1 call, got %d", tracker.callCount)
	}

	// Compare output slices
	if len(funds1) != len(funds2) {
		t.Fatalf("length mismatch: %d != %d", len(funds1), len(funds2))
	}

	// Force cache expiration
	client.mu.Lock()
	client.cachedAt = time.Now().Add(-10 * time.Minute)
	client.mu.Unlock()

	// Third call should trigger HTTP again
	_, err = client.DiscoverETFs(ctx)
	if err != nil {
		t.Fatalf("DiscoverETFs(3) failed: %v", err)
	}
	if tracker.callCount != 2 {
		t.Fatalf("expected 2 calls, got %d", tracker.callCount)
	}
}

func TestHoldingsStubs(t *testing.T) {
	client, err := New("uk")
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	ctx := context.Background()

	t.Run("holdings by identifier", func(t *testing.T) {
		_, err := client.Holdings(ctx, "IE00BK1PV551")
		if err != etfscraper.ErrHoldingsUnavailable {
			t.Fatalf("expected ErrHoldingsUnavailable, got: %v", err)
		}
	})

	t.Run("holdings by fund object", func(t *testing.T) {
		_, err := client.HoldingsForFund(ctx, nil)
		if err != etfscraper.ErrHoldingsUnavailable {
			t.Fatalf("expected ErrHoldingsUnavailable, got: %v", err)
		}
	})
}
