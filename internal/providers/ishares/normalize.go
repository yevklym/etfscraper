package ishares

import (
	"strings"

	"github.com/yevklym/etfscraper"
)

func normalizeCurrency(value string) etfscraper.Currency {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" || trimmed == "-" {
		return ""
	}

	switch strings.ToUpper(trimmed) {
	case "USD", "US DOLLAR", "U.S. DOLLAR":
		return etfscraper.CurrencyUSD
	case "EUR", "EURO":
		return etfscraper.CurrencyEUR
	case "GBP", "BRITISH POUND", "POUND STERLING":
		return etfscraper.CurrencyGBP
	case "JPY", "JAPANESE YEN":
		return etfscraper.CurrencyJPY
	case "CAD", "CANADIAN DOLLAR":
		return etfscraper.CurrencyCAD
	case "AUD", "AUSTRALIAN DOLLAR":
		return etfscraper.CurrencyAUD
	case "CHF", "SWISS FRANC":
		return etfscraper.CurrencyCHF
	case "CNY", "CHINESE YUAN":
		return etfscraper.CurrencyCNY
	case "INR", "INDIAN RUPEE":
		return etfscraper.CurrencyINR
	case "BRL", "BRAZILIAN REAL":
		return etfscraper.CurrencyBRL
	default:
		return etfscraper.Currency(strings.ToUpper(trimmed))
	}
}

// normalizeAssetClass maps an aladdinAssetClass value to a canonical AssetClass
func normalizeAssetClass(value string, mapping map[string]etfscraper.AssetClass) etfscraper.AssetClass {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" || trimmed == "-" {
		return ""
	}

	if ac, ok := mapping[strings.ToLower(trimmed)]; ok {
		return ac
	}

	return etfscraper.AssetClass(trimmed)
}

func normalizeExchange(value string) etfscraper.Exchange {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" || trimmed == "-" {
		return ""
	}

	normalized := strings.ToUpper(trimmed)

	switch normalized {
	case "NYSE", "NEW YORK STOCK EXCHANGE":
		return etfscraper.ExchangeNYSE
	case "NASDAQ":
		return etfscraper.ExchangeNASDAQ
	case "AMEX", "AMERICAN STOCK EXCHANGE":
		return etfscraper.ExchangeAMEX
	case "BATS":
		return etfscraper.ExchangeBATS
	case "LSE", "LONDON STOCK EXCHANGE":
		return etfscraper.ExchangeLSE
	case "EURONEXT", "EURONEXT PARIS", "EURONEXT AMSTERDAM", "EURONEXT BRUSSELS":
		return etfscraper.ExchangeEuronext
	case "XETRA":
		return etfscraper.ExchangeXetra
	case "TSE", "TOKYO STOCK EXCHANGE":
		return etfscraper.ExchangeTSE
	case "HKEX", "HONG KONG EXCHANGE":
		return etfscraper.ExchangeHKEX
	case "SSE", "SHANGHAI STOCK EXCHANGE":
		return etfscraper.ExchangeSSE
	case "SZSE", "SHENZHEN STOCK EXCHANGE":
		return etfscraper.ExchangeSZSE
	case "TSX", "TORONTO STOCK EXCHANGE":
		return etfscraper.ExchangeTSX
	case "ASX", "AUSTRALIAN SECURITIES EXCHANGE":
		return etfscraper.ExchangeASX
	default:
		return etfscraper.Exchange(trimmed)
	}
}
