package xtrackers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/yevklym/etfscraper"
)

// holdingsResponse is the top-level response from the holdings API.
type holdingsResponse struct {
	Tables []holdingsTable `json:"tables"`
}

// holdingsTable represents a single table in the holdings response.
type holdingsTable struct {
	Values []holdingsEntry `json:"values"`
}

// holdingsEntry represents a single holding row.
type holdingsEntry struct {
	Header  holdingsField        `json:"header"`
	Column0 holdingsField        `json:"column_0"`
	Column1 holdingsNumericField `json:"column_1"`
	Column2 holdingsNumericField `json:"column_2"`
	Column3 holdingsField        `json:"column_3"`
	Column4 holdingsField        `json:"column_4"`
	Column5 holdingsField        `json:"column_5"`
}

// holdingsField is a text field in the holdings response.
type holdingsField struct {
	Value string `json:"value"`
}

// holdingsNumericField has both a display value and a numeric sortValue.
type holdingsNumericField struct {
	Value     string  `json:"value"`
	SortValue float64 `json:"sortValue"`
}

// fetchHoldings fetches and parses holdings for a fund using the browser-based GET.
func (c *Client) fetchHoldings(ctx context.Context, fund *etfscraper.Fund) (*etfscraper.HoldingsSnapshot, error) {
	metadata, ok := fund.ProviderMetadata.(fundMetadata)
	if !ok || metadata.ProductURL == "" {
		return nil, fmt.Errorf("xtrackers: holdings: fund %s is missing product URL metadata", fund.ISIN)
	}

	url := c.holdingsURL(metadata.ProductURL)

	if c.httpConfig.Debug {
		c.httpConfig.Logger.Printf("xtrackers: holdings request %s (isin: %s)", url, fund.ISIN)
	}

	var respBody []byte

	if c.skipBrowserFetch {
		resp, err := c.doGet(ctx, url)
		if err != nil {
			return nil, fmt.Errorf("xtrackers: holdings: %w", err)
		}
		defer func() {
			if closeErr := resp.Body.Close(); closeErr != nil {
				c.httpConfig.Logger.Printf("Warning: failed to close response body: %v", closeErr)
			}
		}()

		bodyBytes, err := readAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("xtrackers: holdings: failed to read response: %w", err)
		}
		respBody = bodyBytes
	} else {
		var err error
		respBody, err = c.doGetBrowser(ctx, url)
		if err != nil {
			return nil, fmt.Errorf("xtrackers: holdings: %w", err)
		}
	}

	var payload holdingsResponse
	if err := json.Unmarshal(respBody, &payload); err != nil {
		return nil, fmt.Errorf("xtrackers: holdings: failed to decode response: %w", err)
	}

	if len(payload.Tables) == 0 || len(payload.Tables[0].Values) == 0 {
		return nil, fmt.Errorf("%w: fund %s", etfscraper.ErrHoldingsUnavailable, fund.ISIN)
	}

	holdings := c.convertHoldings(payload.Tables[0].Values)
	if len(holdings) == 0 {
		return nil, fmt.Errorf("%w: fund %s", etfscraper.ErrHoldingsUnavailable, fund.ISIN)
	}

	return &etfscraper.HoldingsSnapshot{
		Fund:          *fund,
		AsOfDate:      fund.LastUpdated,
		Holdings:      holdings,
		LastUpdated:   time.Now(),
		TotalHoldings: len(holdings),
	}, nil
}

// convertHoldings maps raw API entries to the canonical Holding type.
func (c *Client) convertHoldings(entries []holdingsEntry) []etfscraper.Holding {
	holdings := make([]etfscraper.Holding, 0, len(entries))
	for _, e := range entries {
		name := strings.TrimSpace(e.Column0.Value)
		if name == "" {
			continue
		}

		holding := etfscraper.Holding{
			ISIN:        strings.TrimSpace(e.Header.Value),
			Name:        name,
			Weight:      normalizeWeight(e.Column1.SortValue),
			MarketValue: e.Column2.SortValue,
			Location:    etfscraper.Location(strings.TrimSpace(e.Column3.Value)),
			Sector:      etfscraper.Sector(strings.TrimSpace(e.Column4.Value)),
			AssetClass:  normalizeAssetClass(e.Column5.Value, c.config.AssetClassMapping),
		}

		holdings = append(holdings, holding)
	}
	return holdings
}

// holdingsURL builds the full URL for the holdings API endpoint.
func (c *Client) holdingsURL(productURL string) string {
	// productURL is like "/de-de/IE00BJ0KDQ92-msci-world-ucits-etf-1c/"
	slug := strings.TrimPrefix(productURL, "/"+c.config.Locale+"/")
	slug = strings.TrimSuffix(slug, "/")
	return fmt.Sprintf("%s/api/pdp/%s/etf/%s/holdings", c.config.BaseURL, c.config.Locale, slug)
}
