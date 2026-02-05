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

func TestOpenSpec(t *testing.T) {
	provider, err := OpenSpec(Spec{Name: "ishares", Region: "us"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if provider == nil {
		t.Fatal("expected provider")
	}
}

func TestOpenUnknownProvider(t *testing.T) {
	if _, err := Open("unknown:us"); err == nil {
		t.Fatal("expected error for unknown provider")
	}
}
