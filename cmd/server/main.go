package main

import (
	"log"
	"net/http"
	"os"

	"github.com/naoyafurudono/hello-std-webhooks/api"
	"github.com/naoyafurudono/hello-std-webhooks/server"
)

func main() {
	// Get the webhook secret from environment variable
	secret := os.Getenv("WEBHOOK_SECRET")
	if secret == "" {
		// Default secret for development (base64 encoded)
		// In production, always use a secure secret from environment
		secret = "whsec_MfKQ9r8GKYqrTwjUPD8ILPZIo2LaLaSw"
	}

	addr := os.Getenv("SERVER_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	// Create the handler
	handler := server.NewWebhookHandler()

	// Create the ogen server
	srv, err := api.NewServer(handler)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Wrap with signature verification middleware
	middleware, err := server.NewWebhookVerificationMiddleware(secret, srv)
	if err != nil {
		log.Fatalf("Failed to create middleware: %v", err)
	}

	log.Printf("Starting webhook server on %s", addr)
	if err := http.ListenAndServe(addr, middleware); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
