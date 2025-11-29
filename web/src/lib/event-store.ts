export interface WebhookEvent {
  id: string;
  type: string;
  data: Record<string, unknown>;
  receivedAt: string;
}

// In-memory store for webhook events
const events: WebhookEvent[] = [];
const MAX_EVENTS = 100;

export function addEvent(event: Omit<WebhookEvent, "receivedAt">): WebhookEvent {
  const storedEvent: WebhookEvent = {
    ...event,
    receivedAt: new Date().toISOString(),
  };

  events.unshift(storedEvent);

  // Keep only the most recent events
  if (events.length > MAX_EVENTS) {
    events.pop();
  }

  return storedEvent;
}

export function getEvents(): WebhookEvent[] {
  return [...events];
}

export function clearEvents(): void {
  events.length = 0;
}
