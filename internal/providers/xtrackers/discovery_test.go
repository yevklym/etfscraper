package xtrackers

import (
	"context"
	_ "embed"
	"math"
	"net/http"
	"strings"
	"testing"

	"github.com/yevklym/etfscraper"
	"github.com/yevklym/etfscraper/internal/testutil"
)

//go:embed data/datatable-en-gb.json
var discoveryResponseGB string

//go:embed data/datatable-fr-fr.json
var discoveryResponseFR string

func TestDiscoverETFs(t *testing.T) {
	client, err := New("uk",
		WithHTTPClient(&testutil.MockHTTPClient{
			StatusCode:   http.StatusOK,
			ResponseBody: discoveryResponseGB,
		}),
		withSkipBrowserFetch(),
	)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	funds, err := client.DiscoverETFs(context.Background())
	if err != nil {
		t.Fatalf("DiscoverETFs() failed: %v", err)
	}
	if len(funds) != 5 {
		t.Fatalf("expected 5 funds, got %d", len(funds))
	}

	tests := []struct {
		name           string
		isin           string
		wantName       string
		wantCurrency   etfscraper.Currency
		wantAssetClass etfscraper.AssetClass
		wantTER        float64
		wantDist       bool
	}{
		{
			name:           "equity distributing USD",
			isin:           "IE00BK1PV551",
			wantName:       "Xtrackers MSCI World UCITS ETF 1D",
			wantCurrency:   etfscraper.CurrencyUSD,
			wantAssetClass: etfscraper.AssetClassEquity,
			wantTER:        0.0012,
			wantDist:       true,
		},
		{
			name:           "equity capitalizing USD",
			isin:           "IE00BJ0KDQ92",
			wantName:       "Xtrackers MSCI World UCITS ETF 1C",
			wantCurrency:   etfscraper.CurrencyUSD,
			wantAssetClass: etfscraper.AssetClassEquity,
			wantTER:        0.0019,
			wantDist:       false,
		},
		{
			name:           "equity capitalizing EUR",
			isin:           "LU0274208692",
			wantName:       "Xtrackers MSCI Europe UCITS ETF 1C",
			wantCurrency:   etfscraper.CurrencyEUR,
			wantAssetClass: etfscraper.AssetClassEquity,
			wantTER:        0.0012,
			wantDist:       false,
		},
		{
			name:           "fixed income bond",
			isin:           "LU0290357846",
			wantName:       "Xtrackers II EUR Corporate Bond UCITS ETF 1C",
			wantCurrency:   etfscraper.CurrencyEUR,
			wantAssetClass: etfscraper.AssetClassBond,
			wantTER:        0.0016,
			wantDist:       false,
		},
		{
			name:           "commodity GBP",
			isin:           "GB00BLD4ZL17",
			wantName:       "Xtrackers Physical Gold ETC (GBP)",
			wantCurrency:   etfscraper.CurrencyGBP,
			wantAssetClass: etfscraper.AssetClassCommodity,
			wantTER:        0.0025,
			wantDist:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found := findFundByISIN(funds, tt.isin)
			if found == nil {
				t.Fatalf("fund %s not found", tt.isin)
			}
			if found.Name != tt.wantName {
				t.Errorf("name = %q, want %q", found.Name, tt.wantName)
			}
			if found.Provider != etfscraper.ProviderXtrackers {
				t.Errorf("provider = %q, want %q", found.Provider, etfscraper.ProviderXtrackers)
			}
			if found.Currency != tt.wantCurrency {
				t.Errorf("currency = %q, want %q", found.Currency, tt.wantCurrency)
			}
			if found.AssetClass != tt.wantAssetClass {
				t.Errorf("asset class = %q, want %q", found.AssetClass, tt.wantAssetClass)
			}
			if math.Abs(found.ExpenseRatio-tt.wantTER) > 1e-6 {
				t.Errorf("expense ratio = %f, want %f", found.ExpenseRatio, tt.wantTER)
			}
			if found.IsDistributing != tt.wantDist {
				t.Errorf("distributing = %v, want %v", found.IsDistributing, tt.wantDist)
			}
			if found.TotalAssets <= 0 {
				t.Error("expected positive total assets")
			}
			if found.InceptionDate == nil {
				t.Error("expected non-nil inception date")
			}
		})
	}
}

