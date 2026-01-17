package ishares

import (
	"context"
	"net/http"
	"strings"
	"testing"
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

		mockClient := &mockHTTPClient{
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

		mockClient := &mockHTTPClient{
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
