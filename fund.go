package etfscraper

import "time"

// Fund holds comprehensive metadata about an ETF
type Fund struct {
	// Basic identification
	Ticker string `json:"ticker,omitempty"`
	ISIN   string `json:"isin,omitempty"`
	Name   string `json:"name"`

	// Provider information
	Provider ProviderName `json:"provider,omitempty"`

	// Financial details
	Currency Currency `json:"currency,omitempty"`

	// Fund characteristics
	InceptionDate  *time.Time `json:"inceptionDate,omitempty"`
	TotalAssets    float64    `json:"totalAssets,omitempty"`
	ExpenseRatio   float64    `json:"expenseRatio,omitempty"`
	IsDistributing bool       `json:"isDistributing,omitempty"`

	// Classification
	AssetClass AssetClass `json:"assetClass,omitempty"`

	// Trading information
	Exchange Exchange `json:"exchange,omitempty"`

	// Additional metadata
	LastUpdated      time.Time `json:"lastUpdated,omitzero"`
	ProviderMetadata any       `json:"-"`
}
