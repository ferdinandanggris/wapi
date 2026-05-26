package cloud_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/ferdinandanggris/wapi/cloud"
	"github.com/ferdinandanggris/wapi/types"
)

func TestSendTextMessage(t *testing.T) {
	ms := newDefaultMockServer()
	defer ms.Close()

	c := ms.client()
	ctx := context.Background()

	resp, err := c.SendMessage(ctx, "123", types.NewTextMessage("+16505551234", "Hello!", false))
	if err != nil {
		t.Fatalf("SendMessage failed: %v", err)
	}

	if len(resp.Messages) == 0 {
		t.Fatal("expected at least one message in response")
	}
	if resp.Messages[0].ID == "" {
		t.Error("expected non-empty message ID")
	}
}

func TestSendImageMessage(t *testing.T) {
	ms := newDefaultMockServer()
	defer ms.Close()

	c := ms.client()
	_, err := c.SendMessage(context.Background(), "123", types.NewImageMessage("+16505551234", "img-id-123", "caption"))
	if err != nil {
		t.Fatalf("SendMessage image failed: %v", err)
	}
}

func TestSendImageByLink(t *testing.T) {
	ms := newDefaultMockServer()
	defer ms.Close()

	c := ms.client()
	_, err := c.SendMessage(context.Background(), "123", types.NewImageByLink("+16505551234", "https://example.com/img.jpg", "caption"))
	if err != nil {
		t.Fatalf("SendMessage image link failed: %v", err)
	}
}

func TestSendTemplateMessage(t *testing.T) {
	ms := newDefaultMockServer()
	defer ms.Close()

	c := ms.client()
	msg := types.NewTemplateMessage("+16505551234", "hello_world", "en_US",
		types.NewBodyComponent(types.NewTextParameter("John")),
	)
	_, err := c.SendMessage(context.Background(), "123", msg)
	if err != nil {
		t.Fatalf("SendMessage template failed: %v", err)
	}
}

func TestSendInteractiveButton(t *testing.T) {
	ms := newDefaultMockServer()
	defer ms.Close()

	c := ms.client()
	msg := types.NewInteractiveButton("+16505551234", "Confirm your order",
		types.NewButton("yes", "Yes"),
		types.NewButton("no", "No"),
	)
	_, err := c.SendMessage(context.Background(), "123", msg)
	if err != nil {
		t.Fatalf("SendMessage interactive button failed: %v", err)
	}
}

func TestSendInteractiveList(t *testing.T) {
	ms := newDefaultMockServer()
	defer ms.Close()

	c := ms.client()
	msg := types.NewInteractiveList("+16505551234", "View Options", "Choose an option",
		types.NewSection("Section 1",
			types.NewRow("opt1", "Option 1", "Description 1"),
			types.NewRow("opt2", "Option 2", "Description 2"),
		),
	)
	_, err := c.SendMessage(context.Background(), "123", msg)
	if err != nil {
		t.Fatalf("SendMessage interactive list failed: %v", err)
	}
}

func TestSendLocationMessage(t *testing.T) {
	ms := newDefaultMockServer()
	defer ms.Close()

	c := ms.client()
	_, err := c.SendMessage(context.Background(), "123", types.NewLocationMessage("+16505551234", 18.4861, -69.9312, "Office", "Santo Domingo"))
	if err != nil {
		t.Fatalf("SendMessage location failed: %v", err)
	}
}

func TestSendContactMessage(t *testing.T) {
	ms := newDefaultMockServer()
	defer ms.Close()

	c := ms.client()
	contact := &types.Contact{
		Name: &types.Name{FormattedName: "John Doe", FirstName: "John"},
		Phones: []*types.Phone{{Phone: "+16505551234", Type: "WORK"}},
	}
	_, err := c.SendMessage(context.Background(), "123", types.NewContactMessage("+16505551234", contact))
	if err != nil {
		t.Fatalf("SendMessage contact failed: %v", err)
	}
}

func TestSendReactionMessage(t *testing.T) {
	ms := newDefaultMockServer()
	defer ms.Close()

	c := ms.client()
	_, err := c.SendMessage(context.Background(), "123", types.NewReactionMessage("+16505551234", "wamid.abc", "\U0001f44d"))
	if err != nil {
		t.Fatalf("SendMessage reaction failed: %v", err)
	}
}

func TestRemoveReaction(t *testing.T) {
	ms := newDefaultMockServer()
	defer ms.Close()

	c := ms.client()
	_, err := c.SendMessage(context.Background(), "123", types.NewRemoveReactionMessage("+16505551234", "wamid.abc"))
	if err != nil {
		t.Fatalf("Remove reaction failed: %v", err)
	}
}

func TestSendMessageWithContext(t *testing.T) {
	ms := newDefaultMockServer()
	defer ms.Close()

	c := ms.client()
	msg := types.NewTextMessage("+16505551234", "Reply message", false).WithContext("wamid.original-msg")
	_, err := c.SendMessage(context.Background(), "123", msg)
	if err != nil {
		t.Fatalf("SendMessage with context failed: %v", err)
	}
}

func TestMarkAsRead(t *testing.T) {
	ms := newDefaultMockServer()
	defer ms.Close()

	c := ms.client()
	err := c.MarkAsRead(context.Background(), "123", "wamid.abc")
	if err != nil {
		t.Fatalf("MarkAsRead failed: %v", err)
	}
}

func TestSendMessageAPIError(t *testing.T) {
	ms := newMockServer()
	defer ms.Close()

	ms.on("POST", "/123/messages", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, 131030, "Recipient not valid WhatsApp user", 0)
	})

	c := ms.client()
	_, err := c.SendMessage(context.Background(), "123", types.NewTextMessage("+16505551234", "test", false))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestSendMessageToEmpty(t *testing.T) {
	ms := newDefaultMockServer()
	defer ms.Close()

	c := ms.client()
	_, err := c.SendMessage(context.Background(), "123", &types.Message{Type: "text"})
	if err == nil {
		t.Fatal("expected error for empty recipient")
	}
}

func TestSendMessageRateLimitError(t *testing.T) {
	ms := newMockServer()
	defer ms.Close()

	attempts := 0
	ms.on("POST", "/123/messages", func(w http.ResponseWriter, r *http.Request) {
		attempts++
		writeError(w, 130429, "Rate limit hit", 0)
	})

	c := ms.client(
		cloud.WithRetry(2),
	)
	_, err := c.SendMessage(context.Background(), "123", types.NewTextMessage("+16505551234", "test", false))
	if err == nil {
		t.Fatal("expected error")
	}
	if attempts != 2 {
		t.Errorf("expected 2 attempts, got %d", attempts)
	}
}

func TestSendCTAUrlMessage(t *testing.T) {
	ms := newDefaultMockServer()
	defer ms.Close()

	c := ms.client()
	msg := types.NewInteractiveCTA("+16505551234", "Track Order", "https://example.com/track/123", "Your order is ready")
	_, err := c.SendMessage(context.Background(), "123", msg)
	if err != nil {
		t.Fatalf("SendMessage CTA failed: %v", err)
	}
}
