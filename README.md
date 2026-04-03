# ETF Scraper

`etfscraper` is a Go library for discovering ETFs and fetching fund metadata and holdings from diefferent exchange-traded fund and index mutual fund providers.
## Design

The library has two public packages with different responsibilities:

- `etfscraper` contains the core interfaces, domain types, errors, and shared configuration/logging primitives (for example `Provider`, `Fund`, `Holding`, `ErrHoldingsUnavailable`, `Logger`, `NopLogger()`).
- `providers` contains the factory layer for creating concrete provider instances (for example `providers.Open(...)` and `providers.OpenSpec(...)`).

In typical usage, import both packages: use `providers` to construct a provider, then work with shared types, interfaces, and helpers from `etfscraper`.

## Supported Providers

| Provider  | Regions        |
|-----------|----------------|
| iShares   | us, de, uk, fr |
| Amundi    | de, uk, fr     |
| Xtrackers | de, uk         |

## Installation

```bash
go get github.com/yevklym/etfscraper
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/yevklym/etfscraper/providers"
)

func main() {
	provider, err := providers.Open("ishares:us")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// Discover all ETFs
	funds, err := provider.DiscoverETFs(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d ETFs\n", len(funds))

	// Get holdings for a specific fund
	snapshot, err := provider.Holdings(ctx, "IVV")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Holdings: %d (as of %s)\n", snapshot.TotalHoldings, snapshot.AsOfDate.Format("2006-01-02"))
	for _, h := range snapshot.Holdings[:5] {
		fmt.Printf("  %s: %.2f%%\n", h.Name, h.Weight*100)
	}
}
```

### Bulk Holdings (avoiding N+1)

When fetching holdings for multiple funds, use `HoldingsForFund` to skip the internal discovery lookup on each call:

```go
funds, _ := provider.DiscoverETFs(ctx)
for i := range funds[:5] {
	snapshot, err := provider.HoldingsForFund(ctx, &funds[i])
	if err != nil {
		log.Printf("skipping %s: %v", funds[i].Ticker, err)
		continue
	}
	fmt.Printf("%s: %d holdings\n", funds[i].Ticker, snapshot.TotalHoldings)
}
```

## Configuration

```go
// Custom timeout
provider, _ := providers.Open("amundi:de", providers.WithTimeout(30*time.Second))

// Custom HTTP client
provider, _ := providers.Open("ishares:uk", providers.WithHTTPClient(myClient))

// Debug logging
provider, _ := providers.Open("ishares:us", providers.WithDebug(true))

// Custom logger (or silence with etfscraper.NopLogger())
provider, _ := providers.Open("ishares:us", providers.WithLogger(etfscraper.NopLogger()))

// Discovery cache TTL (default 5m, set 0 to disable)
provider, _ := providers.Open("ishares:us", providers.WithCacheTTL(10*time.Minute))

// Typed spec
provider, _ := providers.OpenSpec(providers.Spec{Name: "ishares", Region: "uk"})
```

List supported providers at runtime:

```go
for _, p := range providers.SupportedProviders() {
	fmt.Printf("%s: %v\n", p.Name, p.Regions)
}
```

## Provider Interface

All providers implement the `etfscraper.Provider` interface:

```go
type Provider interface {
	DiscoverETFs(ctx context.Context) ([]Fund, error)
	FundInfo(ctx context.Context, identifier string) (*Fund, error)
	Holdings(ctx context.Context, identifier string) (*HoldingsSnapshot, error)
	HoldingsForFund(ctx context.Context, fund *Fund) (*HoldingsSnapshot, error)
}
```

- `identifier` accepts a fund ticker or ISIN.
- `Holdings` internally calls `FundInfo` to resolve the identifier; use `HoldingsForFund` when you already have the `Fund` to avoid repeated lookups.
- `ErrHoldingsUnavailable` is returned when a provider cannot supply holdings for a given fund. Use `errors.Is(err, etfscraper.ErrHoldingsUnavailable)` to check.

## Architecture

```
etfprovider.go                        Provider interface
fund.go, holding.go                   Fund, Holding, HoldingsSnapshot structs
enums.go                              Currency, Exchange, Sector, AssetClass types
errors.go                             Sentinel errors (ErrHoldingsUnavailable)
config.go                             HTTPClient, HTTPConfig, Logger interfaces
providers/                            Public factory: Open(), OpenSpec(), options
internal/providers/ishares/           iShares provider (CSV holdings, JSON discovery)
internal/providers/amundi/            Amundi provider (JSON API)
internal/testutil/                    Shared HTTP mock for tests
cmd/example/                          Runnable usage example
```

## Testing

```bash
go test -v -race ./...
```
