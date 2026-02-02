package amundi

import (
	"bytes"
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
	if c.httpConfig.Debug {
		log.Printf("amundi: discovery request %s%s", c.config.BaseURL, c.config.DiscoveryPath)
	}
	req, err := http.NewRequestWithContext(ctx, "POST", c.config.BaseURL+c.config.DiscoveryPath, bytes.NewReader(c.config.RequestBody))
	if err != nil {
		return nil, err
	}

	for k, v := range c.config.DefaultHeaders {
		req.Header.Set(k, v)
	}

	resp, err := c.httpConfig.Client.Do(req)
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
		return nil, err
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
			AssetClass:     mapAssetClass(p.Characteristics.AssetClass),
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
