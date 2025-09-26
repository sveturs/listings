'use client';

import { useEffect, useState } from 'react';
import configManager from '@/config';

export default function TestAuthConfig() {
  const [authServiceUrl, setAuthServiceUrl] = useState<string>('');
  const [apiUrl, setApiUrl] = useState<string>('');

  useEffect(() => {
    // Получаем URLs из конфигурации
    setAuthServiceUrl(configManager.getAuthServiceUrl());
    setApiUrl(configManager.getApiUrl());
  }, []);

  return (
    <div className="container mx-auto p-8">
      <h1 className="text-2xl font-bold mb-4">Auth Configuration Test</h1>

      <div className="bg-base-200 p-4 rounded-lg space-y-2">
        <div>
          <strong>Auth Service URL:</strong>
          <code className="ml-2 bg-base-300 px-2 py-1 rounded">
            {authServiceUrl || 'Loading...'}
          </code>
        </div>

        <div>
          <strong>API URL:</strong>
          <code className="ml-2 bg-base-300 px-2 py-1 rounded">
            {apiUrl || 'Loading...'}
          </code>
        </div>

        <div>
          <strong>NEXT_PUBLIC_AUTH_SERVICE_URL (env):</strong>
          <code className="ml-2 bg-base-300 px-2 py-1 rounded">
            {process.env.NEXT_PUBLIC_AUTH_SERVICE_URL || 'Not set'}
          </code>
        </div>
      </div>

      <div className="mt-4">
        <button
          className="btn btn-primary"
          onClick={() => window.location.reload()}
        >
          Refresh Page
        </button>
      </div>
    </div>
  );
}
