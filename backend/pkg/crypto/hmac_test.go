package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testBody = "Hello, World!"

// Helper function to generate correct signature for testing
func generateSignature(secret, body string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(body))
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

func TestGitHubSignatureVerifierVerifySignature(t *testing.T) {
	verifier := NewGitHubSignatureVerifier()
	secret := "my_secret"
	validSignature := generateSignature(secret, testBody)

	tests := []struct {
		name      string
		secret    string
		signature string
		body      []byte
		expected  bool
	}{
		{
			name:      "Valid signature",
			secret:    secret,
			signature: validSignature,
			body:      []byte(testBody),
			expected:  true,
		},
		{
			name:      "Invalid signature",
			secret:    secret,
			signature: "sha256=invalid_signature",
			body:      []byte(testBody),
			expected:  false,
		},
		{
			name:      "Wrong secret",
			secret:    "wrong_secret",
			signature: validSignature,
			body:      []byte(testBody),
			expected:  false,
		},
		{
			name:      "Empty signature",
			secret:    secret,
			signature: "",
			body:      []byte(testBody),
			expected:  false,
		},
		{
			name:      "Malformed signature",
			secret:    secret,
			signature: "invalid_format",
			body:      []byte(testBody),
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := verifier.VerifySignature(tt.secret, tt.signature, tt.body)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewGitHubSignatureVerifier(t *testing.T) {
	verifier := NewGitHubSignatureVerifier()
	assert.NotNil(t, verifier)
}
