'use client';

import { useConfig } from '@/hooks/useConfig';
import configManager from '@/config';

export default function TestConfigPage() {
  const config = useConfig();
  const errors = configManager.getValidationErrors();

  return (
    <div className="p-8">
      <h1 className="text-2xl font-bold mb-4">Configuration Test</h1>

      {errors.length > 0 && (
        <div className="bg-red-100 p-4 mb-4 rounded">
          <h2 className="text-red-800 font-bold">Validation Errors:</h2>
          <ul className="list-disc pl-5">
            {errors.map((error, idx) => (
              <li key={idx} className="text-red-600">
                {error.field}: {error.message}
              </li>
            ))}
          </ul>
        </div>
      )}

      <div className="space-y-4">
        <section>
          <h2 className="text-xl font-semibold">API Configuration</h2>
          <pre className="bg-gray-100 p-2 rounded">
            {JSON.stringify(config.api, null, 2)}
          </pre>
        </section>

        <section>
          <h2 className="text-xl font-semibold">Features</h2>
          <pre className="bg-gray-100 p-2 rounded">
            {JSON.stringify(config.features, null, 2)}
          </pre>
        </section>

        <section>
          <h2 className="text-xl font-semibold">Environment</h2>
          <pre className="bg-gray-100 p-2 rounded">
            {JSON.stringify(config.env, null, 2)}
          </pre>
        </section>

        <section>
          <h2 className="text-xl font-semibold">Storage Configuration</h2>
          <pre className="bg-gray-100 p-2 rounded">
            {JSON.stringify(
              {
                minioUrl: config.storage.minioUrl,
                imagePathPattern: config.storage.imagePathPattern,
                imageHostsCount: config.storage.imageHosts.length,
              },
              null,
              2
            )}
          </pre>
        </section>
      </div>
    </div>
  );
}
