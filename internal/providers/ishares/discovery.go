package ishares

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/yevklym/etfscraper"
)

type isharesDataField struct {
	Display string  `json:"d"`
	Raw     float64 `json:"r"`
}

func (f *isharesDataField) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == '{' {
		type Alias isharesDataField
		var aux Alias
		if err := json.Unmarshal(data, &aux); err != nil {
			return err
		}
		*f = isharesDataField(aux)
	}
	return nil
}

type ISharesETFData struct {
	PortfolioID         int    `json:"portfolioId"`
	FundName            string `json:"fundName"`
	LocalExchangeTicker string `json:"localExchangeTicker"`
	ISIN                string `json:"isin"`
	CUSIP               string `json:"cusip"`
	ProductType         string `json:"productType"`
	InceptionDate       struct {
		Display string `json:"d"`
		Raw     int    `json:"r"`
	} `json:"inceptionDate"`
	Fees               isharesDataField `json:"fees"`
	NetExpenseRatio    isharesDataField `json:"netr"`
	Ter                isharesDataField `json:"ter"`
	TerOcf             isharesDataField `json:"ter_ocf"`
	TotalNetAssets     isharesDataField `json:"totalNetAssets"`
	NavAmount          isharesDataField `json:"navAmount"`
	ProductPageUrl     string           `json:"productPageUrl"`
	AladdinAssetClass  string           `json:"aladdinAssetClass"`
	AladdinCountry     string           `json:"aladdinCountry"`
	AladdinRegion      string           `json:"aladdinRegion"`
	SeriesBaseCurrency string           `json:"seriesBaseCurrency"`
	Exchange           string           `json:"exchange"`
}

type iSharesFundMetadata struct {
	PortfolioID    int
	ProductPageURL string
}

func (c *Client) fetchAndDecodeFunds(ctx context.Context) ([]etfscraper.Fund, error) {
	if c.httpConfig.Debug {
		log.Printf("ishares: discovery request %s", c.config.DiscoveryURL)
	}
	req, err := http.NewRequestWithContext(ctx, "GET", c.config.DiscoveryURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpConfig.Client.Do(req)
	if err != nil {
		return nil, err
	}

	// Only defer if resp is not nil
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Printf("Warning: failed to close response body: %v", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		if c.httpConfig.Debug {
			log.Printf("ishares: discovery response %s", resp.Status)
		}
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// 1. Try the Wrapper Format (US) FIRST
	// We check for 'I' being non-empty to ensure it's actually the wrapped format
	var wrapper struct {
		I map[string]ISharesETFData `json:"i"`
	}
	if err := json.Unmarshal(body, &wrapper); err == nil && len(wrapper.I) > 0 {
		return c.convertToFunds(wrapper.I), nil
	}

	// 2. Fallback to Direct Map Format (DE/Other)
	var result map[string]ISharesETFData
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode JSON with either format: %w", err)
	}

	return c.convertToFunds(result), nil
}

func (c *Client) convertToFunds(etfData map[string]ISharesETFData) []etfscraper.Fund {
	var funds []etfscraper.Fund

	for _, data := range etfData {
		// Skip mutual funds - only process ETFs
		if !c.isValidETF(data) {
			continue
		}

		funds = append(funds, c.convertSingleFund(data))
	}

	return funds
}

func (c *Client) isValidETF(data ISharesETFData) bool {
	if data.ProductType != "ISHARES_FUND_DATA" {
		return false
	}

	// Must have a valid ticker (not dash or empty)
	if data.LocalExchangeTicker == "-" || data.LocalExchangeTicker == "" {
		return false
	}

	// Must have ISIN
	if data.ISIN == "" {
		return false
	}

	return true
}

func (c *Client) convertSingleFund(data ISharesETFData) etfscraper.Fund {
	var expenseRatio float64
	if data.NetExpenseRatio.Raw > 0 {
		expenseRatio = data.NetExpenseRatio.Raw / 100.0
	} else if data.Ter.Raw > 0 {
		expenseRatio = data.Ter.Raw / 100.0
	} else if data.TerOcf.Raw > 0 {
		expenseRatio = data.TerOcf.Raw / 100.0
	}

	currency := c.config.DefaultCurrency
	if normalized := normalizeCurrency(data.SeriesBaseCurrency); normalized != "" {
		currency = normalized
	}

	exchange := c.config.DefaultExchange
	if normalized := normalizeExchange(data.Exchange); normalized != "" {
		exchange = normalized
	}

	fund := etfscraper.Fund{
		Ticker:       data.LocalExchangeTicker,
		Name:         data.FundName,
		ISIN:         data.ISIN,
		Provider:     etfscraper.ProviderIShares,
		Currency:     currency,
		ExpenseRatio: expenseRatio,
		TotalAssets:  data.TotalNetAssets.Raw,
		Exchange:     exchange,
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

	return fund
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
