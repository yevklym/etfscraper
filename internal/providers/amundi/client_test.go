package amundi

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/yevklym/etfscraper/internal/testutil"
)

func TestNew(t *testing.T) {
	t.Run("unsupported region", func(t *testing.T) {
		_, err := New("invalid-region")
		if err == nil {
			t.Fatal("expected error for invalid region")
		}

		expectedMsg := "unsupported region"
		if !strings.Contains(err.Error(), expectedMsg) {
			t.Errorf("expected error containing %q, got %q", expectedMsg, err.Error())
		}
	})

	t.Run("supported regions", func(t *testing.T) {
		regions := []string{"de", "DE"}

		for _, region := range regions {
			_, err := New(region)
			if err != nil {
				t.Errorf("expected region %q to be supported, got error: %v", region, err)
			}
		}
	})
}

// sampleDiscoveryJSON is a minimal Amundi discovery response for cache tests.
const sampleDiscoveryJSON = `{
	"products": [
		{
			"productId": "LU1135865084",
			"productType": "PRODUCT",
			"characteristics": {
				"ISIN": "LU1135865084",
				"SHARE_MARKETING_NAME": "Amundi S&P 500 UCITS ETF",
				"MNEMO": "C500",
				"TER": 0.15,
				"CURRENCY": "EUR",
				"FUND_AUM": 1000,
				"ASSET_CLASS": "Equity",
				"DISTRIBUTION_POLICY": "Capitalisation"
			}
		}
	]
}`

func TestDiscoverCached_ReusesCache(t *testing.T) {
	mock := &testutil.CountingMockHTTPClient{
		MockHTTPClient: testutil.MockHTTPClient{
			ResponseBody: sampleDiscoveryJSON,
			StatusCode:   http.StatusOK,
		},
	}

	c, err := New("de", WithHTTPClient(mock))
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	ctx := context.Background()

	funds1, err := c.DiscoverETFs(ctx)
	if err != nil {
		t.Fatalf("first DiscoverETFs() failed: %v", err)
	}
	if mock.CallCount != 1 {
		t.Fatalf("expected 1 HTTP call after first discover, got %d", mock.CallCount)
	}

	funds2, err := c.DiscoverETFs(ctx)
	if err != nil {
		t.Fatalf("second DiscoverETFs() failed: %v", err)
	}
	if mock.CallCount != 1 {
		t.Fatalf("expected 1 HTTP call after second discover (cache hit), got %d", mock.CallCount)
	}

	if len(funds1) != len(funds2) {
		t.Fatalf("expected same number of funds, got %d and %d", len(funds1), len(funds2))
	}
}

func TestDiscoverCached_FundInfoReusesCacheFromDiscover(t *testing.T) {
	mock := &testutil.CountingMockHTTPClient{
		MockHTTPClient: testutil.MockHTTPClient{
			ResponseBody: sampleDiscoveryJSON,
			StatusCode:   http.StatusOK,
		},
	}

	c, err := New("de", WithHTTPClient(mock))
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	ctx := context.Background()

	_, err = c.DiscoverETFs(ctx)
	if err != nil {
		t.Fatalf("DiscoverETFs() failed: %v", err)
	}

	fund, err := c.FundInfo(ctx, "C500")
	if err != nil {
		t.Fatalf("FundInfo() failed: %v", err)
	}
	if fund.Ticker != "C500" {
		t.Fatalf("expected ticker C500, got %s", fund.Ticker)
	}
	if mock.CallCount != 1 {
		t.Fatalf("expected 1 HTTP call total, got %d", mock.CallCount)
	}
}

