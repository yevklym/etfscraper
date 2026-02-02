package amundi

import (
	_ "embed"
	"sort"
)

//go:embed data/getProductsData-de.request.json
var requestBody []byte

type regionConfig struct {
	BaseURL        string
	DiscoveryPath  string
	RequestBody    []byte
	DefaultHeaders map[string]string
}

var regionConfigs = map[string]regionConfig{
	"de": {
		BaseURL:       "https://www.amundietf.de",
		DiscoveryPath: "/mapi/ProductAPI/getProductsData",
		RequestBody:   requestBody,
		DefaultHeaders: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json, text/plain, */*",
			"Origin":       "https://www.amundietf.de",
			"Referer":      "https://www.amundietf.de/",
		},
	},
}

func SupportedRegions() []string {
	regions := make([]string, 0, len(regionConfigs))
	for region := range regionConfigs {
		regions = append(regions, region)
	}

	sort.Strings(regions)
	return regions
}
