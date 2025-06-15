'use client';

import { useEffect } from 'react';
import { marketplaceApi } from '@/services/api/endpoints';
import { useApi } from '@/hooks/useApi';

export default function TestApiClientPage() {
  const { data, loading, error, execute } = useApi(() =>
    marketplaceApi.getListings({ limit: 5 })
  );

  useEffect(() => {
    execute();
  }, [execute]);

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">API Test - Client Component</h1>

      <div className="bg-base-200 p-4 rounded">
        <h2 className="text-lg font-semibold mb-2">Response:</h2>

        {loading && <div className="text-info">Loading...</div>}

        {error && <div className="text-error">Error: {error.message}</div>}

        {data && (
          <div>
            <pre className="bg-base-300 p-2 rounded mt-2 overflow-auto">
              {JSON.stringify(data, null, 2)}
            </pre>
          </div>
        )}
      </div>

      <button
        onClick={() => execute()}
        disabled={loading}
        className="btn btn-primary mt-4"
      >
        {loading ? 'Loading...' : 'Reload'}
      </button>
    </div>
  );
}
