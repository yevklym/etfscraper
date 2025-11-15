package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/yevklym/etfscraper"
	"github.com/yevklym/etfscraper/internal/providers/ishares"
)

func main() {
	provider, err := getProvider("ishares:us")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Discovering ETFs...")
	funds, err := provider.DiscoverETFs(context.Background())
	if err != nil {
		log.Fatal("Discovery failed:", err)
	}

	fmt.Printf("Successfully discovered %d ETFs!\n\n", len(funds))

	// Show first 5 ETFs
	for i, fund := range funds[:min(5, len(funds))] {
		fmt.Printf("%d. %s (%s)\n", i+1, fund.Name, fund.Ticker)
		fmt.Printf("   ISIN: %s\n", fund.ISIN)
		fmt.Printf("   Expense Ratio: %.2f%%\n", fund.ExpenseRatio*100)
		fmt.Printf("   Assets: $%.1fB\n", fund.TotalAssets/1_000_000_000)
		if fund.InceptionDate != nil {
			fmt.Printf("   Inception: %s\n", fund.InceptionDate.Format("Jan 2, 2006"))
		}
		fmt.Println()
	}
}

func getProvider(providerName string) (etfscraper.Provider, error) {
	parts := strings.SplitN(providerName, ":", 2)
	name := parts[0]
	region := ""
	if len(parts) > 1 {
		region = parts[1]
	}

	switch name {
	case "ishares":
		return ishares.New(region), nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", name)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
