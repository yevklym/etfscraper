package amundi

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

//go:embed data/getProductsData-de.json
var discoveryResponseDE string

func TestDiscoverETFs(t *testing.T) {
	client, err := New("de", WithHTTPClient(&testutil.MockHTTPClient{
		StatusCode:   http.StatusOK,
		ResponseBody: discoveryResponseDE,
	}))
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	funds, err := client.DiscoverETFs(context.Background())
	if err != nil {
		t.Fatalf("DiscoverETFs() failed: %v", err)
	}
	if len(funds) == 0 {
		t.Fatal("expected at least one fund")
	}

	found := findFundByISIN(funds, "DE000ETF7011")
	if found == nil {
		t.Fatal("expected fund DE000ETF7011 to be present")
	}
	if found.Ticker != "F701" {
		t.Errorf("expected ticker F701, got %q", found.Ticker)
	}
	if found.Currency != etfscraper.CurrencyEUR {
		t.Errorf("expected currency EUR, got %q", found.Currency)
	}
	if !found.IsDistributing {
		t.Error("expected fund to be distributing")
	}

	expectedExpense := 0.0042
	if math.Abs(found.ExpenseRatio-expectedExpense) > 1e-6 {
		t.Errorf("expected expense ratio %.4f, got %.6f", expectedExpense, found.ExpenseRatio)
	}
	if found.AssetClass != etfscraper.AssetClassAlternative {
		t.Errorf("expected asset class %q, got %q", etfscraper.AssetClassAlternative, found.AssetClass)
	}

	liquidated := findFundByISIN(funds, "LU2469335025")
	if liquidated != nil {
		t.Errorf("expected liquidated fund LU2469335025 to be filtered")
	}
}

func TestDiscoverETFsHTTPError(t *testing.T) {
	client, err := New("de", WithHTTPClient(&testutil.MockHTTPClient{
		StatusCode:   http.StatusInternalServerError,
		ResponseBody: `{"error":"boom"}`,
	}))
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

func findFundByISIN(funds []etfscraper.Fund, isin string) *etfscraper.Fund {
	for i := range funds {
		if funds[i].ISIN == isin {
			return &funds[i]
		}
	}
	return nil
}
