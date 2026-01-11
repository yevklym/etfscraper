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
	provider, err := getProvider("ishares:de")
	if err != nil {
		log.Fatal(err)
	}

	// --- Example Usage of FundInfo() and Holdings() for a specific ETF ---
	fmt.Println("--- Demonstrating specific fund lookup and holdings ---")
	exampleTicker := "QDVE" // iShares S&P 500 Information Technology Sector UCITS ETF

	fmt.Printf("Fetching FundInfo for %s...\n", exampleTicker)
	specificFund, err := provider.FundInfo(context.Background(), exampleTicker)
	if err != nil {
		log.Printf("Failed to get FundInfo for %s: %v\n", exampleTicker, err)
	} else {
		fmt.Printf("Found specific fund: %s (%s)\n", specificFund.Name, specificFund.Ticker)
		fmt.Printf("   ISIN: %s\n", specificFund.ISIN)
		fmt.Printf("   Expense Ratio: %.2f%%\n", specificFund.ExpenseRatio*100)
		fmt.Printf("   Assets: $%.1fB\n", specificFund.TotalAssets/1_000_000_000)
		if specificFund.InceptionDate != nil {
			fmt.Printf("   Inception: %s\n", specificFund.InceptionDate.Format("Jan 2, 2006"))
		}
		fmt.Println()

		fmt.Printf("Fetching Holdings for %s...\n", exampleTicker)
		specificHoldings, err := provider.Holdings(context.Background(), exampleTicker)
		if err != nil {
			log.Printf("Failed to get Holdings for %s: %v\n", exampleTicker, err)
		} else {
			fmt.Printf("Holdings as of: %s\n", specificHoldings.AsOfDate.Format("Jan 2, 2006"))
			fmt.Printf("Total holdings: %d\n\n", specificHoldings.TotalHoldings)

			fmt.Println("Top 3 Holdings:")
			for k, holding := range specificHoldings.Holdings[:min(3, len(specificHoldings.Holdings))] {
				fmt.Printf("   %d. %s (%.2f%%) - $%.2f\n", k+1, holding.Name, holding.Weight*100, holding.MarketValue)
			}
			fmt.Println("-----------------------------------------------------")
		}
	}

	// --- Example Usage of DiscoverETFs() ---
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

		holdingsSnapshot, err := provider.Holdings(context.Background(), fund.Ticker)
		if err != nil {
			log.Printf("Failed to get holdings for %s (%s): %v\n", fund.Name, fund.Ticker, err)
			continue
		}

		fmt.Printf("Holdings as of:  %s\n", holdingsSnapshot.AsOfDate.Format("Jan 2, 2006"))
		fmt.Printf("Total holdings: %d\n\n", holdingsSnapshot.TotalHoldings)

		fmt.Println("Top 5 Holdings:")
		for j, holding := range holdingsSnapshot.Holdings[:min(5, len(holdingsSnapshot.Holdings))] {
			fmt.Printf("%d. %s (%.2f%%) - $%.2f\n", j+1, holding.Name, holding.Weight*100, holding.MarketValue)
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
		return ishares.New(region)
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
