package client

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	standardwebhooks "github.com/standard-webhooks/standard-webhooks/libraries/go"

	"github.com/naoyafurudono/hello-std-webhooks/api"
)

// WebhookClient wraps the ogen-generated client with standard-webhooks signing.
type WebhookClient struct {
	client *api.Client
	wh     *standardwebhooks.Webhook
}

// NewWebhookClient creates a new webhook client with signature signing capability.
// The secret should be a base64-encoded secret key.
func NewWebhookClient(serverURL string, secret string) (*WebhookClient, error) {
	wh, err := standardwebhooks.NewWebhook(secret)
	if err != nil {
		return nil, err
	}

	// Create a custom HTTP client with signing transport
	transport := &signingTransport{
		wh:   wh,
		base: http.DefaultTransport,
	}
	httpClient := &http.Client{Transport: transport}

	client, err := api.NewClient(serverURL, api.WithClient(httpClient))
	if err != nil {
		return nil, err
	}

	return &WebhookClient{
		client: client,
		wh:     wh,
	}, nil
}

// SendWebhook sends a webhook event with proper standard-webhooks headers.
func (c *WebhookClient) SendWebhook(ctx context.Context, event *api.WebhookEvent) (api.ReceiveWebhookRes, error) {
	return c.client.ReceiveWebhook(ctx, event)
}

// signingTransport adds standard-webhooks headers to outgoing requests.
type signingTransport struct {
	wh   *standardwebhooks.Webhook
	base http.RoundTripper
}

func (t *signingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Only sign webhook requests
	if req.URL.Path == "/webhook" && req.Method == http.MethodPost {
		// Read the body
		var body []byte
		if req.Body != nil {
			var err error
			body, err = io.ReadAll(req.Body)
			if err != nil {
				return nil, err
			}
			req.Body.Close()
			req.Body = io.NopCloser(bytes.NewReader(body))
		}

		// Generate webhook ID and timestamp
		msgID := generateMsgID()
		timestamp := time.Now()

		// Sign the payload
		signature, err := t.wh.Sign(msgID, timestamp, body)
		if err != nil {
			return nil, err
		}

		// Set the standard-webhooks headers
		req.Header.Set("webhook-id", msgID)
		req.Header.Set("webhook-timestamp", formatTimestamp(timestamp))
		req.Header.Set("webhook-signature", signature)
	}

	return t.base.RoundTrip(req)
}

// generateMsgID generates a unique message ID for the webhook.
func generateMsgID() string {
	return "msg_" + uuid.New().String()
}

// formatTimestamp formats the timestamp as a Unix timestamp string.
func formatTimestamp(t time.Time) string {
	return strconv.FormatInt(t.Unix(), 10)
}
