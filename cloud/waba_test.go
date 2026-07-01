package cloud_test

import (
	"context"
	"net/http"
	"testing"

	wapi "github.com/ferdinandanggris/wapi"
)

func TestListWhatsAppBusinessAccounts(t *testing.T) {
	ms := newDefaultMockServer()
	defer ms.Close()

	c := ms.client()
	accounts, err := c.ListWhatsAppBusinessAccounts(context.Background(), "business-789")
	if err != nil {
		t.Fatalf("ListWhatsAppBusinessAccounts failed: %v", err)
	}

	if len(accounts.Data) != 2 {
		t.Fatalf("expected 2 accounts, got %d", len(accounts.Data))
	}

	if accounts.Data[0].Name != "Seller Pulsa" {
		t.Errorf("expected 'Seller Pulsa', got '%s'", accounts.Data[0].Name)
	}
	if accounts.Data[1].Name != "Family Pulsa" {
		t.Errorf("expected 'Family Pulsa', got '%s'", accounts.Data[1].Name)
	}
	if accounts.Paging == nil || accounts.Paging.Cursors == nil {
		t.Fatal("expected paging with cursors")
	}
	if accounts.Paging.Cursors.After != "cursor-after" {
		t.Errorf("expected 'cursor-after', got '%s'", accounts.Paging.Cursors.After)
	}
}

func TestListWhatsAppBusinessAccounts_Pagination(t *testing.T) {
	ms := newMockServer()
	defer ms.Close()

	var capturedLimit, capturedAfter string
	ms.on("GET", "/business-789/owned_whatsapp_business_accounts", func(w http.ResponseWriter, r *http.Request) {
		capturedLimit = r.URL.Query().Get("limit")
		capturedAfter = r.URL.Query().Get("after")
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"data": []map[string]interface{}{},
		})
	})

	c := ms.client()
	_, err := c.ListWhatsAppBusinessAccounts(context.Background(), "business-789",
		wapi.WithLimit(10),
		wapi.WithOffset("cursor-next"),
	)
	if err != nil {
		t.Fatalf("ListWhatsAppBusinessAccounts failed: %v", err)
	}

	if capturedLimit != "10" {
		t.Errorf("expected limit=10, got '%s'", capturedLimit)
	}
	if capturedAfter != "cursor-next" {
		t.Errorf("expected after=cursor-next, got '%s'", capturedAfter)
	}
}
