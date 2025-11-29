package server

import (
	"bytes"
	"io"
	"net/http"

	standardwebhooks "github.com/standard-webhooks/standard-webhooks/libraries/go"
)

// WebhookVerificationMiddleware wraps an HTTP handler to verify standard-webhooks signatures.
type WebhookVerificationMiddleware struct {
	wh   *standardwebhooks.Webhook
	next http.Handler
}

// NewWebhookVerificationMiddleware creates a new verification middleware.
// The secret should be a base64-encoded secret key.
func NewWebhookVerificationMiddleware(secret string, next http.Handler) (*WebhookVerificationMiddleware, error) {
	wh, err := standardwebhooks.NewWebhook(secret)
	if err != nil {
		return nil, err
	}
	return &WebhookVerificationMiddleware{
		wh:   wh,
		next: next,
	}, nil
}

// ServeHTTP implements http.Handler.
func (m *WebhookVerificationMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Only verify POST requests to /webhook
	if r.Method == http.MethodPost && r.URL.Path == "/webhook" {
		// Read the body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, `{"error": "Failed to read request body"}`, http.StatusBadRequest)
			return
		}
		r.Body.Close()

		// Get the webhook headers
		headers := http.Header{}
		headers.Set("webhook-id", r.Header.Get("webhook-id"))
		headers.Set("webhook-timestamp", r.Header.Get("webhook-timestamp"))
		headers.Set("webhook-signature", r.Header.Get("webhook-signature"))

		// Verify the signature
		if err := m.wh.Verify(body, headers); err != nil {
			http.Error(w, `{"error": "Invalid webhook signature"}`, http.StatusUnauthorized)
			return
		}

		// Restore the body for the next handler
		r.Body = io.NopCloser(bytes.NewReader(body))
	}

	m.next.ServeHTTP(w, r)
}
