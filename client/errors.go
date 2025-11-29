package client

import "errors"

// ErrMissingWebhookID is returned when SendWebhook is called without a webhook ID in context.
// Use WithWebhookID to set the message ID before calling SendWebhook.
var ErrMissingWebhookID = errors.New("webhook ID is required: use WithWebhookID to set the message ID")
