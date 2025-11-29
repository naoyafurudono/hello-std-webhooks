import { NextRequest, NextResponse } from "next/server";
import { Webhook } from "standardwebhooks";
import { addEvent } from "@/lib/event-store";

const WEBHOOK_SECRET = process.env.WEBHOOK_SECRET || "whsec_MfKQ9r8GKYqrTwjUPD8ILPZIo2LaLaSw";

export async function POST(request: NextRequest) {
  const wh = new Webhook(WEBHOOK_SECRET);

  // Get raw body as text for signature verification
  const body = await request.text();

  // Get webhook headers
  const headers = {
    "webhook-id": request.headers.get("webhook-id") || "",
    "webhook-timestamp": request.headers.get("webhook-timestamp") || "",
    "webhook-signature": request.headers.get("webhook-signature") || "",
  };

  try {
    // Verify the webhook signature
    const payload = wh.verify(body, headers) as { type: string; data: Record<string, unknown> };

    // Store the event
    const event = addEvent({
      id: headers["webhook-id"],
      type: payload.type,
      data: payload.data,
    });

    console.log(`Received webhook event: type=${event.type}, id=${event.id}`);

    return NextResponse.json({
      success: true,
      message: `Webhook event '${event.type}' processed successfully`,
    });
  } catch (error) {
    console.error("Webhook verification failed:", error);
    return NextResponse.json(
      { error: "Invalid webhook signature" },
      { status: 401 }
    );
  }
}
