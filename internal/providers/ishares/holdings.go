package ishares

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/yevklym/etfscraper"
)

func (c *Client) Holdings(ctx context.Context, identifier string) (*etfscraper.HoldingsSnapshot, error) {
	fund, err := c.FundInfo(ctx, identifier)
	if err != nil {
		return nil, fmt.Errorf("ishares: holdings: %w", err)
	}
	return c.HoldingsForFund(ctx, fund)
}

func (c *Client) HoldingsForFund(ctx context.Context, fund *etfscraper.Fund) (*etfscraper.HoldingsSnapshot, error) {
	if fund == nil {
		return nil, fmt.Errorf("fund cannot be nil")
	}

	url, err := c.generateHoldingsURL(*fund)
	if err != nil {
		return nil, fmt.Errorf("ishares: holdings: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("ishares: holdings: failed to create request: %w", err)
	}

	resp, err := c.httpConfig.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ishares: holdings: request failed: %w", err)
	}

	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			c.httpConfig.Logger.Printf("Warning: failed to close response body: %v", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ishares: holdings: HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	return c.parseHoldings(ctx, resp.Body, fund)
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

func (c *Client) parseHoldings(ctx context.Context, reader io.Reader, fund *etfscraper.Fund) (*etfscraper.HoldingsSnapshot, error) {
	csvReader := csv.NewReader(reader)
	csvReader.LazyQuotes = true
	csvReader.FieldsPerRecord = -1

	asOfDate, err := c.findAndParseDate(csvReader)
	if err != nil {
		return nil, fmt.Errorf("failed to find date header: %w", err)
	}

	headerRow, err := c.findHeaderRow(csvReader)
	if err != nil {
		return nil, fmt.Errorf("failed to find data header: %w", err)
	}

	resolver := newColumnResolver(headerRow)

	holdings, err := c.readHoldingRecords(ctx, csvReader, resolver)
	if err != nil {
		return nil, err
	}

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

// readHoldingRecords iterates over CSV data rows, parsing each into a Holding
// until EOF, an empty row, or a disclaimer row is encountered.
func (c *Client) readHoldingRecords(ctx context.Context, csvReader *csv.Reader, resolver *columnResolver) ([]etfscraper.Holding, error) {
	holdings := make([]etfscraper.Holding, 0, 128)
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			c.httpConfig.Logger.Printf("Warning: CSV parsing error: %v", err)
			continue
		}

		if len(record) == 0 || strings.TrimSpace(record[0]) == "" {
			break
		}

		if isDisclaimerRow(record) {
			break
		}

		holding, err := c.parseHoldingRecord(record, resolver)
		if err != nil {
			c.httpConfig.Logger.Printf("Warning: skipping holding record: %v", err)
			continue
		}

		holdings = append(holdings, holding)
	}
	return holdings, nil
}

// isDisclaimerRow detects footer/disclaimer rows that should terminate CSV parsing.
func isDisclaimerRow(record []string) bool {
	firstField := strings.TrimSpace(record[0])
	if len(firstField) > 100 {
		return true
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
			return true
		}
	}

	return false
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
		{c.config.ColumnMappings.Location, func(v string) { holding.Location = normalizeLocation(v, c.config.LocationMapping) }},
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
			return time.Time{}, fmt.Errorf("csv ended before date header was found")
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
					dateStr = translateMonth(dateStr, c.config.MonthTranslations)
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

// translateMonth applies month name translations longest-first to avoid
// partial matches (e.g. "März" must be replaced before "Mär").
func translateMonth(dateStr string, translations map[string]string) string {
	keys := make([]string, 0, len(translations))
	for k := range translations {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return len(keys[i]) > len(keys[j])
	})
	for _, from := range keys {
		dateStr = strings.ReplaceAll(dateStr, from, translations[from])
	}
	return dateStr
}

func (c *Client) findHeaderRow(csvReader *csv.Reader) ([]string, error) {
	expectedColumns := c.getAllExpectedColumns()

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			return nil, fmt.Errorf("csv ended before data header was found")
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
