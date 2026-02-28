package ishares

import (
	"fmt"
	"strconv"
	"strings"
)

type columnResolver struct {
	columnMap map[string]int
	mapper    ColumnMapper
}

func newColumnResolver(headers []string, mapper ColumnMapper) *columnResolver {
	columnMap := make(map[string]int)
	for i, col := range headers {
		columnMap[strings.TrimSpace(col)] = i
	}

	return &columnResolver{columnMap: columnMap, mapper: mapper}
}

// try all possible column name variations of a field
func (r *columnResolver) findColumn(fieldVariations []string) (int, bool) {
	for _, variation := range fieldVariations {
		if idx, exists := r.columnMap[variation]; exists {
			return idx, true
		}
	}
	return -1, false
}

func (r *columnResolver) getString(record []string, fieldVariations []string) (string, error) {
	idx, found := r.findColumn(fieldVariations)
	if !found {
		return "", fmt.Errorf("none of the column variations found: %v", fieldVariations)
	}

	if idx >= len(record) {
		return "", fmt.Errorf("column index %d out of bounds", idx)
	}

	val := strings.TrimSpace(record[idx])
	val = strings.Trim(val, "\"")
	return val, nil
}

func (r *columnResolver) getFloat(record []string, fieldVariations []string) (float64, error) {
	strVal, err := r.getString(record, fieldVariations)
	if err != nil {
		return 0, err
	}

	if strVal == "" || strVal == "-" {
		return 0, fmt.Errorf("empty value")
	}

	cleanVal := r.normalizeNumber(strVal)

	floatVal, err := strconv.ParseFloat(cleanVal, 64)
	if err != nil {
		return 0, fmt.Errorf("parse error: %w", err)
	}

	return floatVal, nil
}

func (r *columnResolver) normalizeNumber(s string) string {
	s = strings.TrimSpace(s)

	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "\u00A0", "") // non-breaking space (French locale)
	s = strings.ReplaceAll(s, "\u202F", "") // narrow no-break space (French locale)
	s = strings.ReplaceAll(s, "$", "")
	s = strings.ReplaceAll(s, "€", "")
	s = strings.ReplaceAll(s, "'", "") // Swiss thousands separator

	lastComma := strings.LastIndex(s, ",")
	lastDot := strings.LastIndex(s, ".")

	if lastComma > lastDot && lastComma == len(s)-3 {
		// European format: 1.234,56 -> 1234.56
		s = strings.ReplaceAll(s, ".", "")
		s = strings.ReplaceAll(s, ",", ".")
	} else {
		// US format: 1,234.56 -> 1234.56
		s = strings.ReplaceAll(s, ",", "")
	}

	return s
}
