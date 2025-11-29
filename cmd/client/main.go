package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-faster/jx"
	"github.com/google/uuid"
	"github.com/joho/godotenv"

	"github.com/naoyafurudono/hello-std-webhooks/api"
	"github.com/naoyafurudono/hello-std-webhooks/client"
)

func main() {
	// Load env.local if it exists (ignore error if not found)
	_ = godotenv.Load("env.local")

	var targetName string
	flag.StringVar(&targetName, "target", "", "target name (e.g., GO, NEXTJS) - reads WEBHOOK_TARGET_<NAME>_URL and WEBHOOK_TARGET_<NAME>_SECRET")
	flag.Parse()

	targetURL, secret := getTargetConfig(targetName)

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

func getTargetConfig(targetName string) (url, secret string) {
	if targetName != "" {
		// Use named target: WEBHOOK_TARGET_<NAME>_URL and WEBHOOK_TARGET_<NAME>_SECRET
		name := strings.ToUpper(targetName)
		url = os.Getenv(fmt.Sprintf("WEBHOOK_TARGET_%s_URL", name))
		secret = os.Getenv(fmt.Sprintf("WEBHOOK_TARGET_%s_SECRET", name))

		if url == "" {
			log.Fatalf("WEBHOOK_TARGET_%s_URL is not set", name)
		}
		if secret == "" {
			log.Fatalf("WEBHOOK_TARGET_%s_SECRET is not set", name)
		}
		return url, secret
	}

	// Fallback to legacy env vars for backwards compatibility
	url = os.Getenv("WEBHOOK_TARGET_URL")
	if url == "" {
		url = "http://localhost:8080/webhook"
	}

	secret = os.Getenv("WEBHOOK_SECRET")
	if secret == "" {
		// Default secret for development
		secret = "whsec_MfKQ9r8GKYqrTwjUPD8ILPZIo2LaLaSw"
	}

	return url, secret
}

func mustEncodeJSON(v string) jx.Raw {
	var e jx.Encoder
	e.Str(v)
	return e.Bytes()
}
