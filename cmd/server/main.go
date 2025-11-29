package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/naoyafurudono/hello-std-webhooks/api"
	"github.com/naoyafurudono/hello-std-webhooks/server"
)

func main() {
	// Load env.local if it exists (ignore error if not found)
	_ = godotenv.Load("env.local")

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

	// Create the webhook handler
	handler := server.NewWebhookHandler()

	// Create the ogen webhook server
	srv, err := api.NewWebhookServer(handler)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Get the handler for "userEvent" webhook and wrap with signature verification
	webhookHandler := srv.Handler("userEvent")
	middleware, err := server.NewWebhookVerificationMiddleware(secret, webhookHandler)
	if err != nil {
		log.Fatalf("Failed to create middleware: %v", err)
	}

	// Mount at /webhook path
	mux := http.NewServeMux()
	mux.Handle("/webhook", middleware)

	log.Printf("Starting webhook server on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
