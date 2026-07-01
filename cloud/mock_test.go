package cloud_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/ferdinandanggris/wapi/cloud"
)

type mockHandler func(w http.ResponseWriter, r *http.Request)

type mockServer struct {
	*httptest.Server
	handlers map[string]mockHandler
}

func newMockServer() *mockServer {
	ms := &mockServer{handlers: make(map[string]mockHandler)}
	ms.Server = httptest.NewServer(http.HandlerFunc(ms.serve))
	return ms
}

func (ms *mockServer) on(method, path string, h mockHandler) {
	ms.handlers[method+":"+path] = h
}

func (ms *mockServer) serve(w http.ResponseWriter, r *http.Request) {
	key := r.Method + ":" + r.URL.Path
	if h, ok := ms.handlers[key]; ok {
		h(w, r)
		return
	}
	key = r.Method + ":*"
	if h, ok := ms.handlers[key]; ok {
		h(w, r)
		return
	}
	http.NotFound(w, r)
}

func (ms *mockServer) client(opts ...cloud.Option) *cloud.CloudClient {
	all := []cloud.Option{
		cloud.WithBaseURL(ms.URL),
		cloud.WithAccessToken("fake-token"),
		cloud.WithAPIVersion(""),
	}
	all = append(all, opts...)
	return cloud.New(all...)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, metaCode int, msg string, subcode int) {
	httpCode := http.StatusOK
	switch {
	case metaCode == 130429:
		httpCode = http.StatusTooManyRequests
	case metaCode == 131030 || metaCode == 132012 || metaCode == 100:
		httpCode = http.StatusBadRequest
	case metaCode >= 500 && metaCode < 600:
		httpCode = metaCode
	}

	writeJSON(w, httpCode, map[string]interface{}{
		"error": map[string]interface{}{
			"message":        msg,
			"type":           "OAuthException",
			"code":           metaCode,
			"error_subcode":  subcode,
			"fbtrace_id":     "AbC123xYz",
		},
	})
}

func parseBody(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

func defaultSendHandler(w http.ResponseWriter, r *http.Request) {
	var body map[string]interface{}
	if err := parseBody(r, &body); err != nil {
		writeError(w, 100, "invalid JSON", 0)
		return
	}

	to, _ := body["to"].(string)

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"messaging_product": "whatsapp",
		"contacts": []map[string]string{
			{"input": to, "wa_id": strings.TrimPrefix(to, "+")},
		},
		"messages": []map[string]string{
			{"id": "wamid.HBgLMTgwOTEyMzQ1NjcVAgASGBQzRUIwMEY2QjBCNDY2N0YwMzAzMAA="},
		},
	})
}

func newDefaultMockServer() *mockServer {
	ms := newMockServer()
	ms.on("POST", "/123/messages", defaultSendHandler)
	ms.on("POST", "/123/media", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"id": "media-id-123"})
	})
	ms.on("GET", "/media-id-123", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"url":      fmt.Sprintf("%s/download/media-id-123", ms.URL),
			"mime_type": "image/jpeg",
			"sha256":   "abc123",
			"file_size": 1024,
			"id":       "media-id-123",
		})
	})
	ms.on("DELETE", "/media-id-123", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]bool{"success": true})
	})
	ms.on("POST", "/waba-456/message_templates", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"name":     "hello_world",
			"language": "en_US",
			"category": "utility",
		})
	})
	ms.on("GET", "/waba-456/message_templates", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"data": []map[string]interface{}{
				{"name": "hello_world", "language": "en_US", "category": "utility"},
			},
		})
	})
	ms.on("POST", "/123/register", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]bool{"success": true})
	})
	ms.on("GET", "/123", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"id": "123", "display_phone_number": "+16505555555",
			"verified_name": "Test Business", "quality_rating": "GREEN",
		})
	})
	ms.on("GET", "/waba-456/phone_numbers", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "123", "display_phone_number": "+16505555555"},
			},
		})
	})
	ms.on("GET", "/123/whatsapp_business_profile", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"data": []map[string]interface{}{
				{"description": "Test business profile"},
			},
		})
	})
	ms.on("POST", "/123/whatsapp_business_profile", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]bool{"success": true})
	})
	ms.on("GET", "/business-789/owned_whatsapp_business_accounts", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"data": []map[string]interface{}{
				{"id": "906278885743121", "name": "Seller Pulsa", "timezone_id": "66", "message_template_namespace": "dddcb3ee_59b1_4397_99e4_5f2017232c44"},
				{"id": "841222862303866", "name": "Family Pulsa", "currency": "IDR", "timezone_id": "66", "message_template_namespace": "46a16bc1_0dd7_4ca6_95fb_428149ddcb6b"},
			},
			"paging": map[string]interface{}{
				"cursors": map[string]string{
					"before": "cursor-before",
					"after":  "cursor-after",
				},
			},
		})
	})
	return ms
}
