package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ferdinandanggris/wapi/types"
	"github.com/ferdinandanggris/wapi/webhook"
)

func main() {
	handler := &webhook.Handler{
		VerifyToken: os.Getenv("WEBHOOK_VERIFY_TOKEN"),
		AppSecret:   os.Getenv("WABA_APP_SECRET"),
		Logger:      log.Default(),
		OnMessage: func(msg *types.IncomingMsg, meta *types.Metadata, contact *types.WaContact) error {
			name := "unknown"
			if contact != nil {
				name = contact.Profile.Name
			}

			fmt.Printf("[%s][%s] %s: %s %s\n", meta.DisplayPhoneNumber, msg.ID, name, msg.Type, msg.Text.Body)
			switch msg.Type {
			case "image":
				fmt.Printf("[%s][%s] %s sent image: %s\n", meta.DisplayPhoneNumber, msg.ID, name, msg.Image.ID)
			default:
			}

			return nil
		},
		OnStatus: func(status *types.StatusUpdate, meta *types.Metadata) error {
			for _, e := range status.Errors {
				fmt.Printf("[%s] message %s ERROR: %s\n", meta.DisplayPhoneNumber, status.ID, e.Message)
			}
			fmt.Printf("[%s] message %s: %s\n", meta.DisplayPhoneNumber, status.ID, status.Status)
			return nil
		},
	}

	mux := http.NewServeMux()
	mux.Handle("/webhook", handler)

	addr := ":8080"
	fmt.Printf("webhook server listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
