package etfscraper

import "time"

// Fund holds comprehensive metadata about an ETF
type Fund struct {
	// Basic identification
	Ticker string `json:"ticker,omitempty" validate:"required"`
	ISIN   string `json:"isin,omitempty" validate:"len=12"`
	Name   string `json:"name" validate:"required"`

	// Provider information
	Provider ProviderName `json:"provider,omitempty"`

	// Financial details
	Currency Currency `json:"currency,omitempty"`

	// Fund characteristics
	InceptionDate  *time.Time `json:"inceptionDate,omitempty"`
	TotalAssets    float64    `json:"totalAssets,omitempty" validate:"min=0"`
	ExpenseRatio   float64    `json:"expenseRatio,omitempty" validate:"min=0,max=1"`
	DividendYield  float64    `json:"dividendYield,omitempty" validate:"min=0"`
	IsDistributing bool       `json:"isDistributing,omitempty"`

	// Classification
	Category   string     `json:"category,omitempty"`
	Geography  Location   `json:"geography,omitempty"`
	AssetClass AssetClass `json:"assetClass,omitempty"`

	// Trading information
	Exchange Exchange `json:"exchange,omitempty"`

	// Additional metadata
	LastUpdated      time.Time `json:"lastUpdated,omitzero"`
	ProviderMetadata any       `json:"-"`
}
