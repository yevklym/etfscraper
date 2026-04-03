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
	"github.com/yevklym/etfscraper/internal/providers/xtrackers"
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

type providerFactory struct {
	regions     []string
	regionSet   map[string]struct{}
	constructor func(region string, options providerOptions) (etfscraper.Provider, error)
}

var providerRegistry = map[string]providerFactory{
	"amundi": newProviderFactory(
		amundi.SupportedRegions(),
		func(region string, options providerOptions) (etfscraper.Provider, error) {
			amundiOpts := []amundi.ClientOption{amundi.WithHTTPConfig(options.httpConfig)}
			if options.cacheTTL != nil {
				amundiOpts = append(amundiOpts, amundi.WithCacheTTL(*options.cacheTTL))
			}
			return amundi.New(region, amundiOpts...)
		},
	),
	"ishares": newProviderFactory(
		ishares.SupportedRegions(),
		func(region string, options providerOptions) (etfscraper.Provider, error) {
			isharesOpts := []ishares.ClientOption{ishares.WithHTTPConfig(options.httpConfig)}
			if options.cacheTTL != nil {
				isharesOpts = append(isharesOpts, ishares.WithCacheTTL(*options.cacheTTL))
			}
			return ishares.New(region, isharesOpts...)
		},
	),
	"xtrackers": newProviderFactory(
		xtrackers.SupportedRegions(),
		func(region string, options providerOptions) (etfscraper.Provider, error) {
			xOpts := []xtrackers.ClientOption{xtrackers.WithHTTPConfig(options.httpConfig)}
			if options.cacheTTL != nil {
				xOpts = append(xOpts, xtrackers.WithCacheTTL(*options.cacheTTL))
			}
			return xtrackers.New(region, xOpts...)
		},
	),
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
	factory, region, err := validateProviderSelection(name, region)
	if err != nil {
		return nil, err
	}

	options := providerOptions{
		httpConfig: etfscraper.DefaultHTTPConfig(),
	}
	for _, opt := range opts {
		opt(&options)
	}

	return factory.constructor(region, options)
}

func validateProviderSelection(name, region string) (providerFactory, string, error) {
	normalizedName := strings.ToLower(strings.TrimSpace(name))
	if normalizedName == "" {
		return providerFactory{}, "", fmt.Errorf("provider name is required")
	}

	normalizedRegion := strings.ToLower(strings.TrimSpace(region))
	if normalizedRegion == "" {
		return providerFactory{}, "", fmt.Errorf("provider region is required")
	}

	factory, ok := providerRegistry[normalizedName]
	if !ok {
		return providerFactory{}, "", fmt.Errorf("unknown provider: %s", normalizedName)
	}
	if _, ok := factory.regionSet[normalizedRegion]; !ok {
		return providerFactory{}, "", fmt.Errorf("unsupported region %q for provider %q", normalizedRegion, normalizedName)
	}

	return factory, normalizedRegion, nil
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
	specs := make([]ProviderSpec, 0, len(providerRegistry))
	for name, factory := range providerRegistry {
		regions := append([]string(nil), factory.regions...)
		sort.Strings(regions)
		specs = append(specs, ProviderSpec{Name: name, Regions: regions})
	}

	sort.Slice(specs, func(i, j int) bool {
		return specs[i].Name < specs[j].Name
	})

	return specs
}

func newProviderFactory(
	regions []string,
	constructor func(region string, options providerOptions) (etfscraper.Provider, error),
) providerFactory {
	regionList := append([]string(nil), regions...)
	return providerFactory{
		regions:     regionList,
		regionSet:   normalizedRegionSet(regionList),
		constructor: constructor,
	}
}

func normalizedRegionSet(regions []string) map[string]struct{} {
	set := make(map[string]struct{}, len(regions))
	for _, region := range regions {
		normalizedRegion := strings.ToLower(strings.TrimSpace(region))
		if normalizedRegion == "" {
			continue
		}
		set[normalizedRegion] = struct{}{}
	}
	return set
}
