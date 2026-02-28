package etfscraper

import "time"

// Holding represents an individual holding within an ETF
type Holding struct {
	// Basic identification
	Ticker string `json:"ticker,omitempty"`
	ISIN   string `json:"isin,omitempty"`
	Name   string `json:"name"`

	// Financial data
	Weight      float64  `json:"weight"` // decimal: 0.025 = 2.5%, 0.5 = 50%
	Quantity    float64  `json:"quantity,omitempty"`
	MarketValue float64  `json:"marketValue,omitempty"`
	Price       float64  `json:"price,omitempty"`
	Currency    Currency `json:"currency,omitempty"`

	// Classification
	Sector     Sector     `json:"sector,omitempty"`
	AssetClass AssetClass `json:"assetClass,omitempty"`
	Location   Location   `json:"location,omitempty"`
	Exchange   Exchange   `json:"exchange,omitempty"`

	// Timestamps
	LastUpdated time.Time `json:"lastUpdated,omitzero"`
}

// HoldingsSnapshot represents a point-in-time snapshot of fund holdings
type HoldingsSnapshot struct {
	Fund          Fund      `json:"fund"`
	AsOfDate      time.Time `json:"asOfDate"`
	Holdings      []Holding `json:"holdings"`
	TotalHoldings int       `json:"totalHoldings"`

	// Metadata
	LastUpdated time.Time `json:"lastUpdated,omitzero"`
}
