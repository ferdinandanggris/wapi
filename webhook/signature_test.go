package webhook_test

import (
	"testing"

	"github.com/ferdinandanggris/wapi/webhook"
)

func TestVerifySignature(t *testing.T) {
	payload := []byte(`{"test":"data"}`)
	validSig := "sha256=ca742bae4305b1cbf679a7314db99b27ee7b2331e4f5fc1751c0a558914b0723"

	tests := []struct {
		name      string
		payload   []byte
		signature string
		secret    string
		want      bool
	}{
		{"valid signature", payload, validSig, "my-app-secret", true},
		{"invalid signature", payload, "sha256=invalid", "my-app-secret", false},
		{"wrong payload", []byte(`{"other":"data"}`), validSig, "my-app-secret", false},
		{"wrong secret", payload, validSig, "wrong-secret", false},
		{"empty signature", payload, "", "my-app-secret", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := webhook.VerifySignature(tt.payload, tt.signature, tt.secret)
			if got != tt.want {
				t.Errorf("VerifySignature() = %v, want %v", got, tt.want)
			}
		})
	}
}
