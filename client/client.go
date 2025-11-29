package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	standardwebhooks "github.com/standard-webhooks/standard-webhooks/libraries/go"

	"github.com/naoyafurudono/hello-std-webhooks/api"
)

// Default HTTP client with reasonable timeout settings.
var defaultHTTPClient = &http.Client{
	Timeout: 30 * time.Second,
}

// Option is a functional option for configuring WebhookClient.
type Option func(*WebhookClient)

// WithHTTPClient sets a custom HTTP client for the WebhookClient.
// This is useful for setting custom timeouts, transport settings, or for testing.
func WithHTTPClient(client *http.Client) Option {
	return func(wc *WebhookClient) {
		wc.httpClient = client
	}
}

// WebhookClient sends webhook events with standard-webhooks signing.
// Note: This client does not use ogen-generated WebhookClient because
// we need to add standard-webhooks signature headers (webhook-id, webhook-timestamp,
// webhook-signature) to each request, which requires access to the request body
// before signing. Using a custom HTTP client is simpler than using a RoundTripper
// that needs to coordinate the message ID from the caller.
type WebhookClient struct {
	wh         *standardwebhooks.Webhook
	targetURL  string
	httpClient *http.Client
}

// NewWebhookClient creates a new webhook client with signature signing capability.
// The secret should be a base64-encoded secret key.
func NewWebhookClient(targetURL string, secret string, opts ...Option) (*WebhookClient, error) {
	wh, err := standardwebhooks.NewWebhook(secret)
	if err != nil {
		return nil, err
	}

	wc := &WebhookClient{
		wh:         wh,
		targetURL:  targetURL,
		httpClient: defaultHTTPClient,
	}

	for _, opt := range opts {
		opt(wc)
	}

	return wc, nil
}

// SendWebhook sends a webhook event with proper standard-webhooks headers.
// The msgID should be unique per event and remain the same across retries.
// This is used as an idempotency key by consumers.
func (c *WebhookClient) SendWebhook(ctx context.Context, msgID string, event *api.WebhookEvent) (api.UserEventRes, error) {
	// Encode the event to JSON
	body, err := event.MarshalJSON()
	if err != nil {
		return nil, err
	}

	timestamp := time.Now()

	// Sign the payload
	signature, err := c.wh.Sign(msgID, timestamp, body)
	if err != nil {
		return nil, err
	}

	// Create the request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.targetURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("webhook-id", msgID)
	req.Header.Set("webhook-timestamp", formatTimestamp(timestamp))
	req.Header.Set("webhook-signature", signature)

	// Send the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Decode the response based on status code
	switch resp.StatusCode {
	case http.StatusOK:
		var result api.WebhookResponse
		if err := result.UnmarshalJSON(respBody); err != nil {
			return nil, err
		}
		return &result, nil
	case http.StatusBadRequest:
		var result api.UserEventBadRequest
		if err := result.UnmarshalJSON(respBody); err != nil {
			return nil, err
		}
		return &result, nil
	case http.StatusUnauthorized:
		var result api.UserEventUnauthorized
		if err := result.UnmarshalJSON(respBody); err != nil {
			return nil, err
		}
		return &result, nil
	default:
		return nil, &UnexpectedStatusError{StatusCode: resp.StatusCode, Body: respBody}
	}
}

// UnexpectedStatusError is returned when the server returns an unexpected HTTP status code.
type UnexpectedStatusError struct {
	StatusCode int
	Body       []byte
}

func (e *UnexpectedStatusError) Error() string {
	return fmt.Sprintf("unexpected status code: %d", e.StatusCode)
}

// formatTimestamp formats the timestamp as a Unix timestamp string.
func formatTimestamp(t time.Time) string {
	return strconv.FormatInt(t.Unix(), 10)
}
