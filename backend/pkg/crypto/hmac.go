package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// GitHubSignatureVerifier handles GitHub webhook signature verification
type GitHubSignatureVerifier struct{}

// NewGitHubSignatureVerifier creates a new GitHub signature verifier
func NewGitHubSignatureVerifier() *GitHubSignatureVerifier {
	return &GitHubSignatureVerifier{}
}

// VerifySignature verifies GitHub webhook signature using HMAC-SHA256
func (v *GitHubSignatureVerifier) VerifySignature(secret, signature string, body []byte) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}

// SignatureVerifier interface for webhook signature verification
type SignatureVerifier interface {
	VerifySignature(secret, signature string, body []byte) bool
}
