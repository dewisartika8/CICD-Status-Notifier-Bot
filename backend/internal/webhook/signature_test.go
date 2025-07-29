package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func verifyGitHubSignature(secret, signature string, payload []byte) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expectedMAC := mac.Sum(nil)
	signature = signature[len("sha256="):]
	decodedSignature, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}
	return hmac.Equal(expectedMAC, decodedSignature)
}

func TestVerifyGitHubSignature(t *testing.T) {
	secret := "testsecret"
	body := []byte(`{"test":"ok"}`)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	if !verifyGitHubSignature(secret, expected, body) {
		t.Error("Signature verification failed")
	}
	if verifyGitHubSignature(secret, "sha256=invalid", body) {
		t.Error("Signature verification should fail for invalid signature")
	}
}
