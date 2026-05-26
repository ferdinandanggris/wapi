package webhook_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ferdinandanggris/wapi/types"
	"github.com/ferdinandanggris/wapi/webhook"
)

func makeTextPayload(body string) types.WebhookPayload {
	return types.WebhookPayload{
		Object: "whatsapp_business_account",
		Entry: []*types.WebhookEntry{{
			ID: "waba-123",
			Changes: []*types.WebhookChange{{
				Field: "messages",
				Value: &types.WebhookValue{
					MessagingProduct: "whatsapp",
					Metadata:         &types.Metadata{DisplayPhoneNumber: "+16505555555", PhoneNumberID: "pnid-123"},
					Contacts:         []*types.WaContact{{Profile: &types.Profile{Name: "John Doe"}, WaID: "16505551234"}},
					Messages: []*types.IncomingMsg{{
						From: "16505551234", ID: "wamid.abc", Timestamp: "1683229471", Type: "text",
						Text: &types.IncomingText{Body: body},
					}},
				},
			}},
		}},
	}
}

func makeImagePayload() types.WebhookPayload {
	return types.WebhookPayload{
		Object: "whatsapp_business_account",
		Entry: []*types.WebhookEntry{{
			ID: "waba-123",
			Changes: []*types.WebhookChange{{
				Field: "messages",
				Value: &types.WebhookValue{
					MessagingProduct: "whatsapp",
					Metadata:         &types.Metadata{DisplayPhoneNumber: "+16505555555", PhoneNumberID: "pnid-123"},
					Contacts:         []*types.WaContact{{Profile: &types.Profile{Name: "John"}, WaID: "16505551234"}},
					Messages: []*types.IncomingMsg{{
						From: "16505551234", ID: "wamid.img", Timestamp: "1683229471", Type: "image",
						Image: &types.IncomingMedia{ID: "media-123", MimeType: "image/jpeg"},
					}},
				},
			}},
		}},
	}
}

func makeInteractivePayload(interactiveType string) types.WebhookPayload {
	return types.WebhookPayload{
		Object: "whatsapp_business_account",
		Entry: []*types.WebhookEntry{{
			ID: "waba-123",
			Changes: []*types.WebhookChange{{
				Field: "messages",
				Value: &types.WebhookValue{
					MessagingProduct: "whatsapp",
					Metadata:         &types.Metadata{DisplayPhoneNumber: "+16505555555", PhoneNumberID: "pnid-123"},
					Contacts:         []*types.WaContact{{Profile: &types.Profile{Name: "John"}, WaID: "16505551234"}},
					Messages: []*types.IncomingMsg{{
						From: "16505551234", ID: "wamid.int", Timestamp: "1683229471", Type: "interactive",
						Interactive: &types.IncomingInteractive{
							Type:          interactiveType,
							InButtonReply: &types.IncomingButtonReply{ID: "btn_yes", Title: "Yes"},
						},
					}},
				},
			}},
		}},
	}
}

func makeStatusPayload(status string) types.WebhookPayload {
	return types.WebhookPayload{
		Object: "whatsapp_business_account",
		Entry: []*types.WebhookEntry{{
			ID: "waba-123",
			Changes: []*types.WebhookChange{{
				Field: "messages",
				Value: &types.WebhookValue{
					MessagingProduct: "whatsapp",
					Metadata:         &types.Metadata{DisplayPhoneNumber: "+16505555555", PhoneNumberID: "pnid-123"},
					Statuses: []*types.StatusUpdate{{
						ID: "wamid.abc", Status: status, Timestamp: "1683229471", RecipientID: "16505551234",
						Conversation: &types.StatusConversation{
							ID: "conv-123",
							Origin: &types.ConversationOrigin{Type: "business_initiated"},
						},
						Pricing: &types.StatusPricing{Billable: true, PricingModel: "CBP", Category: "marketing"},
					}},
				},
			}},
		}},
	}
}

func makeFailedStatusPayload() types.WebhookPayload {
	return types.WebhookPayload{
		Object: "whatsapp_business_account",
		Entry: []*types.WebhookEntry{{
			ID: "waba-123",
			Changes: []*types.WebhookChange{{
				Field: "messages",
				Value: &types.WebhookValue{
					MessagingProduct: "whatsapp",
					Metadata:         &types.Metadata{DisplayPhoneNumber: "+16505555555", PhoneNumberID: "pnid-123"},
					Statuses: []*types.StatusUpdate{{
						ID: "wamid.abc", Status: "failed", Timestamp: "1683229471", RecipientID: "16505551234",
						Errors: []*types.StatusError{{
							Code: 131047, Title: "Re-engagement message",
							Message: "More than 24 hours have passed",
							ErrorData: &types.StatusErrorData{Details: "Send a template message instead"},
						}},
					}},
				},
			}},
		}},
	}
}

