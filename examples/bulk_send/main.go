package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/ferdinandanggris/wapi/cloud"
	"github.com/ferdinandanggris/wapi/transport"
	"github.com/ferdinandanggris/wapi/types"
)

func main() {
	ctx := context.Background()
	to := os.Getenv("TO_NUMBER")

	if to == "" {
		log.Fatal("TO_NUMBER is required")
	}

	client := cloud.New(
		cloud.WithAccessToken(os.Getenv("WABA_TOKEN")),
		cloud.WithHTTPClient(&http.Client{
			Transport: transport.Chain(
				http.DefaultTransport,
				transport.RateLimit(80, 80),
				transport.DefaultRetry(),
			),
		}),
	)

	phoneID := os.Getenv("PHONE_NUMBER_ID")
	recipients := []string{to}

	var wg sync.WaitGroup
	start := time.Now()

	for _, recipient := range recipients {
		wg.Add(1)
		go func(r string) {
			defer wg.Done()

			msg := types.NewTextMessage(r, "Bulk message test", false)
			resp, err := client.SendMessage(ctx, phoneID, msg)
			if err != nil {
				log.Printf("FAIL %s: %v", r, err)
				return
			}
			log.Printf("OK %s: %s", r, resp.Messages[0].ID)
		}(recipient)
	}

	wg.Wait()
	fmt.Printf("done in %s\n", time.Since(start))
}
