# Copilot Instructions for etfscraper

## Summary
- `etfscraper` is a Go library that discovers ETFs and fetches fund metadata and holdings for providers like iShares and Amundi.
- Small Go module targeting Go 1.26 with a CLI example under `cmd/example`.

## Repo Facts
- Language/runtime: Go 1.26 (module `github.com/yevklym/etfscraper`).
- No Makefile or shell scripts; build/test via `go` tool.
- Lint in CI uses `golangci-lint` with default configuration (no `.golangci.*` file).

## Build, Test, Lint, Run (validated commands)
Validated on macOS (darwin/arm64) with Go 1.26.1.

### Bootstrap
- Always install Go 1.26 (matches `go.mod` and CI).
- `golangci-lint` is required for local lint parity with CI.

### Build
```bash
cd /Users/yevheniiklymenko/GolandProjects/etfscraper
go build -v ./...
```
- Works with no additional setup.

### Test
```bash
cd /Users/yevheniiklymenko/GolandProjects/etfscraper
go test -v -race ./...
```
- Matches CI; passes locally (tests finish in a few seconds).
- `go test -v ./...` also works but CI uses `-race`.

### Lint
```bash
cd /Users/yevheniiklymenko/GolandProjects/etfscraper
golangci-lint run ./...
```
- Uses default golangci-lint configuration (no custom config file). Repository-wide lint may surface pre-existing staticcheck findings in internal provider discovery tests.

### Clean
```bash
cd /Users/yevheniiklymenko/GolandProjects/etfscraper
go clean -testcache
```
- Safe to clear test cache before re-running tests.

### Run (CLI example)
```bash
cd /Users/yevheniiklymenko/GolandProjects/etfscraper
go run ./cmd/example
```
- Not validated here; it makes network requests to ETF providers, so expect external I/O and potential rate limits.

### Notes on order and pitfalls
- Running `go build` or `go test` is sufficient; both auto-download module deps.
- No repo-specific environment variables were required in local validation.

## CI and Validation
GitHub Actions workflow: `.github/workflows/ci.yml`
- Lint job: `golangci-lint` (Go 1.26).
- Test job: `go build -v ./...` then `go test -v -race ./...`.
- Always keep local checks aligned with these commands.

## Project Layout and Architecture
- Core interface: `etfprovider.go` defines `Provider` with `DiscoverETFs`, `FundInfo`, `Holdings`, and `HoldingsForFund`.
- Domain models: `fund.go`, `holding.go`, `enums.go`, `errors.go`, `config.go`.
- Provider factory (public API): `providers/open.go`, `providers/options.go`.
- Concrete providers: `internal/providers/amundi` and `internal/providers/ishares`.
- Test utilities: `internal/testutil`.
- CLI example: `cmd/example/main.go`.
- Go style guidance: `.github/go.instructions.md`.

## Root Contents
- `.agent/`, `.git/`, `.github/`, `.idea/`, `.DS_Store`
- `.gitignore`, `LICENSE`, `README.md`, `go.mod`, `go.sum`
- `cmd/`, `internal/`, `providers/`
- `config.go`, `enums.go`, `errors.go`, `etfprovider.go`, `fund.go`, `holding.go`

## Next-Level Directories
- `cmd/example/main.go`
- `internal/providers/amundi/*` (client, discovery, holdings, options, tests, data fixtures)
- `internal/providers/ishares/*` (client, discovery, holdings, column resolver, options, tests, data fixtures)
- `internal/testutil/httpmock.go`
- `providers/open.go`, `providers/options.go`, `providers/open_test.go`

## README Highlights (summary)
- Usage examples for `providers.Open` and `OpenSpec`.
- Supported regions: iShares (US/DE/UK), Amundi (DE/UK/FR).
- Tests: `go test -v -race ./...`.

## Key File Snippets
- `cmd/example/main.go` (entrypoint) opens a provider and demonstrates `FundInfo`, `Holdings`, and `DiscoverETFs`.
- `providers/open.go` parses `provider:region` specs and uses a registry-backed factory for validation and constructor dispatch.

## Trust These Instructions
- Follow this document first and only search the repo when details are missing or appear incorrect.
- Use Context7 when generating code that involves third-party packages,
  Go standard library features added after Go 1.20.
  Skip it for general logic, algorithms, or well-established patterns.
