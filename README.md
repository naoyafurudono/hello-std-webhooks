# hello-std-webhooks

A demo project for [Standard Webhooks](https://www.standardwebhooks.com/), featuring a Go webhook client and a Next.js webhook server.

## Overview

Webhooks are HTTP callbacks that enable event-driven communication between services. This project demonstrates how to send and verify webhooks securely using the Standard Webhooks specification.

Standard Webhooks provides protection against:
- **Spoofing** - Signature verification ensures the webhook came from a trusted source
- **Replay attacks** - Timestamp validation prevents reuse of captured requests

This project includes:

- **Go Client** (`cmd/client`): Sends signed webhook requests
- **Next.js Server** (`web/`): Receives and verifies webhook signatures
- **Key Generator** (`cmd/keygen`): Generates `whsec_` formatted secrets

## Quick Start

```bash
# Install dependencies
make deps

# Generate env.local files with a shared secret
make setup-env

# Start the Next.js server (in one terminal)
make web-dev

# Send a test webhook (in another terminal)
make send
```

Then open http://localhost:3000 to view the API documentation and received events.

## Project Structure

```
.
├── api/                    # OpenAPI schema and generated code (ogen)
├── cmd/
│   ├── client/            # Go webhook client
│   └── keygen/            # Secret key generator
├── client/                # Webhook client library
├── web/                   # Next.js webhook server
│   └── src/
│       ├── app/
│       │   ├── api/webhook/   # POST /api/webhook endpoint
│       │   ├── api/events/    # GET/DELETE /api/events endpoint
│       │   ├── events/        # Events viewer page
│       │   └── page.tsx       # Home page with API docs
│       └── lib/
│           └── event-store.ts # In-memory event storage
├── Makefile
└── README.md
```

## Makefile Commands

| Command | Description |
|---------|-------------|
| `make deps` | Install Go and npm dependencies |
| `make setup-env` | Generate env.local files with shared secret |
| `make web-dev` | Start Next.js dev server |
| `make web-build` | Build Next.js for production |
| `make send` | Send a test webhook to the server |
| `make keygen` | Generate a new webhook secret |
| `make generate` | Regenerate ogen code from OpenAPI schema |
| `make build` | Build Go binaries |
| `make test` | Run Go tests |
| `make clean` | Remove build artifacts |

## Standard Webhooks Specification

This project follows the [Standard Webhooks](https://github.com/standard-webhooks/standard-webhooks) specification for signing and verifying webhooks.

### Required Headers

| Header | Description |
|--------|-------------|
| `webhook-id` | Unique message identifier (e.g., `msg_abc123`) |
| `webhook-timestamp` | Unix timestamp in seconds |
| `webhook-signature` | `v1,<base64-hmac-sha256>` |

### Signature Calculation

```
signed_content = webhook_id + "." + webhook_timestamp + "." + body
signature = base64(HMAC-SHA256(base64_decode(secret), signed_content))
```

The signature header format is `v1,<signature>`.

## Environment Variables

### Client (`env.local`)

| Variable | Description |
|----------|-------------|
| `WEBHOOK_TARGET_URL` | Target webhook endpoint URL |
| `WEBHOOK_SECRET` | Shared secret (`whsec_...` format) |

### Server (`web/env.local`)

| Variable | Description |
|----------|-------------|
| `WEBHOOK_SECRET` | Shared secret for verifying webhook signatures |

## Deployment

The project includes a `render.yaml` for deploying to [Render](https://render.com/).

## License

MIT
