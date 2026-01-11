package ishares_test

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/yevklym/etfscraper/internal/providers/ishares"
)

func ExampleNew() {
	client, err := ishares.New("de")
	if err != nil {
		log.Fatal(err)
	}

	funds, _ := client.DiscoverETFs(context.Background())
	fmt.Printf("Found %d ETFs\n", len(funds))
}

func ExampleNew_withOptions() {
	client, err := ishares.New("de",
		ishares.WithTimeout(30*time.Second),
		ishares.WithDebug(true),
	)
	if err != nil {
		log.Fatal(err)
	}

	_ = client
}
