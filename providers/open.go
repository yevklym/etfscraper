package providers

import (
	"fmt"
	"sort"
	"strings"

	"github.com/yevklym/etfscraper"
	"github.com/yevklym/etfscraper/internal/providers/amundi"
	"github.com/yevklym/etfscraper/internal/providers/ishares"
)

type ProviderSpec struct {
	Name    string
	Regions []string
}

type Spec struct {
	Name   string
	Region string
}

func Open(spec string, opts ...Option) (etfscraper.Provider, error) {
	name, region, err := ParseProviderSpec(spec)
	if err != nil {
		return nil, err
	}

	return OpenNameRegion(name, region, opts...)
}

func OpenSpec(spec Spec, opts ...Option) (etfscraper.Provider, error) {
	return OpenNameRegion(spec.Name, spec.Region, opts...)
}

func OpenNameRegion(name, region string, opts ...Option) (etfscraper.Provider, error) {
	options := providerOptions{
		httpConfig: etfscraper.DefaultHTTPConfig(),
	}
	for _, opt := range opts {
		opt(&options)
	}

	switch strings.ToLower(name) {
	case "ishares":
		return ishares.New(region, ishares.WithHTTPConfig(options.httpConfig))
	case "amundi":
		return amundi.New(region, amundi.WithHTTPConfig(options.httpConfig))
	default:
		return nil, fmt.Errorf("unknown provider: %s", name)
	}
}

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
