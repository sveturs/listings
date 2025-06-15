'use client';

import { env } from 'next-runtime-env';
import { publicEnv } from '@/utils/env';
import { useEffect, useState } from 'react';

export default function TestEnvPage() {
  const [windowEnv, setWindowEnv] = useState<any>(null);

  useEffect(() => {
    // Check window.__ENV after hydration
    if (typeof window !== 'undefined') {
      setWindowEnv((window as any).__ENV);
    }
  }, []);

  return (
    <div className="p-8">
      <h1 className="text-2xl font-bold mb-4">Runtime Environment Test</h1>

      <div className="space-y-2">
        <h2 className="text-xl font-semibold">Direct access:</h2>
        <p>API URL: {env('NEXT_PUBLIC_API_URL')}</p>
        <p>Minio URL: {env('NEXT_PUBLIC_MINIO_URL')}</p>

        <h2 className="text-xl font-semibold mt-4">Typed access:</h2>
        <p>API URL: {publicEnv.API_URL}</p>
        <p>Minio URL: {publicEnv.MINIO_URL}</p>
        <p>Payments enabled: {publicEnv.ENABLE_PAYMENTS ? 'Yes' : 'No'}</p>

        <h2 className="text-xl font-semibold mt-4">window.__ENV object:</h2>
        <pre className="bg-gray-100 p-2 rounded overflow-auto">
          {windowEnv ? JSON.stringify(windowEnv, null, 2) : 'Loading...'}
        </pre>

        <h2 className="text-xl font-semibold mt-4">
          Check for non-public vars:
        </h2>
        <p>
          SECRET_KEY in window.__ENV:{' '}
          {windowEnv && windowEnv.SECRET_KEY ? 'LEAKED!' : 'Not present (good)'}
        </p>
      </div>
    </div>
  );
}
