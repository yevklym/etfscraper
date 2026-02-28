package ishares

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/yevklym/etfscraper"
)

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

	url, err := c.generateHoldingsURL(*fund)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
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
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	return c.parseHoldings(resp.Body, fund)
}

func (c *Client) generateHoldingsURL(fund etfscraper.Fund) (string, error) {
	metadata, ok := fund.ProviderMetadata.(iSharesFundMetadata)
	if !ok {
		return "", fmt.Errorf("internal error: fund %s is missing iShares metadata", fund.Ticker)
	}

	if metadata.ProductPageURL != "" {
		if strings.TrimSpace(c.config.HoldingsURLTemplate) == "" {
			return "", fmt.Errorf("holdings not configured for region %s", c.region)
		}
		return fmt.Sprintf(c.config.HoldingsURLTemplate,
			c.config.BaseURL,
			metadata.ProductPageURL), nil
	}

	return "", fmt.Errorf("unable to generate holding URL for %s", fund.Ticker)
}

func (c *Client) parseHoldings(reader io.Reader, fund *etfscraper.Fund) (*etfscraper.HoldingsSnapshot, error) {
	csvReader := csv.NewReader(reader)
	csvReader.LazyQuotes = true
	csvReader.FieldsPerRecord = -1

	// Find and parse the "as of" date
	asOfDate, err := c.findAndParseDate(csvReader)
	if err != nil {
		return nil, err
	}

	// Find and parse the data header row
	headerRow, err := c.findHeaderRow(csvReader)
	if err != nil {
		return nil, err
	}

	resolver := newColumnResolver(headerRow, c.config.ColumnMappings)

	// Parse holdings data using column map
	var holdings []etfscraper.Holding
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Warning: CSV parsing error: %v", err)
			continue
		}

		// Stop at empty rows
		if len(record) == 0 || record[0] == "" {
			break
		}

		if !c.isValidHoldingRow(record, resolver) {
			break
		}

		holding, err := c.parseHoldingRecord(record, resolver)
		if err != nil {
			continue
		}

		holdings = append(holdings, holding)
	}

	if len(holdings) == 0 {
		return nil, fmt.Errorf("no holdings found for fund %s", fund.Ticker)
	}

	return &etfscraper.HoldingsSnapshot{
		Fund:          *fund,
		AsOfDate:      asOfDate,
		Holdings:      holdings,
		LastUpdated:   time.Now(),
		TotalHoldings: len(holdings),
	}, nil
}

func (c *Client) isValidHoldingRow(record []string, resolver *columnResolver) bool {

	// Check if we can get a name
	name, err := resolver.getString(record, c.config.ColumnMappings.Name)
	if err != nil || strings.TrimSpace(name) == "" {
		return false
	}

	// Check if a market value is numeric
	_, err = resolver.getFloat(record, c.config.ColumnMappings.MarketValue)
	if err != nil {
		return false
	}

	// Check if weight is numeric
	_, err = resolver.getFloat(record, c.config.ColumnMappings.Weight)
	if err != nil {
		return false
	}

	// Disclaimer check
	firstField := strings.TrimSpace(record[0])
	if len(firstField) > 100 {
		return false
	}

	disqualifiers := []string{
		"©", "http://", "https://", "www.",
		"The content", "This information",
		"Holdings subject to change",
		"carefully consider",
	}

	firstFieldLower := strings.ToLower(firstField)
	for _, disqualifier := range disqualifiers {
		if strings.Contains(firstFieldLower, strings.ToLower(disqualifier)) {
			return false
		}
	}

	return true
}

// parseHoldingRecord extracts a Holding from a CSV record using dynamic column mapping
func (c *Client) parseHoldingRecord(record []string, resolver *columnResolver) (etfscraper.Holding, error) {
	holding := etfscraper.Holding{}

	// Extract required fields
	name, err := resolver.getString(record, c.config.ColumnMappings.Name)
	if err != nil {
		return holding, fmt.Errorf("missing required field 'Name': %w", err)
	}
	holding.Name = name

	marketValue, err := resolver.getFloat(record, c.config.ColumnMappings.MarketValue)
	if err != nil {
		return holding, fmt.Errorf("invalid market value for %q: %w", name, err)
	}
	holding.MarketValue = marketValue

	weight, err := resolver.getFloat(record, c.config.ColumnMappings.Weight)
	if err != nil {
		return holding, fmt.Errorf("invalid weight for %q: %w", name, err)
	}
	// Normalize weight to 0-1 range
	if weight > 1 {
		holding.Weight = weight / 100.0
	} else {
		holding.Weight = weight
	}

	c.extractOptionalFields(&holding, record, resolver)

	return holding, nil
}

