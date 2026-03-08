// Package providers is the public factory for creating ETF data providers.
// Use Open, OpenSpec, or OpenNameRegion to create a Provider for a given
// provider name and region (e.g. "ishares:us", "amundi:de").
package providers

import (
	"fmt"
	"sort"
	"strings"

	"github.com/yevklym/etfscraper"
	"github.com/yevklym/etfscraper/internal/providers/amundi"
	"github.com/yevklym/etfscraper/internal/providers/ishares"
)

// ProviderSpec describes a provider and its supported regions.
type ProviderSpec struct {
	Name    string
	Regions []string
}

// Spec represents a provider selection with explicit name and region.
type Spec struct {
	Name   string
	Region string
}

// Open creates a provider from a "provider:region" spec string (e.g.
// "ishares:us", "amundi:de"). The provider name is case-insensitive.
func Open(spec string, opts ...Option) (etfscraper.Provider, error) {
	name, region, err := ParseProviderSpec(spec)
	if err != nil {
		return nil, fmt.Errorf("invalid provider spec %q: %w", spec, err)
	}

	return OpenNameRegion(name, region, opts...)
}

// OpenSpec creates a provider from a Spec value.
func OpenSpec(spec Spec, opts ...Option) (etfscraper.Provider, error) {
	return OpenNameRegion(spec.Name, spec.Region, opts...)
}

// OpenNameRegion creates a provider from explicit name and region values.
// The name is case-insensitive (e.g. "iShares" and "ishares" are equivalent).
func OpenNameRegion(name, region string, opts ...Option) (etfscraper.Provider, error) {
	name, region, err := validateProviderSelection(name, region)
	if err != nil {
		return nil, err
	}

	options := providerOptions{
		httpConfig: etfscraper.DefaultHTTPConfig(),
	}
	for _, opt := range opts {
		opt(&options)
	}

	switch strings.ToLower(name) {
	case "ishares":
		isharesOpts := []ishares.ClientOption{ishares.WithHTTPConfig(options.httpConfig)}
		if options.cacheTTL != nil {
			isharesOpts = append(isharesOpts, ishares.WithCacheTTL(*options.cacheTTL))
		}
		return ishares.New(region, isharesOpts...)
	case "amundi":
		amundiOpts := []amundi.ClientOption{amundi.WithHTTPConfig(options.httpConfig)}
		if options.cacheTTL != nil {
			amundiOpts = append(amundiOpts, amundi.WithCacheTTL(*options.cacheTTL))
		}
		return amundi.New(region, amundiOpts...)
	default:
		return nil, fmt.Errorf("unknown provider: %s", name)
	}
}

func validateProviderSelection(name, region string) (string, string, error) {
	normalizedName := strings.ToLower(strings.TrimSpace(name))
	if normalizedName == "" {
		return "", "", fmt.Errorf("provider name is required")
	}

	normalizedRegion := strings.ToLower(strings.TrimSpace(region))
	if normalizedRegion == "" {
		return "", "", fmt.Errorf("provider region is required")
	}

	supported := make(map[string]map[string]struct{})
	for _, spec := range SupportedProviders() {
		normalizedProvider := strings.ToLower(strings.TrimSpace(spec.Name))
		regions := make(map[string]struct{}, len(spec.Regions))
		for _, supportedRegion := range spec.Regions {
			normalizedSupportedRegion := strings.ToLower(strings.TrimSpace(supportedRegion))
			if normalizedSupportedRegion == "" {
				continue
			}
			regions[normalizedSupportedRegion] = struct{}{}
		}
		supported[normalizedProvider] = regions
	}

	regions, ok := supported[normalizedName]
	if !ok {
		return "", "", fmt.Errorf("unknown provider: %s", normalizedName)
	}
	if _, ok := regions[normalizedRegion]; !ok {
		return "", "", fmt.Errorf("unsupported region %q for provider %q", normalizedRegion, normalizedName)
	}

	return normalizedName, normalizedRegion, nil
}

// ParseProviderSpec parses a "provider:region" string into name and region.
// If no colon is present, region is returned as an empty string.
func ParseProviderSpec(spec string) (name string, region string, err error) {
	trimmed := strings.TrimSpace(spec)
	if trimmed == "" {
		return "", "", fmt.Errorf("provider spec cannot be empty")
	}

	parts := strings.SplitN(trimmed, ":", 2)
	name = strings.TrimSpace(parts[0])
	if name == "" {
		return "", "", fmt.Errorf("provider name cannot be empty")
	}
	if len(parts) > 1 {
		region = strings.TrimSpace(parts[1])
	}
	return name, region, nil
}

// SupportedProviders returns all providers and their supported regions.
func SupportedProviders() []ProviderSpec {
	specs := []ProviderSpec{
		{Name: "amundi", Regions: amundi.SupportedRegions()},
		{Name: "ishares", Regions: ishares.SupportedRegions()},
	}

	sort.Slice(specs, func(i, j int) bool {
		return specs[i].Name < specs[j].Name
	})

	return specs
}
