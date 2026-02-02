package amundi

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

func (c *Client) FundInfo(ctx context.Context, identifier string) (*etfscraper.Fund, error) {
	identifier = strings.TrimSpace(identifier)
	if identifier == "" {
		return nil, fmt.Errorf("identifier cannot be empty")
	}

	funds, err := c.DiscoverETFs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch funds: %w", err)
	}

	for _, fund := range funds {
		if strings.EqualFold(fund.Ticker, identifier) || strings.EqualFold(fund.ISIN, identifier) {
			return &fund, nil
		}
	}

	return nil, fmt.Errorf("fund not found with identifier: %s", identifier)
}

func (c *Client) Holdings(_ context.Context, _ string) (*etfscraper.HoldingsSnapshot, error) {
	return nil, fmt.Errorf("%w: amundi", etfscraper.ErrHoldingsUnavailable)
}