func TestDiscoverCached_ExpiredCacheRefetches(t *testing.T) {
	mock := &testutil.CountingMockHTTPClient{
		MockHTTPClient: testutil.MockHTTPClient{
			ResponseBody: sampleDiscoveryJSON,
			StatusCode:   http.StatusOK,
		},
	}

	c, err := New("de", WithHTTPClient(mock), WithCacheTTL(1*time.Millisecond))
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	ctx := context.Background()

	_, err = c.DiscoverETFs(ctx)
	if err != nil {
		t.Fatalf("first DiscoverETFs() failed: %v", err)
	}
	if mock.CallCount != 1 {
		t.Fatalf("expected 1 call, got %d", mock.CallCount)
	}

	time.Sleep(5 * time.Millisecond)

	_, err = c.DiscoverETFs(ctx)
	if err != nil {
		t.Fatalf("second DiscoverETFs() failed: %v", err)
	}
	if mock.CallCount != 2 {
		t.Fatalf("expected 2 calls after cache expiry, got %d", mock.CallCount)
	}
}

func TestDiscoverCached_ZeroTTLAlwaysFetches(t *testing.T) {
	mock := &testutil.CountingMockHTTPClient{
		MockHTTPClient: testutil.MockHTTPClient{
			ResponseBody: sampleDiscoveryJSON,
			StatusCode:   http.StatusOK,
		},
	}

	c, err := New("de", WithHTTPClient(mock), WithCacheTTL(0))
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	ctx := context.Background()

	for i := 0; i < 3; i++ {
		_, err = c.DiscoverETFs(ctx)
		if err != nil {
			t.Fatalf("DiscoverETFs() call %d failed: %v", i+1, err)
		}
	}

	if mock.CallCount != 3 {
		t.Fatalf("expected 3 HTTP calls with TTL=0, got %d", mock.CallCount)
	}
}

func TestDiscoverETFs_ReturnsCopy(t *testing.T) {
	mock := &testutil.CountingMockHTTPClient{
		MockHTTPClient: testutil.MockHTTPClient{
			ResponseBody: sampleDiscoveryJSON,
			StatusCode:   http.StatusOK,
		},
	}

	c, err := New("de", WithHTTPClient(mock))
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	ctx := context.Background()

	funds1, err := c.DiscoverETFs(ctx)
	if err != nil {
		t.Fatalf("DiscoverETFs() failed: %v", err)
	}

	if len(funds1) > 0 {
		funds1[0].Ticker = "MUTATED"
	}

	funds2, err := c.DiscoverETFs(ctx)
	if err != nil {
		t.Fatalf("second DiscoverETFs() failed: %v", err)
	}

	if len(funds2) > 0 && funds2[0].Ticker == "MUTATED" {
		t.Fatal("DiscoverETFs should return a copy; mutation leaked to cache")
	}
}

func TestWithCacheTTL(t *testing.T) {
	ttl := 10 * time.Minute
	c, err := New("de", WithCacheTTL(ttl))
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	if c.cacheTTL != ttl {
		t.Fatalf("expected cacheTTL %v, got %v", ttl, c.cacheTTL)
	}
}

func TestBuildFundIndex(t *testing.T) {
	mock := &testutil.MockHTTPClient{
		ResponseBody: sampleDiscoveryJSON,
		StatusCode:   http.StatusOK,
	}

	c, err := New("de", WithHTTPClient(mock))
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	ctx := context.Background()
	_, err = c.DiscoverETFs(ctx)
	if err != nil {
		t.Fatalf("DiscoverETFs() failed: %v", err)
	}

	// Lookup by ticker (case insensitive).
	fund, err := c.FundInfo(ctx, "c500")
	if err != nil {
		t.Fatalf("FundInfo by lowercase ticker failed: %v", err)
	}
	if fund.Ticker != "C500" {
		t.Fatalf("expected C500, got %s", fund.Ticker)
	}

	// Lookup by ISIN.
	fund, err = c.FundInfo(ctx, "LU1135865084")
	if err != nil {
		t.Fatalf("FundInfo by ISIN failed: %v", err)
	}
	if fund.ISIN != "LU1135865084" {
		t.Fatalf("expected LU1135865084, got %s", fund.ISIN)
	}
}

func TestHoldingsForFund_NilFund(t *testing.T) {
	c, err := New("de")
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	_, err = c.HoldingsForFund(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for nil fund")
	}
	if !strings.Contains(err.Error(), "cannot be nil") {
		t.Fatalf("expected 'cannot be nil' error, got: %v", err)
	}
}
