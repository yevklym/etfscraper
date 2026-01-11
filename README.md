# ETF Scraper

`etfscraper` is a Go library and CLI tool designed to discover and extract detailed data from Exchange-Traded Fund (ETF) providers. 

Currently, it supports the **iShares** (US) provider, offering capabilities to scrape fund metadata and granular holdings information. Other regions and providers are planned for future releases.

## Features

*   **Fund Discovery**: Automatically discover all available ETFs from a provider.
*   **Detailed Metadata**: Extract key fund information including:
    *   Ticker, Name, and ISIN
    *   Total Assets (AUM)
    *   Expense Ratio
    *   Inception Date
*   **Deep Holdings Analysis**: Download and parse full holdings for specific funds.
    *   Extracts Ticker, Name, Sector, Asset Class, Weight, and Market Value.

## Installation

### As a Library
To use `etfscraper` in your own Go project:

```bash
go get github.com/yevklym/etfscraper
```

### As a CLI Tool
In development.

To run the included command-line interface directly:

## Usage

### Run basic example
```bash
git clone https://github.com/yevklym/etfscraper.git
cd etfscraper
go run cmd/cli/main.go
```

### Library Example

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/yevklym/etfscraper/internal/providers/ishares"
)

func main() {
    // 1. Initialize the provider (pass nil to use default HTTP client)
    provider := ishares.New("us")

    // 2. Get specific Fund Information
    fund, err := provider.FundInfo(context.Background(), "IVV")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Found Fund: %s (%s)\n", fund.Name, fund.Ticker)

    // 3. Get Full Holdings
    holdings, err := provider.Holdings(context.Background(), "IVV")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Total Holdings: %d\n", holdings.TotalHoldings)
    for _, h := range holdings.Holdings[:5] {
        fmt.Printf("- %s: %.2f%%\n", h.Name, h.Weight*100)
    }
}
```

## Architecture

The project follows the following layout:

*   **`etfprovider.go`**: Defines the core `Provider` interface.
*   **`fund.go` & `holding.go`**: Domain models representing Funds and Holdings.
*   **`internal/providers/`**: Contains concrete implementations of the `Provider` interface.
    *   **`ishares/`**: 
        *   `discovery.go`: Handles fetching the list of all ETFs.
        *   `holdings.go`: Handles downloading and parsing CSV holdings files.
        *   `client.go`: Main entry point for the iShares provider.

## Testing
The project is test-driven.

To run all tests:

```bash
go test -v ./...
```