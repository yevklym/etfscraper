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
	ExchangeXetra    Exchange = "Xetra"    // Xetra
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

const (
	ProviderIShares ProviderName = "iShares"
	ProviderAmundi  ProviderName = "Amundi"
)
