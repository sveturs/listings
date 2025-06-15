import { marketplaceApi } from '@/services/api/endpoints';

export default async function TestApiServerPage() {
  // Server Component - автоматически использует внутренний URL
  const response = await marketplaceApi.getListings({ limit: 5 });

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">API Test - Server Component</h1>

      <div className="bg-base-200 p-4 rounded">
        <h2 className="text-lg font-semibold mb-2">Response:</h2>
        {response.error ? (
          <div className="text-error">Error: {response.error.message}</div>
        ) : (
          <div>
            <p>Status: {response.status}</p>
            <pre className="bg-base-300 p-2 rounded mt-2 overflow-auto">
              {JSON.stringify(response.data, null, 2)}
            </pre>
          </div>
        )}
      </div>
    </div>
  );
}
