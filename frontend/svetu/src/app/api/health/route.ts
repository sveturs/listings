import { NextResponse } from 'next/server';
import configManager from '@/config';

export async function GET() {
  const config = configManager.getConfig();
  const errors = configManager.getValidationErrors();

  // Проверяем критичные сервисы
  const checks = {
    api: await checkApiHealth(config.api.url),
    storage: await checkStorageHealth(config.storage.minioUrl),
    config: errors.length === 0,
  };

  const isHealthy = Object.values(checks).every((v) => v === true);

  return NextResponse.json(
    {
      status: isHealthy ? 'healthy' : 'unhealthy',
      timestamp: new Date().toISOString(),
      checks,
      config: {
        environment: config.env.isProduction ? 'production' : 'development',
        features: config.features,
      },
      errors: errors.length > 0 ? errors : undefined,
    },
    {
      status: isHealthy ? 200 : 503,
    }
  );
}

async function checkApiHealth(url: string): Promise<boolean> {
  try {
    const response = await fetch(`${url}/health`, {
      method: 'GET',
      signal: AbortSignal.timeout(5000),
    });
    return response.ok;
  } catch {
    return false;
  }
}

async function checkStorageHealth(url: string): Promise<boolean> {
  try {
    const response = await fetch(`${url}/minio/health/live`, {
      method: 'GET',
      signal: AbortSignal.timeout(5000),
    });
    return response.ok;
  } catch {
    return false;
  }
}
