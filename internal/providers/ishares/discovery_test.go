package ishares

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/yevklym/etfscraper"
)

type mockHTTPClient struct {
	ResponseBody string
	StatusCode   int
	Error        error
	Delay        time.Duration
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	select {
	case <-req.Context().Done():
		return nil, req.Context().Err()
	default:
	}

	if m.Delay > 0 {
		select {
		case <-time.After(m.Delay):
		case <-req.Context().Done():
			return nil, req.Context().Err()
		}
	}

	if m.Error != nil {
		return nil, m.Error
	}

	response := &http.Response{
		StatusCode: m.StatusCode,
		Body:       io.NopCloser(bytes.NewReader([]byte(m.ResponseBody))),
		Header:     make(http.Header),
	}
	response.Header.Set("Content-Type", "application/json")

	return response, nil
}

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
	c := &Client{}

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

func TestContextCancellation(t *testing.T) {
	t.Run("immediate cancellation", func(t *testing.T) {
		slowMock := &mockHTTPClient{
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
		slowMock := &mockHTTPClient{
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
