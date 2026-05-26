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

	appID := os.Getenv("APP_ID")
	callbackURL := os.Getenv("CALLBACK_URL")
	verifyToken := os.Getenv("VERIFY_TOKEN")
	if callbackURL == "" {
		callbackURL = "https://c906-125-164-232-173.ngrok-free.app/webhook"
	}
	if verifyToken == "" {
		verifyToken = "verify_token"
	}

	fmt.Println("=== Setting callback URL ===")
	fmt.Printf("callback: %s\n", callbackURL)
	fmt.Printf("verify token: %s\n", verifyToken)

	err := client.SetWebhookCallback(ctx, appID, callbackURL, verifyToken,
		"messages", "message_template_status")
	if err != nil {
		fmt.Fprintf(os.Stderr, "set callback failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("callback set successfully!")
}
