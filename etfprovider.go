// Package etfscraper provides types and interfaces for discovering ETFs
// and fetching fund metadata and holdings from providers like iShares and Amundi.
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

	// Holdings retrieves the holdings of a specific fund by ticker or ISIN.
	// This method performs an internal fund lookup on each call.
	// Use HoldingsForFund to skip the lookup when you already have a Fund.
	Holdings(ctx context.Context, identifier string) (*HoldingsSnapshot, error)

	// HoldingsForFund retrieves the holdings using a previously fetched Fund.
	// This avoids the internal discovery lookup, making it the preferred method
	// when you already have a Fund from DiscoverETFs or FundInfo.
	HoldingsForFund(ctx context.Context, fund *Fund) (*HoldingsSnapshot, error)
}
