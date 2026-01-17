package ishares

type regionConfig struct {
	DiscoveryURL        string
	BaseURL             string
	HoldingsURLTemplate string
	ColumnMappings      ColumnMapper
	MonthTranslations   map[string]string
	DateFormats         []string
	DateHeaderPatterns  []string
}

type ColumnMapper struct {
	Name          []string
	Ticker        []string
	ISIN          []string
	MarketValue   []string
	Weight        []string
	Quantity      []string
	ParValue      []string // For bonds
	NotionalValue []string
	Price         []string
	Sector        []string
	AssetClass    []string
	Location      []string
	Exchange      []string
	Currency      []string
	FXRate        []string
	Type          []string
}

var regionConfigs = map[string]regionConfig{
	"us": {
		DiscoveryURL:        "https://www.ishares.com/us/product-screener/product-screener-v3.1.jsn?dcrPath=/templatedata/config/product-screener-v3/data/en/us-ishares/ishares-product-screener-backend-config&siteEntryPassthrough=true",
		BaseURL:             "https://www.ishares.com",
		HoldingsURLTemplate: "%s%s/1467271812596.ajax?fileType=csv",
		MonthTranslations:   nil,
		DateFormats:         []string{"Jan 2, 2006"},
		DateHeaderPatterns:  []string{"Fund Holdings as of"},
		ColumnMappings: ColumnMapper{
			Name:          []string{"Name"},
			Ticker:        []string{"Ticker"},
			ISIN:          []string{"ISIN"},
			MarketValue:   []string{"Market Value"},
			Weight:        []string{"Weight (%)", "Market Weight"},
			NotionalValue: []string{"Notional Value", "Notional Weight"},
			Quantity:      []string{"Quantity"},
			Price:         []string{"Price"},
			Sector:        []string{"Sector"},
			AssetClass:    []string{"Asset Class"},
			Location:      []string{"Location"},
			Exchange:      []string{"Exchange"},
			Currency:      []string{"Currency", "Market Currency"},
			FXRate:        []string{"FX Rate"},
			Type:          []string{"Type"},
		},
	},
	"de": {
		DiscoveryURL:        "https://www.ishares.com/de/privatanleger/de/product-screener/product-screener-v3.1.jsn?dcrPath=/templatedata/config/product-screener-v3/data/de/germany/product-screener/ishares-product-screener-backend-config&siteEntryPassthrough=true",
		BaseURL:             "https://www.ishares.com",
		HoldingsURLTemplate: "%s%s/fund/1478358465952.ajax?fileType=csv",
		MonthTranslations: map[string]string{
			"Jan": "Jan",
			"Feb": "Feb",
			"Mär": "Mar",
			"Apr": "Apr",
			"Mai": "May",
			"Jun": "Jun",
			"Jul": "Jul",
			"Aug": "Aug",
			"Sep": "Sep",
			"Okt": "Oct",
			"Nov": "Nov",
			"Dez": "Dec",
		},
		DateFormats:        []string{"02.Jan.2006"},
		DateHeaderPatterns: []string{"Fondsposition per", "Fondsbestände am"},
		ColumnMappings: ColumnMapper{
			Name:          []string{"Name"},
			Ticker:        []string{"Emittententicker"},
			ISIN:          []string{"ISIN"},
			MarketValue:   []string{"Marktwert"},
			Weight:        []string{"Gewichtung (%)"},
			NotionalValue: []string{"Nominalwert"},
			Quantity:      []string{"Nominale"},
			Price:         []string{"Kurs"},
			Sector:        []string{"Sektor"},
			AssetClass:    []string{"Anlageklasse"},
			Location:      []string{"Standort"},
			Exchange:      []string{"Börse"},
			Currency:      []string{"Währung", "Marktwährung"},
			FXRate:        []string{"Wechselkurs"},
		},
	},
}
