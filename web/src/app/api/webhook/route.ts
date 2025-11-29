import { NextRequest, NextResponse } from "next/server";
import { Webhook } from "standardwebhooks";
import { addEvent } from "@/lib/event-store";

function getWebhookSecret(): string {
  const secret = process.env.WEBHOOK_SECRET;
  if (!secret) {
    throw new Error("WEBHOOK_SECRET is not set. Run 'make setup-env' to generate env.local files.");
  }
  return secret;
}

export async function POST(request: NextRequest) {
  const wh = new Webhook(getWebhookSecret());

  // Get raw body as text for signature verification
  const body = await request.text();

  // Get webhook headers
  const headers = {
    id: request.headers.get("webhook-id") || "",
    timestamp: request.headers.get("webhook-timestamp") || "",
    signature: request.headers.get("webhook-signature") || "",
  };

  const stdHeaders = {
    "webhook-id": headers.id,
    "webhook-timestamp": headers.timestamp,
    "webhook-signature": headers.signature,
  };

  try {
    // Verify the webhook signature and parse as JSON
    const payload = wh.verify(body, stdHeaders) as Record<string, unknown>;

    // Store the verified event
    const event = addEvent({
      headers,
      payload,
      rawBody: body,
      verified: true,
    });

    console.log(`Received webhook (verified): id=${event.headers.id}`);

    return NextResponse.json({
      success: true,
      message: "Webhook received and verified successfully",
    });
  } catch (err) {
    const errorMessage = err instanceof Error ? err.message : "Unknown error";

    // Try to parse body as JSON for display
    let payload: Record<string, unknown> | null = null;
    try {
      payload = JSON.parse(body);
    } catch {
      // Body is not valid JSON
    }

    // Store the failed event
    const event = addEvent({
      headers,
      payload,
      rawBody: body,
      verified: false,
      error: errorMessage,
    });

    console.error(`Webhook verification failed: id=${event.headers.id}, error=${errorMessage}`);

    return NextResponse.json(
      { error: "Invalid webhook signature" },
      { status: 401 }
    );
  }
}
