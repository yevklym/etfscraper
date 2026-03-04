package etfscraper

import "errors"

// ErrHoldingsUnavailable is returned by Holdings and HoldingsForFund when a
// provider cannot supply holdings data for a given fund. Use errors.Is to
// check for this sentinel:
//
//	if errors.Is(err, etfscraper.ErrHoldingsUnavailable) { ... }
var ErrHoldingsUnavailable = errors.New("holdings unavailable")
