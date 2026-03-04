package etfscraper

import "time"

// Fund holds comprehensive metadata about an ETF.
type Fund struct {
	// Ticker is the exchange trading symbol (e.g. "IVV", "EUNL").
	Ticker string `json:"ticker,omitempty"`
	// ISIN is the International Securities Identification Number (e.g. "IE00B5BMR087").
	ISIN string `json:"isin,omitempty"`
	// Name is the full fund name as reported by the provider.
	Name string `json:"name"`

	// Provider identifies which data provider this fund came from.
	Provider ProviderName `json:"provider,omitempty"`

	// Currency is the fund's base currency (ISO 4217 code).
	Currency Currency `json:"currency,omitempty"`

	// InceptionDate is the fund's launch date. It is nil when the provider
	// does not supply this information.
	InceptionDate *time.Time `json:"inceptionDate,omitempty"`
	// TotalAssets is the fund's total net assets in the fund's base Currency.
	TotalAssets float64 `json:"totalAssets,omitempty"`
	// ExpenseRatio is the annual expense ratio as a decimal (0.0003 = 0.03%).
	ExpenseRatio float64 `json:"expenseRatio,omitempty"`
	// IsDistributing is true for funds that pay dividends, false for
	// accumulating funds that reinvest income.
	IsDistributing bool `json:"isDistributing,omitempty"`

	// AssetClass is the fund's primary asset class (e.g. Equity, Bond).
	AssetClass AssetClass `json:"assetClass,omitempty"`

	// Exchange is the primary exchange where the fund is listed.
	Exchange Exchange `json:"exchange,omitempty"`

	// LastUpdated records when this Fund data was fetched.
	LastUpdated time.Time `json:"lastUpdated,omitzero"`
	// ProviderMetadata holds provider-specific data not represented by the
	// common fields above. Its concrete type varies by provider and is
	// intended for internal use; callers should not rely on its contents.
	ProviderMetadata any `json:"-"`
}
