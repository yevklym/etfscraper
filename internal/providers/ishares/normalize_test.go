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

func TestNormalizeSector_AllRegions(t *testing.T) {
	tests := []struct {
		name     string
		region   string
		input    string
		expected etfscraper.Sector
	}{
		// English (US)
		{name: "energy us", region: "us", input: "Energy", expected: etfscraper.SectorEnergy},
		{name: "materials us", region: "us", input: "Materials", expected: etfscraper.SectorMaterials},
		{name: "industrials us", region: "us", input: "Industrials", expected: etfscraper.SectorIndustrials},
		{name: "consumer discretionary us", region: "us", input: "Consumer Discretionary", expected: etfscraper.SectorConsumerDiscretionary},
		{name: "consumer staples us", region: "us", input: "Consumer Staples", expected: etfscraper.SectorConsumerStaples},
		{name: "health care us", region: "us", input: "Health Care", expected: etfscraper.SectorHealthcare},
		{name: "financials us", region: "us", input: "Financials", expected: etfscraper.SectorFinancials},
		{name: "information technology us", region: "us", input: "Information Technology", expected: etfscraper.SectorInformationTechnology},
		{name: "communication us", region: "us", input: "Communication", expected: etfscraper.SectorTelecommunication},
		{name: "utilities us", region: "us", input: "Utilities", expected: etfscraper.SectorUtilities},
		{name: "real estate us", region: "us", input: "Real Estate", expected: etfscraper.SectorRealEstate},

		// English (UK)
		{name: "energy uk", region: "uk", input: "Energy", expected: etfscraper.SectorEnergy},
		{name: "health care uk", region: "uk", input: "Health Care", expected: etfscraper.SectorHealthcare},
		{name: "communication uk", region: "uk", input: "Communication", expected: etfscraper.SectorTelecommunication},

		// German (DE)
		{name: "energy german", region: "de", input: "Energie", expected: etfscraper.SectorEnergy},
		{name: "materials german", region: "de", input: "Materialien", expected: etfscraper.SectorMaterials},
		{name: "industrials german", region: "de", input: "Industrie", expected: etfscraper.SectorIndustrials},
		{name: "consumer discretionary german", region: "de", input: "Zyklische Konsumgüter", expected: etfscraper.SectorConsumerDiscretionary},
		{name: "consumer staples german", region: "de", input: "Nichtzyklische Konsumgüter", expected: etfscraper.SectorConsumerStaples},
		{name: "healthcare german", region: "de", input: "Gesundheitsversorgung", expected: etfscraper.SectorHealthcare},
		{name: "financials german", region: "de", input: "Financials", expected: etfscraper.SectorFinancials},
		{name: "it german", region: "de", input: "IT", expected: etfscraper.SectorInformationTechnology},
		{name: "communication german", region: "de", input: "Kommunikation", expected: etfscraper.SectorTelecommunication},
		{name: "utilities german", region: "de", input: "Versorger", expected: etfscraper.SectorUtilities},
		{name: "real estate german", region: "de", input: "Immobilien", expected: etfscraper.SectorRealEstate},

		// French (FR)
		{name: "energy french", region: "fr", input: "Energie", expected: etfscraper.SectorEnergy},
		{name: "materials french", region: "fr", input: "Matériaux", expected: etfscraper.SectorMaterials},
		{name: "industrials french", region: "fr", input: "Industries", expected: etfscraper.SectorIndustrials},
		{name: "consumer discretionary french", region: "fr", input: "Biens de consommation cycliques", expected: etfscraper.SectorConsumerDiscretionary},
		{name: "consumer staples french", region: "fr", input: "Biens de consommation de base", expected: etfscraper.SectorConsumerStaples},
		{name: "healthcare french", region: "fr", input: "Santé", expected: etfscraper.SectorHealthcare},
		{name: "financials french", region: "fr", input: "Finance", expected: etfscraper.SectorFinancials},
		{name: "information technology french", region: "fr", input: "Technologie de l'information", expected: etfscraper.SectorInformationTechnology},
		{name: "communication french", region: "fr", input: "La communication", expected: etfscraper.SectorTelecommunication},
		{name: "utilities french", region: "fr", input: "Services publics", expected: etfscraper.SectorUtilities},
		{name: "real estate french", region: "fr", input: "Immobilier", expected: etfscraper.SectorRealEstate},

		// Edge cases
		{name: "empty string", region: "us", input: "", expected: ""},
		{name: "dash", region: "us", input: "-", expected: ""},
		{name: "whitespace", region: "us", input: "  Energy  ", expected: etfscraper.SectorEnergy},
		{name: "case insensitive", region: "us", input: "HEALTH CARE", expected: etfscraper.SectorHealthcare},
		{name: "unknown passthrough", region: "us", input: "Cash and/or Derivatives", expected: etfscraper.Sector("Cash and/or Derivatives")},
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

func TestNormalizeLocation_AllRegions(t *testing.T) {
	tests := []struct {
		name     string
		region   string
		input    string
		expected etfscraper.Location
	}{
		// English (US, UK)
		{name: "us english", region: "us", input: "United States", expected: etfscraper.LocationUnitedStates},
		{name: "uk english", region: "uk", input: "United Kingdom", expected: etfscraper.LocationUnitedKingdom},
		{name: "japan english", region: "us", input: "Japan", expected: etfscraper.LocationJapan},

		// German (DE)
		{name: "us german", region: "de", input: "Vereinigte Staaten", expected: etfscraper.LocationUnitedStates},
		{name: "germany german", region: "de", input: "Deutschland", expected: etfscraper.LocationGermany},
		{name: "japan german", region: "de", input: "Japan", expected: etfscraper.LocationJapan},

		// French (FR)
		{name: "us french", region: "fr", input: "États-unis", expected: etfscraper.LocationUnitedStates},
		{name: "france french", region: "fr", input: "France", expected: etfscraper.LocationFrance},
		{name: "japan french", region: "fr", input: "Japon", expected: etfscraper.LocationJapan},

		// Edge cases
		{name: "empty string", region: "us", input: "", expected: ""},
		{name: "dash", region: "us", input: "-", expected: ""},
		{name: "whitespace", region: "us", input: "  Japan  ", expected: etfscraper.LocationJapan},
		{name: "case insensitive", region: "us", input: "UNITED STATES", expected: etfscraper.LocationUnitedStates},
		{name: "unknown passthrough", region: "us", input: "Vietnam", expected: etfscraper.Location("Vietnam")},
		{name: "nil mapping", region: "", input: "Japan", expected: etfscraper.Location("Japan")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var mapping map[string]etfscraper.Location
			if cfg, ok := regionConfigs[tt.region]; ok {
				mapping = cfg.LocationMapping
			}
			got := normalizeLocation(tt.input, mapping)
			if got != tt.expected {
				t.Errorf("normalizeLocation(%q) = %q, want %q", tt.input, got, tt.expected)
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
