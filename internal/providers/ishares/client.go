package ishares

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/yevklym/etfscraper"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	httpClient HTTPClient
	region     string
}

func New(region string, client HTTPClient) *Client {
	if client == nil {
		client = &http.Client{
			Timeout: time.Second * 15,
		}
	}
	return &Client{
		httpClient: client,
		region:     region,
	}
}

func (c *Client) DiscoverETFs(ctx context.Context) ([]etfscraper.Fund, error) {
	return c.fetchAndDecodeFunds(ctx)
}

// FundInfo retrieves detailed information about a specific fund by ticker
func (c *Client) FundInfo(ctx context.Context, identifier string) (*etfscraper.Fund, error) {
	url, err := c.buildFundURL(identifier)
	if err != nil {
		return nil, fmt.Errorf("failed to build fund URL: %w", err)
	}

	funds, err := c.fetchAndDecodeFunds(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch or decode funds from %s: %w", url, err)
	}

	for _, fund := range funds {
		if strings.EqualFold(fund.Ticker, identifier) {
			return &fund, nil
		}
	}

	return nil, fmt.Errorf("fund not found with identifier: %s", identifier)
}

func (c *Client) buildFundURL(identifier string) (string, error) {
	return fmt.Sprintf(
		"https://www.ishares.com/us/product-screener/product-screener-v3.1.jsn?dcrPath=/templatedata/config/product-screener-v3/data/en/us-ishares/ishares-product-screener-backend-config&siteEntryPassthrough=true&ticker=%s",
		identifier,
	), nil
}
