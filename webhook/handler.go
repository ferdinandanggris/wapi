// Package webhook provides an HTTP handler for WhatsApp Cloud API webhook events
// with signature verification, OnMessage and OnStatus callbacks.
package webhook

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/ferdinandanggris/wapi/types"
)

// Handler implements http.Handler for WhatsApp webhook verification and event delivery.
// Configure VerifyToken and AppSecret for production use.
type Handler struct {
	// VerifyToken must match the token set in the Meta App dashboard during webhook setup.
	VerifyToken string
	// AppSecret is used to verify X-Hub-Signature-256. Leave empty to skip verification.
	AppSecret string
	// OnMessage is called for each incoming message. Return an error to log it.
	OnMessage func(msg *types.IncomingMsg, meta *types.Metadata, contact *types.WaContact) error
	// OnStatus is called for each status update (sent, delivered, read, failed).
	OnStatus func(status *types.StatusUpdate, meta *types.Metadata) error
	// Logger is used for debug logging. If nil, no logging is performed.
	Logger *log.Logger
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.verify(w, r)
	case http.MethodPost:
		h.handle(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) verify(w http.ResponseWriter, r *http.Request) {
	mode := r.URL.Query().Get("hub.mode")
	token := r.URL.Query().Get("hub.verify_token")
	challenge := r.URL.Query().Get("hub.challenge")

	if mode != "subscribe" || token != h.VerifyToken {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(challenge)) //nolint:errcheck // net/http handles write errors
}

func (h *Handler) handle(w http.ResponseWriter, r *http.Request) {
	if h.AppSecret != "" {
		sig := r.Header.Get("X-Hub-Signature-256")
		if sig == "" {
			h.log("webhook: missing X-Hub-Signature-256")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			h.log("webhook: read body: %v", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		r.Body.Close()

		if !VerifySignature(body, sig, h.AppSecret) {
			h.log("webhook: invalid signature")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		r.Body = io.NopCloser(bytes.NewReader(body))
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`)) //nolint:errcheck // net/http handles write errors

	var payload types.WebhookPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.log("webhook: decode payload: %v", err)
		return
	}

	for _, entry := range payload.Entry {
		for _, change := range entry.Changes {
			if change.Value == nil {
				continue
			}

			v := change.Value

			if len(v.Messages) > 0 && h.OnMessage != nil {
				msg := v.Messages[0]
				meta := v.Metadata
				var contact *types.WaContact
				if len(v.Contacts) > 0 {
					contact = v.Contacts[0]
				}
				if err := h.OnMessage(msg, meta, contact); err != nil {
					h.log("webhook: onMessage: %v", err)
				}
			}

			if len(v.Statuses) > 0 && h.OnStatus != nil {
				for _, status := range v.Statuses {
					if err := h.OnStatus(status, v.Metadata); err != nil {
						h.log("webhook: onStatus: %v", err)
					}
				}
			}
		}
	}
}

func (h *Handler) log(format string, args ...interface{}) {
	if h.Logger != nil {
		h.Logger.Printf(format, args...)
	}
}
