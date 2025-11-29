package server

import (
	"context"
	"fmt"
	"log"

	"github.com/naoyafurudono/hello-std-webhooks/api"
)

// WebhookHandler implements api.Handler for processing webhook events.
type WebhookHandler struct{}

// NewWebhookHandler creates a new WebhookHandler.
func NewWebhookHandler() *WebhookHandler {
	return &WebhookHandler{}
}

// ReceiveWebhook handles incoming webhook events.
// Note: Signature verification is done in the middleware before this handler is called.
func (h *WebhookHandler) ReceiveWebhook(ctx context.Context, req *api.WebhookEvent) (api.ReceiveWebhookRes, error) {
	log.Printf("Received webhook event: type=%s", req.Type)

	// Process the webhook event based on type
	switch req.Type {
	case "user.created":
		log.Printf("Processing user.created event with data: %v", req.Data)
	case "user.updated":
		log.Printf("Processing user.updated event with data: %v", req.Data)
	case "user.deleted":
		log.Printf("Processing user.deleted event with data: %v", req.Data)
	default:
		log.Printf("Unknown event type: %s", req.Type)
	}

	return &api.WebhookResponse{
		Success: true,
		Message: fmt.Sprintf("Webhook event '%s' processed successfully", req.Type),
	}, nil
}
