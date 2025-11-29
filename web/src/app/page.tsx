import Link from "next/link";

export default function Home() {
  return (
    <div className="min-h-screen bg-gray-100 p-8">
      <div className="max-w-2xl mx-auto">
        <h1 className="text-3xl font-bold text-gray-900 mb-6">
          Webhook Server
        </h1>
        <p className="text-gray-600 mb-8">
          A Next.js webhook server with standard-webhooks signature verification.
        </p>

        <div className="bg-white rounded-lg shadow p-6 mb-6">
          <h2 className="text-xl font-semibold text-gray-800 mb-4">Endpoints</h2>

          {/* POST /api/webhook */}
          <div className="border-b pb-6 mb-6">
            <code className="bg-blue-100 px-2 py-1 rounded text-sm font-bold">
              POST /api/webhook
            </code>
            <p className="text-gray-600 text-sm mt-2 mb-4">
              Receive webhook events with standard-webhooks signature verification.
            </p>

            <h3 className="text-sm font-semibold text-gray-700 mb-2">Required Headers</h3>
            <table className="w-full text-sm mb-4">
              <thead>
                <tr className="border-b">
                  <th className="text-left py-2 text-gray-600">Header</th>
                  <th className="text-left py-2 text-gray-600">Description</th>
                </tr>
              </thead>
              <tbody className="font-mono">
                <tr className="border-b">
                  <td className="py-2"><code className="bg-gray-100 px-1 rounded">webhook-id</code></td>
                  <td className="py-2 font-sans text-gray-600">Unique message identifier (e.g., <code className="bg-gray-100 px-1 rounded">msg_abc123</code>)</td>
                </tr>
                <tr className="border-b">
                  <td className="py-2"><code className="bg-gray-100 px-1 rounded">webhook-timestamp</code></td>
                  <td className="py-2 font-sans text-gray-600">Unix timestamp in seconds (e.g., <code className="bg-gray-100 px-1 rounded">1732000000</code>)</td>
                </tr>
                <tr className="border-b">
                  <td className="py-2"><code className="bg-gray-100 px-1 rounded">webhook-signature</code></td>
                  <td className="py-2 font-sans text-gray-600">Signature: <code className="bg-gray-100 px-1 rounded">v1,&lt;base64-hmac-sha256&gt;</code></td>
                </tr>
              </tbody>
            </table>

            <h3 className="text-sm font-semibold text-gray-700 mb-2">Signature Calculation</h3>
            <div className="bg-gray-50 p-3 rounded text-sm mb-4">
              <p className="mb-2 text-gray-600">1. Construct the signed content:</p>
              <pre className="bg-gray-100 p-2 rounded mb-3 overflow-x-auto">
{`signed_content = webhook_id + "." + webhook_timestamp + "." + body`}
              </pre>
              <p className="mb-2 text-gray-600">2. Compute HMAC-SHA256 with the secret (base64-decode the secret first):</p>
              <pre className="bg-gray-100 p-2 rounded mb-3 overflow-x-auto">
{`signature = base64(HMAC-SHA256(base64_decode(secret), signed_content))`}
              </pre>
              <p className="text-gray-600">3. Prefix with version: <code className="bg-gray-100 px-1 rounded">v1,&lt;signature&gt;</code></p>
            </div>

            <h3 className="text-sm font-semibold text-gray-700 mb-2">Request Body</h3>
            <p className="text-gray-600 text-sm mb-2">JSON payload (any structure):</p>
            <pre className="bg-gray-50 p-3 rounded text-sm mb-4 overflow-x-auto">
{`{
  "type": "user.created",
  "data": {
    "id": "user_123",
    "email": "user@example.com"
  }
}`}
            </pre>

            <h3 className="text-sm font-semibold text-gray-700 mb-2">Response</h3>
            <div className="space-y-2 text-sm">
              <div className="flex items-start gap-2">
                <span className="bg-green-100 text-green-800 px-2 py-0.5 rounded font-mono text-xs">200</span>
                <span className="text-gray-600">Signature verified successfully</span>
              </div>
              <div className="flex items-start gap-2">
                <span className="bg-red-100 text-red-800 px-2 py-0.5 rounded font-mono text-xs">401</span>
                <span className="text-gray-600">Invalid signature</span>
              </div>
            </div>
          </div>

          {/* GET /api/events */}
          <div className="pb-4">
            <code className="bg-green-100 px-2 py-1 rounded text-sm font-bold">
              GET /api/events
            </code>
            <p className="text-gray-600 text-sm mt-2 mb-4">
              List all received webhook events (both verified and failed).
            </p>

            <h3 className="text-sm font-semibold text-gray-700 mb-2">Response</h3>
            <pre className="bg-gray-50 p-3 rounded text-sm overflow-x-auto">
{`[
  {
    "headers": {
      "id": "msg_abc123",
      "timestamp": "1732000000",
      "signature": "v1,..."
    },
    "payload": { ... },
    "rawBody": "...",
    "verified": true,
    "receivedAt": "2024-11-19T..."
  }
]`}
            </pre>
          </div>

          {/* DELETE /api/events */}
          <div className="pt-4 border-t">
            <code className="bg-red-100 px-2 py-1 rounded text-sm font-bold">
              DELETE /api/events
            </code>
            <p className="text-gray-600 text-sm mt-2">
              Clear all stored events.
            </p>
          </div>
        </div>

        <div className="flex gap-4 items-center">
          <Link
            href="/events"
            className="inline-block px-6 py-3 bg-blue-500 text-white rounded-lg hover:bg-blue-600"
          >
            View Events
          </Link>
          <a
            href="https://github.com/naoyafurudono/hello-std-webhooks"
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center gap-2 px-4 py-3 text-gray-700 hover:text-gray-900"
          >
            <svg
              className="w-5 h-5"
              fill="currentColor"
              viewBox="0 0 24 24"
              aria-hidden="true"
            >
              <path
                fillRule="evenodd"
                d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z"
                clipRule="evenodd"
              />
            </svg>
            GitHub
          </a>
        </div>
      </div>
    </div>
  );
}
