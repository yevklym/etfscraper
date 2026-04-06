package etfscraper

// Currency represents an ISO 4217 currency code.
type Currency string

// Exchange represents a stock exchange.
type Exchange string

// AssetClass represents the type of financial asset.
type AssetClass string

// Sector represents the economic sector using GICS classification.
type Sector string

// Location represents a country or region. Providers normalize localized
// location names to these English strings for major economies. Providers
// may return values not in this list.
type Location string

// ProviderName identifies an ETF data provider.
type ProviderName string

// Common currency constants (ISO 4217). Providers may return values not in
// this list; use these constants for comparison rather than raw strings.
const (
	CurrencyUSD Currency = "USD" // US Dollar
	CurrencyEUR Currency = "EUR" // Euro
	CurrencyGBP Currency = "GBP" // British Pound
	CurrencyJPY Currency = "JPY" // Japanese Yen
	CurrencyCAD Currency = "CAD" // Canadian Dollar
	CurrencyAUD Currency = "AUD" // Australian Dollar
	CurrencyCHF Currency = "CHF" // Swiss Franc
	CurrencyCNY Currency = "CNY" // Chinese Yuan
	CurrencyINR Currency = "INR" // Indian Rupee
	CurrencyBRL Currency = "BRL" // Brazilian Real
)

// Common asset class constants. Providers normalize localized values to
// these English strings. Providers may return values not in this list.
const (
	AssetClassEquity         AssetClass = "Equity"
	AssetClassBond           AssetClass = "Bond"
	AssetClassCash           AssetClass = "Cash"
	AssetClassCommodity      AssetClass = "Commodity"
	AssetClassRealEstate     AssetClass = "Real Estate"
	AssetClassAlternative    AssetClass = "Alternative"
	AssetClassCryptocurrency AssetClass = "Cryptocurrency"
)

// Common exchange constants. Providers may return values not in this list.
const (
	ExchangeNYSE     Exchange = "NYSE"     // New York Stock Exchange
	ExchangeNASDAQ   Exchange = "NASDAQ"   // NASDAQ
	ExchangeAMEX     Exchange = "AMEX"     // American Stock Exchange
	ExchangeBATS     Exchange = "BATS"     // BATS Global Markets
	ExchangeLSE      Exchange = "LSE"      // London Stock Exchange
	ExchangeEuronext Exchange = "Euronext" // Euronext
	ExchangeXetra    Exchange = "Xetra"    // Xetra
	ExchangeTSE      Exchange = "TSE"      // Tokyo Stock Exchange
	ExchangeHKEX     Exchange = "HKEX"     // Hong Kong Exchange
	ExchangeSSE      Exchange = "SSE"      // Shanghai Stock Exchange
	ExchangeSZSE     Exchange = "SZSE"     // Shenzhen Stock Exchange
	ExchangeTSX      Exchange = "TSX"      // Toronto Stock Exchange
	ExchangeASX      Exchange = "ASX"      // Australian Securities Exchange
)

// Common GICS sector constants. Providers normalize localized sector names
// to these English strings. Providers may return values not in this list.
const (
	SectorEnergy                Sector = "Energy"
	SectorMaterials             Sector = "Materials"
	SectorIndustrials           Sector = "Industrials"
	SectorConsumerDiscretionary Sector = "Consumer Discretionary"
	SectorConsumerStaples       Sector = "Consumer Staples"
	SectorHealthcare            Sector = "Healthcare"
	SectorFinancials            Sector = "Financials"
	SectorInformationTechnology Sector = "Information Technology"
	SectorTelecommunication     Sector = "Telecommunication Services"
	SectorUtilities             Sector = "Utilities"
	SectorRealEstate            Sector = "Real Estate"
)

// Common Location constants. Providers normalize localized location names
// to these English strings.
const (
	LocationUnitedStates  Location = "United States"
	LocationUnitedKingdom Location = "United Kingdom"
	LocationJapan         Location = "Japan"
	LocationGermany       Location = "Germany"
	LocationFrance        Location = "France"
	LocationSwitzerland   Location = "Switzerland"
	LocationCanada        Location = "Canada"
	LocationAustralia     Location = "Australia"
	LocationChina         Location = "China"
	LocationTaiwan        Location = "Taiwan"
	LocationSouthKorea    Location = "South Korea"
	LocationIndia         Location = "India"
	LocationBrazil        Location = "Brazil"
	LocationNetherlands   Location = "Netherlands"
	LocationSweden        Location = "Sweden"
	LocationItaly         Location = "Italy"
	LocationSpain         Location = "Spain"
	LocationIreland       Location = "Ireland"
	LocationDenmark       Location = "Denmark"
	LocationFinland       Location = "Finland"
	LocationTurkey        Location = "Turkey"
	LocationBelgium       Location = "Belgium"
	LocationAustria       Location = "Austria"
	LocationLuxembourg    Location = "Luxembourg"
	LocationSingapore     Location = "Singapore"
	LocationNorway        Location = "Norway"
	LocationIsrael        Location = "Israel"
	LocationCash          Location = "Cash and/or Derivatives"
	LocationEurope        Location = "Europe"
	LocationAsia          Location = "Asia"
	LocationGlobal        Location = "Global"
)

// ProviderName constants for supported providers.
const (
	ProviderIShares   ProviderName = "iShares"
	ProviderAmundi    ProviderName = "Amundi"
	ProviderXtrackers ProviderName = "Xtrackers"
)
