package ishares

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/yevklym/etfscraper"
)

const (
	isharesUSBaseURL = "https://www.ishares.com/us"
	isharesUKBaseURL = "https://www.ishares.com/uk"
	isharesDEBaseURL = "https://www.ishares.com/de"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	region     string
}

func New(region string) *Client {
	baseURL := isharesUSBaseURL
	switch region {
	case "uk":
		baseURL = isharesUKBaseURL
	case "de":
		baseURL = isharesDEBaseURL
	}

	return &Client{
		httpClient: &http.Client{
			Timeout: time.Second * 15,
		},
		baseURL: baseURL,
		region:  region,
	}
}

func (c *Client) DiscoverETFs(ctx context.Context) ([]etfscraper.Fund, error) {
	return c.discoverFromJSON(ctx)
}

// FundInfo retrieves detailed information about a specific fund.
// TODO: This is inefficient. It should be a direct API call to get fund details.
func (c *Client) FundInfo(ctx context.Context, identifier string) (*etfscraper.Fund, error) {
	funds, err := c.DiscoverETFs(ctx)
	if err != nil {
		return nil, err
	}

	for _, fund := range funds {
		if strings.EqualFold(fund.Ticker, identifier) || strings.EqualFold(fund.ISIN, identifier) {
			return &fund, nil
		}
	}

	return nil, fmt.Errorf("fund not found: %s", identifier)
}

func (c *Client) Holdings(ctx context.Context, identifier string) (*etfscraper.HoldingsSnapshot, error) {
	fund, err := c.FundInfo(ctx, identifier)
	if err != nil {
		return nil, err
	}

	url, err := c.generateHoldingsURL(*fund)
	if err != nil {
		return nil, err
	}

	// TODO: Inefficient to read the whole file into memory. Should stream it.
	csvData, err := c.downloadCSV(ctx, url)
	if err != nil {
		return nil, err
	}

	return c.parseHoldingsCSV(csvData, fund)
}

func (c *Client) generateHoldingsURL(fund etfscraper.Fund) (string, error) {
	if fund.ProductPageURL != "" {
		return fmt.Sprintf("https://www.ishares.com%s/1467271812596.ajax?fileType=csv", fund.ProductPageURL), nil
	}

	if fund.ISIN == "" {
		return "", fmt.Errorf("unable to generate holdings URL for %s without ISIN", fund.Ticker)
	}

	switch c.region {
	case "us":
		return fmt.Sprintf("%s/products/%s/1467271812596.ajax?fileType=csv", c.baseURL, fund.ISIN), nil
	case "uk":
		return fmt.Sprintf("%s/individual/en/products/%s/fund/1506575576011.ajax?fileType=csv", c.baseURL, fund.ISIN), nil
	case "de":
		return fmt.Sprintf("%s/de/privatanleger/de/produkte/%s/fund/1506575576011.ajax?fileType=csv", c.baseURL, fund.ISIN), nil
	default:
		return "", fmt.Errorf("unsupported region for holdings URL: %s", c.region)
	}
}

func (c *Client) downloadCSV(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "text/csv,application/csv,*/*")
	req.Header.Set("Referer", c.baseURL)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	data := make([]byte, 0, resp.ContentLength)
	buf := make([]byte, 4096)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			data = append(data, buf[:n]...)
		}
		if err != nil {
			break
		}
	}

	return data, nil
}

// TODO: Implement CSV parsing logic
func (c *Client) parseHoldingsCSV(csvData []byte, fund *etfscraper.Fund) (*etfscraper.HoldingsSnapshot, error) {
	return &etfscraper.HoldingsSnapshot{
		Fund:        *fund,
		AsOfDate:    time.Now(),
		Holdings:    []etfscraper.Holding{},
		Source:      "ishares",
		LastUpdated: time.Now(),
	}, fmt.Errorf("CSV parsing not implemented yet")
}
