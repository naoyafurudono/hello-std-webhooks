package main

import (
	"context"
	"log"
	"os"

	"github.com/go-faster/jx"

	"github.com/naoyafurudono/hello-std-webhooks/api"
	"github.com/naoyafurudono/hello-std-webhooks/client"
)

func main() {
	// Get the webhook secret from environment variable
	secret := os.Getenv("WEBHOOK_SECRET")
	if secret == "" {
		// Default secret for development (base64 encoded)
		// In production, always use a secure secret from environment
		secret = "whsec_MfKQ9r8GKYqrTwjUPD8ILPZIo2LaLaSw"
	}

	serverURL := os.Getenv("WEBHOOK_SERVER_URL")
	if serverURL == "" {
		serverURL = "http://localhost:8080"
	}

	// Create the webhook client
	wc, err := client.NewWebhookClient(serverURL, secret)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Create a sample webhook event
	event := &api.WebhookEvent{
		Type: "user.created",
		Data: api.WebhookEventData{
			"id":    mustEncodeJSON("user_123"),
			"email": mustEncodeJSON("user@example.com"),
			"name":  mustEncodeJSON("John Doe"),
		},
	}

	// Send the webhook
	ctx := context.Background()
	res, err := wc.SendWebhook(ctx, event)
	if err != nil {
		log.Fatalf("Failed to send webhook: %v", err)
	}

	// Handle the response
	switch r := res.(type) {
	case *api.WebhookResponse:
		log.Printf("Webhook sent successfully: success=%v, message=%s", r.Success, r.Message)
	case *api.ReceiveWebhookBadRequest:
		log.Printf("Bad request: %s", r.Error)
	case *api.ReceiveWebhookUnauthorized:
		log.Printf("Unauthorized: %s", r.Error)
	default:
		log.Printf("Unknown response type: %T", r)
	}
}

func mustEncodeJSON(v string) jx.Raw {
	var e jx.Encoder
	e.Str(v)
	return e.Bytes()
}
