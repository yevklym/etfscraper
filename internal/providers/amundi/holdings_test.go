package amundi

import (
	"encoding/json"
	"testing"
	"time"
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

	holdings := convertHoldings(items, 1000)
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
