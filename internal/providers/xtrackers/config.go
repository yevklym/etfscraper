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
	LocationMapping   map[string]etfscraper.Location
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
		LocationMapping: map[string]etfscraper.Location{
			"vereinigte staaten von amerika": etfscraper.LocationUnitedStates,
			"vereinigte staaten":             etfscraper.LocationUnitedStates,
			"united states":                  etfscraper.LocationUnitedStates,
			"vereinigtes königreich":         etfscraper.LocationUnitedKingdom,
			"großbritannien (uk)":            etfscraper.LocationUnitedKingdom,
			"united kingdom":                 etfscraper.LocationUnitedKingdom,
			"japan":                          etfscraper.LocationJapan,
			"deutschland":                    etfscraper.LocationGermany,
			"germany":                        etfscraper.LocationGermany,
			"frankreich":                     etfscraper.LocationFrance,
			"france":                         etfscraper.LocationFrance,
			"schweiz":                        etfscraper.LocationSwitzerland,
			"switzerland":                    etfscraper.LocationSwitzerland,
			"kanada":                         etfscraper.LocationCanada,
			"canada":                         etfscraper.LocationCanada,
			"australien":                     etfscraper.LocationAustralia,
			"australia":                      etfscraper.LocationAustralia,
			"china":                          etfscraper.LocationChina,
			"taiwan":                         etfscraper.LocationTaiwan,
			"südkorea":                       etfscraper.LocationSouthKorea,
			"south korea":                    etfscraper.LocationSouthKorea,
			"indien":                         etfscraper.LocationIndia,
			"india":                          etfscraper.LocationIndia,
			"brasilien":                      etfscraper.LocationBrazil,
			"brazil":                         etfscraper.LocationBrazil,
			"niederlande":                    etfscraper.LocationNetherlands,
			"netherlands":                    etfscraper.LocationNetherlands,
			"schweden":                       etfscraper.LocationSweden,
			"sweden":                         etfscraper.LocationSweden,
			"italien":                        etfscraper.LocationItaly,
			"italy":                          etfscraper.LocationItaly,
			"spanien":                        etfscraper.LocationSpain,
			"spain":                          etfscraper.LocationSpain,
			"irland":                         etfscraper.LocationIreland,
			"ireland":                        etfscraper.LocationIreland,
			"dänemark":                       etfscraper.LocationDenmark,
			"denmark":                        etfscraper.LocationDenmark,
			"finnland":                       etfscraper.LocationFinland,
			"finland":                        etfscraper.LocationFinland,
			"türkei":                         etfscraper.LocationTurkey,
			"turkey":                         etfscraper.LocationTurkey,
			"belgien":                        etfscraper.LocationBelgium,
			"belgium":                        etfscraper.LocationBelgium,
			"österreich":                     etfscraper.LocationAustria,
			"austria":                        etfscraper.LocationAustria,
			"luxemburg":                      etfscraper.LocationLuxembourg,
			"luxembourg":                     etfscraper.LocationLuxembourg,
			"singapur":                       etfscraper.LocationSingapore,
			"singapore":                      etfscraper.LocationSingapore,
			"norwegen":                       etfscraper.LocationNorway,
			"norway":                         etfscraper.LocationNorway,
			"israel":                         etfscraper.LocationIsrael,
			"europäische union":              etfscraper.LocationEurope,
			"european union":                 etfscraper.LocationEurope,
			"cash und/oder derivate":         etfscraper.LocationCash,
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
		LocationMapping: map[string]etfscraper.Location{
			"united states":           etfscraper.LocationUnitedStates,
			"united kingdom":          etfscraper.LocationUnitedKingdom,
			"japan":                   etfscraper.LocationJapan,
			"germany":                 etfscraper.LocationGermany,
			"france":                  etfscraper.LocationFrance,
			"switzerland":             etfscraper.LocationSwitzerland,
			"canada":                  etfscraper.LocationCanada,
			"australia":               etfscraper.LocationAustralia,
			"china":                   etfscraper.LocationChina,
			"taiwan":                  etfscraper.LocationTaiwan,
			"south korea":             etfscraper.LocationSouthKorea,
			"india":                   etfscraper.LocationIndia,
			"brazil":                  etfscraper.LocationBrazil,
			"netherlands":             etfscraper.LocationNetherlands,
			"sweden":                  etfscraper.LocationSweden,
			"italy":                   etfscraper.LocationItaly,
			"spain":                   etfscraper.LocationSpain,
			"ireland":                 etfscraper.LocationIreland,
			"denmark":                 etfscraper.LocationDenmark,
			"finland":                 etfscraper.LocationFinland,
			"turkey":                  etfscraper.LocationTurkey,
			"belgium":                 etfscraper.LocationBelgium,
			"austria":                 etfscraper.LocationAustria,
			"luxembourg":              etfscraper.LocationLuxembourg,
			"singapore":               etfscraper.LocationSingapore,
			"norway":                  etfscraper.LocationNorway,
			"israel":                  etfscraper.LocationIsrael,
			"cash and/or derivatives": etfscraper.LocationCash,
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
