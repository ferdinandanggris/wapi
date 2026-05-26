package cloud_test

import (
	"bytes"
	"context"
	"net/http"
	"testing"
)

func TestUploadMedia(t *testing.T) {
	ms := newDefaultMockServer()
	defer ms.Close()

	c := ms.client()
	content := bytes.NewReader([]byte("fake-image-data"))
	resp, err := c.UploadMedia(context.Background(), "123", "photo.jpg", content, "image/jpeg")
	if err != nil {
		t.Fatalf("UploadMedia failed: %v", err)
	}
	if resp.ID == "" {
		t.Error("expected non-empty media ID")
	}
}

func TestGetMediaURL(t *testing.T) {
	ms := newDefaultMockServer()
	defer ms.Close()

	c := ms.client()
	info, err := c.GetMediaURL(context.Background(), "media-id-123")
	if err != nil {
		t.Fatalf("GetMediaURL failed: %v", err)
	}
	if info.URL == "" {
		t.Error("expected non-empty URL")
	}
	if info.MimeType != "image/jpeg" {
		t.Errorf("expected image/jpeg, got %s", info.MimeType)
	}
}

func TestDownloadMedia(t *testing.T) {
	ms := newDefaultMockServer()
	defer ms.Close()

	ms.on("GET", "/download/media-id-123", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("binary-data"))
	})

	c := ms.client()
	rc, err := c.DownloadMedia(context.Background(), "media-id-123")
	if err != nil {
		t.Fatalf("DownloadMedia failed: %v", err)
	}
	defer rc.Close()

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(rc)
	if buf.String() != "binary-data" {
		t.Errorf("expected binary-data, got %s", buf.String())
	}
}

func TestDeleteMedia(t *testing.T) {
	ms := newDefaultMockServer()
	defer ms.Close()

	c := ms.client()
	err := c.DeleteMedia(context.Background(), "media-id-123")
	if err != nil {
		t.Fatalf("DeleteMedia failed: %v", err)
	}
}
