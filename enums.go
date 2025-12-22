package etfscraper

// Currency represents a currency code
type Currency string

// Exchange represents a stock exchange
type Exchange string

// AssetClass represents the type of financial asset
type AssetClass string

// Sector represents the economic sector (GICS)
type Sector string

// Location represents a country or region
type Location string

type ProviderName string

// Currency constants
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

// AssetClass constants
const (
	AssetClassEquity         AssetClass = "Equity"
	AssetClassBond           AssetClass = "Bond"
	AssetClassCash           AssetClass = "Cash"
	AssetClassCommodity      AssetClass = "Commodity"
	AssetClassRealEstate     AssetClass = "Real Estate"
	AssetClassAlternative    AssetClass = "Alternative"
	AssetClassCryptocurrency AssetClass = "Cryptocurrency"
)

// Exchange constants - major stock exchanges worldwide
const (
	ExchangeNYSE     Exchange = "NYSE"     // New York Stock Exchange
	ExchangeNASDAQ   Exchange = "NASDAQ"   // NASDAQ
	ExchangeAMEX     Exchange = "AMEX"     // American Stock Exchange
	ExchangeBATS     Exchange = "BATS"     // BATS Global Markets
	ExchangeLSE      Exchange = "LSE"      // London Stock Exchange
	ExchangeEuronext Exchange = "Euronext" // Euronext
	ExchangeTSE      Exchange = "TSE"      // Tokyo Stock Exchange
	ExchangeHKEX     Exchange = "HKEX"     // Hong Kong Exchange
	ExchangeSSE      Exchange = "SSE"      // Shanghai Stock Exchange
	ExchangeSZSE     Exchange = "SZSE"     // Shenzhen Stock Exchange
	ExchangeTSX      Exchange = "TSX"      // Toronto Stock Exchange
	ExchangeASX      Exchange = "ASX"      // Australian Securities Exchange
)

// Sector constants - GICS sector classification
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

// Location constants
const (
	LocationUSA         Location = "USA"
	LocationCanada      Location = "Canada"
	LocationUK          Location = "United Kingdom"
	LocationGermany     Location = "Germany"
	LocationFrance      Location = "France"
	LocationItaly       Location = "Italy"
	LocationSpain       Location = "Spain"
	LocationNetherlands Location = "Netherlands"
	LocationSwitzerland Location = "Switzerland"
	LocationJapan       Location = "Japan"
	LocationChina       Location = "China"
	LocationHongKong    Location = "Hong Kong"
	LocationSouthKorea  Location = "South Korea"
	LocationIndia       Location = "India"
	LocationAustralia   Location = "Australia"
	LocationBrazil      Location = "Brazil"
	LocationMexico      Location = "Mexico"
	LocationSouthAfrica Location = "South Africa"
	LocationEmerging    Location = "Emerging Markets"
	LocationDeveloped   Location = "Developed Markets"
	LocationGlobal      Location = "Global"
)

const (
	ProviderIShares ProviderName = "iShares"
)
