package cloud_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/ferdinandanggris/wapi/types"
)

func TestRegisterPhone(t *testing.T) {
	ms := newDefaultMockServer()
	defer ms.Close()

	c := ms.client()
	err := c.RegisterPhone(context.Background(), "123", "123456")
	if err != nil {
		t.Fatalf("RegisterPhone failed: %v", err)
	}
}

func TestDeregisterPhone(t *testing.T) {
	ms := newMockServer()
	defer ms.Close()

	deregistered := false
	ms.on("POST", "/123/deregister", func(w http.ResponseWriter, r *http.Request) {
		deregistered = true
		writeJSON(w, http.StatusOK, map[string]bool{"success": true})
	})

	c := ms.client()
	err := c.DeregisterPhone(context.Background(), "123")
	if err != nil {
		t.Fatalf("DeregisterPhone failed: %v", err)
	}
	if !deregistered {
		t.Error("expected deregister to be called")
	}
}

func TestGetPhoneNumber(t *testing.T) {
	ms := newDefaultMockServer()
	defer ms.Close()

	c := ms.client()
	pn, err := c.GetPhoneNumber(context.Background(), "123")
	if err != nil {
		t.Fatalf("GetPhoneNumber failed: %v", err)
	}
	if pn.DisplayPhoneNumber != "+16505555555" {
		t.Errorf("expected +16505555555, got %s", pn.DisplayPhoneNumber)
	}
	if pn.QualityRating != "GREEN" {
		t.Errorf("expected GREEN, got %s", pn.QualityRating)
	}
}

func TestListPhoneNumbers(t *testing.T) {
	ms := newDefaultMockServer()
	defer ms.Close()

	c := ms.client()
	numbers, err := c.ListPhoneNumbers(context.Background(), "waba-456")
	if err != nil {
		t.Fatalf("ListPhoneNumbers failed: %v", err)
	}
	if len(numbers) == 0 {
		t.Fatal("expected at least one phone number")
	}
	if numbers[0].DisplayPhoneNumber != "+16505555555" {
		t.Errorf("expected +16505555555, got %s", numbers[0].DisplayPhoneNumber)
	}
}

func TestSetTwoStepPIN(t *testing.T) {
	ms := newMockServer()
	defer ms.Close()

	var reqBody map[string]string
	ms.on("POST", "/123", func(w http.ResponseWriter, r *http.Request) {
		_ = parseBody(r, &reqBody)
		writeJSON(w, http.StatusOK, map[string]bool{"success": true})
	})

	c := ms.client()
	err := c.SetTwoStepPIN(context.Background(), "123", "654321")
	if err != nil {
		t.Fatalf("SetTwoStepPIN failed: %v", err)
	}
	if reqBody["pin"] != "654321" {
		t.Errorf("expected pin 654321, got %s", reqBody["pin"])
	}
}

func TestGetBusinessProfile(t *testing.T) {
	ms := newDefaultMockServer()
	defer ms.Close()

	c := ms.client()
	profile, err := c.GetBusinessProfile(context.Background(), "123")
	if err != nil {
		t.Fatalf("GetBusinessProfile failed: %v", err)
	}
	if profile.Description != "Test business profile" {
		t.Errorf("expected test description, got %s", profile.Description)
	}
}

func TestUpdateBusinessProfile(t *testing.T) {
	ms := newDefaultMockServer()
	defer ms.Close()

	c := ms.client()
	profile := &types.BusinessProfile{Description: "Updated description"}
	err := c.UpdateBusinessProfile(context.Background(), "123", profile)
	if err != nil {
		t.Fatalf("UpdateBusinessProfile failed: %v", err)
	}
}
