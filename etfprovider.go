// Package etfscraper provides types and interfaces for discovering ETFs
// and fetching fund metadata and holdings from providers like iShares and Amundi.
package etfscraper

import (
	"context"
)

// Provider defines the interface for an ETF data provider.
type Provider interface {
	// DiscoverETFs returns all ETFs available from the provider.
	// Results are cached for the duration configured by WithCacheTTL
	// (default 5 minutes). Returns an error on network or parsing failures.
	DiscoverETFs(ctx context.Context) ([]Fund, error)

	// FundInfo retrieves detailed information about a specific fund.
	// The identifier can be a fund ticker (e.g. "IVV") or ISIN
	// (e.g. "IE00B5BMR087"). Returns an error if the fund is not found.
	FundInfo(ctx context.Context, identifier string) (*Fund, error)

	// Holdings retrieves the holdings of a specific fund by ticker or ISIN.
	// This method performs an internal fund lookup (via FundInfo) on each call.
	// Use HoldingsForFund to skip the lookup when you already have a Fund.
	// Returns ErrHoldingsUnavailable if the provider cannot supply holdings
	// for the given fund.
	Holdings(ctx context.Context, identifier string) (*HoldingsSnapshot, error)

	// HoldingsForFund retrieves the holdings using a previously fetched Fund.
	// This avoids the internal discovery lookup, making it the preferred method
	// when iterating over funds from DiscoverETFs. The fund must not be nil.
	// Returns ErrHoldingsUnavailable if the provider cannot supply holdings.
	HoldingsForFund(ctx context.Context, fund *Fund) (*HoldingsSnapshot, error)
}