func TestVerifyWebhook(t *testing.T) {
	tests := []struct {
		name  string
		query string
		code  int
		body  string
	}{
		{"valid token", "/webhook?hub.mode=subscribe&hub.verify_token=mytoken&hub.challenge=ch123", http.StatusOK, "ch123"},
		{"invalid token", "/webhook?hub.mode=subscribe&hub.verify_token=wrong&hub.challenge=ch123", http.StatusForbidden, ""},
		{"wrong mode", "/webhook?hub.mode=unsubscribe&hub.verify_token=mytoken&hub.challenge=ch123", http.StatusForbidden, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &webhook.Handler{VerifyToken: "mytoken"}
			req := httptest.NewRequest("GET", tt.query, nil)
			w := httptest.NewRecorder()
			h.ServeHTTP(w, req)

			if w.Code != tt.code {
				t.Fatalf("expected %d, got %d", tt.code, w.Code)
			}
			if tt.body != "" && w.Body.String() != tt.body {
				t.Errorf("expected '%s', got '%s'", tt.body, w.Body.String())
			}
		})
	}
}

func TestHandleIncomingTextMessage(t *testing.T) {
	var gotMsg *types.IncomingMsg
	var gotMeta *types.Metadata

	h := &webhook.Handler{
		OnMessage: func(msg *types.IncomingMsg, meta *types.Metadata, contact *types.WaContact) error {
			gotMsg = msg
			gotMeta = meta
			return nil
		},
	}

	payload := makeTextPayload("Hello!")
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/webhook", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if gotMsg == nil {
		t.Fatal("expected OnMessage to be called")
	}
	if gotMsg.Type != "text" {
		t.Errorf("expected text, got %s", gotMsg.Type)
	}
	if gotMsg.Text.Body != "Hello!" {
		t.Errorf("expected 'Hello!', got '%s'", gotMsg.Text.Body)
	}
	if gotMeta.PhoneNumberID != "pnid-123" {
		t.Errorf("expected pnid-123, got %s", gotMeta.PhoneNumberID)
	}
}

func TestHandleIncomingImageMessage(t *testing.T) {
	var gotMsg *types.IncomingMsg
	h := &webhook.Handler{
		OnMessage: func(msg *types.IncomingMsg, meta *types.Metadata, contact *types.WaContact) error {
			gotMsg = msg
			return nil
		},
	}

	payload := makeImagePayload()
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/webhook", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	_ = w

	if gotMsg.Type != "image" {
		t.Errorf("expected image, got %s", gotMsg.Type)
	}
	if gotMsg.Image.ID != "media-123" {
		t.Errorf("expected media-123, got %s", gotMsg.Image.ID)
	}
}

func TestHandleInteractiveButtonReply(t *testing.T) {
	var gotMsg *types.IncomingMsg
	h := &webhook.Handler{
		OnMessage: func(msg *types.IncomingMsg, meta *types.Metadata, contact *types.WaContact) error {
			gotMsg = msg
			return nil
		},
	}

	payload := makeInteractivePayload("button_reply")
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/webhook", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if gotMsg.Interactive.Type != "button_reply" {
		t.Errorf("expected button_reply, got %s", gotMsg.Interactive.Type)
	}
	if gotMsg.Interactive.InButtonReply.ID != "btn_yes" {
		t.Errorf("expected btn_yes, got %s", gotMsg.Interactive.InButtonReply.ID)
	}
}

func TestHandleStatusUpdate(t *testing.T) {
	var gotStatus *types.StatusUpdate
	h := &webhook.Handler{
		OnStatus: func(status *types.StatusUpdate, meta *types.Metadata) error {
			gotStatus = status
			return nil
		},
	}

	payload := makeStatusPayload("delivered")
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/webhook", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if gotStatus == nil {
		t.Fatal("expected OnStatus to be called")
	}
	if gotStatus.Status != "delivered" {
		t.Errorf("expected delivered, got %s", gotStatus.Status)
	}
	if gotStatus.Pricing.Category != "marketing" {
		t.Errorf("expected marketing, got %s", gotStatus.Pricing.Category)
	}
}

func TestHandleFailedStatus(t *testing.T) {
	var gotStatus *types.StatusUpdate
	h := &webhook.Handler{
		OnStatus: func(status *types.StatusUpdate, meta *types.Metadata) error {
			gotStatus = status
			return nil
		},
	}

	payload := makeFailedStatusPayload()
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/webhook", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if len(gotStatus.Errors) == 0 {
		t.Fatal("expected errors")
	}
	if gotStatus.Errors[0].Code != 131047 {
		t.Errorf("expected 131047, got %d", gotStatus.Errors[0].Code)
	}
}

func TestHandleWebhookReturns200Immediately(t *testing.T) {
	called := false
	h := &webhook.Handler{
		OnMessage: func(msg *types.IncomingMsg, meta *types.Metadata, contact *types.WaContact) error {
			called = true
			return nil
		},
	}

	payload := makeTextPayload("Hello")
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/webhook", bytes.NewReader(body))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if !called {
		t.Error("expected OnMessage to be called")
	}
}

func TestMethodNotAllowed(t *testing.T) {
	h := &webhook.Handler{}
	req := httptest.NewRequest("PUT", "/webhook", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", w.Code)
	}
}
