// Command example demonstrates the etfscraper SDK by discovering ETFs,
// looking up a specific fund, and fetching its holdings.
//
// Usage:
//
//	go run ./cmd/example
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/yevklym/etfscraper/providers"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Create a provider for iShares US region.
	provider, err := providers.Open("ishares:us")
	if err != nil {
		log.Fatal(err)
	}

	// Discover all ETFs offered by the provider.
	funds, err := provider.DiscoverETFs(ctx)
	if err != nil {
		log.Fatal("Discovery failed:", err)
	}
	fmt.Printf("Discovered %d ETFs\n\n", len(funds))

	// Show the first 5 discovered funds.
	for i, fund := range funds[:min(5, len(funds))] {
		fmt.Printf("%d. %s (%s) — %s\n", i+1, fund.Name, fund.Ticker, fund.ISIN)
	}
	fmt.Println()

	// Look up a specific fund by ticker and print its metadata.
	fund, err := provider.FundInfo(ctx, "IVV")
	if err != nil {
		log.Fatalf("FundInfo failed: %v", err)
	}
	fmt.Printf("Fund: %s (%s)\n", fund.Name, fund.Ticker)
	fmt.Printf("  ISIN:          %s\n", fund.ISIN)
	fmt.Printf("  Currency:      %s\n", fund.Currency)
	fmt.Printf("  Exchange:      %s\n", fund.Exchange)
	fmt.Printf("  Asset Class:   %s\n", fund.AssetClass)
	fmt.Printf("  Expense Ratio: %.2f%%\n", fund.ExpenseRatio*100)
	if fund.InceptionDate != nil {
		fmt.Printf("  Inception:     %s\n", fund.InceptionDate.Format("Jan 2, 2006"))
	}
	fmt.Println()

	// Fetch holdings for the fund.
	snapshot, err := provider.HoldingsForFund(ctx, fund)
	if err != nil {
		log.Fatalf("Holdings failed: %v", err)
	}
	fmt.Printf("Holdings as of %s (%d total):\n", snapshot.AsOfDate.Format("2006-01-02"), snapshot.TotalHoldings)
	for i, h := range snapshot.Holdings[:min(10, len(snapshot.Holdings))] {
		fmt.Printf("  %2d. %-30s %6.2f%%  $%14.2f\n", i+1, h.Name, h.Weight*100, h.MarketValue)
	}
}
