// Package xtrackers provides a client for fetching Xtrackers (DWS) ETF data.
package xtrackers

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/yevklym/etfscraper"
)

// Client is the Xtrackers ETF data provider.
type Client struct {
	httpConfig etfscraper.HTTPConfig
	region     string
	config     regionConfig

	mu       sync.Mutex
	cache    []etfscraper.Fund
	index    map[string]*etfscraper.Fund // lowercase ISIN -> fund
	cachedAt time.Time
	cacheTTL time.Duration

	// skipBrowserFetch is an internal flag used by unit tests to avoid
	// launching headless Chromium during mocked executions.
	skipBrowserFetch bool
}

// New creates a new Xtrackers client for the given region.
func New(region string, opts ...ClientOption) (*Client, error) {
	config, ok := regionConfigs[strings.ToLower(region)]
	if !ok {
		return nil, fmt.Errorf("unsupported region '%s'", region)
	}

	c := &Client{
		httpConfig: etfscraper.DefaultHTTPConfig(),
		region:     strings.ToLower(region),
		config:     config,
		cacheTTL:   5 * time.Minute,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

// DiscoverETFs returns all Xtrackers ETFs available from the provider.
func (c *Client) DiscoverETFs(ctx context.Context) ([]etfscraper.Fund, error) {
	funds, err := c.discoverCached(ctx)
	if err != nil {
		return nil, err
	}
	// Return a copy so callers cannot mutate the cache.
	out := make([]etfscraper.Fund, len(funds))
	copy(out, funds)
	return out, nil
}

// discoverCached returns the cached discovery result or fetches fresh data.
func (c *Client) discoverCached(ctx context.Context) ([]etfscraper.Fund, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cache != nil && time.Since(c.cachedAt) < c.cacheTTL {
		return c.cache, nil
	}

	funds, err := c.discoverFresh(ctx)
	if err != nil {
		return nil, err
	}

	c.cache = funds
	c.cachedAt = time.Now()
	c.index = buildFundIndex(funds)
	return funds, nil
}

// buildFundIndex creates a lookup map keyed by lowercase ISIN.
func buildFundIndex(funds []etfscraper.Fund) map[string]*etfscraper.Fund {
	idx := make(map[string]*etfscraper.Fund, len(funds))
	for i := range funds {
		f := &funds[i]
		if f.ISIN != "" {
			idx[strings.ToLower(f.ISIN)] = f
		}
	}
	return idx
}

// FundInfo retrieves information about a specific fund by ISIN.
func (c *Client) FundInfo(ctx context.Context, identifier string) (*etfscraper.Fund, error) {
	identifier = strings.TrimSpace(identifier)
	if identifier == "" {
		return nil, fmt.Errorf("identifier cannot be empty")
	}

	if _, err := c.discoverCached(ctx); err != nil {
		return nil, fmt.Errorf("failed to fetch funds: %w", err)
	}

	c.mu.Lock()
	fund := c.index[strings.ToLower(identifier)]
	c.mu.Unlock()

	if fund != nil {
		return fund, nil
	}

	return nil, fmt.Errorf("fund not found with identifier: %s", identifier)
}

// Holdings retrieves the holdings of a specific fund by ISIN.
func (c *Client) Holdings(ctx context.Context, identifier string) (*etfscraper.HoldingsSnapshot, error) {
	fund, err := c.FundInfo(ctx, identifier)
	if err != nil {
		return nil, fmt.Errorf("xtrackers: holdings: %w", err)
	}
	return c.HoldingsForFund(ctx, fund)
}

// HoldingsForFund retrieves the holdings using a previously fetched Fund.
func (c *Client) HoldingsForFund(ctx context.Context, fund *etfscraper.Fund) (*etfscraper.HoldingsSnapshot, error) {
	if fund == nil {
		return nil, fmt.Errorf("fund cannot be nil")
	}
	return c.fetchHoldings(ctx, fund)
}
