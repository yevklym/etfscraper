package amundi

import (
	"fmt"
	"sort"
)

type regionConfig struct {
	BaseURL        string
	DiscoveryPath  string
	HoldingsPath   string
	CountryCode    string
	CountryName    string
	LanguageCode   string
	LanguageName   string
	DefaultHeaders map[string]string
}

var regionConfigs = map[string]regionConfig{
	"de": {
		BaseURL:       "https://www.amundietf.de",
		DiscoveryPath: "/mapi/ProductAPI/getProductsData",
		HoldingsPath:  "/mapi/ProductAPI/getProductsData",
		CountryCode:   "DEU",
		CountryName:   "Germany",
		LanguageCode:  "en",
		LanguageName:  "English",
		DefaultHeaders: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json, text/plain, */*",
			"Origin":       "https://www.amundietf.de",
			"Referer":      "https://www.amundietf.de/",
		},
	},
	"uk": {
		BaseURL:       "https://www.amundietf.co.uk",
		DiscoveryPath: "/mapi/ProductAPI/getProductsData",
		HoldingsPath:  "/mapi/ProductAPI/getProductsData",
		CountryCode:   "GBR",
		LanguageCode:  "en",
		LanguageName:  "English",
		DefaultHeaders: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json, text/plain, */*",
			"Origin":       "https://www.amundietf.co.uk",
			"Referer":      "https://www.amundietf.co.uk/",
		},
	},
	"fr": {
		BaseURL:       "https://www.amundietf.fr",
		DiscoveryPath: "/mapi/ProductAPI/getProductsData",
		HoldingsPath:  "/mapi/ProductAPI/getProductsData",
		CountryCode:   "FRA",
		LanguageCode:  "en",
		LanguageName:  "English",
		DefaultHeaders: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json, text/plain, */*",
			"Origin":       "https://www.amundietf.fr",
			"Referer":      "https://www.amundietf.fr/",
		},
	},
}

// buildDiscoveryRequest creates a discovery request for a specific region
func buildDiscoveryRequest(region string) (map[string]any, error) {
	cfg, ok := regionConfigs[region]
	if !ok {
		return nil, fmt.Errorf("unsupported region: %s (supported: %v)", region, SupportedRegions())
	}

	return map[string]any{
		"context": map[string]string{
			"countryCode":     cfg.CountryCode,
			"countryName":     cfg.CountryName,
			"languageCode":    cfg.LanguageCode,
			"languageName":    cfg.LanguageName,
			"userProfileName": "RETAIL",
		},
		"productType": "PRODUCT",
		"characteristics": []string{
			// Identity
			"ISIN",
			"SHARE_MARKETING_NAME",
			"MNEMO",
			"WKN",

			// Core Data
			"TER",
			"CURRENCY",
			"FUND_AUM",
			"NAV",
			"NAV_DATE_DISPLAYED",

			// Classification
			"ASSET_CLASS",
			"SUBASSET_CLASS",
			"DISTRIBUTION_POLICY",
			"FUND_ISSUER",
			"FUND_DOMICILIATION_COUNTRY",

			// ESG & Strategy
			"IS_ESG",
			"IS_CLIMATE",
			"IS_THEMATIC",
			"FUND_SFDR_CLASSIFICATION",
			"ESG_SCOPE",
			"IMPACT",
			"STRATEGY",

			// Listings
			"MAIN_LISTINGS",
			"LISTING_PLACES",

			// Metadata
			"INCEPTION_DATE",
			"FUND_REPLICATION_METHODOLOGY",
			"SRRI",
		},
		"metrics": []map[string]string{
			{"indicator": "shareCumulativePerformance", "period": "ONE_YEAR"},
			{"indicator": "shareCumulativePerformance", "period": "THREE_YEARS"},
		},
		"filters":       []any{},
		"sortCriterias": []any{},
		"historics":     []any{},
	}, nil
}

func buildHoldingsRequest(region, isin string) (map[string]any, error) {
	cfg, ok := regionConfigs[region]
	if !ok {
		return nil, fmt.Errorf("unsupported region: %s (supported: %v)", region, SupportedRegions())
	}

	return map[string]any{
		"context": map[string]string{
			"countryCode":     cfg.CountryCode,
			"countryName":     cfg.CountryName,
			"languageCode":    cfg.LanguageCode,
			"languageName":    cfg.LanguageName,
			"userProfileName": "RETAIL",
		},
		"productIds": []string{isin},
		"characteristics": []string{
			"ISIN",
			"SHARE_MARKETING_NAME",
			"FUND_AUM",
			"CURRENCY",
			"POSITION_AS_OF_DATE",
			"FUND_BREAKDOWNS_AS_OF_DATE",
		},
		"breakDown": map[string]any{
			"aggregationFields": []string{"FUND_TOP10"},
		},
		"composition": map[string]any{
			"compositionFields": []string{
				"name",
				"isin",
				"bbg",
				"weight",
				"quantity",
				"currency",
				"sector",
				"type",
				"countryOfRisk",
			},
		},
		"productType": "PRODUCT",
		"historics":   []any{},
	}, nil
}

func SupportedRegions() []string {
	regions := make([]string, 0, len(regionConfigs))
	for region := range regionConfigs {
		regions = append(regions, region)
	}
	sort.Strings(regions)
	return regions
}
