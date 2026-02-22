package ishares

import (
	"testing"

	"github.com/yevklym/etfscraper"
)

func TestNormalizeCurrency_FullNameMappings(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected etfscraper.Currency
	}{
		{name: "british pound", input: "BRITISH POUND", expected: etfscraper.CurrencyGBP},
		{name: "pound sterling", input: "Pound Sterling", expected: etfscraper.CurrencyGBP},
		{name: "us dollar", input: "U.S. DOLLAR", expected: etfscraper.CurrencyUSD},
		{name: "euro", input: "Euro", expected: etfscraper.CurrencyEUR},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := normalizeCurrency(test.input); got != test.expected {
				t.Fatalf("expected %q, got %q", test.expected, got)
			}
		})
	}
}

func TestNormalizeExchange_CommonMappings(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected etfscraper.Exchange
	}{
		{name: "lse code", input: "LSE", expected: etfscraper.ExchangeLSE},
		{name: "lse full", input: "London Stock Exchange", expected: etfscraper.ExchangeLSE},
		{name: "xetra", input: "Xetra", expected: etfscraper.ExchangeXetra},
		{name: "euronext", input: "Euronext Paris", expected: etfscraper.ExchangeEuronext},
		{name: "nyse", input: "New York Stock Exchange", expected: etfscraper.ExchangeNYSE},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := normalizeExchange(test.input); got != test.expected {
				t.Fatalf("expected %q, got %q", test.expected, got)
			}
		})
	}
}
