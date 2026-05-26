package cloud_test

import (
	"context"
	"net/http"
	"testing"
)

func TestSubscribeToWebhooks(t *testing.T) {
	ms := newMockServer()
	defer ms.Close()

	ms.on("POST", "/waba-456/subscribed_apps", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"id": "98765", "name": "My App"})
	})

	c := ms.client()
	app, err := c.SubscribeToWebhooks(context.Background(), "waba-456")
	if err != nil {
		t.Fatalf("SubscribeToWebhooks failed: %v", err)
	}
	if app.ID != "98765" {
		t.Errorf("expected 98765, got %s", app.ID)
	}
}

func TestUnsubscribeFromWebhooks(t *testing.T) {
	ms := newMockServer()
	defer ms.Close()

	deleted := false
	ms.on("DELETE", "/waba-456/subscribed_apps", func(w http.ResponseWriter, r *http.Request) {
		deleted = true
		writeJSON(w, http.StatusOK, map[string]bool{"success": true})
	})

	c := ms.client()
	err := c.UnsubscribeFromWebhooks(context.Background(), "waba-456")
	if err != nil {
		t.Fatalf("UnsubscribeFromWebhooks failed: %v", err)
	}
	if !deleted {
		t.Error("expected DELETE to be called")
	}
}

func TestGetWebhookSubscription(t *testing.T) {
	ms := newMockServer()
	defer ms.Close()

	ms.on("GET", "/waba-456/subscribed_apps", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"id": "98765", "name": "My App"})
	})

	c := ms.client()
	app, err := c.GetWebhookSubscription(context.Background(), "waba-456")
	if err != nil {
		t.Fatalf("GetWebhookSubscription failed: %v", err)
	}
	if app.ID != "98765" {
		t.Errorf("expected 98765, got %s", app.ID)
	}
}

func TestSetWebhookFields(t *testing.T) {
	ms := newMockServer()
	defer ms.Close()

	var fields string
	ms.on("POST", "/app-789/subscriptions", func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		fields = r.Form.Get("fields")
		writeJSON(w, http.StatusOK, map[string]bool{"success": true})
	})

	c := ms.client()
	err := c.SetWebhookFields(context.Background(), "app-789", "messages", "message_template_status")
	if err != nil {
		t.Fatalf("SetWebhookFields failed: %v", err)
	}
	if fields != "messages,message_template_status" {
		t.Errorf("expected 'messages,message_template_status', got '%s'", fields)
	}
}

func TestSetWebhookCallback(t *testing.T) {
	ms := newMockServer()
	defer ms.Close()

	var gotURL, gotToken, gotFields string
	ms.on("POST", "/app-789/subscriptions", func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		gotURL = r.Form.Get("callback_url")
		gotToken = r.Form.Get("verify_token")
		gotFields = r.Form.Get("fields")
		writeJSON(w, http.StatusOK, map[string]bool{"success": true})
	})

	c := ms.client()
	err := c.SetWebhookCallback(context.Background(), "app-789",
		"https://example.com/webhook", "mytoken",
		"messages", "message_template_status")
	if err != nil {
		t.Fatalf("SetWebhookCallback failed: %v", err)
	}
	if gotURL != "https://example.com/webhook" {
		t.Errorf("expected callback url, got %s", gotURL)
	}
	if gotToken != "mytoken" {
		t.Errorf("expected mytoken, got %s", gotToken)
	}
	if gotFields != "messages,message_template_status" {
		t.Errorf("expected messages,message_template_status, got %s", gotFields)
	}
}
