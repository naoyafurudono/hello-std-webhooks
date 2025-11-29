.PHONY: generate build run-server run-client test clean

# Generate ogen code from OpenAPI schema
generate:
	go generate ./api/...

# Build all binaries
build:
	go build -o bin/server ./cmd/server
	go build -o bin/client ./cmd/client

# Run the webhook server
run-server:
	go run ./cmd/server/main.go

# Run the webhook client
run-client:
	go run ./cmd/client/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Install dependencies
deps:
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	go vet ./...
