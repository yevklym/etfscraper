package amundi

import (
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("unsupported region", func(t *testing.T) {
		_, err := New("invalid-region")
		if err == nil {
			t.Fatal("expected error for invalid region")
		}

		expectedMsg := "unsupported region"
		if !strings.Contains(err.Error(), expectedMsg) {
			t.Errorf("expected error containing %q, got %q", expectedMsg, err.Error())
		}
	})

	t.Run("supported regions", func(t *testing.T) {
		regions := []string{"de", "DE"}

		for _, region := range regions {
			_, err := New(region)
			if err != nil {
				t.Errorf("expected region %q to be supported, got error: %v", region, err)
			}
		}
	})
}