func (c *Client) extractOptionalFields(holding *etfscraper.Holding, record []string, resolver *columnResolver) {
	// String fields
	stringFields := []struct {
		mapping []string
		setter  func(string)
	}{
		{c.config.ColumnMappings.Ticker, func(v string) { holding.Ticker = v }},
		{c.config.ColumnMappings.ISIN, func(v string) { holding.ISIN = v }},
		{c.config.ColumnMappings.Sector, func(v string) { holding.Sector = normalizeSector(v, c.config.SectorMapping) }},
		{c.config.ColumnMappings.AssetClass, func(v string) { holding.AssetClass = normalizeAssetClass(v, c.config.AssetClassMapping) }},
		{c.config.ColumnMappings.Location, func(v string) { holding.Location = etfscraper.Location(v) }},
		{c.config.ColumnMappings.Exchange, func(v string) { holding.Exchange = normalizeExchange(v) }},
	}

	for _, field := range stringFields {
		if len(field.mapping) > 0 {
			if val, err := resolver.getString(record, field.mapping); err == nil {
				field.setter(val)
			}
		}
	}

	if currency := c.resolveHoldingCurrency(record, resolver); currency != "" {
		holding.Currency = currency
	}

	// Float fields
	floatFields := []struct {
		mapping []string
		setter  func(float64)
	}{
		{c.config.ColumnMappings.Price, func(v float64) { holding.Price = v }},
		{c.config.ColumnMappings.Quantity, func(v float64) { holding.Quantity = v }},
	}

	for _, field := range floatFields {
		if len(field.mapping) > 0 {
			if val, err := resolver.getFloat(record, field.mapping); err == nil {
				field.setter(val)
			}
		}
	}

	// Par Value as fallback for Quantity (bonds)
	if holding.Quantity == 0 && len(c.config.ColumnMappings.ParValue) > 0 {
		if parValue, err := resolver.getFloat(record, c.config.ColumnMappings.ParValue); err == nil {
			holding.Quantity = parValue
		}
	}
}

func (c *Client) resolveHoldingCurrency(record []string, resolver *columnResolver) etfscraper.Currency {
	marketCurrency := c.readNormalizedCurrency(record, resolver, c.config.ColumnMappings.MarketCurrency)
	if marketCurrency != "" {
		return marketCurrency
	}
	return c.readNormalizedCurrency(record, resolver, c.config.ColumnMappings.Currency)
}

func (c *Client) readNormalizedCurrency(record []string, resolver *columnResolver, mapping []string) etfscraper.Currency {
	if len(mapping) == 0 {
		return ""
	}
	val, err := resolver.getString(record, mapping)
	if err != nil {
		return ""
	}
	return normalizeCurrency(val)
}

func containsAny(slice []string, targets []string) bool {
	for _, item := range slice {
		if slices.Contains(targets, strings.TrimSpace(item)) {
			return true
		}
	}
	return false
}

func (c *Client) findAndParseDate(csvReader *csv.Reader) (time.Time, error) {
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			return time.Time{}, fmt.Errorf("CSV ended before date header was found")
		}
		if err != nil {
			return time.Time{}, fmt.Errorf("error reading CSV header: %w", err)
		}

		if len(record) < 2 {
			continue
		}

		for _, pattern := range c.config.DateHeaderPatterns {
			firstField := strings.TrimSpace(record[0])
			if strings.Contains(firstField, pattern) {
				dateStr := strings.TrimSpace(record[1])

				if c.config.MonthTranslations != nil {
					for from, to := range c.config.MonthTranslations {
						dateStr = strings.ReplaceAll(dateStr, from, to)
					}
				}

				var asOfDate time.Time
				var parseErr error

				for _, format := range c.config.DateFormats {
					asOfDate, parseErr = time.Parse(format, dateStr)
					if parseErr == nil {
						return asOfDate, nil
					}
				}

				return time.Time{}, fmt.Errorf("failed to parse date %q with any of the configured formats: %w",
					dateStr, parseErr)
			}
		}
	}
}

func (c *Client) findHeaderRow(csvReader *csv.Reader) ([]string, error) {
	expectedColumns := c.getAllExpectedColumns()

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			return nil, fmt.Errorf("CSV ended before data header was found")
		}
		if err != nil {
			return nil, fmt.Errorf("error while searching for CSV header: %w", err)
		}

		if len(record) == 0 || (len(record) == 1 && record[0] == "") {
			continue
		}

		if containsAny(record, expectedColumns) {
			return record, nil
		}
	}
}

func (c *Client) getAllExpectedColumns() []string {
	var columns []string

	columns = append(columns, c.config.ColumnMappings.Name...)
	columns = append(columns, c.config.ColumnMappings.Ticker...)
	columns = append(columns, c.config.ColumnMappings.ISIN...)
	columns = append(columns, c.config.ColumnMappings.MarketValue...)
	columns = append(columns, c.config.ColumnMappings.Weight...)

	if len(c.config.ColumnMappings.Sector) > 0 {
		columns = append(columns, c.config.ColumnMappings.Sector...)
	}
	if len(c.config.ColumnMappings.Location) > 0 {
		columns = append(columns, c.config.ColumnMappings.Location...)
	}

	return columns
}
