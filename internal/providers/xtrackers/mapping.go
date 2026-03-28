package xtrackers

import (
	"strconv"
	"strings"
	"time"

	"github.com/yevklym/etfscraper"
)

// fundMetadata holds Xtrackers-specific data stored in Fund.ProviderMetadata.
type fundMetadata struct {
	ProductURL string
}

// mapCurrency maps a currency string to the canonical Currency type.
func mapCurrency(value string) etfscraper.Currency {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return ""
	}
	return etfscraper.Currency(strings.ToUpper(trimmed))
}

// normalizeAssetClass maps a Xtrackers asset class string to the canonical AssetClass.
func normalizeAssetClass(value string, mapping map[string]etfscraper.AssetClass) etfscraper.AssetClass {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return ""
	}

	if ac, ok := mapping[strings.ToLower(trimmed)]; ok {
		return ac
	}

	return etfscraper.AssetClass(trimmed)
}

// isDistributing returns whether the fund distributes based on the UseOfProfit field.
func isDistributing(useOfProfit string) bool {
	lower := strings.ToLower(strings.TrimSpace(useOfProfit))
	return strings.Contains(lower, "distributing") || strings.Contains(lower, "ausschüttend")
}

// parseAUM parses an AUM value from the API's sortValue field.
// The sortValue is a numeric value representing total assets.
func parseAUM(sortValue float64) float64 {
	return sortValue
}

// parseTER parses the TER value string (e.g. "0.12" or "0,12") and converts to decimal (0.0012).
func parseTER(value string) float64 {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return 0
	}

	// Handle European number formats (0,30 -> 0.30)
	trimmed = strings.ReplaceAll(trimmed, ",", ".")

	f, err := strconv.ParseFloat(trimmed, 64)
	if err != nil {
		return 0
	}
	return f / 100.0
}

// parseDateStr attempts to parse a date string from known regional formats.
func parseDateStr(value string) time.Time {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return time.Time{}
	}

	// Try German format: DD.MM.YYYY
	if t, err := time.Parse("02.01.2006", trimmed); err == nil {
		return t
	}
	// Try UK format: DD/MM/YYYY
	if t, err := time.Parse("02/01/2006", trimmed); err == nil {
		return t
	}
	// Try ISO format: YYYY-MM-DD
	if t, err := time.Parse("2006-01-02", trimmed); err == nil {
		return t
	}

	return time.Time{}
}

// parseLaunchDate parses a date string into a pointer.
func parseLaunchDate(value string) *time.Time {
	t := parseDateStr(value)
	if t.IsZero() {
		return nil
	}
	return &t
}

// parseLastUpdated parses the performance date.
func parseLastUpdated(value string) time.Time {
	return parseDateStr(value)
}
