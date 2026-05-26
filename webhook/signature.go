package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// VerifySignature validates the X-Hub-Signature-256 header against the request body using HMAC-SHA256.
func VerifySignature(payload []byte, signature, appSecret string) bool {
	mac := hmac.New(sha256.New, []byte(appSecret))
	mac.Write(payload)
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}
