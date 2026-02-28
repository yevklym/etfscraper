package amundi

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/yevklym/etfscraper"
)

type holdingsResponse struct {
	Products []holdingsProduct `json:"products"`
}

type holdingsProduct struct {
	ProductID       string                  `json:"productId"`
	ProductType     string                  `json:"productType"`
	Characteristics holdingsCharacteristics `json:"characteristics"`
	Composition     json.RawMessage         `json:"composition"`
}

type holdingsCharacteristics struct {
	ISIN                   string    `json:"ISIN"`
	PositionAsOfDate       dateValue `json:"POSITION_AS_OF_DATE"`
	FundBreakdownsAsOfDate dateValue `json:"FUND_BREAKDOWNS_AS_OF_DATE"`
}

type compositionItem struct {
	Name          string  `json:"name"`
	ISIN          string  `json:"isin"`
	Bloomberg     string  `json:"bbg"`
	Weight        float64 `json:"weight"`
	Quantity      float64 `json:"quantity"`
	MarketValue   float64 `json:"marketValue"`
	Value         float64 `json:"value"`
	Currency      string  `json:"currency"`
	Sector        string  `json:"sector"`
	Type          string  `json:"type"`
	CountryOfRisk string  `json:"countryOfRisk"`
}

type compositionEntry struct {
	CompositionCharacteristics compositionItem `json:"compositionCharacteristics"`
	Weight                     float64         `json:"weight"`
}

type compositionResponse struct {
	TotalNumberOfInstruments int                `json:"totalNumberOfInstruments"`
	CompositionData          []compositionEntry `json:"compositionData"`
}

func (c *Client) Holdings(ctx context.Context, identifier string) (*etfscraper.HoldingsSnapshot, error) {
	fund, err := c.FundInfo(ctx, identifier)
	if err != nil {
		return nil, err
	}
	return c.HoldingsForFund(ctx, fund)
}

func (c *Client) HoldingsForFund(ctx context.Context, fund *etfscraper.Fund) (*etfscraper.HoldingsSnapshot, error) {
	if fund == nil {
		return nil, fmt.Errorf("fund cannot be nil")
	}

	requestBody, err := buildHoldingsRequest(c.region, fund.ISIN)
	if err != nil {
		return nil, fmt.Errorf("failed to build holdings request: %w", err)
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	url := c.config.BaseURL + c.config.HoldingsPath

	if c.httpConfig.Debug {
		log.Printf("amundi: holdings request %s (isin: %s)", url, fund.ISIN)
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
			log.Printf("amundi: holdings response %s", resp.Status)
		}
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	var payload holdingsResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	product, err := selectHoldingsProduct(payload.Products, fund.ISIN)
	if err != nil {
		return nil, err
	}

	asOfDate, err := parseHoldingsDate(product.Characteristics)
	if err != nil {
		return nil, err
	}

	composition, err := parseComposition(product.Composition)
	if err != nil {
		return nil, err
	}

	holdings := c.convertHoldings(composition, fund.TotalAssets)
	if len(holdings) == 0 {
		return nil, fmt.Errorf("%w: fund %s", etfscraper.ErrHoldingsUnavailable, fund.Ticker)
	}

	return &etfscraper.HoldingsSnapshot{
		Fund:          *fund,
		AsOfDate:      asOfDate,
		Holdings:      holdings,
		LastUpdated:   time.Now(),
		TotalHoldings: len(holdings),
	}, nil
}

func selectHoldingsProduct(products []holdingsProduct, isin string) (*holdingsProduct, error) {
	if len(products) == 0 {
		return nil, fmt.Errorf("no holdings data returned")
	}

	for i := range products {
		if strings.EqualFold(products[i].Characteristics.ISIN, isin) {
			return &products[i], nil
		}
	}

	return &products[0], nil
}

func parseHoldingsDate(ch holdingsCharacteristics) (time.Time, error) {
	if ch.PositionAsOfDate.Valid {
		return ch.PositionAsOfDate.Time, nil
	}
	if ch.FundBreakdownsAsOfDate.Valid {
		return ch.FundBreakdownsAsOfDate.Time, nil
	}

	return time.Time{}, fmt.Errorf("holdings date missing")
}

func parseComposition(raw json.RawMessage) ([]compositionItem, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return nil, fmt.Errorf("composition data missing")
	}

	var direct []compositionItem
	if err := json.Unmarshal(raw, &direct); err == nil {
		return direct, nil
	}

	var response compositionResponse
	if err := json.Unmarshal(raw, &response); err == nil && len(response.CompositionData) > 0 {
		items := make([]compositionItem, 0, len(response.CompositionData))
		for _, entry := range response.CompositionData {
			item := entry.CompositionCharacteristics
			if item.Weight == 0 {
				item.Weight = entry.Weight
			}
			items = append(items, item)
		}
		return items, nil
	}

	var wrapper struct {
		Composition []compositionItem `json:"composition"`
		Holdings    []compositionItem `json:"holdings"`
		Items       []compositionItem `json:"items"`
	}
	if err := json.Unmarshal(raw, &wrapper); err != nil {
		return nil, fmt.Errorf("unexpected composition format")
	}

	if len(wrapper.Composition) > 0 {
		return wrapper.Composition, nil
	}
	if len(wrapper.Holdings) > 0 {
		return wrapper.Holdings, nil
	}
	if len(wrapper.Items) > 0 {
		return wrapper.Items, nil
	}

	return nil, fmt.Errorf("composition data missing")
}

func (c *Client) convertHoldings(items []compositionItem, fundTotalAssets float64) []etfscraper.Holding {
	holdings := make([]etfscraper.Holding, 0, len(items))
	for _, item := range items {
		name := strings.TrimSpace(item.Name)
		if name == "" {
			continue
		}

		weight := normalizeWeight(item.Weight)
		marketValue := item.MarketValue
		if marketValue == 0 && item.Value != 0 {
			marketValue = item.Value
		}
		if marketValue == 0 && fundTotalAssets > 0 && weight > 0 {
			marketValue = fundTotalAssets * weight
		}

		holding := etfscraper.Holding{
			Name:        name,
			ISIN:        strings.TrimSpace(item.ISIN),
			Ticker:      strings.TrimSpace(item.Bloomberg),
			Weight:      weight,
			Quantity:    item.Quantity,
			MarketValue: marketValue,
			Currency:    mapCurrency(item.Currency),
			Sector:      normalizeSector(item.Sector, c.config.SectorMapping),
			AssetClass:  normalizeAssetClass(item.Type, c.config.AssetClassMapping),
			Location:    etfscraper.Location(item.CountryOfRisk),
		}

		holdings = append(holdings, holding)
	}

	return holdings
}

func normalizeWeight(weight float64) float64 {
	if weight > 1 {
		return weight / 100.0
	}
	return weight
}

type dateValue struct {
	Time  time.Time
	Valid bool
}

func (d *dateValue) UnmarshalJSON(data []byte) error {
	str := strings.TrimSpace(string(data))
	if str == "" || str == "null" {
		return nil
	}

	if str[0] == '"' {
		unquoted := strings.Trim(str, "\"")
		if unquoted == "" {
			return nil
		}
		parsed, err := time.Parse("2006-01-02", unquoted)
		if err != nil {
			return err
		}
		d.Time = parsed.UTC()
		d.Valid = true
		return nil
	}

	var value float64
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	if value <= 0 {
		return nil
	}

	d.Time = time.UnixMilli(int64(value)).UTC()
	d.Valid = true
	return nil
}
