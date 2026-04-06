package xtrackers

import (
	"context"
	_ "embed"
	"math"
	"net/http"
	"testing"

	"github.com/yevklym/etfscraper"
	"github.com/yevklym/etfscraper/internal/testutil"
)

//go:embed data/holdings-de-de.json
var holdingsResponseDE string

//go:embed data/holdings-en-gb.json
var holdingsResponseGB string

func TestHoldingsForFund(t *testing.T) {
	client, err := New("de",
		WithHTTPClient(&testutil.MockHTTPClient{
			StatusCode:   http.StatusOK,
			ResponseBody: holdingsResponseDE,
		}),
		withSkipBrowserFetch(),
	)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	fund := &etfscraper.Fund{
		ISIN:     "IE00BJ0KDQ92",
		Name:     "Xtrackers MSCI World UCITS ETF 1C",
		Provider: etfscraper.ProviderXtrackers,
		ProviderMetadata: fundMetadata{
			ProductURL: "/de-de/IE00BJ0KDQ92-msci-world-ucits-etf-1c/",
		},
	}

	snapshot, err := client.HoldingsForFund(context.Background(), fund)
	if err != nil {
		t.Fatalf("HoldingsForFund() failed: %v", err)
	}

	if snapshot.TotalHoldings != 7 {
		t.Fatalf("expected 7 holdings, got %d", snapshot.TotalHoldings)
	}

	tests := []struct {
		name            string
		isin            string
		wantName        string
		wantWeight      float64
		wantMarketValue float64
		wantLocation    string
		wantSector      string
		wantAssetClass  etfscraper.AssetClass
	}{
		{
			name:            "nvidia",
			isin:            "US67066G1040",
			wantName:        "NVIDIA CORP",
			wantWeight:      0.0523670013,
			wantMarketValue: 1261876635.72,
			wantLocation:    "United States",
			wantSector:      "Information Technology",
			wantAssetClass:  etfscraper.AssetClassEquity,
		},
		{
			name:            "jpmorgan finance sector",
			isin:            "US46625H1005",
			wantName:        "JPMORGAN CHASE",
			wantWeight:      0.0099871635,
			wantMarketValue: 240658582.44,
			wantLocation:    "United States",
			wantSector:      "Financials",
			wantAssetClass:  etfscraper.AssetClassEquity,
		},
		{
			name:            "toyota japan",
			isin:            "JP3633400001",
			wantName:        "TOYOTA MOTOR CORP",
			wantWeight:      0.003512,
			wantMarketValue: 84610000.00,
			wantLocation:    "Japan",
			wantSector:      "Consumer Discretionary",
			wantAssetClass:  etfscraper.AssetClassEquity,
		},
		{
			name:            "cash position empty sector",
			isin:            "_CURRENCYUSD",
			wantName:        "US DOLLAR",
			wantWeight:      0.00012,
			wantMarketValue: 2890000.00,
			wantLocation:    "",
			wantSector:      "",
			wantAssetClass:  etfscraper.AssetClass("Cash und/oder Derivate"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found := findHoldingByISIN(snapshot.Holdings, tt.isin)
			if found == nil {
				t.Fatalf("holding %s not found", tt.isin)
			}
			if found.Name != tt.wantName {
				t.Errorf("name = %q, want %q", found.Name, tt.wantName)
			}
			if math.Abs(found.Weight-tt.wantWeight) > 1e-8 {
				t.Errorf("weight = %v, want %v", found.Weight, tt.wantWeight)
			}
			if math.Abs(found.MarketValue-tt.wantMarketValue) > 0.01 {
				t.Errorf("market value = %v, want %v", found.MarketValue, tt.wantMarketValue)
			}
			if string(found.Location) != tt.wantLocation {
				t.Errorf("location = %q, want %q", found.Location, tt.wantLocation)
			}
			if string(found.Sector) != tt.wantSector {
				t.Errorf("sector = %q, want %q", found.Sector, tt.wantSector)
			}
			if found.AssetClass != tt.wantAssetClass {
				t.Errorf("asset class = %q, want %q", found.AssetClass, tt.wantAssetClass)
			}
		})
	}
}

func TestHoldingsForFund_NilFund(t *testing.T) {
	client, err := New("de")
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	_, err = client.HoldingsForFund(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for nil fund")
	}
}

func TestHoldingsForFund_MissingMetadata(t *testing.T) {
	client, err := New("de",
		WithHTTPClient(&testutil.MockHTTPClient{
			StatusCode:   http.StatusOK,
			ResponseBody: holdingsResponseDE,
		}),
		withSkipBrowserFetch(),
	)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	fund := &etfscraper.Fund{
		ISIN:     "IE00BJ0KDQ92",
		Provider: etfscraper.ProviderXtrackers,
		// No ProviderMetadata — should fail
	}

	_, err = client.HoldingsForFund(context.Background(), fund)
	if err == nil {
		t.Fatal("expected error for fund without metadata")
	}
}

