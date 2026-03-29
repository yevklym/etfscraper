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
		wantAssetClass  etfscraper.AssetClass
	}{
		{
			name:            "nvidia",
			isin:            "US67066G1040",
			wantName:        "NVIDIA CORP",
			wantWeight:      0.0523670013,
			wantMarketValue: 1261876635.72,
			wantLocation:    "Vereinigte Staaten von Amerika",
			wantAssetClass:  etfscraper.AssetClassEquity,
		},
		{
			name:            "apple",
			isin:            "US0378331005",
			wantName:        "APPLE INC",
			wantWeight:      0.0467771656,
			wantMarketValue: 1127179537.77,
			wantLocation:    "Vereinigte Staaten von Amerika",
			wantAssetClass:  etfscraper.AssetClassEquity,
		},
		{
			name:            "jpmorgan finance sector",
			isin:            "US46625H1005",
			wantName:        "JPMORGAN CHASE",
			wantWeight:      0.0099871635,
			wantMarketValue: 240658582.44,
			wantLocation:    "Vereinigte Staaten von Amerika",
			wantAssetClass:  etfscraper.AssetClassEquity,
		},
		{
			name:            "toyota japan",
			isin:            "JP3633400001",
			wantName:        "TOYOTA MOTOR CORP",
			wantWeight:      0.003512,
			wantMarketValue: 84610000.00,
			wantLocation:    "Japan",
			wantAssetClass:  etfscraper.AssetClassEquity,
		},
		{
			name:            "cash position",
			isin:            "_CURRENCYUSD",
			wantName:        "US DOLLAR",
			wantWeight:      0.00012,
			wantMarketValue: 2890000.00,
			wantLocation:    "",
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

func TestHoldingsURL(t *testing.T) {
	client, err := New("de")
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	tests := []struct {
		name       string
		productURL string
		want       string
	}{
		{
			name:       "normal URL",
			productURL: "/de-de/IE00BJ0KDQ92-msci-world-ucits-etf-1c/",
			want:       "https://etf.dws.com/api/pdp/de-de/etf/IE00BJ0KDQ92-msci-world-ucits-etf-1c/holdings",
		},
		{
			name:       "no trailing slash",
			productURL: "/de-de/IE00BJ0KDQ92-msci-world-ucits-etf-1c",
			want:       "https://etf.dws.com/api/pdp/de-de/etf/IE00BJ0KDQ92-msci-world-ucits-etf-1c/holdings",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
