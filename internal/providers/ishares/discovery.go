package ishares

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yevklym/etfscraper"
)

type ISharesETFData struct {
	PortfolioID         int    `json:"portfolioId"`
	FundName            string `json:"fundName"`
	LocalExchangeTicker string `json:"localExchangeTicker"`
	ISIN                string `json:"isin"`
	CUSIP               string `json:"cusip"`
	InceptionDate       struct {
		Display string `json:"d"`
		Raw     int    `json:"r"`
	} `json:"inceptionDate"`
	Fees struct {
		Display string  `json:"d"`
		Raw     float64 `json:"r"`
	} `json:"fees"`
	NetExpenseRatio struct {
		Display string  `json:"d"`
		Raw     float64 `json:"r"`
	} `json:"netr"`
	TotalNetAssets struct {
		Display string  `json:"d"`
		Raw     float64 `json:"r"`
	} `json:"totalNetAssets"`
	NavAmount struct {
		Display string  `json:"d"`
		Raw     float64 `json:"r"`
	} `json:"navAmount"`
	ProductPageUrl    string `json:"productPageUrl"`
	AladdinAssetClass string `json:"aladdinAssetClass"`
	AladdinCountry    string `json:"aladdinCountry"`
	AladdinRegion     string `json:"aladdinRegion"`
}

type iSharesFundMetadata struct {
	PortfolioID    int
	ProductPageURL string
}

const usETFDiscoveryURL = "https://www.ishares.com/us/product-screener/product-screener-v3.1.jsn?dcrPath=/templatedata/config/product-screener-v3/data/en/us-ishares/ishares-product-screener-backend-config&siteEntryPassthrough=true"

func (c *Client) discoverFromJSON(ctx context.Context) ([]etfscraper.Fund, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", usETFDiscoveryURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	var etfData map[string]ISharesETFData
	if err := json.NewDecoder(resp.Body).Decode(&etfData); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return c.convertToFunds(etfData), nil
}

func (c *Client) convertToFunds(etfData map[string]ISharesETFData) []etfscraper.Fund {
	var funds []etfscraper.Fund

	for _, data := range etfData {
		if data.FundName == "" || data.LocalExchangeTicker == "" {
			continue
		}

		fund := etfscraper.Fund{
			Ticker:       data.LocalExchangeTicker,
			Name:         data.FundName,
			ISIN:         data.ISIN,
			Provider:     etfscraper.ProviderIShares,
			Currency:     etfscraper.CurrencyUSD,
			ExpenseRatio: data.NetExpenseRatio.Raw / 100.0,
			TotalAssets:  data.TotalNetAssets.Raw,
			Exchange:     etfscraper.ExchangeNYSE,
		}

		fund.ProviderMetadata = iSharesFundMetadata{
			PortfolioID:    data.PortfolioID,
			ProductPageURL: data.ProductPageUrl,
		}

		// Parse inception date
		if data.InceptionDate.Raw > 0 {
			if date := parseISharesDate(data.InceptionDate.Raw); date != nil {
				fund.InceptionDate = date
			}
		}

		funds = append(funds, fund)
	}

	return funds
}

func parseISharesDate(dateInt int) *time.Time {
	dateStr := fmt.Sprintf("%d", dateInt)
	if len(dateStr) != 8 {
		return nil
	}

	if date, err := time.Parse("20060102", dateStr); err == nil {
		return &date
	}
	return nil
}
