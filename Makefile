.PHONY: generate build run-server run-client test clean deps fmt lint \
        web-install web-build web-dev web-send keygen send-go send-nextjs

# Generate ogen code from OpenAPI schema
generate:
	go tool ogen --config api/ogen.yaml --target api --package api --clean api/openapi.yaml

# Build all binaries
build:
	go build -o bin/server ./cmd/server
	go build -o bin/client ./cmd/client
	go build -o bin/keygen ./cmd/keygen

# Generate a new webhook secret
keygen:
	@go run ./cmd/keygen/

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

# Send webhook to Next.js server (legacy, use send-nextjs instead)
web-send:
	WEBHOOK_TARGET_URL=http://localhost:3000/api/webhook go run ./cmd/client/

# === Multi-target commands (requires .env with WEBHOOK_TARGET_<NAME>_* vars) ===

# Send webhook to Go server target
send-go:
	go run ./cmd/client/ -target go

# Send webhook to Next.js server target
send-nextjs:
	go run ./cmd/client/ -target nextjs
