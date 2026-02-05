package ishares

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/yevklym/etfscraper"
	"github.com/yevklym/etfscraper/internal/testutil"
)

func TestParseISharesDate(t *testing.T) {
	t.Run("valid date", func(t *testing.T) {
		dateInt := 20231225
		expected := time.Date(2023, time.December, 25, 0, 0, 0, 0, time.UTC)

		result := parseISharesDate(dateInt)

		if result == nil {
			t.Fatalf("Expected date, but got nil")
		}
		if !result.Equal(expected) {
			t.Errorf("Expected date %v, but got %v", expected, *result)
		}
	})

	t.Run("invalid date format", func(t *testing.T) {
		dateInt := 202312250 // Extra digit

		result := parseISharesDate(dateInt)

		if result != nil {
			t.Errorf("Expected nil for invalid date, but got %v", *result)
		}
	})

	t.Run("zero date", func(t *testing.T) {
		dateInt := 0

		result := parseISharesDate(dateInt)

		if result != nil {
			t.Errorf("Expected nil for zero date, but got %v", *result)
		}
	})
}

func TestConvertSingleFund(t *testing.T) {
	c := &Client{config: regionConfigs["us"]}

	input := ISharesETFData{
		PortfolioID:         12345,
		FundName:            "Test Fund",
		LocalExchangeTicker: "TEST",
		ISIN:                "US1234567890",
		InceptionDate: struct {
			Display string `json:"d"`
			Raw     int    `json:"r"`
		}{Raw: 20200115},
		NetExpenseRatio: struct {
			Display string  `json:"d"`
			Raw     float64 `json:"r"`
		}{Raw: 7.5},
		TotalNetAssets: struct {
			Display string  `json:"d"`
			Raw     float64 `json:"r"`
		}{Raw: 1000000.0},
		ProductPageUrl: "/us/products/12345/test-fund",
	}

	inceptionDate := time.Date(2020, time.January, 15, 0, 0, 0, 0, time.UTC)
	expected := etfscraper.Fund{
		Ticker:        "TEST",
		Name:          "Test Fund",
		ISIN:          "US1234567890",
		Provider:      etfscraper.ProviderIShares,
		Currency:      etfscraper.CurrencyUSD,
		InceptionDate: &inceptionDate,
		ExpenseRatio:  0.075,
		TotalAssets:   1000000.0,
		Exchange:      etfscraper.ExchangeNYSE,
		ProviderMetadata: iSharesFundMetadata{
			PortfolioID:    12345,
			ProductPageURL: "/us/products/12345/test-fund",
		},
	}

	result := c.convertSingleFund(input)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Fund conversion mismatch:\nGot:  %+v\nWant: %+v", result, expected)
	}
}

func TestDiscoverETFs_WrapperFormat(t *testing.T) {
	sampleJSON := `{
		"i": {
			"239619": {
				"fundName": "iShares MSCI China ETF",
				"localExchangeTicker": "MCHI",
				"isin": "US4642874659",
				"productType": "ISHARES_FUND_DATA",
				"totalNetAssets": {"r": 7779083697.85},
				"netr": {"r": 0.59},
				"portfolioId": 239619,
				"productPageUrl": ":/us/products/239619/test"
			}
		}
	}`

	mockClient := &testutil.MockHTTPClient{
		ResponseBody: sampleJSON,
		StatusCode:   http.StatusOK,
	}

	c, _ := New("us", WithHTTPClient(mockClient))

	funds, err := c.DiscoverETFs(context.Background())
	if err != nil {
		t.Fatalf("DiscoverETFs failed: %v", err)
	}
	if len(funds) != 1 {
		t.Fatalf("expected 1 fund, got %d", len(funds))
	}
	if funds[0].Ticker != "MCHI" {
		t.Fatalf("expected ticker MCHI, got %q", funds[0].Ticker)
	}
}

func TestConvertSingleFund_RegionDefaults(t *testing.T) {
	input := ISharesETFData{
		PortfolioID:         12345,
		FundName:            "Test Fund",
		LocalExchangeTicker: "TEST",
		ISIN:                "DE1234567890",
		ProductType:         "ISHARES_FUND_DATA",
	}

	deClient := &Client{config: regionConfigs["de"]}
	deFund := deClient.convertSingleFund(input)
	if deFund.Currency != etfscraper.CurrencyEUR {
		t.Fatalf("expected EUR currency, got %q", deFund.Currency)
	}
	if deFund.Exchange != etfscraper.ExchangeXetra {
		t.Fatalf("expected Xetra exchange, got %q", deFund.Exchange)
	}

	ukClient := &Client{config: regionConfigs["uk"]}
	ukFund := ukClient.convertSingleFund(input)
	if ukFund.Currency != etfscraper.CurrencyGBP {
		t.Fatalf("expected GBP currency, got %q", ukFund.Currency)
	}
	if ukFund.Exchange != etfscraper.ExchangeLSE {
		t.Fatalf("expected LSE exchange, got %q", ukFund.Exchange)
	}
}

func TestContextCancellation(t *testing.T) {
	t.Run("immediate cancellation", func(t *testing.T) {
		slowMock := &testutil.MockHTTPClient{
			ResponseBody: `{}`,
			StatusCode:   http.StatusOK,
		}

		c, _ := New("us", WithHTTPClient(slowMock))

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		_, err := c.DiscoverETFs(ctx)
		if err == nil {
			t.Fatal("expected error from cancelled context")
		}

		if !errors.Is(err, context.Canceled) {
			t.Errorf("expected context.Canceled error, got: %v", err)
		}
	})

	t.Run("cancellation during request", func(t *testing.T) {
		slowMock := &testutil.MockHTTPClient{
			ResponseBody: `{}`,
			StatusCode:   http.StatusOK,
			Delay:        100 * time.Millisecond, // Simulate slow response
		}

		c, _ := New("us", WithHTTPClient(slowMock))

		ctx, cancel := context.WithCancel(context.Background())

		// Cancel after 50ms (before the 100ms delay completes)
		go func() {
			time.Sleep(50 * time.Millisecond)
			cancel()
		}()

		_, err := c.DiscoverETFs(ctx)
		if err == nil {
			t.Fatal("expected error from cancelled context")
		}

		if !errors.Is(err, context.Canceled) {
			t.Errorf("expected context.Canceled error, got: %v", err)
		}
	})
}
