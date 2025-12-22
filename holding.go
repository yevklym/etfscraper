package etfscraper

import "time"

// Holding represents an individual holding within an ETF
type Holding struct {
	// Basic identification
	Ticker      string `json:"ticker,omitempty" validate:"required_without=ISIN"`
	ISIN        string `json:"isin,omitempty" validate:"len=12"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description,omitempty"`

	// Financial data
	Weight      float64  `json:"weight" validate:"min=0,max=1"` // decimal: 0.025 = 2.5%, 0.5 = 50%
	Quantity    float64  `json:"quantity,omitempty" validate:"min=0"`
	MarketValue float64  `json:"marketValue,omitempty" validate:"min=0"`
	Price       float64  `json:"price,omitempty" validate:"min=0"`
	Currency    Currency `json:"currency,omitempty"`

	// Classification
	Sector     Sector     `json:"sector,omitempty"`
	AssetClass AssetClass `json:"assetClass,omitempty"`
	Location   Location   `json:"location,omitempty"`
	Exchange   Exchange   `json:"exchange,omitempty"`

	// Timestamps
	LastUpdated time.Time `json:"lastUpdated,omitempty"`
}

// HoldingsSnapshot represents a point-in-time snapshot of fund holdings
type HoldingsSnapshot struct {
	Fund          Fund      `json:"fund" validate:"required"`
	AsOfDate      time.Time `json:"asOfDate" validate:"required"`
	Holdings      []Holding `json:"holdings" validate:"dive"`
	TotalHoldings int       `json:"totalHoldings"`
	TopHoldings   int       `json:"topHoldings,omitempty"` // Number of holdings included if truncated

	// Summary statistics
	TotalWeight float64 `json:"totalWeight,omitempty"`
	CashWeight  float64 `json:"cashWeight,omitempty"`

	// Metadata
	LastUpdated time.Time `json:"lastUpdated,omitempty"`
}
