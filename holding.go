package etfscraper

import "time"

// Holding represents an individual holding within an ETF.
type Holding struct {
	// Ticker is the holding's exchange symbol, if available.
	Ticker string `json:"ticker,omitempty"`
	// ISIN is the holding's International Securities Identification Number.
	ISIN string `json:"isin,omitempty"`
	// Name is the holding's full name as reported by the provider.
	Name string `json:"name"`

	// Weight is the holding's portfolio weight as a decimal (0.025 = 2.5%).
	Weight float64 `json:"weight"`
	// Quantity is the number of shares or units held.
	Quantity float64 `json:"quantity,omitempty"`
	// MarketValue is the holding's total market value in Currency.
	MarketValue float64 `json:"marketValue,omitempty"`
	// Price is the per-unit price in Currency.
	Price float64 `json:"price,omitempty"`
	// Currency is the holding's trading currency (ISO 4217 code).
	Currency Currency `json:"currency,omitempty"`

	// Sector is the GICS economic sector.
	Sector Sector `json:"sector,omitempty"`
	// AssetClass is the holding's asset type (e.g. Equity, Bond).
	AssetClass AssetClass `json:"assetClass,omitempty"`
	// Location is the holding's country or region as reported by the provider.
	// Values are provider-specific free-form strings (e.g. "United States",
	// "Deutschland", "Etats-Unis").
	Location Location `json:"location,omitempty"`
	// Exchange is the exchange where the holding is traded.
	Exchange Exchange `json:"exchange,omitempty"`

	// LastUpdated records when this Holding data was fetched.
	LastUpdated time.Time `json:"lastUpdated,omitzero"`
}

// HoldingsSnapshot represents a point-in-time snapshot of fund holdings.
type HoldingsSnapshot struct {
	// Fund is a copy of the fund metadata that these holdings belong to.
	Fund Fund `json:"fund"`
	// AsOfDate is the date the holdings data is effective, as reported by
	// the provider (not the time the request was made).
	AsOfDate time.Time `json:"asOfDate"`
	// Holdings is the list of individual holdings in the fund.
	Holdings []Holding `json:"holdings"`
	// TotalHoldings is the number of holdings returned (always len(Holdings)).
	TotalHoldings int `json:"totalHoldings"`

	// LastUpdated records when this snapshot was fetched.
	LastUpdated time.Time `json:"lastUpdated,omitzero"`
}