func TestDiscoverETFs_FR(t *testing.T) {
	client, err := New("fr",
		WithHTTPClient(&testutil.MockHTTPClient{
			StatusCode:   http.StatusOK,
			ResponseBody: discoveryResponseFR,
		}),
		withSkipBrowserFetch(),
	)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	funds, err := client.DiscoverETFs(context.Background())
	if err != nil {
		t.Fatalf("DiscoverETFs() failed: %v", err)
	}
	if len(funds) != 5 {
		t.Fatalf("expected 5 funds, got %d", len(funds))
	}

	tests := []struct {
		name           string
		isin           string
		wantName       string
		wantCurrency   etfscraper.Currency
		wantAssetClass etfscraper.AssetClass
		wantTER        float64
		wantDist       bool
	}{
		{
			name:           "french equity distributing",
			isin:           "IE00BK1PV551",
			wantName:       "Xtrackers MSCI World UCITS ETF 1D",
			wantCurrency:   etfscraper.CurrencyUSD,
			wantAssetClass: etfscraper.AssetClassEquity,
			wantTER:        0.0012,
			wantDist:       true,
		},
		{
			name:           "french equity capitalisation",
			isin:           "IE00BJ0KDQ92",
			wantName:       "Xtrackers MSCI World UCITS ETF 1C",
			wantCurrency:   etfscraper.CurrencyUSD,
			wantAssetClass: etfscraper.AssetClassEquity,
			wantTER:        0.0019,
			wantDist:       false,
		},
		{
			name:           "french obligations bond",
			isin:           "LU0290357846",
			wantName:       "Xtrackers II EUR Corporate Bond UCITS ETF 1C",
			wantCurrency:   etfscraper.CurrencyEUR,
			wantAssetClass: etfscraper.AssetClassBond,
			wantTER:        0.0016,
			wantDist:       false,
		},
		{
			name:           "french matières premières commodity",
			isin:           "DE000A0S9GB0",
			wantName:       "Xtrackers Physical Gold ETC (EUR)",
			wantCurrency:   etfscraper.CurrencyEUR,
			wantAssetClass: etfscraper.AssetClassCommodity,
			wantTER:        0.0036,
			wantDist:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found := findFundByISIN(funds, tt.isin)
			if found == nil {
				t.Fatalf("fund %s not found", tt.isin)
			}
			if found.Name != tt.wantName {
				t.Errorf("name = %q, want %q", found.Name, tt.wantName)
			}
			if found.Provider != etfscraper.ProviderXtrackers {
				t.Errorf("provider = %q, want %q", found.Provider, etfscraper.ProviderXtrackers)
			}
			if found.Currency != tt.wantCurrency {
				t.Errorf("currency = %q, want %q", found.Currency, tt.wantCurrency)
			}
			if found.AssetClass != tt.wantAssetClass {
				t.Errorf("asset class = %q, want %q", found.AssetClass, tt.wantAssetClass)
			}
			if math.Abs(found.ExpenseRatio-tt.wantTER) > 1e-6 {
				t.Errorf("expense ratio = %f, want %f", found.ExpenseRatio, tt.wantTER)
			}
			if found.IsDistributing != tt.wantDist {
				t.Errorf("distributing = %v, want %v", found.IsDistributing, tt.wantDist)
			}
			if found.TotalAssets <= 0 {
				t.Error("expected positive total assets")
			}
			if found.InceptionDate == nil {
				t.Error("expected non-nil inception date")
			}
		})
	}
}

func TestDiscoverETFsHTTPError(t *testing.T) {
	client, err := New("uk",
		WithHTTPClient(&testutil.MockHTTPClient{
			StatusCode:   http.StatusInternalServerError,
			ResponseBody: `{"error":"boom"}`,
		}),
		withSkipBrowserFetch(),
	)
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	_, err = client.DiscoverETFs(context.Background())
	if err == nil {
		t.Fatal("expected error for non-200 response")
	}
	if !strings.Contains(err.Error(), "HTTP") {
		t.Errorf("expected HTTP error, got %q", err.Error())
	}
}

func TestNewUnsupportedRegion(t *testing.T) {
	_, err := New("xx")
	if err == nil {
		t.Fatal("expected error for unsupported region")
	}
	if !strings.Contains(err.Error(), "unsupported") {
		t.Errorf("expected unsupported region error, got %q", err.Error())
	}
}

func findFundByISIN(funds []etfscraper.Fund, isin string) *etfscraper.Fund {
	for i := range funds {
		if funds[i].ISIN == isin {
			return &funds[i]
		}
	}
	return nil
}
