.PHONY: generate build run-server run-client test clean deps fmt lint \
        web-install web-build web-dev web-send

# Generate ogen code from OpenAPI schema
generate:
	go tool ogen --config api/ogen.yaml --target api --package api --clean api/openapi.yaml

# Build all binaries
build:
	go build -o bin/server ./cmd/server
	go build -o bin/client ./cmd/client

# Run the Go webhook server
run-server:
	go run ./cmd/server/main.go

# Run the Go webhook client (sends to Go server by default)
run-client:
	go run ./cmd/client/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf web/.next/

# Install dependencies
deps:
	go mod tidy
	cd web && npm install

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	go vet ./...

# === Next.js webhook server ===

# Install web dependencies
web-install:
	cd web && npm install

# Build Next.js app
web-build:
	cd web && npm run build

# Run Next.js dev server
web-dev:
	cd web && npm run dev

# Send webhook to Next.js server
web-send:
	WEBHOOK_TARGET_URL=http://localhost:3000/api/webhook go run ./cmd/client/
