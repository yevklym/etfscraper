package xtrackers

import (
	"fmt"
	"sort"

	"github.com/yevklym/etfscraper"
)

type regionConfig struct {
	BaseURL           string
	DiscoveryPath     string
	Locale            string
	DefaultHeaders    map[string]string
	AssetClassMapping map[string]etfscraper.AssetClass
	SectorMapping     map[string]etfscraper.Sector
}

var regionConfigs = map[string]regionConfig{
	"de": {
		BaseURL:       "https://etf.dws.com",
		DiscoveryPath: "/api/fundfinder/%s/datatable",
		Locale:        "de-de",
		DefaultHeaders: map[string]string{
			"Content-Type":   "application/json",
			"Accept":         "application/json",
			"client-id":      "passive-frontend",
			"User-Agent":     "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.0 Safari/605.1.15",
			"Origin":         "https://etf.dws.com",
			"Referer":        "https://etf.dws.com/de-de/produktfinder/",
			"Sec-Fetch-Site": "same-origin",
			"Sec-Fetch-Mode": "cors",
			"Sec-Fetch-Dest": "empty",
		},
		AssetClassMapping: map[string]etfscraper.AssetClass{
			"aktien":    etfscraper.AssetClassEquity,
			"renten":    etfscraper.AssetClassBond,
			"rohstoffe": etfscraper.AssetClassCommodity,
			"equity":    etfscraper.AssetClassEquity,
		},
		SectorMapping: map[string]etfscraper.Sector{
			"energie":                 etfscraper.SectorEnergy,
			"material":                etfscraper.SectorMaterials,
			"industrieunternehmen":    etfscraper.SectorIndustrials,
			"nicht-basiskonsumgüter":  etfscraper.SectorConsumerDiscretionary,
			"basiskonsumgüter":        etfscraper.SectorConsumerStaples,
			"gesundheitswesen":        etfscraper.SectorHealthcare,
			"finanzen":                etfscraper.SectorFinancials,
			"informationstechnologie": etfscraper.SectorInformationTechnology,
			"kommunikationsdienste":   etfscraper.SectorTelecommunication,
			"versorgungsunternehmen":  etfscraper.SectorUtilities,
			"immobilien":              etfscraper.SectorRealEstate,
		},
	},
	"uk": {
		BaseURL:       "https://etf.dws.com",
		DiscoveryPath: "/api/fundfinder/%s/datatable",
		Locale:        "en-gb",
		DefaultHeaders: map[string]string{
			"Content-Type":   "application/json",
			"Accept":         "application/json",
			"client-id":      "passive-frontend",
			"User-Agent":     "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.0 Safari/605.1.15",
			"Origin":         "https://etf.dws.com",
			"Referer":        "https://etf.dws.com/en-gb/product-finder/",
			"Sec-Fetch-Site": "same-origin",
			"Sec-Fetch-Mode": "cors",
			"Sec-Fetch-Dest": "empty",
		},
		AssetClassMapping: map[string]etfscraper.AssetClass{
			"equities":     etfscraper.AssetClassEquity,
			"fixed income": etfscraper.AssetClassBond,
			"commodities":  etfscraper.AssetClassCommodity,
			"multi asset":  etfscraper.AssetClassAlternative,
			"alternatives": etfscraper.AssetClassAlternative,
			"equity":       etfscraper.AssetClassEquity,
		},
		SectorMapping: map[string]etfscraper.Sector{
			"energy":                 etfscraper.SectorEnergy,
			"materials":              etfscraper.SectorMaterials,
			"industrials":            etfscraper.SectorIndustrials,
			"consumer discretionary": etfscraper.SectorConsumerDiscretionary,
			"consumer staples":       etfscraper.SectorConsumerStaples,
			"health care":            etfscraper.SectorHealthcare,
			"financials":             etfscraper.SectorFinancials,
			"information technology": etfscraper.SectorInformationTechnology,
			"communication services": etfscraper.SectorTelecommunication,
			"utilities":              etfscraper.SectorUtilities,
			"real estate":            etfscraper.SectorRealEstate,
		},
	},
}

// discoveryRequestBody returns the exact JSON payload expected by the /datatable API.
func discoveryRequestBody() map[string]any {
	return map[string]any{
		"searchTerm": "",
		"filters":    []any{},
	}
}

// SupportedRegions returns all supported region codes.
func SupportedRegions() []string {
	regions := make([]string, 0, len(regionConfigs))
	for region := range regionConfigs {
		regions = append(regions, region)
	}
	sort.Strings(regions)
	return regions
}

// discoveryURL returns the full URL for the discovery endpoint.
func (c *Client) discoveryURL() string {
	path := fmt.Sprintf(c.config.DiscoveryPath, c.config.Locale)
	return c.config.BaseURL + path
}
