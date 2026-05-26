package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

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
	filePath := os.Args[1]

	f, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	fmt.Println("uploading media...")
	mimeType := "image/jpeg"
	if ext := filepath.Ext(filePath); ext == ".png" {
		mimeType = "image/png"
	} else if ext == ".mp4" {
		mimeType = "video/mp4"
	}

	upload, err := client.UploadMedia(ctx, phoneID, filepath.Base(filePath), f, mimeType)
	if err != nil {
		fmt.Fprintf(os.Stderr, "upload failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("uploaded! media id: %s\n", upload.ID)

	fmt.Println("sending by media ID...")
	msg := types.NewImageMessage(to, upload.ID, "Here's your image")
	resp, err := client.SendMessage(ctx, phoneID, msg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "send failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("sent! message id: %s\n", resp.Messages[0].ID)
}
