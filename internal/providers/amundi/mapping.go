package amundi

import (
	"strings"

	"github.com/yevklym/etfscraper"
)

type amundiFundMetadata struct {
	ProductID     string
	MainListings  map[string]string
	FundIssuer    string
	FundAUMInEuro float64
}

func pickTicker(c characteristics) string {
	if value := strings.TrimSpace(c.Mnemo); value != "" {
		return value
	}
	if value := strings.TrimSpace(c.Ticker); value != "" {
		return value
	}
	return pickListingTicker(c.MainListings)
}

func pickListingTicker(listings map[string]string) string {
	if listings == nil {
		return ""
	}
	if value, ok := listings["DEU"]; ok {
		return normalizeListing(value)
	}
	for _, value := range listings {
		if normalized := normalizeListing(value); normalized != "" {
			return normalized
		}
	}
	return ""
}

func normalizeListing(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return ""
	}
	fields := strings.Fields(trimmed)
	if len(fields) > 0 {
		return fields[0]
	}
	return trimmed
}

func pickISIN(p product) string {
	if value := strings.TrimSpace(p.Characteristics.ISIN); value != "" {
		return value
	}
	productID := strings.TrimSpace(p.ProductID)
	if len(productID) == 12 {
		return productID
	}
	return ""
}

func mapCurrency(value string) etfscraper.Currency {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return ""
	}

	switch strings.ToUpper(trimmed) {
	case "USD":
		return etfscraper.CurrencyUSD
	case "EUR":
		return etfscraper.CurrencyEUR
	case "GBP":
		return etfscraper.CurrencyGBP
	case "JPY":
		return etfscraper.CurrencyJPY
	case "CAD":
		return etfscraper.CurrencyCAD
	case "AUD":
		return etfscraper.CurrencyAUD
	case "CHF":
		return etfscraper.CurrencyCHF
	case "CNY":
		return etfscraper.CurrencyCNY
	case "INR":
		return etfscraper.CurrencyINR
	case "BRL":
		return etfscraper.CurrencyBRL
	default:
		return etfscraper.Currency(strings.ToUpper(trimmed))
	}
}

// normalizeAssetClass maps an asset class or holdings type value to a canonical AssetClass
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

// normalizeSector maps a sector value to a canonical Sector
func normalizeSector(value string, mapping map[string]etfscraper.Sector) etfscraper.Sector {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return ""
	}

	if s, ok := mapping[strings.ToLower(trimmed)]; ok {
		return s
	}

	return etfscraper.Sector(trimmed)
}

func isDistributing(policy string) bool {
	trimmed := strings.TrimSpace(policy)
	if trimmed == "" {
		return false
	}
	return strings.Contains(strings.ToLower(trimmed), "distribution")
}
