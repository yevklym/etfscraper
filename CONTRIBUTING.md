# Contributing

Contributions are welcome. Please read this guide before opening a pull request.

## Prerequisites

- Go 1.26.1 or later
- External dependencies should be kept to an absolute minimum (e.g., `go-rod` is used only to bypass Akamai bot protection). Heavy reliance on the Go standard library is strongly preferred.

## Getting started

```bash
git clone https://github.com/yevklym/etfscraper
cd etfscraper
go build ./...
go test -v -race ./...
```

## Pull requests

- Target the main branch
- Include tests for any new or changed behaviour
- Keep changes focused — one concern per PR

## Adding a provider
New providers go under `internal/providers/<name>/` and must implement the `etfscraper.Provider` interface. A new provider
should also be registered in the factory in `providers/open.go` and listed in `providers/open_test.go`.

## Scope
To keep the library focused, the following are unlikely to be accepted:
- Provider additions without tests and fixture data
- Heavy or unnecessary external dependencies