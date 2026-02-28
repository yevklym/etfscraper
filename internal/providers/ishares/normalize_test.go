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

func TestNormalizeAssetClass_AllRegions(t *testing.T) {
	tests := []struct {
		name     string
		region   string
		input    string
		expected etfscraper.AssetClass
	}{
		// English (US)
		{name: "equity us", region: "us", input: "Equity", expected: etfscraper.AssetClassEquity},
		{name: "fixed income us", region: "us", input: "Fixed Income", expected: etfscraper.AssetClassBond},
		{name: "cash us", region: "us", input: "Cash", expected: etfscraper.AssetClassCash},
		{name: "commodity us", region: "us", input: "Commodity", expected: etfscraper.AssetClassCommodity},
		{name: "real estate us", region: "us", input: "Real Estate", expected: etfscraper.AssetClassRealEstate},
		{name: "digital assets us", region: "us", input: "Digital Assets", expected: etfscraper.AssetClassCryptocurrency},
		{name: "multi asset us", region: "us", input: "Multi Asset", expected: etfscraper.AssetClassAlternative},

		// English (UK)
		{name: "equity uk", region: "uk", input: "Equity", expected: etfscraper.AssetClassEquity},
		{name: "fixed income uk", region: "uk", input: "Fixed Income", expected: etfscraper.AssetClassBond},

		// German (DE)
		{name: "equity german", region: "de", input: "Aktien", expected: etfscraper.AssetClassEquity},
		{name: "bonds german", region: "de", input: "Anleihen", expected: etfscraper.AssetClassBond},
		{name: "cash german", region: "de", input: "Barmittel", expected: etfscraper.AssetClassCash},
		{name: "commodity german", region: "de", input: "Rohstoffe", expected: etfscraper.AssetClassCommodity},
		{name: "real estate german", region: "de", input: "Immobilien", expected: etfscraper.AssetClassRealEstate},
		{name: "multi-asset german", region: "de", input: "Multi-Asset", expected: etfscraper.AssetClassAlternative},

		// French (FR)
		{name: "equity french", region: "fr", input: "Actions", expected: etfscraper.AssetClassEquity},
		{name: "bonds french", region: "fr", input: "Obligations", expected: etfscraper.AssetClassBond},
		{name: "cash french", region: "fr", input: "Liquidités", expected: etfscraper.AssetClassCash},
		{name: "commodity french", region: "fr", input: "Matières premières", expected: etfscraper.AssetClassCommodity},
		{name: "real estate french", region: "fr", input: "Immobilier", expected: etfscraper.AssetClassRealEstate},
		{name: "multi-actifs french", region: "fr", input: "Multi-actifs", expected: etfscraper.AssetClassAlternative},
		{name: "private markets french", region: "fr", input: "Marchés privés", expected: etfscraper.AssetClassAlternative},

		// Edge cases
		{name: "empty string", region: "us", input: "", expected: ""},
		{name: "dash", region: "us", input: "-", expected: ""},
		{name: "whitespace", region: "us", input: "  Equity  ", expected: etfscraper.AssetClassEquity},
		{name: "case insensitive", region: "us", input: "FIXED INCOME", expected: etfscraper.AssetClassBond},
		{name: "unknown passthrough", region: "us", input: "SomethingNew", expected: etfscraper.AssetClass("SomethingNew")},
		{name: "nil mapping", region: "", input: "Equity", expected: etfscraper.AssetClass("Equity")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var mapping map[string]etfscraper.AssetClass
			if cfg, ok := regionConfigs[tt.region]; ok {
				mapping = cfg.AssetClassMapping
			}
			got := normalizeAssetClass(tt.input, mapping)
			if got != tt.expected {
				t.Errorf("normalizeAssetClass(%q) = %q, want %q", tt.input, got, tt.expected)
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
