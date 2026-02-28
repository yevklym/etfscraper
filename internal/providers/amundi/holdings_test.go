package amundi

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/yevklym/etfscraper"
	"github.com/yevklym/etfscraper/internal/testutil"
)

func TestDateValueUnmarshal_String(t *testing.T) {
	var dv dateValue
	if err := json.Unmarshal([]byte(`"2026-01-29"`), &dv); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if !dv.Valid {
		t.Fatal("expected valid date")
	}

	expected := time.Date(2026, time.January, 29, 0, 0, 0, 0, time.UTC)
	if !dv.Time.Equal(expected) {
		t.Fatalf("expected %v, got %v", expected, dv.Time)
	}
}

func TestDateValueUnmarshal_Millis(t *testing.T) {
	var dv dateValue
	expectedMillis := int64(1769625600000)
	if err := json.Unmarshal([]byte(`1769625600000`), &dv); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if !dv.Valid {
		t.Fatal("expected valid date")
	}
	if dv.Time.UnixMilli() != expectedMillis {
		t.Fatalf("expected %d, got %d", expectedMillis, dv.Time.UnixMilli())
	}
}

func TestParseHoldingsDate_PrioritizesPositionDate(t *testing.T) {
	position := dateValue{
		Time:  time.Date(2026, time.January, 29, 0, 0, 0, 0, time.UTC),
		Valid: true,
	}
	fallback := dateValue{
		Time:  time.Date(2025, time.December, 31, 0, 0, 0, 0, time.UTC),
		Valid: true,
	}

	got, err := parseHoldingsDate(holdingsCharacteristics{
		PositionAsOfDate:       position,
		FundBreakdownsAsOfDate: fallback,
	})
	if err != nil {
		t.Fatalf("parseHoldingsDate failed: %v", err)
	}
	if !got.Equal(position.Time) {
		t.Fatalf("expected %v, got %v", position.Time, got)
	}
}

