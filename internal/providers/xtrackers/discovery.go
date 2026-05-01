package xtrackers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yevklym/etfscraper"
)

// datatableResponse is the top-level response from the discovery API.
type datatableResponse struct {
	Values []datatableEntry `json:"values"`
}

// datatableEntry represents a single ETF in the datatable response.
type datatableEntry struct {
	ID                   entryField   `json:"ID"`
	Column0              column0      `json:"column_0"`
	AssetClass           entryField   `json:"AssetClass"`
	AssetUnderManagement numericField `json:"AssetUnderManagement"`
	Currency             entryField   `json:"Currency"`
	TotalExpenseRatio    entryField   `json:"TotalExpenseRatio"`
	UseOfProfit          entryField   `json:"UseOfProfit"`
	FundLaunchDate       entryField   `json:"FundLaunchDate"`
	PerformanceDate      entryField   `json:"PerformanceDate"`
}

// entryField is a generic field with a string value.
type entryField struct {
	Value string `json:"value"`
}

// numericField is a field with both a display value and a numeric sortValue.
type numericField struct {
	Value     string  `json:"value"`
	SortValue float64 `json:"sortValue"`
}

// column0 holds the fund name and ISIN sub-columns.
type column0 struct {
	Column00 column0Name `json:"column_0_0"`
}

// column0Name has a nested value object containing the fund name and URL.
// Real response: {"value": {"text": "...", "name": "...", "url": "..."}, "sortValue": "...", "type": "link"}
type column0Name struct {
	Value column0NameValue `json:"value"`
}

// column0NameValue holds the fund name and product page URL inside column_0_0.value.
type column0NameValue struct {
	Text string `json:"text"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

func (c *Client) discoverFresh(ctx context.Context) ([]etfscraper.Fund, error) {
	body, err := json.Marshal(discoveryRequestBody())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	url := c.discoveryURL()

	if c.httpConfig.Debug {
		c.httpConfig.Logger.Printf("xtrackers: discovery request %s (region: %s)", url, c.region)
	}

	var respBody []byte

	if c.skipBrowserFetch {
		// Used by unit tests to mock standard HTTP
		resp, err := c.doPost(ctx, url, body)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch discovery data: %w", err)
		}
		defer func() {
			if closeErr := resp.Body.Close(); closeErr != nil {
				c.httpConfig.Logger.Printf("Warning: failed to close response body: %v", closeErr)
			}
		}()

		if resp.StatusCode != http.StatusOK {
			if c.httpConfig.Debug {
				c.httpConfig.Logger.Printf("xtrackers: discovery response %s", resp.Status)
			}
			return nil, fmt.Errorf("xtrackers: discovery: HTTP %d: %s", resp.StatusCode, resp.Status)
		}

		// Read body so we can decode it exactly like the browser fetch
		var buf bytes.Buffer
		if _, err := buf.ReadFrom(resp.Body); err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}
		respBody = buf.Bytes()
	} else {
		// Used in real environments to bypass Akamai
		var err error
		respBody, err = c.doPostBrowser(ctx, url, body)
		if err != nil {
			return nil, err
		}
	}

	if c.httpConfig.Debug {
		c.httpConfig.Logger.Printf("xtrackers: received response payload (length %d)", len(respBody))
	}

	var payload datatableResponse
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if c.httpConfig.Debug {
		c.httpConfig.Logger.Printf("xtrackers: discovered %d ETFs", len(payload.Values))
	}

	return c.convertToFunds(payload.Values), nil
}

func (c *Client) convertToFunds(entries []datatableEntry) []etfscraper.Fund {
	out := make([]etfscraper.Fund, 0, len(entries))
	for _, e := range entries {
		isin := e.ID.Value
		name := e.Column0.Column00.Value.Name
		if name == "" {
			name = e.Column0.Column00.Value.Text
		}

		if isin == "" || name == "" {
			continue
		}

		f := etfscraper.Fund{
			ISIN:           isin,
			Name:           name,
			Provider:       etfscraper.ProviderXtrackers,
			Currency:       mapCurrency(e.Currency.Value),
			ExpenseRatio:   parseTER(e.TotalExpenseRatio.Value),
			TotalAssets:    parseAUM(e.AssetUnderManagement.SortValue),
			AssetClass:     normalizeAssetClass(e.AssetClass.Value, c.config.AssetClassMapping),
			IsDistributing: isDistributing(e.UseOfProfit.Value, c.config.DistributionTerms),
			InceptionDate:  parseLaunchDate(e.FundLaunchDate.Value),
			LastUpdated:    parseLastUpdated(e.PerformanceDate.Value),
		}
		f.ProviderMetadata = fundMetadata{
			ProductURL: e.Column0.Column00.Value.URL,
		}

		out = append(out, f)
	}
	return out
}
