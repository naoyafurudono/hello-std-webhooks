package client

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strconv"
	"time"

	standardwebhooks "github.com/standard-webhooks/standard-webhooks/libraries/go"

	"github.com/naoyafurudono/hello-std-webhooks/api"
)

// WebhookClient wraps the ogen-generated WebhookClient with standard-webhooks signing.
type WebhookClient struct {
	wh        *standardwebhooks.Webhook
	targetURL string
}

// NewWebhookClient creates a new webhook client with signature signing capability.
// The secret should be a base64-encoded secret key.
func NewWebhookClient(targetURL string, secret string) (*WebhookClient, error) {
	wh, err := standardwebhooks.NewWebhook(secret)
	if err != nil {
		return nil, err
	}

	return &WebhookClient{
		wh:        wh,
		targetURL: targetURL,
	}, nil
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
	resp, err := http.DefaultClient.Do(req)
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
		return nil, &unexpectedStatusError{StatusCode: resp.StatusCode, Body: respBody}
	}
}

type unexpectedStatusError struct {
	StatusCode int
	Body       []byte
}

func (e *unexpectedStatusError) Error() string {
	return "unexpected status code: " + strconv.Itoa(e.StatusCode)
}

// formatTimestamp formats the timestamp as a Unix timestamp string.
func formatTimestamp(t time.Time) string {
	return strconv.FormatInt(t.Unix(), 10)
}
