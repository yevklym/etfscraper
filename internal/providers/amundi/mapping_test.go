package amundi

import (
	"testing"

	"github.com/yevklym/etfscraper"
)

func TestNormalizeAssetClass_AllRegions(t *testing.T) {
	tests := []struct {
		name     string
		region   string
		input    string
		expected etfscraper.AssetClass
	}{
		// Discovery ASSET_CLASS values
		{name: "equity", region: "de", input: "Equity", expected: etfscraper.AssetClassEquity},
		{name: "fixed income", region: "de", input: "Fixed Income", expected: etfscraper.AssetClassBond},
		{name: "commodities", region: "de", input: "Commodities", expected: etfscraper.AssetClassCommodity},
		{name: "multi asset", region: "de", input: "Multi Asset", expected: etfscraper.AssetClassAlternative},
		{name: "alternatives", region: "de", input: "Alternatives", expected: etfscraper.AssetClassAlternative},

		// Holdings type values
		{name: "equity_ordinary", region: "de", input: "EQUITY_ORDINARY", expected: etfscraper.AssetClassEquity},
		{name: "equity_preferred", region: "de", input: "EQUITY_PREFERRED", expected: etfscraper.AssetClassEquity},
		{name: "bond type", region: "de", input: "Bond", expected: etfscraper.AssetClassBond},
		{name: "cash type", region: "de", input: "Cash", expected: etfscraper.AssetClassCash},

		// Same mappings across regions
		{name: "equity uk", region: "uk", input: "Equity", expected: etfscraper.AssetClassEquity},
		{name: "equity fr", region: "fr", input: "Equity", expected: etfscraper.AssetClassEquity},
		{name: "multi asset uk", region: "uk", input: "Multi Asset", expected: etfscraper.AssetClassAlternative},
		{name: "multi asset fr", region: "fr", input: "Multi Asset", expected: etfscraper.AssetClassAlternative},

		// Edge cases
		{name: "empty string", region: "de", input: "", expected: ""},
		{name: "whitespace", region: "de", input: "  Equity  ", expected: etfscraper.AssetClassEquity},
		{name: "case insensitive", region: "de", input: "FIXED INCOME", expected: etfscraper.AssetClassBond},
		{name: "unknown passthrough", region: "de", input: "SomethingNew", expected: etfscraper.AssetClass("SomethingNew")},
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

func TestNormalizeSector_AllRegions(t *testing.T) {
	tests := []struct {
		name     string
		region   string
		input    string
		expected etfscraper.Sector
	}{
		// All 11 GICS sectors
		{name: "energy", region: "de", input: "Energy", expected: etfscraper.SectorEnergy},
		{name: "materials", region: "de", input: "Materials", expected: etfscraper.SectorMaterials},
		{name: "industrials", region: "de", input: "Industrials", expected: etfscraper.SectorIndustrials},
		{name: "consumer discretionary", region: "de", input: "Consumer Discretionary", expected: etfscraper.SectorConsumerDiscretionary},
		{name: "consumer staples", region: "de", input: "Consumer Staples", expected: etfscraper.SectorConsumerStaples},
		{name: "health care", region: "de", input: "Health Care", expected: etfscraper.SectorHealthcare},
		{name: "financials", region: "de", input: "Financials", expected: etfscraper.SectorFinancials},
		{name: "information technology", region: "de", input: "Information Technology", expected: etfscraper.SectorInformationTechnology},
		{name: "communication", region: "de", input: "Communication", expected: etfscraper.SectorTelecommunication},
		{name: "communication services", region: "de", input: "Communication Services", expected: etfscraper.SectorTelecommunication},
		{name: "utilities", region: "de", input: "Utilities", expected: etfscraper.SectorUtilities},
		{name: "real estate", region: "de", input: "Real Estate", expected: etfscraper.SectorRealEstate},

		// Same mappings across regions
		{name: "energy uk", region: "uk", input: "Energy", expected: etfscraper.SectorEnergy},
		{name: "energy fr", region: "fr", input: "Energy", expected: etfscraper.SectorEnergy},
		{name: "information technology uk", region: "uk", input: "Information Technology", expected: etfscraper.SectorInformationTechnology},

		// Edge cases
		{name: "empty string", region: "de", input: "", expected: ""},
		{name: "whitespace", region: "de", input: "  Energy  ", expected: etfscraper.SectorEnergy},
		{name: "case insensitive", region: "de", input: "HEALTH CARE", expected: etfscraper.SectorHealthcare},
		{name: "unknown passthrough", region: "de", input: "Other", expected: etfscraper.Sector("Other")},
		{name: "nil mapping", region: "", input: "Energy", expected: etfscraper.Sector("Energy")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var mapping map[string]etfscraper.Sector
			if cfg, ok := regionConfigs[tt.region]; ok {
				mapping = cfg.SectorMapping
			}
			got := normalizeSector(tt.input, mapping)
			if got != tt.expected {
				t.Errorf("normalizeSector(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestMapCurrency(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  etfscraper.Currency
	}{
		{name: "USD", input: "USD", want: etfscraper.CurrencyUSD},
		{name: "EUR", input: "EUR", want: etfscraper.CurrencyEUR},
		{name: "GBP", input: "GBP", want: etfscraper.CurrencyGBP},
		{name: "JPY", input: "JPY", want: etfscraper.CurrencyJPY},
		{name: "CAD", input: "CAD", want: etfscraper.CurrencyCAD},
		{name: "AUD", input: "AUD", want: etfscraper.CurrencyAUD},
		{name: "CHF", input: "CHF", want: etfscraper.CurrencyCHF},
		{name: "CNY", input: "CNY", want: etfscraper.CurrencyCNY},
		{name: "INR", input: "INR", want: etfscraper.CurrencyINR},
		{name: "BRL", input: "BRL", want: etfscraper.CurrencyBRL},
		{name: "lowercase", input: "usd", want: etfscraper.CurrencyUSD},
		{name: "mixed case", input: "Eur", want: etfscraper.CurrencyEUR},
		{name: "whitespace", input: "  GBP  ", want: etfscraper.CurrencyGBP},
		{name: "unknown passthrough", input: "SEK", want: etfscraper.Currency("SEK")},
		{name: "empty", input: "", want: ""},
		{name: "whitespace only", input: "   ", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapCurrency(tt.input)
			if got != tt.want {
				t.Errorf("mapCurrency(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestPickTicker(t *testing.T) {
	tests := []struct {
		name  string
		chars characteristics
		want  string
	}{
		{
			name:  "mnemo preferred",
			chars: characteristics{Mnemo: "GSPX", Ticker: "OTHER"},
			want:  "GSPX",
		},
		{
			name:  "ticker fallback",
			chars: characteristics{Ticker: "GSPX"},
			want:  "GSPX",
		},
		{
			name:  "listing fallback",
			chars: characteristics{MainListings: map[string]string{"DEU": "GSPX GY Equity"}},
			want:  "GSPX",
		},
		{
			name:  "empty",
			chars: characteristics{},
			want:  "",
		},
		{
			name:  "whitespace mnemo",
			chars: characteristics{Mnemo: "  ", Ticker: "TICK"},
			want:  "TICK",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pickTicker(tt.chars)
			if got != tt.want {
				t.Errorf("pickTicker() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPickISIN(t *testing.T) {
	tests := []struct {
		name string
		prod product
		want string
	}{
		{
			name: "from characteristics",
			prod: product{Characteristics: characteristics{ISIN: "IE00B5BMR087"}},
			want: "IE00B5BMR087",
		},
		{
			name: "from product ID",
			prod: product{ProductID: "IE00B5BMR087"},
			want: "IE00B5BMR087",
		},
		{
			name: "product ID wrong length",
			prod: product{ProductID: "SHORT"},
			want: "",
		},
		{
			name: "empty",
			prod: product{},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pickISIN(tt.prod)
			if got != tt.want {
				t.Errorf("pickISIN() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestIsDistributing(t *testing.T) {
	tests := []struct {
		name   string
		policy string
		want   bool
	}{
		{name: "distribution", policy: "Distribution", want: true},
		{name: "accumulation", policy: "Accumulation", want: false},
		{name: "empty", policy: "", want: false},
		{name: "contains distribution", policy: "Quarterly Distribution", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isDistributing(tt.policy)
			if got != tt.want {
				t.Errorf("isDistributing(%q) = %v, want %v", tt.policy, got, tt.want)
			}
		})
	}
}
