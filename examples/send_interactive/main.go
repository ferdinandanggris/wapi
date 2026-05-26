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

	fmt.Println("=== Reply Buttons ===")
	btnMsg := types.NewInteractiveButton(to, "Confirm your appointment?",
		types.NewButton("confirm", "Confirm"),
		types.NewButton("reschedule", "Reschedule"),
		types.NewButton("cancel", "Cancel"),
	)
	btnMsg.Interactive.
		WithHeader("text", "Appointment Reminder").
		WithFooter("Tap to respond")

	resp, err := client.SendMessage(ctx, phoneID, btnMsg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "send buttons failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("buttons sent! id: %s\n", resp.Messages[0].ID)

	fmt.Println("\n=== List Message ===")
	listMsg := types.NewInteractiveList(to, "View Services", "Choose a service:",
		types.NewSection("Consulting",
			types.NewRow("strategy", "Strategy Session", "1-hour consultation"),
			types.NewRow("audit", "Technical Audit", "Full architecture review"),
		),
		types.NewSection("Development",
			types.NewRow("mvp", "MVP Build", "4-week prototype"),
			types.NewRow("custom", "Custom Project", "Tailored development"),
		),
	)

	resp, err = client.SendMessage(ctx, phoneID, listMsg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "send list failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("list sent! id: %s\n", resp.Messages[0].ID)
}
