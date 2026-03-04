package ishares

import (
	"sort"

	"github.com/yevklym/etfscraper"
)

type regionConfig struct {
	DiscoveryURL        string
	BaseURL             string
	HoldingsURLTemplate string
	DefaultCurrency     etfscraper.Currency
	DefaultExchange     etfscraper.Exchange
	AssetClassMapping   map[string]etfscraper.AssetClass
	SectorMapping       map[string]etfscraper.Sector
	ColumnMappings      ColumnMapper
	MonthTranslations   map[string]string
	DateFormats         []string
	DateHeaderPatterns  []string
}

type ColumnMapper struct {
	Name           []string
	Ticker         []string
	ISIN           []string
	MarketValue    []string
	Weight         []string
	Quantity       []string
	ParValue       []string // For bonds
	Price          []string
	Sector         []string
	AssetClass     []string
	Location       []string
	Exchange       []string
	Currency       []string
	MarketCurrency []string
}

var regionConfigs = map[string]regionConfig{
	"us": {
		DiscoveryURL:        "https://www.ishares.com/us/product-screener/product-screener-v3.1.jsn?dcrPath=/templatedata/config/product-screener-v3/data/en/us-ishares/ishares-product-screener-backend-config&siteEntryPassthrough=true",
		BaseURL:             "https://www.ishares.com",
		HoldingsURLTemplate: "%s%s/1467271812596.ajax?fileType=csv",
		DefaultCurrency:     etfscraper.CurrencyUSD,
		DefaultExchange:     etfscraper.ExchangeNYSE,
		AssetClassMapping: map[string]etfscraper.AssetClass{
			"equity":         etfscraper.AssetClassEquity,
			"fixed income":   etfscraper.AssetClassBond,
			"cash":           etfscraper.AssetClassCash,
			"commodity":      etfscraper.AssetClassCommodity,
			"real estate":    etfscraper.AssetClassRealEstate,
			"digital assets": etfscraper.AssetClassCryptocurrency,
			"multi asset":    etfscraper.AssetClassAlternative,
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
			"communication":          etfscraper.SectorTelecommunication,
			"utilities":              etfscraper.SectorUtilities,
			"real estate":            etfscraper.SectorRealEstate,
		},
		MonthTranslations:  nil,
		DateFormats:        []string{"Jan 2, 2006"},
		DateHeaderPatterns: []string{"Fund Holdings as of"},
		ColumnMappings: ColumnMapper{
			Name:           []string{"Name"},
			Ticker:         []string{"Ticker"},
			ISIN:           []string{"ISIN"},
			MarketValue:    []string{"Market Value"},
			Weight:         []string{"Weight (%)", "Market Weight"},
			Quantity:       []string{"Quantity"},
			Price:          []string{"Price"},
			Sector:         []string{"Sector"},
			AssetClass:     []string{"Asset Class"},
			Location:       []string{"Location"},
			Exchange:       []string{"Exchange"},
			Currency:       []string{"Currency"},
			MarketCurrency: []string{"Market Currency"},
		},
	},
	"de": {
		DiscoveryURL:        "https://www.ishares.com/de/privatanleger/de/product-screener/product-screener-v3.1.jsn?dcrPath=/templatedata/config/product-screener-v3/data/de/germany/product-screener/ishares-product-screener-backend-config&siteEntryPassthrough=true",
		BaseURL:             "https://www.ishares.com",
		HoldingsURLTemplate: "%s%s/fund/1478358465952.ajax?fileType=csv",
		DefaultCurrency:     etfscraper.CurrencyEUR,
		DefaultExchange:     etfscraper.ExchangeXetra,
		AssetClassMapping: map[string]etfscraper.AssetClass{
			"aktien":         etfscraper.AssetClassEquity,
			"anleihen":       etfscraper.AssetClassBond,
			"barmittel":      etfscraper.AssetClassCash,
			"rohstoffe":      etfscraper.AssetClassCommodity,
			"immobilien":     etfscraper.AssetClassRealEstate,
			"digital assets": etfscraper.AssetClassCryptocurrency,
			"multi-asset":    etfscraper.AssetClassAlternative,
		},
		SectorMapping: map[string]etfscraper.Sector{
			"energie":                    etfscraper.SectorEnergy,
			"materialien":                etfscraper.SectorMaterials,
			"industrie":                  etfscraper.SectorIndustrials,
			"zyklische konsumgüter":      etfscraper.SectorConsumerDiscretionary,
			"nichtzyklische konsumgüter": etfscraper.SectorConsumerStaples,
			"gesundheitsversorgung":      etfscraper.SectorHealthcare,
			"financials":                 etfscraper.SectorFinancials,
			"it":                         etfscraper.SectorInformationTechnology,
			"kommunikation":              etfscraper.SectorTelecommunication,
			"versorger":                  etfscraper.SectorUtilities,
			"immobilien":                 etfscraper.SectorRealEstate,
		},
		MonthTranslations: map[string]string{
			// Full German month names (must be matched before abbreviations)
			"Januar":    "Jan",
			"Februar":   "Feb",
			"März":      "Mar",
			"April":     "Apr",
			"Mai":       "May",
			"Juni":      "Jun",
			"Juli":      "Jul",
			"August":    "Aug",
			"September": "Sep",
			"Oktober":   "Oct",
			"November":  "Nov",
			"Dezember":  "Dec",
			// Abbreviated forms
			"Jan": "Jan",
			"Feb": "Feb",
			"Mär": "Mar",
			"Apr": "Apr",
			"Jun": "Jun",
			"Jul": "Jul",
			"Aug": "Aug",
			"Sep": "Sep",
			"Okt": "Oct",
			"Nov": "Nov",
			"Dez": "Dec",
		},
		DateFormats:        []string{"02.Jan.2006", "02.Jan2006"},
		DateHeaderPatterns: []string{"Fondsposition per", "Fondsbestände am"},
		ColumnMappings: ColumnMapper{
			Name:           []string{"Name"},
			Ticker:         []string{"Emittententicker"},
			ISIN:           []string{"ISIN"},
			MarketValue:    []string{"Marktwert"},
			Weight:         []string{"Gewichtung (%)"},
			Quantity:       []string{"Nominale"},
			Price:          []string{"Kurs"},
			Sector:         []string{"Sektor"},
			AssetClass:     []string{"Anlageklasse"},
			Location:       []string{"Standort"},
			Exchange:       []string{"Börse"},
			Currency:       []string{"Währung"},
			MarketCurrency: []string{"Marktwährung"},
		},
	},
	"uk": {
		DiscoveryURL:        "https://www.ishares.com/uk/individual/en/product-screener/product-screener-v3.1.jsn?dcrPath=/templatedata/config/product-screener-v3/data/en/uk/product-screener/ishares-product-screener-backend-config&siteEntryPassthrough=true",
		BaseURL:             "https://www.ishares.com",
		HoldingsURLTemplate: "%s%s/fund/1506575576011.ajax?fileType=csv",
		DefaultCurrency:     etfscraper.CurrencyGBP,
		DefaultExchange:     etfscraper.ExchangeLSE,
		AssetClassMapping: map[string]etfscraper.AssetClass{
			"equity":         etfscraper.AssetClassEquity,
			"fixed income":   etfscraper.AssetClassBond,
			"cash":           etfscraper.AssetClassCash,
			"commodity":      etfscraper.AssetClassCommodity,
			"real estate":    etfscraper.AssetClassRealEstate,
			"digital assets": etfscraper.AssetClassCryptocurrency,
			"multi asset":    etfscraper.AssetClassAlternative,
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
			"communication":          etfscraper.SectorTelecommunication,
			"utilities":              etfscraper.SectorUtilities,
			"real estate":            etfscraper.SectorRealEstate,
		},
		MonthTranslations:  nil,
		DateFormats:        []string{"02/Jan/2006"},
		DateHeaderPatterns: []string{"Fund Holdings as of"},
		ColumnMappings: ColumnMapper{
			Name:           []string{"Name"},
			Ticker:         []string{"Ticker"},
			ISIN:           []string{"ISIN"},
			MarketValue:    []string{"Market Value"},
			Weight:         []string{"Weight (%)", "Market Weight"},
			Quantity:       []string{"Quantity"},
			Price:          []string{"Price"},
			Sector:         []string{"Sector"},
			AssetClass:     []string{"Asset Class"},
			Location:       []string{"Location"},
			Exchange:       []string{"Exchange"},
			Currency:       []string{"Currency"},
			MarketCurrency: []string{"Market Currency"},
		},
	},
	"fr": {
		DiscoveryURL:        "https://www.blackrock.com/fr/particuliers/product-screener/product-screener-v3.1.jsn?dcrPath=/templatedata/config/product-screener-v3/data/fr/France/product-screener-backend-config&siteEntryPassthrough=true",
		BaseURL:             "https://www.blackrock.com",
		HoldingsURLTemplate: "%s%s/1499538099380.ajax?fileType=csv",
		DefaultCurrency:     etfscraper.CurrencyEUR,
		DefaultExchange:     etfscraper.ExchangeEuronext,
		AssetClassMapping: map[string]etfscraper.AssetClass{
			"actions":            etfscraper.AssetClassEquity,
			"obligations":        etfscraper.AssetClassBond,
			"liquidités":         etfscraper.AssetClassCash,
			"matières premières": etfscraper.AssetClassCommodity,
			"immobilier":         etfscraper.AssetClassRealEstate,
			"digital assets":     etfscraper.AssetClassCryptocurrency,
			"multi-actifs":       etfscraper.AssetClassAlternative,
			"marchés privés":     etfscraper.AssetClassAlternative,
		},
		SectorMapping: map[string]etfscraper.Sector{
			"energie":                         etfscraper.SectorEnergy,
			"matériaux":                       etfscraper.SectorMaterials,
			"industries":                      etfscraper.SectorIndustrials,
			"biens de consommation cycliques": etfscraper.SectorConsumerDiscretionary,
			"biens de consommation de base":   etfscraper.SectorConsumerStaples,
			"santé":                           etfscraper.SectorHealthcare,
			"finance":                         etfscraper.SectorFinancials,
			"technologie de l'information":    etfscraper.SectorInformationTechnology,
			"la communication":                etfscraper.SectorTelecommunication,
			"services publics":                etfscraper.SectorUtilities,
			"immobilier":                      etfscraper.SectorRealEstate,
		},
		MonthTranslations: map[string]string{
			"janv.": "Jan",
			"févr.": "Feb",
			"mars":  "Mar",
			"avr.":  "Apr",
			"mai":   "May",
			"juin":  "Jun",
			"juil.": "Jul",
			"août":  "Aug",
			"sept.": "Sep",
			"oct.":  "Oct",
			"nov.":  "Nov",
			"déc.":  "Dec",
		},
		DateFormats:        []string{"02/Jan/2006"},
		DateHeaderPatterns: []string{"Fund Holdings as of"},
		ColumnMappings: ColumnMapper{
			Name:           []string{"Name"},
			Ticker:         []string{"Ticker"},
			ISIN:           []string{"ISIN"},
			MarketValue:    []string{"Market Value"},
			Weight:         []string{"Weight (%)", "Market Weight"},
			Quantity:       []string{"Shares", "Quantity"},
			Price:          []string{"Price"},
			Sector:         []string{"Sector"},
			AssetClass:     []string{"Asset Class"},
			Location:       []string{"Location"},
			Exchange:       []string{"Exchange"},
			Currency:       []string{"Currency"},
			MarketCurrency: []string{"Market Currency"},
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
