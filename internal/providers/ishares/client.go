// Package ishares provides a client for fetching iShares ETF data.
//
// The client supports multiple regions (US, DE) and allows configuration
// through functional options.
//
// Example usage:
//
//	client, err := ishares.New("de")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	funds, err := client.DiscoverETFs(context.Background())
//
// With custom configuration:
//
//	client, err := ishares.New("de",
//	    ishares.WithTimeout(30*time.Second),
//	    ishares.WithDebug(true),
//	)
package ishares

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/yevklym/etfscraper"
)

type Client struct {
	httpConfig etfscraper.HTTPConfig
	region     string
	config     regionConfig

	mu       sync.Mutex
	cache    []etfscraper.Fund
	index    map[string]*etfscraper.Fund // lowercase ticker -> fund
	cachedAt time.Time
	cacheTTL time.Duration
}

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

	funds, err := c.fetchAndDecodeFunds(ctx)
	if err != nil {
		return nil, err
	}

	c.cache = funds
	c.cachedAt = time.Now()
	c.index = buildFundIndex(funds)
	return funds, nil
}

// buildFundIndex creates a lookup map keyed by lowercase ticker and ISIN.
func buildFundIndex(funds []etfscraper.Fund) map[string]*etfscraper.Fund {
	idx := make(map[string]*etfscraper.Fund, len(funds)*2)
	for i := range funds {
		f := &funds[i]
		if f.Ticker != "" {
			idx[strings.ToLower(f.Ticker)] = f
		}
		if f.ISIN != "" {
			idx[strings.ToLower(f.ISIN)] = f
		}
	}
	return idx
}

// FundInfo retrieves detailed information about a specific fund by ticker or ISIN.
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
