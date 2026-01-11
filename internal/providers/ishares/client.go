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

	"github.com/yevklym/etfscraper"
)

type Client struct {
	httpConfig etfscraper.HTTPConfig
	region     string
	config     regionConfig
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
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

func (c *Client) DiscoverETFs(ctx context.Context) ([]etfscraper.Fund, error) {
	return c.fetchAndDecodeFunds(ctx)
}

// FundInfo retrieves detailed information about a specific fund by ticker
func (c *Client) FundInfo(ctx context.Context, identifier string) (*etfscraper.Fund, error) {
	identifier = strings.TrimSpace(identifier)
	if identifier == "" {
		return nil, fmt.Errorf("identifier cannot be empty")
	}

	funds, err := c.fetchAndDecodeFunds(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch funds: %w", err)
	}

	for _, fund := range funds {
		if strings.EqualFold(fund.Ticker, identifier) {
			return &fund, nil
		}
	}

	return nil, fmt.Errorf("fund not found with identifier: %s", identifier)
}
