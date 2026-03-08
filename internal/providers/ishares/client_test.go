package ishares

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
		regions := []string{"us", "de", "US", "DE"}

		for _, region := range regions {
			_, err := New(region)
			if err != nil {
				t.Errorf("expected region %q to be supported, got error: %v", region, err)
			}
		}
	})
}

func TestFundInfo(t *testing.T) {
	t.Run("fund found", func(t *testing.T) {
		sampleJSON := `{
        "239619": {
            "fundName": "iShares MSCI China ETF",
            "localExchangeTicker": "MCHI",
            "isin": "US4642874659",
            "productType": "ISHARES_FUND_DATA",
            "inceptionDate": {"r": 20181329},
            "fees": {"r": 0.59},
            "totalNetAssets": {"r": 7779083697.85},
            "portfolioId": 239619,
            "productPageUrl": ":/us/products/239619/test"
        }
    }`

		mockClient := &testutil.MockHTTPClient{
			ResponseBody: sampleJSON,
			StatusCode:   http.StatusOK,
		}

		c, _ := New("us", WithHTTPClient(mockClient))

		fund, err := c.FundInfo(context.Background(), "MCHI")
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if fund.Ticker != "MCHI" {
			t.Errorf("expected ticker MCHI, got %s", fund.Ticker)
		}
	})

	t.Run("fund not found", func(t *testing.T) {
		sampleJSON := `{"239619": {"localExchangeTicker": "MCHI"}}`

		mockClient := &testutil.MockHTTPClient{
			ResponseBody: sampleJSON,
			StatusCode:   http.StatusOK,
		}

		c, _ := New("us", WithHTTPClient(mockClient))

		_, err := c.FundInfo(context.Background(), "NOTFOUND")
		if err == nil {
			t.Fatal("expected error for non-existent fund")
		}
	})

	t.Run("empty identifier", func(t *testing.T) {
		c, _ := New("us")

		_, err := c.FundInfo(context.Background(), "")
		if err == nil {
			t.Fatal("expected error for empty identifier")
		}

		if !strings.Contains(err.Error(), "cannot be empty") {
			t.Errorf("expected 'cannot be empty' error, got: %v", err)
		}
	})

	t.Run("whitespace identifier", func(t *testing.T) {
		c, _ := New("us")

		_, err := c.FundInfo(context.Background(), "   ")
		if err == nil {
			t.Fatal("expected error for whitespace identifier")
		}
	})
}

// sampleDiscoveryJSON is a minimal iShares discovery response for cache tests.
const sampleDiscoveryJSON = `{
	"239619": {
		"fundName": "iShares MSCI China ETF",
		"localExchangeTicker": "MCHI",
		"isin": "US4642874659",
		"productType": "ISHARES_FUND_DATA",
		"inceptionDate": {"r": 20181329},
		"fees": {"r": 0.59},
		"totalNetAssets": {"r": 7779083697.85},
		"portfolioId": 239619,
		"productPageUrl": ":/us/products/239619/test"
	}
}`

func TestDiscoverCached_ReusesCache(t *testing.T) {
	mock := &testutil.CountingMockHTTPClient{
		MockHTTPClient: testutil.MockHTTPClient{
			ResponseBody: sampleDiscoveryJSON,
			StatusCode:   http.StatusOK,
		},
	}

	c, err := New("us", WithHTTPClient(mock))
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	ctx := context.Background()

	// First call should hit HTTP.
	funds1, err := c.DiscoverETFs(ctx)
	if err != nil {
		t.Fatalf("first DiscoverETFs() failed: %v", err)
	}
	if mock.CallCount != 1 {
		t.Fatalf("expected 1 HTTP call after first discover, got %d", mock.CallCount)
	}

	// Second call should use cache.
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

	c, err := New("us", WithHTTPClient(mock))
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	ctx := context.Background()

	_, err = c.DiscoverETFs(ctx)
	if err != nil {
		t.Fatalf("DiscoverETFs() failed: %v", err)
	}

	// FundInfo should not trigger another HTTP call.
	fund, err := c.FundInfo(ctx, "MCHI")
	if err != nil {
		t.Fatalf("FundInfo() failed: %v", err)
	}
	if fund.Ticker != "MCHI" {
		t.Fatalf("expected ticker MCHI, got %s", fund.Ticker)
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

	c, err := New("us", WithHTTPClient(mock), WithCacheTTL(1*time.Millisecond))
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

	// Wait for cache to expire.
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

	c, err := New("us", WithHTTPClient(mock), WithCacheTTL(0))
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	ctx := context.Background()

	for i := range 3 {
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

	c, err := New("us", WithHTTPClient(mock))
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	ctx := context.Background()

	funds1, err := c.DiscoverETFs(ctx)
	if err != nil {
		t.Fatalf("DiscoverETFs() failed: %v", err)
	}

	// Mutate the returned slice.
	if len(funds1) > 0 {
		funds1[0].Ticker = "MUTATED"
	}

	// Second call should still return the original data.
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
	c, err := New("us", WithCacheTTL(ttl))
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}
	if c.cacheTTL != ttl {
		t.Fatalf("expected cacheTTL %v, got %v", ttl, c.cacheTTL)
	}
}

func TestBuildFundIndex(t *testing.T) {
	t.Run("lookup by ticker and ISIN", func(t *testing.T) {
		mock := &testutil.MockHTTPClient{
			ResponseBody: sampleDiscoveryJSON,
			StatusCode:   http.StatusOK,
		}

		c, err := New("us", WithHTTPClient(mock))
		if err != nil {
			t.Fatalf("New() failed: %v", err)
		}

		ctx := context.Background()
		_, err = c.DiscoverETFs(ctx)
		if err != nil {
			t.Fatalf("DiscoverETFs() failed: %v", err)
		}

		// Lookup by ticker (case insensitive).
		fund, err := c.FundInfo(ctx, "mchi")
		if err != nil {
			t.Fatalf("FundInfo by lowercase ticker failed: %v", err)
		}
		if fund.Ticker != "MCHI" {
			t.Fatalf("expected MCHI, got %s", fund.Ticker)
		}

		// Lookup by ISIN.
		fund, err = c.FundInfo(ctx, "US4642874659")
		if err != nil {
			t.Fatalf("FundInfo by ISIN failed: %v", err)
		}
		if fund.ISIN != "US4642874659" {
			t.Fatalf("expected US4642874659, got %s", fund.ISIN)
		}
	})
}

func TestHoldingsForFund_NilFund(t *testing.T) {
	c, err := New("us")
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