func TestHoldingsForFund_UK(t *testing.T) {
	client, err := New("uk",
		WithHTTPClient(&testutil.MockHTTPClient{
			StatusCode:   http.StatusOK,
			ResponseBody: holdingsResponseGB,
		}),
		withSkipBrowserFetch(),
	)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	fund := &etfscraper.Fund{
		ISIN:     "IE00BK1PV551",
		Name:     "Xtrackers MSCI World UCITS ETF 1D",
		Provider: etfscraper.ProviderXtrackers,
		ProviderMetadata: fundMetadata{
			ProductURL: "/en-gb/IE00BK1PV551-msci-world-ucits-etf-1d/",
		},
	}

	snapshot, err := client.HoldingsForFund(context.Background(), fund)
	if err != nil {
		t.Fatalf("HoldingsForFund() failed: %v", err)
	}

	if snapshot.TotalHoldings != 5 {
		t.Fatalf("expected 5 holdings, got %d", snapshot.TotalHoldings)
	}

	tests := []struct {
		name           string
		isin           string
		wantName       string
		wantLocation   string
		wantSector     string
		wantAssetClass etfscraper.AssetClass
	}{
		{
			name:           "nvidia english locale",
			isin:           "US67066G1040",
			wantName:       "NVIDIA CORP",
			wantLocation:   "United States",
			wantSector:     "Information Technology",
			wantAssetClass: etfscraper.AssetClassEquity,
		},
		{
			name:           "amazon consumer discretionary",
			isin:           "US0231351067",
			wantName:       "AMAZON COM INC",
			wantLocation:   "United States",
			wantSector:     "Consumer Discretionary",
			wantAssetClass: etfscraper.AssetClassEquity,
		},
		{
			name:           "toyota japan english",
			isin:           "JP3633400001",
			wantName:       "TOYOTA MOTOR CORP",
			wantLocation:   "Japan",
			wantSector:     "Consumer Discretionary",
			wantAssetClass: etfscraper.AssetClassEquity,
		},
		{
			name:           "cash english",
			isin:           "_CURRENCYUSD",
			wantName:       "US DOLLAR",
			wantLocation:   "",
			wantSector:     "",
			wantAssetClass: etfscraper.AssetClass("Cash and/or Derivatives"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found := findHoldingByISIN(snapshot.Holdings, tt.isin)
			if found == nil {
				t.Fatalf("holding %s not found", tt.isin)
			}
			if found.Name != tt.wantName {
				t.Errorf("name = %q, want %q", found.Name, tt.wantName)
			}
			if string(found.Location) != tt.wantLocation {
				t.Errorf("location = %q, want %q", found.Location, tt.wantLocation)
			}
			if string(found.Sector) != tt.wantSector {
				t.Errorf("sector = %q, want %q", found.Sector, tt.wantSector)
			}
			if found.AssetClass != tt.wantAssetClass {
				t.Errorf("asset class = %q, want %q", found.AssetClass, tt.wantAssetClass)
			}
		})
	}
}

func TestHoldingsURL(t *testing.T) {
	tests := []struct {
		name       string
		region     string
		productURL string
		want       string
	}{
		{
			name:       "DE normal URL",
			region:     "de",
			productURL: "/de-de/IE00BJ0KDQ92-msci-world-ucits-etf-1c/",
			want:       "https://etf.dws.com/api/pdp/de-de/etf/IE00BJ0KDQ92-msci-world-ucits-etf-1c/holdings",
		},
		{
			name:       "UK normal URL",
			region:     "uk",
			productURL: "/en-gb/IE00BK1PV551-msci-world-ucits-etf-1d/",
			want:       "https://etf.dws.com/api/pdp/en-gb/etf/IE00BK1PV551-msci-world-ucits-etf-1d/holdings",
		},
		{
			name:       "no trailing slash",
			region:     "de",
			productURL: "/de-de/IE00BJ0KDQ92-msci-world-ucits-etf-1c",
			want:       "https://etf.dws.com/api/pdp/de-de/etf/IE00BJ0KDQ92-msci-world-ucits-etf-1c/holdings",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := New(tt.region)
			if err != nil {
				t.Fatalf("New() failed: %v", err)
			}
			got := client.holdingsURL(tt.productURL)
			if got != tt.want {
				t.Errorf("holdingsURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNormalizeWeight(t *testing.T) {
	tests := []struct {
		name  string
		input float64
		want  float64
	}{
		{"percentage", 5.24, 0.0524},
		{"sub one percentage", 0.95, 0.0095},
		{"zero", 0.0, 0.0},
		{"one hundred percent", 100.0, 1.0},
		{"small fraction", 0.012, 0.00012},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeWeight(tt.input)
			if math.Abs(got-tt.want) > 1e-10 {
				t.Errorf("normalizeWeight(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func findHoldingByISIN(holdings []etfscraper.Holding, isin string) *etfscraper.Holding {
	for i := range holdings {
		if holdings[i].ISIN == isin {
			return &holdings[i]
		}
	}
	return nil
}
