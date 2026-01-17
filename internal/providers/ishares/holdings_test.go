package ishares

import (
	"strings"
	"testing"
	"time"

	"github.com/yevklym/etfscraper"
)

func TestParseHoldings(t *testing.T) {
	csvData := `iShares S&P Mid-Cap 400 Value ETF
Fund Holdings as of,"Nov 14, 2025"
Inception Date,"Jul 24, 2000"
Shares Outstanding,"60,350,000.00"
Stock,"-"
Bond,"-"
Cash,"-"
Other,"-"

Ticker,Name,Type,Sector,Asset Class,Market Value,Notional Value,Quantity,Price,Location,Exchange,Currency,FX Rate,Market Currency,Accrual Date,Notional Weight,Market Weight
"FLEX","FLEX LTD","EQUITY","Information Technology","Equity","119,376,740.56","119,376,740.56","1,989,944.00","59.99","United States","NASDAQ","USD","1.00","USD","-","1.55","1.55"
"TLN","TALEN ENERGY CORP","EQUITY","Utilities","Equity","87,460,299.92","87,460,299.92","242,326.00","360.92","United States","NASDAQ","USD","1.00","USD","-","1.14","1.14"
"FAZ5","S&P MID 400 EMINI DEC 25","INDEX","Cash and/or Derivatives","Futures","0.00","10,282,880.00","32.00","3,213.40","-","Index And Options Market","USD","1.00","USD","-","0.13","-"

"The content contained herein..."`

	c := &Client{}
	fund := &etfscraper.Fund{Ticker: "IJJ", Name: "iShares S&P Mid-Cap 400 Value ETF"}

	snapshot, err := c.parseHoldings(strings.NewReader(csvData), fund)
	if err != nil {
		t.Errorf("parseHoldings failed: %v", err)
	}

	expectedDate := time.Date(2025, time.November, 14, 0, 0, 0, 0, time.UTC)
	if !snapshot.AsOfDate.Equal(expectedDate) {
		t.Errorf("Expected AsOfDate %v, got %v", expectedDate, snapshot.AsOfDate)
	}

	if snapshot.TotalHoldings != 3 {
		t.Errorf("Expected 3 holdings, got %d", snapshot.TotalHoldings)
	}

	flex := snapshot.Holdings[0]
	if flex.Ticker != "FLEX" {
		t.Errorf("Expected first ticker FLEX, got %s", flex.Ticker)
	}
	if flex.Sector != "Information Technology" {
		t.Errorf("Expected sector Information Technology, got %s", flex.Sector)
	}
	if flex.Weight != 0.0155 {
		t.Errorf("Expected weight 0.0155, got %f", flex.Weight)
	}

	faz := snapshot.Holdings[2]
	if faz.Ticker != "FAZ5" {
		t.Errorf("Expected ticker FAZ5, got %s", faz.Ticker)
	}
	if faz.AssetClass != "Futures" {
		t.Errorf("Expected Asset Class Futures, got %s", faz.AssetClass)
	}
}

func TestParseHoldings_InternationalFormat(t *testing.T) {
	csvData := `Fund Holdings as of,"Sep 30, 2025"

   Name,Market Value,Weight (%),Quantity
   "ASML HOLDING NV","246,946,976.41","1.93","253,795.00"
   "SAP","180,198,609.94","1.41","672,929.00"
   "JPY/USD","-109,819.70","0.00","-2,579,800,311.00"

   "The content contained herein..."`

	c := &Client{}
	fund := &etfscraper.Fund{Ticker: "IEUR", Name: "iShares Core MSCI Europe ETF"}

	snapshot, err := c.parseHoldings(strings.NewReader(csvData), fund)
	if err != nil {
		t.Fatalf("parseHoldings failed: %v", err)
	}

	expectedDate := time.Date(2025, time.September, 30, 0, 0, 0, 0, time.UTC)
	if !snapshot.AsOfDate.Equal(expectedDate) {
		t.Errorf("Expected AsOfDate %v, got %v", expectedDate, snapshot.AsOfDate)
	}

	if snapshot.TotalHoldings != 3 {
		t.Errorf("Expected 3 holdings, got %d", snapshot.TotalHoldings)
	}

	asml := snapshot.Holdings[0]
	if asml.Name != "ASML HOLDING NV" {
		t.Errorf("Expected name ASML HOLDING NV, got %s", asml.Name)
	}
	if asml.Ticker != "" {
		t.Errorf("Expected empty ticker, got %s", asml.Ticker)
	}
	epsilon := 0.00001
	if diff := asml.Weight - 0.0193; diff > epsilon || diff < -epsilon {
		t.Errorf("Expected weight 0.0193, got %f", asml.Weight)
	}

	if asml.Quantity != 253795.0 {
		t.Errorf("Expected quantity 253795.0, got %f", asml.Quantity)
	}
}
