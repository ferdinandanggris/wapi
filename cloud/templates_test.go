package cloud_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/ferdinandanggris/wapi/types"
)

func TestCreateTemplate(t *testing.T) {
	ms := newDefaultMockServer()
	defer ms.Close()

	c := ms.client()
	tpl := &types.Template{
		Name:     "order_confirmation",
		Language: "en_US",
		Category: "utility",
		Components: []*types.TemplateComponent{
			{Type: "body", Text: "Hi {{1}}, your order is confirmed."},
		},
	}

	created, err := c.CreateTemplate(context.Background(), "waba-456", tpl)
	if err != nil {
		t.Fatalf("CreateTemplate failed: %v", err)
	}
	if created.Name != "hello_world" {
		t.Errorf("expected hello_world, got %s", created.Name)
	}
}

func TestListTemplates(t *testing.T) {
	ms := newDefaultMockServer()
	defer ms.Close()

	c := ms.client()
	list, err := c.ListTemplates(context.Background(), "waba-456")
	if err != nil {
		t.Fatalf("ListTemplates failed: %v", err)
	}
	if len(list.Data) == 0 {
		t.Fatal("expected at least one template")
	}
	if list.Data[0].Name != "hello_world" {
		t.Errorf("expected hello_world, got %s", list.Data[0].Name)
	}
}

func TestDeleteTemplate(t *testing.T) {
	ms := newMockServer()
	defer ms.Close()

	deleted := false
	ms.on("DELETE", "/template-123", func(w http.ResponseWriter, r *http.Request) {
		deleted = true
		writeJSON(w, http.StatusOK, map[string]bool{"success": true})
	})

	c := ms.client()
	err := c.DeleteTemplate(context.Background(), "template-123")
	if err != nil {
		t.Fatalf("DeleteTemplate failed: %v", err)
	}
	if !deleted {
		t.Error("expected delete to be called")
	}
}

func TestGetTemplate(t *testing.T) {
	ms := newMockServer()
	defer ms.Close()

	ms.on("GET", "/tpl-123", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"name": "hello_world", "language": "en_US", "category": "utility",
		})
	})

	c := ms.client()
	tpl, err := c.GetTemplate(context.Background(), "tpl-123")
	if err != nil {
		t.Fatalf("GetTemplate failed: %v", err)
	}
	if tpl.Name != "hello_world" {
		t.Errorf("expected hello_world, got %s", tpl.Name)
	}
}

func TestEditTemplate(t *testing.T) {
	ms := newDefaultMockServer()
	defer ms.Close()

	c := ms.client()
	tpl := &types.Template{Name: "hello_world", Language: "en_US", Category: "utility"}
	err := c.EditTemplate(context.Background(), "waba-456", "", tpl)
	if err != nil {
		t.Fatalf("EditTemplate failed: %v", err)
	}
}
