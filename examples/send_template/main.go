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

	msg := types.NewTemplateMessage(to, "order_confirmation", "en_US",
		types.NewHeaderComponent(types.NewTextParameter("ORD-12345")),
		types.NewBodyComponent(
			types.NewTextParameter("John"),
			types.NewTextParameter("ORD-12345"),
			types.NewTextParameter("49.99"),
		),
		types.NewURLButtonComponent(0, "ORD-12345"),
	)

	resp, err := client.SendMessage(ctx, phoneID, msg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "send template failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("template sent! message id: %s\n", resp.Messages[0].ID)
}
