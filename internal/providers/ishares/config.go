package ishares

type regionConfig struct {
	DiscoveryURL       string
	BaseURL            string
	HoldingsDownloadID string
}

var regionConfigs = map[string]regionConfig{
	"us": {
		DiscoveryURL:       "https://www.ishares.com/us/product-screener/product-screener-v3.1.jsn?dcrPath=/templatedata/config/product-screener-v3/data/en/us-ishares/ishares-product-screener-backend-config&siteEntryPassthrough=true",
		BaseURL:            "https://www.ishares.com",
		HoldingsDownloadID: "1467271812596",
	},
	"de": {
		DiscoveryURL:       "https://www.ishares.com/de/privatanleger/de/product-screener/product-screener-v3.1.jsn?dcrPath=/templatedata/config/product-screener-v3/data/de/germany/product-screener/ishares-product-screener-backend-config&siteEntryPassthrough=true",
		BaseURL:            "https://www.ishares.com",
		HoldingsDownloadID: "1506575576011",
	},
}
