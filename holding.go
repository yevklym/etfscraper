package etfscraper

import (
	"time"
)

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

// Holding represents an individual holding within an ETF
type Holding struct {
	// Basic identification
	Ticker      string `json:"ticker,omitempty" validate:"required_without=ISIN"`
	ISIN        string `json:"isin,omitempty" validate:"len=12"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description,omitempty"`

	// Financial data
	Weight      float64  `json:"weight" validate:"min=0,max=1"` // decimal: 0.025 = 2.5%, 0.5 = 50%
	Quantity    float64  `json:"quantity,omitempty" validate:"min=0"`
	MarketValue float64  `json:"marketValue,omitempty" validate:"min=0"`
	Price       float64  `json:"price,omitempty" validate:"min=0"`
	Currency    Currency `json:"currency,omitempty"`

	// Classification
	Sector     Sector     `json:"sector,omitempty"`
	AssetClass AssetClass `json:"assetClass,omitempty"`
	Location   Location   `json:"location,omitempty"`
	Exchange   Exchange   `json:"exchange,omitempty"`

	// Timestamps
	LastUpdated time.Time `json:"lastUpdated,omitempty"`
}

// Fund holds comprehensive metadata about an ETF
type Fund struct {
	// Basic identification
	Ticker   string `json:"ticker,omitempty" validate:"required"`
	ISIN     string `json:"isin,omitempty" validate:"len=12"`
	Name     string `json:"name" validate:"required"`
	FullName string `json:"fullName,omitempty"`

	// Provider information
	Provider ProviderName `json:"provider,omitempty"`
	Family   string       `json:"family,omitempty"`

	// Financial details
	Currency     Currency `json:"currency,omitempty"`
	BaseCurrency Currency `json:"baseCurrency,omitempty"`

	// Fund characteristics
	InceptionDate  *time.Time `json:"inceptionDate,omitempty"`
	TotalAssets    float64    `json:"totalAssets,omitempty" validate:"min=0"`
	ExpenseRatio   float64    `json:"expenseRatio,omitempty" validate:"min=0,max=1"`
	DividendYield  float64    `json:"dividendYield,omitempty" validate:"min=0"`
	IsDistributing bool       `json:"isDistributing,omitempty"`

	// Classification
	Category   string     `json:"category,omitempty"`
	Geography  Location   `json:"geography,omitempty"`
	AssetClass AssetClass `json:"assetClass,omitempty"`

	// Trading information
	Exchange      Exchange `json:"exchange,omitempty"`
	TradingSymbol string   `json:"tradingSymbol,omitempty"`

	// Additional metadata
	CUSIP            string      `json:"cusip,omitempty" validate:"len=9"`
	SEDOL            string      `json:"sedol,omitempty" validate:"len=7"`
	Description      string      `json:"description,omitempty"`
	Objective        string      `json:"objective,omitempty"`
	ProductPageURL   string      `json:"productPageUrl,omitempty"`
	LastUpdated      time.Time   `json:"lastUpdated,omitempty"`
	ProviderMetadata interface{} `json:"-"`
}

// HoldingsSnapshot represents a point-in-time snapshot of fund holdings
type HoldingsSnapshot struct {
	Fund          Fund      `json:"fund" validate:"required"`
	AsOfDate      time.Time `json:"asOfDate" validate:"required"`
	Holdings      []Holding `json:"holdings" validate:"dive"`
	TotalHoldings int       `json:"totalHoldings"`
	TopHoldings   int       `json:"topHoldings,omitempty"` // Number of holdings included if truncated

	// Summary statistics
	TotalWeight float64 `json:"totalWeight,omitempty"`
	CashWeight  float64 `json:"cashWeight,omitempty"`

	// Metadata
	LastUpdated time.Time `json:"lastUpdated,omitempty"`
}
