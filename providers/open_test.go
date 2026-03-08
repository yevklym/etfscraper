package providers

import "testing"

func TestParseProviderSpec(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		if _, _, err := ParseProviderSpec(""); err == nil {
			t.Fatal("expected error for empty spec")
		}
	})

	t.Run("name only", func(t *testing.T) {
		name, region, err := ParseProviderSpec("ishares")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if name != "ishares" {
			t.Fatalf("expected name ishares, got %q", name)
		}
		if region != "" {
			t.Fatalf("expected empty region, got %q", region)
		}
	})

	t.Run("name and region", func(t *testing.T) {
		name, region, err := ParseProviderSpec("ishares:uk")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if name != "ishares" || region != "uk" {
			t.Fatalf("expected ishares:uk, got %q:%q", name, region)
		}
	})
}

func TestSupportedProviders(t *testing.T) {
	specs := SupportedProviders()
	if len(specs) == 0 {
		t.Fatal("expected supported providers")
	}
	if specs[0].Name != "amundi" {
		t.Fatalf("expected amundi first, got %q", specs[0].Name)
	}
	found := false
	for _, spec := range specs {
		if spec.Name == "ishares" && len(spec.Regions) > 0 {
			found = true
		}
	}
	if !found {
		t.Fatal("expected ishares provider with regions")
	}
}

func TestSupportedProvidersCanBeOpened(t *testing.T) {
	for _, spec := range SupportedProviders() {
		if len(spec.Regions) == 0 {
			t.Fatalf("expected regions for provider %q", spec.Name)
		}
		for _, region := range spec.Regions {
			t.Run(spec.Name+":"+region, func(t *testing.T) {
				provider, err := OpenNameRegion(spec.Name, region)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if provider == nil {
					t.Fatal("expected provider")
				}
			})
		}
	}
}

func TestOpenSpec(t *testing.T) {
	provider, err := OpenSpec(Spec{Name: "ishares", Region: "us"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if provider == nil {
		t.Fatal("expected provider")
	}
}

func TestOpenValidation(t *testing.T) {
	tests := []struct {
		name    string
		spec    string
		wantErr string
	}{
		{
			name:    "missing region",
			spec:    "ishares",
			wantErr: "provider region is required",
		},
		{
			name:    "unknown provider",
			spec:    "unknown:us",
			wantErr: "unknown provider: unknown",
		},
		{
			name:    "unsupported region",
			spec:    "ishares:xx",
			wantErr: "unsupported region \"xx\" for provider \"ishares\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Open(tt.spec)
			if err == nil {
				t.Fatal("expected error")
			}
			if err.Error() != tt.wantErr {
				t.Fatalf("expected %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}

func TestOpenSpecValidation(t *testing.T) {
	tests := []struct {
		name    string
		spec    Spec
		wantErr string
	}{
		{
			name:    "missing name",
			spec:    Spec{Region: "us"},
			wantErr: "provider name is required",
		},
		{
			name:    "missing region",
			spec:    Spec{Name: "ishares"},
			wantErr: "provider region is required",
		},
		{
			name:    "unsupported region",
			spec:    Spec{Name: "amundi", Region: "xx"},
			wantErr: "unsupported region \"xx\" for provider \"amundi\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := OpenSpec(tt.spec)
			if err == nil {
				t.Fatal("expected error")
			}
			if err.Error() != tt.wantErr {
				t.Fatalf("expected %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}

func TestOpenNameRegionValidation(t *testing.T) {
	tests := []struct {
		name     string
		provider string
		region   string
		wantErr  string
	}{
		{
			name:     "missing name",
			provider: "",
			region:   "us",
			wantErr:  "provider name is required",
		},
		{
			name:     "missing region",
			provider: "ishares",
			region:   "",
			wantErr:  "provider region is required",
		},
		{
			name:     "unknown provider",
			provider: "unknown",
			region:   "us",
			wantErr:  "unknown provider: unknown",
		},
		{
			name:     "unsupported region",
			provider: "ishares",
			region:   "xx",
			wantErr:  "unsupported region \"xx\" for provider \"ishares\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := OpenNameRegion(tt.provider, tt.region)
			if err == nil {
				t.Fatal("expected error")
			}
			if err.Error() != tt.wantErr {
				t.Fatalf("expected %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}

func TestOpenCasingNormalization(t *testing.T) {
	provider, err := Open("IShares:US")
	if err != nil {
		t.Fatalf("unexpected error for Open: %v", err)
	}
	if provider == nil {
		t.Fatal("expected provider from Open")
	}

	provider, err = OpenSpec(Spec{Name: "AmUnDi", Region: "Fr"})
	if err != nil {
		t.Fatalf("unexpected error for OpenSpec: %v", err)
	}
	if provider == nil {
		t.Fatal("expected provider from OpenSpec")
	}

	provider, err = OpenNameRegion("iShares", "Uk")
	if err != nil {
		t.Fatalf("unexpected error for OpenNameRegion: %v", err)
	}
	if provider == nil {
		t.Fatal("expected provider from OpenNameRegion")
	}
}
