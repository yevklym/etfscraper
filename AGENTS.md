# AGENTS.md

Guidance for AI coding agents operating in this repository.
Also see `.github/copilot-instructions.md` and `.github/go.instructions.md` for additional context.

## Project Overview

`etfscraper` is a Go 1.26 library (module `github.com/yevklym/etfscraper`) that discovers ETFs and fetches fund metadata and holdings for providers like iShares and Amundi. It has **zero external dependencies** -- stdlib only.

## Build / Test / Lint Commands

```bash
# Build everything
go build -v ./...

# Run all tests (matches CI)
go test -v -race ./...

# Run tests in a single package
go test -v -race ./internal/providers/ishares/

# Run a single test by name (regex match)
go test -v -race ./internal/providers/ishares/ -run TestParseHoldings_FrenchFormat

# Run tests matching a pattern across all packages
go test -v -race ./... -run TestColumnResolver

# Lint (uses default config, no .golangci.yml)
golangci-lint run ./...

# Clear test cache
go clean -testcache

# Run CLI example (makes real network requests)
go run ./cmd/example
```

CI (`.github/workflows/ci.yml`) runs `golangci-lint` then `go build -v ./...` and `go test -v -race ./...` on Go 1.26.

## Project Layout

```
etfprovider.go, fund.go, holding.go  Root package: Provider interface, domain types
enums.go                              Custom string types: Currency, Exchange, Sector, etc.
errors.go                             Sentinel errors (ErrHoldingsUnavailable)
config.go                             HTTPClient interface, HTTPConfig struct
providers/                            Public factory: Open(), OpenSpec(), options
internal/providers/ishares/           iShares provider (CSV holdings, JSON discovery)
internal/providers/amundi/            Amundi provider (JSON API)
internal/testutil/                    Shared MockHTTPClient
cmd/example/                          CLI demo app
```

Each provider package follows a consistent file structure:
- `client.go` -- Client struct, `New()` constructor, interface methods
- `config.go` -- Region-specific configuration (URLs, column mappings, defaults)
- `discovery.go` -- ETF discovery logic, JSON parsing
- `holdings.go` -- Holdings fetching and parsing
- `options.go` -- Functional options (`WithTimeout`, `WithDebug`, etc.)
- `*_test.go` -- Co-located tests (same package for white-box testing)
- `data/` -- Embedded test fixture JSON files

## Code Style

### Imports

Two groups separated by a blank line: stdlib first, then module-internal imports.

```go
import (
    "context"
    "fmt"
    "net/http"

    "github.com/yevklym/etfscraper"
)
```

No external dependencies exist. If added, use three groups: stdlib, external, internal.

### Naming

- **No stuttering**: package `ishares` has type `Client`, not `ISharesClient`
- **Acronyms uppercase**: `ISIN`, `URL`, `HTTP` (e.g., `ProductPageURL`, `HTTPConfig`)
- **Receivers**: single letter (`c` for Client, `r` for resolver)
- **Exported**: PascalCase (`Fund`, `CurrencyUSD`, `WithTimeout`)
- **Unexported**: camelCase (`regionConfig`, `normalizeCurrency`, `pickTicker`)
- **Test vars**: `tt` for table entries, `got`/`want`/`expected` for assertions

### Types

- Custom string types for domain enums: `type Currency string`, `type Exchange string`
- Constants use PascalCase: `CurrencyUSD`, `ExchangeNYSE`, `ProviderIShares`
- Interface-based HTTP client: `HTTPClient` interface with `Do(*http.Request)` for test mocking
- `Fund.ProviderMetadata` is `any`, cast with type assertions in provider code
- JSON tags: `json:"camelCase,omitempty"` or `json:"field,omitzero"`

### Functions and Methods

- `context.Context` is always the first parameter
- Constructors: `New(region string, opts ...ClientOption) (*Client, error)`
- Functional options: `type ClientOption func(*Client)` with `With*` prefix
- Returns: `(result, error)` -- error is always last

### Error Handling

- Wrap with `%w` when propagating: `fmt.Errorf("failed to fetch funds: %w", err)`
- Plain `fmt.Errorf` for leaf errors: `fmt.Errorf("unsupported region '%s'", region)`
- Sentinel errors: `var ErrHoldingsUnavailable = errors.New("holdings unavailable")`
- Messages: lowercase, no trailing punctuation
- Check errors immediately; don't log and return (choose one)
- Response body close errors: log in deferred func, don't return

### Comments

- Godoc style: start with the symbol name (`// Fund holds comprehensive metadata...`)
- Package comments: `// Package ishares provides a client for...`
- Explain "why", not "what", unless the logic is complex
- No emoji in code or comments

## Testing Conventions

- **Standard `testing` package only** -- no testify or external test libraries
- **Table-driven tests** with `t.Run` subtests are the default pattern
- **Naming**: `TestFunctionName` or `TestFunctionName_Scenario` (e.g., `TestParseHoldings_FrenchFormat`)
- **Subtest names**: lowercase descriptive strings (e.g., `"fund found"`, `"European format"`)
- **Float comparison**: use epsilon tolerance, not exact equality
- **HTTP mocking**: use `testutil.MockHTTPClient` for single responses
- **Test data**: inline CSV/JSON strings for small data, `//go:embed data/file.json` for large fixtures
- **White-box**: test files use the same package name to access unexported symbols
- **Black-box**: `example_test.go` uses `package ishares_test` for godoc examples
- **Assertions**: plain `t.Errorf`/`t.Fatalf`, no assertion helpers

```go
// Standard table-driven test pattern
func TestNormalize_Scenarios(t *testing.T) {
    tests := []struct {
        name  string
        input string
        want  string
    }{
        {name: "simple case", input: "foo", want: "bar"},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := normalize(tt.input)
            if got != tt.want {
                t.Errorf("normalize() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

## Key Architecture Decisions

- Public API is in `providers/` package; concrete implementations are in `internal/`
- Factory wiring in `providers/open.go` is registry-backed (supported regions + constructor dispatch share one source of truth)
- The `Provider` interface has four methods: `DiscoverETFs`, `FundInfo`, `Holdings`, `HoldingsForFund`
- Provider spec format: `"ishares:us"`, `"amundi:fr"` (parsed by `providers.Open` / `ParseProviderSpec`; `providers.OpenSpec` accepts structured `Spec`)
- Each region has its own config with locale-specific column names, date formats, and month translations
- Number normalization handles US (`1,234.56`), European (`1.234,56`), Swiss (`1'234.56`), and French space-separated (`10 569 831,35`) formats
