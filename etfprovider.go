package etfscraper

import (
	"context"
)

// Provider defines the interface for an ETF data provider.
type Provider interface {
	// DiscoverETFs discovers all ETFs from the provider.
	DiscoverETFs(ctx context.Context) ([]Fund, error)

	// FundInfo retrieves detailed information about a specific fund.
	FundInfo(ctx context.Context, identifier string) (*Fund, error)

	// Holdings retrieves the holdings of a specific fund.
	Holdings(ctx context.Context, identifier string) (*HoldingsSnapshot, error)
}