func TestParseComposition_CompositionData(t *testing.T) {
	raw := []byte(`{
		"totalNumberOfInstruments": 2,
		"compositionData": [
			{
				"compositionCharacteristics": {
					"name": "NVIDIA CORP",
					"isin": "US67066G1040",
					"bbg": "NVDA UW",
					"weight": 0.1,
					"quantity": 123,
					"currency": "USD",
					"type": "EQUITY_ORDINARY",
					"sector": "Information Technology",
					"countryOfRisk": "United States"
				},
				"weight": 0.1
			},
			{
				"compositionCharacteristics": {
					"name": "AMAZON.COM INC",
					"isin": "US0231351067",
					"bbg": "AMZN UW",
					"weight": 0.2,
					"quantity": 456,
					"currency": "USD",
					"type": "EQUITY_ORDINARY",
					"sector": "Consumer Discretionary",
					"countryOfRisk": "United States"
				},
				"weight": 0.2
			}
		]
	}`)

	items, err := parseComposition(raw)
	if err != nil {
		t.Fatalf("parseComposition failed: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
	if items[0].Name != "NVIDIA CORP" {
		t.Fatalf("unexpected first item name: %s", items[0].Name)
	}
	if items[1].ISIN != "US0231351067" {
		t.Fatalf("unexpected second item ISIN: %s", items[1].ISIN)
	}
}

func TestConvertHoldings_DerivesMarketValue(t *testing.T) {
	items := []compositionItem{
		{Name: "Example", Weight: 0.1},
		{Name: "Percent", Weight: 5},
		{Name: "Provided", Weight: 0.2, MarketValue: 42},
	}

	c, err := New("de")
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	holdings := c.convertHoldings(items, 1000)
	if len(holdings) != 3 {
		t.Fatalf("expected 3 holdings, got %d", len(holdings))
	}

	if holdings[0].MarketValue != 100 {
		t.Fatalf("expected derived market value 100, got %f", holdings[0].MarketValue)
	}
	if holdings[1].Weight != 0.05 {
		t.Fatalf("expected normalized weight 0.05, got %f", holdings[1].Weight)
	}
	if holdings[2].MarketValue != 42 {
		t.Fatalf("expected provided market value 42, got %f", holdings[2].MarketValue)
	}
}

func TestHoldings_EndToEnd(t *testing.T) {
	discoveryResponse := `{
		"products": [
			{
				"productId": "LU1135865084",
				"productType": "PRODUCT",
				"characteristics": {
					"ISIN": "LU1135865084",
					"SHARE_MARKETING_NAME": "Amundi Core S&P 500 Swap UCITS ETF Acc",
					"MNEMO": "C500",
					"TER": 0.15,
					"CURRENCY": "EUR",
					"FUND_AUM": 1000,
					"ASSET_CLASS": "Equity",
					"DISTRIBUTION_POLICY": "Capitalisation"
				}
			}
		]
	}`

	holdingsResponse := `{
		"products": [
			{
				"productId": "LU1135865084",
				"productType": "PRODUCT",
				"characteristics": {
					"ISIN": "LU1135865084",
					"POSITION_AS_OF_DATE": "2026-01-29",
					"FUND_BREAKDOWNS_AS_OF_DATE": "2026-01-29"
				},
				"composition": {
					"totalNumberOfInstruments": 2,
					"compositionData": [
						{
							"compositionCharacteristics": {
								"quantity": 10,
								"bbg": "NVDA UW",
								"name": "NVIDIA CORP",
								"weight": 0.1,
								"currency": "USD",
								"type": "EQUITY_ORDINARY",
								"sector": "Information Technology",
								"isin": "US67066G1040",
								"countryOfRisk": "United States"
							},
							"weight": 0.1
						}
					]
				}
			}
		]
	}`

	client, err := New("de", WithHTTPClient(&sequenceHTTPClient{
		responses: []mockResponse{
			{status: http.StatusOK, body: discoveryResponse},
			{status: http.StatusOK, body: holdingsResponse},
		},
	}))
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	snapshot, err := client.Holdings(context.Background(), "C500")
	if err != nil {
		t.Fatalf("Holdings failed: %v", err)
	}

	if snapshot.TotalHoldings != 1 {
		t.Fatalf("expected 1 holding, got %d", snapshot.TotalHoldings)
	}
	if snapshot.AsOfDate.Format("2006-01-02") != "2026-01-29" {
		t.Fatalf("unexpected AsOfDate: %v", snapshot.AsOfDate)
	}
	if snapshot.Holdings[0].MarketValue != 100 {
		t.Fatalf("expected market value 100, got %f", snapshot.Holdings[0].MarketValue)
	}
	if snapshot.Holdings[0].Sector != etfscraper.SectorInformationTechnology {
		t.Errorf("expected sector %q, got %q", etfscraper.SectorInformationTechnology, snapshot.Holdings[0].Sector)
	}
	if snapshot.Holdings[0].AssetClass != etfscraper.AssetClassEquity {
		t.Errorf("expected asset class %q, got %q", etfscraper.AssetClassEquity, snapshot.Holdings[0].AssetClass)
	}
}

func TestHoldings_HTTPError(t *testing.T) {
	client, err := New("de", WithHTTPClient(&sequenceHTTPClient{
		responses: []mockResponse{
			{status: http.StatusOK, body: `{"products":[]}`},
			{status: http.StatusInternalServerError, body: `{"error":"boom"}`},
		},
	}))
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	_, err = client.Holdings(context.Background(), "UNKNOWN")
	if err == nil {
		t.Fatal("expected error for HTTP failure")
	}
}

func TestFundInfo_ByISIN(t *testing.T) {
	discoveryResponse := `{"products":[{"productId":"LU1135865084","productType":"PRODUCT","characteristics":{"ISIN":"LU1135865084","SHARE_MARKETING_NAME":"Test Fund","MNEMO":"C500","TER":0.1,"CURRENCY":"EUR","FUND_AUM":1000,"ASSET_CLASS":"Equity","DISTRIBUTION_POLICY":"Capitalisation"}}]}`

	client, err := New("de", WithHTTPClient(&testutil.MockHTTPClient{
		StatusCode:   http.StatusOK,
		ResponseBody: discoveryResponse,
	}))
	if err != nil {
		t.Fatalf("New() failed: %v", err)
	}

	fund, err := client.FundInfo(context.Background(), "LU1135865084")
	if err != nil {
		t.Fatalf("FundInfo failed: %v", err)
	}
	if fund.ISIN != "LU1135865084" {
		t.Fatalf("unexpected ISIN %s", fund.ISIN)
	}
}

type mockResponse struct {
	status int
	body   string
	err    error
}

type sequenceHTTPClient struct {
	responses []mockResponse
}

func (s *sequenceHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if len(s.responses) == 0 {
		return nil, errors.New("no mock responses available")
	}
	response := s.responses[0]
	s.responses = s.responses[1:]

	if response.err != nil {
		return nil, response.err
	}

	return &http.Response{
		StatusCode: response.status,
		Body:       io.NopCloser(strings.NewReader(response.body)),
		Header:     make(http.Header),
	}, nil
}
