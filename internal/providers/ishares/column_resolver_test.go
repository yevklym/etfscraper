package ishares

import "testing"

func TestColumnResolver_FindColumn(t *testing.T) {
	headers := []string{"Name", "Ticker", "Market Value", "Weight (%)"}
	resolver := newColumnResolver(headers, ColumnMapper{})

	tests := []struct {
		name       string
		variations []string
		wantIdx    int
		wantFound  bool
	}{
		{
			name:       "exact match",
			variations: []string{"Name"},
			wantIdx:    0,
			wantFound:  true,
		},
		{
			name:       "first variation matches",
			variations: []string{"Ticker Symbol", "Ticker"},
			wantIdx:    1,
			wantFound:  true,
		},
		{
			name:       "column with special chars",
			variations: []string{"Weight (%)", "Weight"},
			wantIdx:    3,
			wantFound:  true,
		},
		{
			name:       "no match",
			variations: []string{"ISIN", "Identifier"},
			wantIdx:    -1,
			wantFound:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx, found := resolver.findColumn(tt.variations)
			if idx != tt.wantIdx {
				t.Errorf("findColumn() idx = %v, want %v", idx, tt.wantIdx)
			}
			if found != tt.wantFound {
				t.Errorf("findColumn() found = %v, want %v", found, tt.wantFound)
			}
		})
	}
}

func TestColumnResolver_GetString(t *testing.T) {
	headers := []string{"Name", "Ticker", "Description"}
	resolver := newColumnResolver(headers, ColumnMapper{})

	tests := []struct {
		name       string
		record     []string
		variations []string
		want       string
		wantErr    bool
	}{
		{
			name:       "normal string",
			record:     []string{"Apple Inc", "AAPL", "Tech company"},
			variations: []string{"Name"},
			want:       "Apple Inc",
			wantErr:    false,
		},
		{
			name:       "quoted string",
			record:     []string{`"Microsoft Corp"`, "MSFT", "Software"},
			variations: []string{"Name"},
			want:       "Microsoft Corp",
			wantErr:    false,
		},
		{
			name:       "string with whitespace",
			record:     []string{"  Tesla  ", "TSLA", "EV"},
			variations: []string{"Name"},
			want:       "Tesla",
			wantErr:    false,
		},
		{
			name:       "column not found",
			record:     []string{"Apple", "AAPL", "Tech"},
			variations: []string{"NonExistent"},
			want:       "",
			wantErr:    true,
		},
		{
			name:       "index out of bounds",
			record:     []string{"Apple"},
			variations: []string{"Description"},
			want:       "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolver.getString(tt.record, tt.variations)
			if (err != nil) != tt.wantErr {
				t.Errorf("getString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColumnResolver_GetFloat(t *testing.T) {
	headers := []string{"Name", "Market Value", "Weight"}
	resolver := newColumnResolver(headers, ColumnMapper{})

	tests := []struct {
		name       string
		record     []string
		variations []string
		want       float64
		wantErr    bool
	}{
		{
			name:       "US format",
			record:     []string{"Apple", "1,234,567.89", "5.5"},
			variations: []string{"Market Value"},
			want:       1234567.89,
			wantErr:    false,
		},
		{
			name:       "European format",
			record:     []string{"Apple", "1.234.567,89", "5.5"},
			variations: []string{"Market Value"},
			want:       1234567.89,
			wantErr:    false,
		},
		{
			name:       "with dollar sign",
			record:     []string{"Apple", "$1,234.56", "5.5"},
			variations: []string{"Market Value"},
			want:       1234.56,
			wantErr:    false,
		},
		{
			name:       "with euro sign",
			record:     []string{"Apple", "€1.234,56", "5.5"},
			variations: []string{"Market Value"},
			want:       1234.56,
			wantErr:    false,
		},
		{
			name:       "dash value",
			record:     []string{"Apple", "-", "5.5"},
			variations: []string{"Market Value"},
			want:       0,
			wantErr:    true,
		},
		{
			name:       "empty value",
			record:     []string{"Apple", "", "5.5"},
			variations: []string{"Market Value"},
			want:       0,
			wantErr:    true,
		},
		{
			name:       "invalid number",
			record:     []string{"Apple", "invalid", "5.5"},
			variations: []string{"Market Value"},
			want:       0,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolver.getFloat(tt.record, tt.variations)
			if (err != nil) != tt.wantErr {
				t.Errorf("getFloat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				epsilon := 0.001
				if diff := got - tt.want; diff > epsilon || diff < -epsilon {
					t.Errorf("getFloat() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestColumnResolver_NormalizeNumber(t *testing.T) {
	resolver := &columnResolver{}

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "US format with commas",
			input: "1,234,567.89",
			want:  "1234567.89",
		},
		{
			name:  "European format",
			input: "1.234.567,89",
			want:  "1234567.89",
		},
		{
			name:  "Swiss format",
			input: "1'234'567.89",
			want:  "1234567.89",
		},
		{
			name:  "with dollar sign",
			input: "$1,234.56",
			want:  "1234.56",
		},
		{
			name:  "with euro sign",
			input: "€1.234,56",
			want:  "1234.56",
		},
		{
			name:  "with spaces",
			input: "1 234 567.89",
			want:  "1234567.89",
		},
		{
			name:  "simple number",
			input: "123.45",
			want:  "123.45",
		},
		{
			name:  "negative number",
			input: "-123.45",
			want:  "-123.45",
		},
		{
			name:  "French format with spaces and comma decimal",
			input: "10 569 831 271,35",
			want:  "10569831271.35",
		},
		{
			name:  "non-breaking spaces",
			input: "10\u00A0569\u00A0831\u00A0271,35",
			want:  "10569831271.35",
		},
		{
			name:  "narrow no-break spaces",
			input: "10\u202F569\u202F831\u202F271,35",
			want:  "10569831271.35",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolver.normalizeNumber(tt.input)
			if got != tt.want {
				t.Errorf("normalizeNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}
