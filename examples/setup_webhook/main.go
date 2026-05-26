package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ferdinandanggris/wapi/cloud"
)

func main() {
	ctx := context.Background()

	client := cloud.New(
		cloud.WithAccessToken(os.Getenv("WABA_TOKEN")),
	)

	wabaID := os.Getenv("WHATSAPP_WABA_ID")
	appID := os.Getenv("APP_ID")

	fmt.Println("=== Subscribing WABA to webhooks ===")
	app, err := client.SubscribeToWebhooks(ctx, wabaID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "subscribe failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("subscribed! app id: %s, name: %s\n", app.ID, app.Name)

	fmt.Println("\n=== Setting callback URL ===")
	callbackURL := "https://c906-125-164-232-173.ngrok-free.app/webhook"
	verifyToken := "verify_token"

	err = client.SetWebhookCallback(ctx, appID, callbackURL, verifyToken,
		"messages", "message_template_status")
	if err != nil {
		fmt.Fprintf(os.Stderr, "set callback failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("callback set to %s\n", callbackURL)

	fmt.Println("\n=== Done! ===")
}
