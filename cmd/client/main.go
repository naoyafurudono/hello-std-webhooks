package main

import (
	"context"
	"log"
	"os"

	"github.com/go-faster/jx"
	"github.com/google/uuid"
	"github.com/joho/godotenv"

	"github.com/naoyafurudono/hello-std-webhooks/api"
	"github.com/naoyafurudono/hello-std-webhooks/client"
)

func main() {
	// Load env.local if it exists (ignore error if not found)
	_ = godotenv.Load("env.local")

	targetURL := os.Getenv("WEBHOOK_TARGET_URL")
	if targetURL == "" {
		log.Fatal("WEBHOOK_TARGET_URL is not set. Run 'make setup-env' first.")
	}

	secret := os.Getenv("WEBHOOK_SECRET")
	if secret == "" {
		log.Fatal("WEBHOOK_SECRET is not set. Run 'make setup-env' first.")
	}

	// Create the webhook client
	wc, err := client.NewWebhookClient(targetURL, secret)
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

	// Generate a unique message ID for this event.
	// In production, this should be derived from the event itself
	// (e.g., "msg_" + event ID) and stored for retries.
	msgID := "msg_" + uuid.New().String()

	// Send the webhook with the message ID
	log.Printf("Sending webhook to %s", targetURL)
	res, err := wc.SendWebhook(context.Background(), msgID, event)
	if err != nil {
		log.Fatalf("Failed to send webhook: %v", err)
	}

	// Handle the response
	switch r := res.(type) {
	case *api.WebhookResponse:
		log.Printf("Webhook sent successfully: success=%v, message=%s", r.Success, r.Message)
	case *api.UserEventBadRequest:
		log.Printf("Bad request: %s", r.Error)
	case *api.UserEventUnauthorized:
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
