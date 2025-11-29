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
          <ul className="space-y-3">
            <li>
              <code className="bg-gray-100 px-2 py-1 rounded text-sm">
                POST /api/webhook
              </code>
              <p className="text-gray-600 text-sm mt-1">
                Receive webhook events with signature verification
              </p>
            </li>
            <li>
              <code className="bg-gray-100 px-2 py-1 rounded text-sm">
                GET /api/events
              </code>
              <p className="text-gray-600 text-sm mt-1">
                List all received events
              </p>
            </li>
          </ul>
        </div>

        <Link
          href="/events"
          className="inline-block px-6 py-3 bg-blue-500 text-white rounded-lg hover:bg-blue-600"
        >
          View Events
        </Link>
      </div>
    </div>
  );
}
