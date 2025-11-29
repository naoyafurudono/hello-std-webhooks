"use client";

import { useEffect, useState } from "react";

interface WebhookHeaders {
  id: string;
  timestamp: string;
  signature: string;
}

interface WebhookEvent {
  headers: WebhookHeaders;
  payload: Record<string, unknown>;
  receivedAt: string;
}

export default function EventsPage() {
  const [events, setEvents] = useState<WebhookEvent[]>([]);
  const [loading, setLoading] = useState(true);

  const fetchEvents = async () => {
    try {
      const res = await fetch("/api/events");
      const data = await res.json();
      setEvents(data);
    } catch (error) {
      console.error("Failed to fetch events:", error);
    } finally {
      setLoading(false);
    }
  };

  const clearEvents = async () => {
    try {
      await fetch("/api/events", { method: "DELETE" });
      setEvents([]);
    } catch (error) {
      console.error("Failed to clear events:", error);
    }
  };

  useEffect(() => {
    fetchEvents();
    // Poll for new events every 2 seconds
    const interval = setInterval(fetchEvents, 2000);
    return () => clearInterval(interval);
  }, []);

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-100 p-8">
        <div className="max-w-4xl mx-auto">
          <p className="text-gray-600">Loading...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-100 p-8">
      <div className="max-w-4xl mx-auto">
        <div className="flex justify-between items-center mb-6">
          <h1 className="text-2xl font-bold text-gray-900">Webhook Events</h1>
          <div className="space-x-2">
            <button
              onClick={fetchEvents}
              className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
            >
              Refresh
            </button>
            <button
              onClick={clearEvents}
              className="px-4 py-2 bg-red-500 text-white rounded hover:bg-red-600"
            >
              Clear
            </button>
          </div>
        </div>

        {events.length === 0 ? (
          <div className="bg-white rounded-lg shadow p-6">
            <p className="text-gray-500 text-center">No events received yet.</p>
            <p className="text-gray-400 text-center text-sm mt-2">
              Send a webhook to /api/webhook to see events here.
            </p>
          </div>
        ) : (
          <div className="space-y-4">
            {events.map((event) => (
              <div
                key={`${event.headers.id}-${event.receivedAt}`}
                className="bg-white rounded-lg shadow p-4"
              >
                <div className="flex justify-between items-start mb-3">
                  <span className="text-sm text-gray-500">
                    {new Date(event.receivedAt).toLocaleString()}
                  </span>
                </div>

                <div className="mb-3">
                  <h3 className="text-sm font-semibold text-gray-700 mb-2">Headers</h3>
                  <div className="bg-gray-50 p-3 rounded text-sm font-mono space-y-1">
                    <div><span className="text-gray-500">webhook-id:</span> {event.headers.id}</div>
                    <div><span className="text-gray-500">webhook-timestamp:</span> {event.headers.timestamp}</div>
                    <div className="break-all"><span className="text-gray-500">webhook-signature:</span> {event.headers.signature}</div>
                  </div>
                </div>

                <div>
                  <h3 className="text-sm font-semibold text-gray-700 mb-2">Payload</h3>
                  <pre className="bg-gray-50 p-3 rounded text-sm overflow-x-auto">
                    {JSON.stringify(event.payload, null, 2)}
                  </pre>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
