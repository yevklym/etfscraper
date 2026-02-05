package etfscraper

import "errors"

// ErrHoldingsUnavailable is returned when a provider cannot supply holdings data.
var ErrHoldingsUnavailable = errors.New("holdings unavailable")
