package amundi

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/yevklym/etfscraper"
)

type productsResponse struct {
	Products []product `json:"products"`
}

type product struct {
	ProductID       string          `json:"productId"`
	ProductType     string          `json:"productType"`
	Characteristics characteristics `json:"characteristics"`
}

type characteristics struct {
	ISIN               string            `json:"ISIN"`
	ShareName          string            `json:"SHARE_MARKETING_NAME"`
	Mnemo              string            `json:"MNEMO"`
	Ticker             string            `json:"TICKER"`
	Ter                float64           `json:"TER"`
	Currency           string            `json:"CURRENCY"`
	FundAUM            float64           `json:"FUND_AUM"`
	FundAUMInEuro      float64           `json:"FUND_AUM_IN_EURO"`
	AssetClass         string            `json:"ASSET_CLASS"`
	DistributionPolicy string            `json:"DISTRIBUTION_POLICY"`
	MainListings       map[string]string `json:"MAIN_LISTINGS"`
	FundIssuer         string            `json:"FUND_ISSUER"`
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

func (c *Client) discoverFresh(ctx context.Context) ([]etfscraper.Fund, error) {
	requestBody, err := buildDiscoveryRequest(c.region)
	if err != nil {
		return nil, fmt.Errorf("failed to build discovery request: %w", err)
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	url := c.config.BaseURL + c.config.DiscoveryPath

	if c.httpConfig.Debug {
		log.Printf("amundi: discovery request %s (region: %s)", url, c.region)
	}

	resp, err := c.doPost(ctx, url, body)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Printf("Warning: failed to close response body: %v", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		if c.httpConfig.Debug {
			log.Printf("amundi: discovery response %s", resp.Status)
		}
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	var payload productsResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if c.httpConfig.Debug {
		log.Printf("amundi: discovered %d ETFs", len(payload.Products))
	}

	return c.convertToFunds(payload.Products), nil
}

func (c *Client) convertToFunds(products []product) []etfscraper.Fund {
	var out []etfscraper.Fund
	for _, p := range products {
		if p.ProductType != "PRODUCT" {
			continue
		}
		f := etfscraper.Fund{
			Ticker:         pickTicker(p.Characteristics),
			ISIN:           pickISIN(p),
			Name:           p.Characteristics.ShareName,
			Provider:       etfscraper.ProviderAmundi,
			Currency:       mapCurrency(p.Characteristics.Currency),
			ExpenseRatio:   p.Characteristics.Ter / 100.0,
			TotalAssets:    p.Characteristics.FundAUM,
			AssetClass:     normalizeAssetClass(p.Characteristics.AssetClass, c.config.AssetClassMapping),
			IsDistributing: isDistributing(p.Characteristics.DistributionPolicy),
		}
		f.ProviderMetadata = amundiFundMetadata{
			ProductID:     p.ProductID,
			MainListings:  p.Characteristics.MainListings,
			FundIssuer:    p.Characteristics.FundIssuer,
			FundAUMInEuro: p.Characteristics.FundAUMInEuro,
		}
		if f.Ticker != "" && f.Name != "" && f.ISIN != "" {
			out = append(out, f)
		}
	}
	return out
}
