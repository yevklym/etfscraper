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
"FLEX","FLEX LTD","EQUITY","Information Technology","Equity","119,376,740.56","119,376,740.56","1,989,944.00","59.99","United States","NASDAQ","USD","1.00","EUR","-","1.55","1.55"
"TLN","TALEN ENERGY CORP","EQUITY","Utilities","Equity","87,460,299.92","87,460,299.92","242,326.00","360.92","United States","NASDAQ","USD","1.00","USD","-","1.14","1.14"
"FAZ5","S&P MID 400 EMINI DEC 25","INDEX","Cash and/or Derivatives","Futures","0.00","10,282,880.00","32.00","3,213.40","-","Index And Options Market","USD","1.00","USD","-","0.13","0.00"

"The content contained herein..."`

	c, err := New("us")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	fund := &etfscraper.Fund{Ticker: "IJJ", Name: "iShares S&P Mid-Cap 400 Value ETF"}

	snapshot, err := c.parseHoldings(strings.NewReader(csvData), fund)
	if err != nil {
		t.Fatalf("parseHoldings failed: %v", err)
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
	if flex.Currency != etfscraper.CurrencyEUR {
		t.Errorf("Expected currency EUR, got %s", flex.Currency)
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

	c, err := New("us")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

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

func TestParseHoldings_FrenchFormat(t *testing.T) {
	csvData := "iShares S&P 500 (Acc)\n" +
		"Fund Holdings as of,\"26/f\u00e9vr./2026\"\n" +
		"Inception Date,\"15/sept./2010\"\n" +
		"Shares Outstanding,\"100 000 000,00\"\n" +
		"\n" +
		"Ticker,Name,Sector,Asset Class,Market Value,Weight (%),Notional Value,Shares,Price,Location,Exchange,Market Currency\n" +
		"\"NVDA\",\"NVIDIA CORP\",\"Technologie de l'information\",\"Actions\",\"10 569 831 271,35\",\"7,60\",\"10 569 831 271,35\",\"57 168 215,00\",\"184,89\",\"Etats-Unis\",\"NASDAQ\",\"USD\"\n" +
		"\"AAPL\",\"APPLE INC\",\"Technologie de l'information\",\"Actions\",\"9 488 522 637,00\",\"6,82\",\"9 488 522 637,00\",\"34 762 860,00\",\"272,95\",\"Etats-Unis\",\"NASDAQ\",\"USD\"\n" +
		"\"MSFT\",\"MICROSOFT CORP\",\"Technologie de l'information\",\"Actions\",\"7 024 245 332,72\",\"5,05\",\"7 024 245 332,72\",\"17 485 426,00\",\"401,72\",\"Etats-Unis\",\"NASDAQ\",\"USD\"\n" +
		"\n" +
		"\"The content contained herein...\""

	c, err := New("fr")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	fund := &etfscraper.Fund{Ticker: "CSPX", Name: "iShares S&P 500 (Acc)"}

	snapshot, err := c.parseHoldings(strings.NewReader(csvData), fund)
	if err != nil {
		t.Fatalf("parseHoldings failed: %v", err)
	}

	expectedDate := time.Date(2026, time.February, 26, 0, 0, 0, 0, time.UTC)
	if !snapshot.AsOfDate.Equal(expectedDate) {
		t.Errorf("Expected AsOfDate %v, got %v", expectedDate, snapshot.AsOfDate)
	}

	if snapshot.TotalHoldings != 3 {
		t.Errorf("Expected 3 holdings, got %d", snapshot.TotalHoldings)
	}

	nvidia := snapshot.Holdings[0]
	if nvidia.Ticker != "NVDA" {
		t.Errorf("Expected first ticker NVDA, got %s", nvidia.Ticker)
	}
	if nvidia.Name != "NVIDIA CORP" {
		t.Errorf("Expected name NVIDIA CORP, got %s", nvidia.Name)
	}
	if nvidia.Sector != "Technologie de l'information" {
		t.Errorf("Expected sector Technologie de l'information, got %s", nvidia.Sector)
	}
	if nvidia.AssetClass != "Actions" {
		t.Errorf("Expected asset class Actions, got %s", nvidia.AssetClass)
	}
	if nvidia.Location != "Etats-Unis" {
		t.Errorf("Expected location Etats-Unis, got %s", nvidia.Location)
	}
	if nvidia.Currency != etfscraper.CurrencyUSD {
		t.Errorf("Expected currency USD, got %s", nvidia.Currency)
	}

	epsilon := 0.01
	if diff := nvidia.MarketValue - 10569831271.35; diff > epsilon || diff < -epsilon {
		t.Errorf("Expected market value 10569831271.35, got %f", nvidia.MarketValue)
	}

	if diff := nvidia.Weight - 0.076; diff > 0.001 || diff < -0.001 {
		t.Errorf("Expected weight ~0.076, got %f", nvidia.Weight)
	}

	if diff := nvidia.Quantity - 57168215.0; diff > epsilon || diff < -epsilon {
		t.Errorf("Expected quantity 57168215.0, got %f", nvidia.Quantity)
	}

	if diff := nvidia.Price - 184.89; diff > epsilon || diff < -epsilon {
		t.Errorf("Expected price 184.89, got %f", nvidia.Price)
	}

	msft := snapshot.Holdings[2]
	if msft.Ticker != "MSFT" {
		t.Errorf("Expected third ticker MSFT, got %s", msft.Ticker)
	}
}
