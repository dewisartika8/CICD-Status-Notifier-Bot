package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestWebhookEndpoint(t *testing.T) {
	os.Setenv("GITHUB_WEBHOOK_SECRET", "testsecret")
	app := fiber.New()
	// Daftarkan endpoint seperti di main.go
	app.Post("/api/v1/webhooks/github/:projectId", func(c *fiber.Ctx) error {
		sig := c.Get("X-Hub-Signature-256")
		if sig == "" {
			return c.Status(fiber.StatusUnauthorized).SendString("Missing signature header")
		}
		body := c.Body()
		secret := os.Getenv("GITHUB_WEBHOOK_SECRET")
		if !verifyGitHubSignature(secret, sig, body) {
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid signature")
		}
		var payload GitHubActionsPayload
		if err := json.Unmarshal(body, &payload); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid payload")
		}
		return c.SendStatus(fiber.StatusAccepted)
	})

	// Mock payload
	payload := GitHubActionsPayload{
		Action:   "completed",
		Workflow: "CI",
	}
	payload.Repository.Name = "repo"
	payload.Sender.Login = "user"
	body, _ := json.Marshal(payload)
	mac := hmac.New(sha256.New, []byte("testsecret"))
	mac.Write(body)
	signature := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	req := httptest.NewRequest("POST", "/api/v1/webhooks/github/test", bytes.NewReader(body))
	req.Header.Set("X-Hub-Signature-256", signature)
	resp, _ := app.Test(req)
	if resp.StatusCode != http.StatusAccepted {
		t.Errorf("Expected 202, got %d", resp.StatusCode)
	}
}
