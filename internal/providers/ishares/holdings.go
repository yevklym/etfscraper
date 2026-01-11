package ishares

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/yevklym/etfscraper"
)

func (c *Client) Holdings(ctx context.Context, identifier string) (*etfscraper.HoldingsSnapshot, error) {
	fund, err := c.FundInfo(ctx, identifier)
	if err != nil {
		return nil, err
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
	var asOfDate time.Time
	const layout = "Jan 2, 2006"

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			return nil, fmt.Errorf("CSV ended before date header was found")
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV header: %w", err)
		}

		if len(record) < 2 {
			continue
		}

		if strings.Contains(record[0], "Fund Holdings as of") || strings.Contains(record[0], "as of") {
			asOfDate, err = time.Parse(layout, record[1])
			if err != nil {
				return nil, fmt.Errorf("failed to parse date %q: %w", record[1], err)
			}
			break
		}
	}

	// Find and parse the data header row
	var headerRow []string
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

		// Detect header row by checking for known column names
		if containsAny(record, []string{"Name", "Ticker", "Market Value", "Weight (%)", "Market Weight"}) {
			headerRow = record
			break
		}
	}

	// Create column index map for dynamic field access
	columnMap := make(map[string]int)
	for i, col := range headerRow {
		columnMap[strings.TrimSpace(col)] = i
	}

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

		// Stop at empty rows or legal disclaimer text
		if len(record) == 0 || record[0] == "" || isDisclaimerRow(record[0]) {
			break
		}

		holding, err := c.parseHoldingRecord(record, columnMap)
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

// isDisclaimerRow detects if a row contains legal disclaimer text
func isDisclaimerRow(firstField string) bool {
	disclaimerPatterns := []string{
		"The content contained herein",
		"The Funds are distributed by",
		"Holdings subject to change",
		"CAREFULLY CONSIDER THE FUNDS",
		"Past performance does not guarantee",
		"©20", // Copyright notices
		"The iShares Funds are not sponsored",
	}

	for _, pattern := range disclaimerPatterns {
		if strings.Contains(firstField, pattern) {
			return true
		}
	}

	return false
}

// parseHoldingRecord extracts a Holding from a CSV record using dynamic column mapping
func (c *Client) parseHoldingRecord(record []string, columnMap map[string]int) (etfscraper.Holding, error) {
	holding := etfscraper.Holding{}

	// Extract Name (required field)
	name, err := getStringField(record, columnMap, "Name")
	if err != nil {
		return holding, fmt.Errorf("missing required field 'Name': %w", err)
	}
	holding.Name = name

	// Extract Ticker (optional)
	if ticker, err := getStringField(record, columnMap, "Ticker"); err == nil {
		holding.Ticker = ticker
	}

	// Extract ISIN (optional)
	if isin, err := getStringField(record, columnMap, "ISIN"); err == nil {
		holding.ISIN = isin
	}

	// Extract Market Value (required)
	marketValue, err := getFloatField(record, columnMap, "Market Value")
	if err != nil {
		return holding, fmt.Errorf("invalid market value for %q: %w", name, err)
	}
	holding.MarketValue = marketValue

	// Extract Weight (required)
	weight, err := getFloatFieldWithAlternatives(record, columnMap, []string{"Weight (%)", "Market Weight", "Notional Weight"})
	if err != nil {
		return holding, fmt.Errorf("invalid weight for %q: %w", name, err)
	}
	// If column is "Weight (%)", it's already in percentage; if "Market Weight", it might be too
	// Check if value is > 1 to determine if it needs division by 100
	if weight > 1 {
		holding.Weight = weight / 100.0
	} else {
		holding.Weight = weight
	}

	// Extract Quantity (optional - not present in all CSV formats)
	if quantity, err := getFloatField(record, columnMap, "Quantity"); err == nil {
		holding.Quantity = quantity
	} else if parValue, err := getFloatField(record, columnMap, "Par Value"); err == nil {
		// Bonds use Par Value instead of Quantity
		holding.Quantity = parValue
	}

	// Extract Price (optional)
	if price, err := getFloatField(record, columnMap, "Price"); err == nil {
		holding.Price = price
	}

	// Extract Sector (optional)
	if sector, err := getStringField(record, columnMap, "Sector"); err == nil {
		holding.Sector = etfscraper.Sector(sector)
	}

	// Extract Asset Class (optional)
	if assetClass, err := getStringField(record, columnMap, "Asset Class"); err == nil {
		holding.AssetClass = etfscraper.AssetClass(assetClass)
	}

	// Extract Location (optional)
	if location, err := getStringField(record, columnMap, "Location"); err == nil {
		holding.Location = etfscraper.Location(location)
	}

	// Extract Exchange (optional)
	if exchange, err := getStringField(record, columnMap, "Exchange"); err == nil {
		holding.Exchange = etfscraper.Exchange(exchange)
	}

	// Extract Currency (optional)
	if currency, err := getStringField(record, columnMap, "Currency"); err == nil {
		holding.Currency = etfscraper.Currency(currency)
	}

	return holding, nil
}

// getFloatFieldWithAlternatives tries multiple column names in order
func getFloatFieldWithAlternatives(record []string, columnMap map[string]int, fieldNames []string) (float64, error) {
	var lastErr error
	for _, fieldName := range fieldNames {
		val, err := getFloatField(record, columnMap, fieldName)
		if err == nil {
			return val, nil
		}
		lastErr = err
	}
	return 0, fmt.Errorf("none of the alternative columns found: %w", lastErr)
}

// Helper functions for safe field extraction
func getStringField(record []string, columnMap map[string]int, fieldName string) (string, error) {
	idx, exists := columnMap[fieldName]
	if !exists {
		return "", fmt.Errorf("column %q not found", fieldName)
	}
	if idx >= len(record) {
		return "", fmt.Errorf("column %q index %d out of bounds", fieldName, idx)
	}

	val := strings.TrimSpace(record[idx])
	val = strings.Trim(val, "\"")
	return val, nil
}

func getFloatField(record []string, columnMap map[string]int, fieldName string) (float64, error) {
	strVal, err := getStringField(record, columnMap, fieldName)
	if err != nil {
		return 0, err
	}

	// Handle empty values
	if strVal == "" || strVal == "-" {
		return 0, fmt.Errorf("empty value")
	}

	// Remove commas and parse
	cleanVal := strings.ReplaceAll(strVal, ",", "")
	floatVal, err := strconv.ParseFloat(cleanVal, 64)
	if err != nil {
		return 0, fmt.Errorf("parse error: %w", err)
	}

	return floatVal, nil
}

func containsAny(slice []string, targets []string) bool {
	for _, item := range slice {
		for _, target := range targets {
			if strings.TrimSpace(item) == target {
				return true
			}
		}
	}
	return false
}
