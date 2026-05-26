package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ferdinandanggris/wapi/cloud"
	"github.com/ferdinandanggris/wapi/types"
)

func main() {
	ctx := context.Background()

	client := cloud.New(
		cloud.WithAccessToken(os.Getenv("WABA_TOKEN")),
	)

	to := os.Getenv("TO_NUMBER")
	phoneID := os.Getenv("PHONE_NUMBER_ID")

	msg := types.NewTextMessage(to, "Hello from wapi!", false)
	resp, err := client.SendMessage(ctx, phoneID, msg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "send failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("sent! message id: %s\n", resp.Messages[0].ID)
}
