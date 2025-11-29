package main

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/go-faster/jx"
	standardwebhooks "github.com/standard-webhooks/standard-webhooks/libraries/go"

	"github.com/naoyafurudono/hello-std-webhooks/api"
	"github.com/naoyafurudono/hello-std-webhooks/client"
	"github.com/naoyafurudono/hello-std-webhooks/server"
)

const testSecret = "whsec_MfKQ9r8GKYqrTwjUPD8ILPZIo2LaLaSw"

func setupTestServer(t *testing.T) *httptest.Server {
	t.Helper()

	handler := server.NewWebhookHandler()
	srv, err := api.NewWebhookServer(handler)
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}

	webhookHandler := srv.Handler("userEvent")
	middleware, err := server.NewWebhookVerificationMiddleware(testSecret, webhookHandler)
	if err != nil {
		t.Fatalf("failed to create middleware: %v", err)
	}

	return httptest.NewServer(middleware)
}

func TestIntegration_ValidSignature(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()

	wc, err := client.NewWebhookClient(ts.URL, testSecret)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	event := &api.WebhookEvent{
		Type: "user.created",
		Data: api.WebhookEventData{
			"id":    mustEncodeJSON("user_123"),
			"email": mustEncodeJSON("test@example.com"),
		},
	}

	res, err := wc.SendWebhook(context.Background(), event)
	if err != nil {
		t.Fatalf("failed to send webhook: %v", err)
	}

	resp, ok := res.(*api.WebhookResponse)
	if !ok {
		t.Fatalf("unexpected response type: %T", res)
	}

	if !resp.Success {
		t.Errorf("expected success=true, got false")
	}
	if resp.Message != "Webhook event 'user.created' processed successfully" {
		t.Errorf("unexpected message: %s", resp.Message)
	}
}

func TestIntegration_InvalidSignature(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()

	// Use a different secret for the client (valid base64 format)
	wrongSecret := "whsec_C2FtcGxlLXdyb25nLXNlY3JldC1rZXk="
	wh, err := standardwebhooks.NewWebhook(wrongSecret)
	if err != nil {
		t.Fatalf("failed to create webhook: %v", err)
	}

	payload := []byte(`{"type":"user.created","data":{"id":"user_123"}}`)
	msgID := "msg_test123"
	timestamp := time.Now()

	signature, err := wh.Sign(msgID, timestamp, payload)
	if err != nil {
		t.Fatalf("failed to sign: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, ts.URL, bytes.NewReader(payload))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("webhook-id", msgID)
	req.Header.Set("webhook-timestamp", strconv.FormatInt(timestamp.Unix(), 10))
	req.Header.Set("webhook-signature", signature)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", resp.StatusCode)
	}
}

func TestIntegration_MissingHeaders(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()

	payload := []byte(`{"type":"user.created","data":{"id":"123"}}`)
	req, err := http.NewRequest(http.MethodPost, ts.URL, bytes.NewReader(payload))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	// No webhook headers set

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", resp.StatusCode)
	}
}

func TestIntegration_ExpiredTimestamp(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()

	wh, err := standardwebhooks.NewWebhook(testSecret)
	if err != nil {
		t.Fatalf("failed to create webhook: %v", err)
	}

	payload := []byte(`{"type":"user.created","data":{"id":"123"}}`)
	msgID := "msg_test123"
	// Timestamp from 10 minutes ago (beyond the default tolerance)
	expiredTimestamp := time.Now().Add(-10 * time.Minute)

	signature, err := wh.Sign(msgID, expiredTimestamp, payload)
	if err != nil {
		t.Fatalf("failed to sign: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, ts.URL, bytes.NewReader(payload))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("webhook-id", msgID)
	req.Header.Set("webhook-timestamp", strconv.FormatInt(expiredTimestamp.Unix(), 10))
	req.Header.Set("webhook-signature", signature)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		body, _ := io.ReadAll(resp.Body)
		t.Errorf("expected status 401, got %d: %s", resp.StatusCode, body)
	}
}

func TestIntegration_TamperedPayload(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()

	wh, err := standardwebhooks.NewWebhook(testSecret)
	if err != nil {
		t.Fatalf("failed to create webhook: %v", err)
	}

	originalPayload := []byte(`{"type":"user.created","data":{"id":"123"}}`)
	tamperedPayload := []byte(`{"type":"user.created","data":{"id":"456"}}`)
	msgID := "msg_test123"
	timestamp := time.Now()

	// Sign with original payload
	signature, err := wh.Sign(msgID, timestamp, originalPayload)
	if err != nil {
		t.Fatalf("failed to sign: %v", err)
	}

	// Send tampered payload with original signature
	req, err := http.NewRequest(http.MethodPost, ts.URL, bytes.NewReader(tamperedPayload))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("webhook-id", msgID)
	req.Header.Set("webhook-timestamp", strconv.FormatInt(timestamp.Unix(), 10))
	req.Header.Set("webhook-signature", signature)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", resp.StatusCode)
	}
}

func TestIntegration_DifferentEventTypes(t *testing.T) {
	ts := setupTestServer(t)
	defer ts.Close()

	wc, err := client.NewWebhookClient(ts.URL, testSecret)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	eventTypes := []string{"user.created", "user.updated", "user.deleted", "order.placed"}

	for _, eventType := range eventTypes {
		t.Run(eventType, func(t *testing.T) {
			event := &api.WebhookEvent{
				Type: eventType,
				Data: api.WebhookEventData{
					"id": mustEncodeJSON("test_123"),
				},
			}

			res, err := wc.SendWebhook(context.Background(), event)
			if err != nil {
				t.Fatalf("failed to send webhook: %v", err)
			}

			resp, ok := res.(*api.WebhookResponse)
			if !ok {
				t.Fatalf("unexpected response type: %T", res)
			}

			if !resp.Success {
				t.Errorf("expected success=true, got false")
			}
		})
	}
}

func mustEncodeJSON(v string) jx.Raw {
	var e jx.Encoder
	e.Str(v)
	return e.Bytes()
}
